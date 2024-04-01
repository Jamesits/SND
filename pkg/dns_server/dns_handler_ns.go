package dns_server

import (
	"github.com/miekg/dns"
	"log"
)

func handleNS(handler *Handler, r, msg *dns.Msg) {
	if handler.config.Debug {
		log.Printf("NS %s\n", msg.Question[0].Name)
	}

	// TODO: check if domain exists
	// same for root zone
	for _, ns := range handler.config.DefaultNSes {
		msg.Answer = append(msg.Answer, &dns.NS{
			Hdr: dns.RR_Header{Name: msg.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: r.Question[0].Qclass, Ttl: handler.config.DefaultSOARecord.TTL},
			Ns:  *ns,
		})
	}
}
