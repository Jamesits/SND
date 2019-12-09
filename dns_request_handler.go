package main

import (
	"fmt"
	"github.com/miekg/dns"
	"log"
	"net"
	"strings"
)

func IPToArpaDomain(ip net.IP, reverse bool, ipv6ConversionMode IPv6NotationMode) string {
	ipv4 := ip.To4()
	var ret []string

	if ipv4 == nil {
		// ipv6
		// TODO: edge cases?
		for _, elem := range ip {
			s := fmt.Sprintf("%x", elem) // 2 characters per iteration
			if len(s) == 2 {
				ret = append(ret, s[0:1])
				ret = append(ret, s[1:2])
			} else {
				ret = append(ret, "0")
				ret = append(ret, s[0:1])
			}
		}
	} else {
		// ipv4
		for _, elem := range ipv4 {
			ret = append(ret, fmt.Sprintf("%d", elem))
		}
	}

	switch ipv6ConversionMode {
	case ARPA_NOTATION:
		break
	case FOUR_HEXS_NOTATION:
		reverse = !reverse // in this mode, ret is processed in reverse, so we need to reverse it again before returning
		var ret2 []string
		for i := len(ret) - 1; i >= 0; i -= 4 {
			var b strings.Builder
			var isLeadingZero = true
			for j := 3; j >= 0; j-- {
				if len(ret) < i-j || 0 > i-j {
					log.Panicf("Assertion for IP length can be divided in 4 failed")
				}
				//noinspection GoNilness
				if isLeadingZero {
					if ret[i-j] != "0" {
						isLeadingZero = false
						b.WriteString(ret[i-j])
					}
				} else {
					b.WriteString(ret[i-j])
				}
			}
			if b.Len() == 0 {
				b.WriteString("0")
			}

			ret2 = append(ret2, b.String())
		}

		ret = ret2
	default:
		break
	}

	if reverse {
		var ret2 []string
		for i := len(ret) - 1; i >= 0; i-- {
			ret2 = append(ret2, ret[i])
		}

		ret = ret2
	}

	return strings.Join(ret, ".")
}

type handler struct{}

func newDnsReplyMsg() *dns.Msg {
	msg := dns.Msg{}

	msg.Compress = conf.CompressDNSMessages

	// this is an authoritative DNS server
	msg.Authoritative = true
	msg.RecursionAvailable = false

	// DNSSEC disabled for now
	// TODO: fix DNSSEC
	msg.AuthenticatedData = false
	msg.CheckingDisabled = true

	return &msg
}

func finishAnswer(w *dns.ResponseWriter, r *dns.Msg) {
	err := (*w).WriteMsg(r)
	if err != nil {
		softFailIf(err)

		// reply with SERVFAIL
		msg := newDnsReplyMsg()
		msg.SetReply(r)
		msg.Rcode = dns.RcodeServerFailure
		err = (*w).WriteMsg(msg)
		softFailIf(err)
	}
}

func (this *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := newDnsReplyMsg()
	msg.SetReply(r)

	defer finishAnswer(&w, msg)

	// sanity check
	if len(r.Question) != 1 {
		msg.Rcode = dns.RcodeRefused
		return
	}

	switch r.Question[0].Qtype {
	case dns.TypeSOA:
		log.Printf("SOA %s\n", msg.Question[0].Name)

		// TODO: check if domain exists, and use the actual SOA for that domain
		// if not our case, we should reply with authority section containing root zone
		// See: http://www-inf.int-evry.fr/~hennequi/CoursDNS/NOTES-COURS_eng/msg.html
		// Get root zone from https://www.internic.net/domain/named.root
		msg.Answer = append(msg.Answer, &dns.SOA{
			Hdr:     dns.RR_Header{Name: msg.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: dns.ClassINET, Ttl: conf.DefaultSOARecord.TTL},
			Ns:      *conf.DefaultSOARecord.MName,
			Mbox:    *conf.DefaultSOARecord.RName,
			Serial:  conf.DefaultSOARecord.Serial,
			Refresh: conf.DefaultSOARecord.Refresh,
			Retry:   conf.DefaultSOARecord.Retry,
			Expire:  conf.DefaultSOARecord.Expire,
			Minttl:  conf.DefaultSOARecord.TTL,
		})

		return

	case dns.TypeNS:
		log.Printf("NS %s\n", msg.Question[0].Name)

		// TODO: check if domain exists
		// same for root zone
		for _, ns := range conf.DefaultNSes {
			msg.Answer = append(msg.Answer, &dns.NS{
				Hdr: dns.RR_Header{Name: msg.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: dns.ClassINET, Ttl: conf.DefaultSOARecord.TTL},
				Ns:  *ns,
			})
		}

		return

	case dns.TypePTR:
		nameBreakout := strings.Split(msg.Question[0].Name, ".")
		index := len(nameBreakout) - 1

		// sanity check
		if index < 3 || nameBreakout[index] != "" || nameBreakout[index-1] != "arpa" {
			log.Printf("PTR %s not rational\n", msg.Question[0].Name)
			return
		}

		var split string

		// parse IP address out of the request
		index -= 3
		var b strings.Builder
		switch nameBreakout[index+1] {
		case "in-addr": // IPv4
			split = "."
			for ; index >= 0; index-- {
				b.WriteString(nameBreakout[index])
				b.WriteString(split)
			}
		case "ip6": // IPv6
			split = ":"
			for i := 0; index >= 0; {
				b.WriteString(nameBreakout[index])
				index--
				i++
				if i%4 == 0 {
					b.WriteString(split)
				}
			}
		default:
			log.Printf("PTR %s unable to parse IP address\n", msg.Question[0].Name)
			return
		}
		ipaddr := net.ParseIP(strings.TrimRight(b.String(), split))
		if split == "." {
			ipaddr = ipaddr.To4()
		}

		// find a matching config
		// TODO: optimize to less then O(n)
		found := false
		for _, netBlock := range conf.PerNetConfigs {
			if netBlock.IPNet.Contains(ipaddr) {
				found = true

				// construct ptr
				var p strings.Builder
				if netBlock.DomainPrefix != nil {
					p.WriteString(*netBlock.DomainPrefix)
				}

				switch netBlock.PtrGenerationMode {
				case FIXED:
					p.WriteString(*netBlock.Domain)
				case PREPEND_LEFT_TO_RIGHT:
					p.WriteString(IPToArpaDomain(ipaddr, false, netBlock.IPv6NotationMode))
					p.WriteString(".")
					p.WriteString(*netBlock.Domain)
				case PREPEND_RIGHT_TO_LEFT:
					p.WriteString(IPToArpaDomain(ipaddr, true, netBlock.IPv6NotationMode))
					p.WriteString(".")
					p.WriteString(*netBlock.Domain)
				default:
					return
				}

				log.Printf("PTR %s => %s", msg.Question[0].Name, p.String())

				// generate an answer
				msg.Answer = append(msg.Answer, &dns.PTR{
					Hdr: dns.RR_Header{Name: msg.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: dns.ClassINET, Ttl: 1919},
					Ptr: p.String(),
				})
				break
			}
		}

		if !found {
			log.Printf("PTR %s unknown net", ipaddr.String())
		}

	default:
		msg.Rcode = dns.RcodeNotImplemented
		log.Printf("%d %s not implemented\n", msg.Question[0].Qtype, msg.Question[0].Name)
		return
	}
}
