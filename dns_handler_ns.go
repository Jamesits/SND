package main

import (
	"github.com/miekg/dns"
	"log"
)

func handleNS(this *handler, r *dns.Msg, msg *dns.Msg) {
	log.Printf("NS %s\n", msg.Question[0].Name)

	// TODO: check if domain exists
	// same for root zone
	for _, ns := range conf.DefaultNSes {
		msg.Answer = append(msg.Answer, &dns.NS{
			Hdr: dns.RR_Header{Name: msg.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: r.Question[0].Qclass, Ttl: conf.DefaultSOARecord.TTL},
			Ns:  *ns,
		})
	}
}
