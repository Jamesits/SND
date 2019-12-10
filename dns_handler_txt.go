package main

import (
	"github.com/miekg/dns"
	"log"
)

// replies a TXT record containing server name and version
func handleTXTVersionRequest(this *handler, r, msg *dns.Msg) {
	log.Printf("TXT %s\n", msg.Question[0].Name)

	if !conf.AllowVersionReporting {
		msg.Rcode = dns.RcodeRefused
		return
	}

	var versionString string
	if len(conf.OverrideVersionString) == 0 {
		versionString = getVersionFullString()
	} else {
		versionString = conf.OverrideVersionString
	}

	msg.Answer = append(msg.Answer, &dns.TXT{
		Hdr: dns.RR_Header{Name: msg.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: r.Question[0].Qclass, Ttl: conf.DefaultSOARecord.TTL},
		Txt: []string{versionString},
	})
}
