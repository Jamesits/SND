package dns_server

import (
	"github.com/miekg/dns"
	"log"
)

// simply replies NOTIMPL
func handleDefault(handler *Handler, r, msg *dns.Msg) {
	if handler.config.Debug {
		log.Printf("%d %s not implemented\n", msg.Question[0].Qtype, msg.Question[0].Name)
	}
	msg.Rcode = dns.RcodeNotImplemented
}
