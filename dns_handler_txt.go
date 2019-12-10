package main

import (
	"github.com/miekg/dns"
	"log"
)

func handleTXTVersionRequest(this *handler, r *dns.Msg, msg *dns.Msg) {
	log.Printf("TXT %s\n", msg.Question[0].Name)

	msg.Answer = append(msg.Answer, &dns.TXT{
		Hdr: dns.RR_Header{Name: msg.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: dns.ClassINET, Ttl: conf.DefaultSOARecord.TTL},
		Txt: []string{getVersionFullString()},
	})
}
