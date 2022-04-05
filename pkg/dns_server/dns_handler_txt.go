package dns_server

import (
	"github.com/jamesits/snd/pkg/version"
	"github.com/miekg/dns"
	"log"
)

// replies a TXT record containing server name and version
func handleTXTVersionRequest(this *handler, r, msg *dns.Msg) {
	log.Printf("TXT %s\n", msg.Question[0].Name)

	if !this.config.AllowVersionReporting {
		msg.Rcode = dns.RcodeRefused
		return
	}

	var versionString string
	if len(this.config.OverrideVersionString) == 0 {
		versionString = version.GetVersionFullString()
	} else {
		versionString = this.config.OverrideVersionString
	}

	msg.Answer = append(msg.Answer, &dns.TXT{
		Hdr: dns.RR_Header{Name: msg.Question[0].Name, Rrtype: r.Question[0].Qtype, Class: r.Question[0].Qclass, Ttl: this.config.DefaultSOARecord.TTL},
		Txt: []string{versionString},
	})
}
