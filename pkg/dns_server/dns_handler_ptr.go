package dns_server

import (
	"fmt"
	"github.com/jamesits/snd/pkg/config"
	"log"
	"net"
	"strings"

	"github.com/miekg/dns"
)

func handlePTR(handler *Handler, r, msg *dns.Msg) {
	msg.Question[0].Name = strings.ToLower(msg.Question[0].Name)
	nameBreakout := strings.Split(msg.Question[0].Name, ".")
	index := len(nameBreakout) - 1

	// sanity check
	if index < 3 || nameBreakout[index] != "" || nameBreakout[index-1] != "arpa" {
		if handler.config.Debug {
			log.Printf("PTR %s not rational\n", msg.Question[0].Name)
		}
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
		if handler.config.Debug {
			log.Printf("PTR %s unable to parse IP address\n", msg.Question[0].Name)
		}
		return
	}
	ipaddr := net.ParseIP(strings.TrimRight(b.String(), split))
	if split == "." {
		ipaddr = ipaddr.To4()
	}

	// find a matching Config
	// TODO: optimize to less then O(n)
	found := false
	for _, netBlock := range handler.config.PerNetConfigs {
		if netBlock.IPNet.Contains(ipaddr) {
			found = true

			// construct ptr
			var p strings.Builder
			if netBlock.DomainPrefix != nil {
				p.WriteString(*netBlock.DomainPrefix)
			}

			switch netBlock.PtrGenerationMode {
			case config.Fixed:
				p.WriteString(*netBlock.Domain)
			case config.PrependLeftToRight:
				p.WriteString(IPToArpaDomain(ipaddr, false, netBlock.IPv6NotationMode))
				p.WriteString(".")
				p.WriteString(*netBlock.Domain)
			case config.PrependRightToLeft:
				p.WriteString(IPToArpaDomain(ipaddr, true, netBlock.IPv6NotationMode))
				p.WriteString(".")
				p.WriteString(*netBlock.Domain)
			case config.PrependLeftToRightDash:
				IPGenerate := IPToArpaDomain(ipaddr, false, netBlock.IPv6NotationMode)
				p.WriteString(strings.Replace(IPGenerate, ".", "-", -1))
				p.WriteString(".")
				p.WriteString(*netBlock.Domain)
			case config.PrependRightToLeftDash:
				IPGenerate := IPToArpaDomain(ipaddr, true, netBlock.IPv6NotationMode)
				p.WriteString(strings.Replace(IPGenerate, ".", "-", -1))
				p.WriteString(".")
				p.WriteString(*netBlock.Domain)
			case config.PrependRightToLeftOnlyip:
				IPGenerate := IPToArpaDomain(ipaddr, true, netBlock.IPv6NotationMode)
				p.WriteString(strings.Replace(IPGenerate, ".", "", -1))
				p.WriteString(".")
				p.WriteString(*netBlock.Domain)
			case config.PrependLeftToRightOnlyip:
				IPGenerate := IPToArpaDomain(ipaddr, false, netBlock.IPv6NotationMode)
				p.WriteString(strings.Replace(IPGenerate, ".", "", -1))
				p.WriteString(".")
				p.WriteString(*netBlock.Domain)
			default:
				return
			}

			if handler.config.Debug {
				log.Printf("PTR %s => %s", ipaddr.String(), p.String())
			}

			// generate an answer
			msg.Answer = append(msg.Answer, &dns.PTR{
				Hdr: dns.RR_Header{Name: msg.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: r.Question[0].Qclass, Ttl: netBlock.TTL},
				Ptr: p.String(),
			})
			break
		}
	}

	if !found && handler.config.Debug {
		log.Printf("PTR %s unknown net", ipaddr.String())
	}
}

func IPToArpaDomain(ip net.IP, reverse bool, ipv6ConversionMode config.IPv6NotationMode) string {
	ipv4 := ip.To4()
	var ret []string

	if ipv4 == nil {
		// ipv6
		// TODO: edge cases?
		for _, elem := range ip {
			s := fmt.Sprintf("%x", elem) // 2 characters per iteration
			if len(s) == 2 {
				ret = append(ret, s[0:1], s[1:2])
			} else {
				ret = append(ret, "0", s[0:1])
			}
		}
	} else {
		// ipv4
		for _, elem := range ipv4 {
			ret = append(ret, fmt.Sprintf("%d", elem))
		}
	}

	switch ipv6ConversionMode {
	case config.ArpaNotation:
		break
	case config.FourHexsNotation:
		reverse = !reverse // in this mode, ret is processed in reverse, so we need to reverse it again before returning
		var ret2 []string
		for i := len(ret) - 1; i >= 0; i -= 4 {
			var b strings.Builder
			var isLeadingZero = true
			for j := 3; j >= 0; j-- {
				if i-j < 0 || i-j > len(ret) {
					log.Panicf("Assertion for IP length can be divided in 4 failed")
				}
				// noinspection GoNilness
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
