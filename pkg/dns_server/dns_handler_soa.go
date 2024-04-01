package dns_server

import (
	"github.com/miekg/dns"
	"log"
)

func handleSOA(this *Handler, r, msg *dns.Msg) {
	log.Printf("SOA %s\n", msg.Question[0].Name)

	// TODO: check if domain exists, and use the actual SOA for that domain
	// if not our case, we should reply with authority section containing root zone
	// See: http://www-inf.int-evry.fr/~hennequi/CoursDNS/NOTES-COURS_eng/msg.html
	// Get root zone from https://www.internic.net/domain/named.root
	msg.Answer = append(msg.Answer, &dns.SOA{
		Hdr:     dns.RR_Header{Name: msg.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: r.Question[0].Qclass, Ttl: this.config.DefaultSOARecord.TTL},
		Ns:      *this.config.DefaultSOARecord.MName,
		Mbox:    *this.config.DefaultSOARecord.RName,
		Serial:  this.config.DefaultSOARecord.Serial,
		Refresh: this.config.DefaultSOARecord.Refresh,
		Retry:   this.config.DefaultSOARecord.Retry,
		Expire:  this.config.DefaultSOARecord.Expire,
		Minttl:  this.config.DefaultSOARecord.TTL,
	})
}
