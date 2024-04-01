package dns_server

import (
	"github.com/miekg/dns"
	"log"
)

func handleNS(this *Handler, r, msg *dns.Msg) {
	log.Printf("NS %s\n", msg.Question[0].Name)

	// TODO: check if domain exists
	// same for root zone
	for _, ns := range this.config.DefaultNSes {
		msg.Answer = append(msg.Answer, &dns.NS{
			Hdr: dns.RR_Header{Name: msg.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: r.Question[0].Qclass, Ttl: this.config.DefaultSOARecord.TTL},
			Ns:  *ns,
		})
	}
}
