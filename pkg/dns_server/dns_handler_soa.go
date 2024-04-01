package dns_server

import (
	"github.com/miekg/dns"
	"log"
)

func handleSOA(handler *Handler, r, msg *dns.Msg) {
	if handler.config.Debug {
		log.Printf("SOA %s\n", msg.Question[0].Name)
	}

	// TODO: check if domain exists, and use the actual SOA for that domain
	// if not our case, we should reply with authority section containing root zone
	// See: http://www-inf.int-evry.fr/~hennequi/CoursDNS/NOTES-COURS_eng/msg.html
	// Get root zone from https://www.internic.net/domain/named.root
	msg.Answer = append(msg.Answer, &dns.SOA{
		Hdr:     dns.RR_Header{Name: msg.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: r.Question[0].Qclass, Ttl: handler.config.DefaultSOARecord.TTL},
		Ns:      *handler.config.DefaultSOARecord.MName,
		Mbox:    *handler.config.DefaultSOARecord.RName,
		Serial:  handler.config.DefaultSOARecord.Serial,
		Refresh: handler.config.DefaultSOARecord.Refresh,
		Retry:   handler.config.DefaultSOARecord.Retry,
		Expire:  handler.config.DefaultSOARecord.Expire,
		Minttl:  handler.config.DefaultSOARecord.TTL,
	})
}
