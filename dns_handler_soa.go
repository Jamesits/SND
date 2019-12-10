package main

import (
	"github.com/miekg/dns"
	"log"
)

func handleSOA(this *handler, r, msg *dns.Msg) {
	log.Printf("SOA %s\n", msg.Question[0].Name)

	// TODO: check if domain exists, and use the actual SOA for that domain
	// if not our case, we should reply with authority section containing root zone
	// See: http://www-inf.int-evry.fr/~hennequi/CoursDNS/NOTES-COURS_eng/msg.html
	// Get root zone from https://www.internic.net/domain/named.root
	msg.Answer = append(msg.Answer, &dns.SOA{
		Hdr:     dns.RR_Header{Name: msg.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: r.Question[0].Qclass, Ttl: conf.DefaultSOARecord.TTL},
		Ns:      *conf.DefaultSOARecord.MName,
		Mbox:    *conf.DefaultSOARecord.RName,
		Serial:  conf.DefaultSOARecord.Serial,
		Refresh: conf.DefaultSOARecord.Refresh,
		Retry:   conf.DefaultSOARecord.Retry,
		Expire:  conf.DefaultSOARecord.Expire,
		Minttl:  conf.DefaultSOARecord.TTL,
	})
}
