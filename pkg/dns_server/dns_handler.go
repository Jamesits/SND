package dns_server

import (
	"github.com/jamesits/libiferr/exception"
	"github.com/jamesits/snd/pkg/config"
	"github.com/miekg/dns"
	"strings"
)

type Handler struct {
	config *config.Config
}

func NewHandler(config *config.Config) *Handler {
	return &Handler{config: config}
}

func (handler *Handler) newDNSReplyMsg() *dns.Msg {
	msg := dns.Msg{}

	msg.Compress = handler.config.CompressDNSMessages

	// this is an authoritative DNS server
	msg.Authoritative = true
	msg.RecursionAvailable = false

	// DNSSEC disabled for now
	// TODO: fix DNSSEC
	msg.AuthenticatedData = false
	msg.CheckingDisabled = true

	return &msg
}

// send out the generated answer, and if the answer is not correct, send out a SERVFAIL
func (handler *Handler) finishAnswer(w *dns.ResponseWriter, r *dns.Msg) {
	err := (*w).WriteMsg(r)
	if err != nil {
		exception.SoftFailWithReason("failed to send primary DNS answer", err)

		// if answer sanity check (miekg/dns automatically does handler) fails, reply with SERVFAIL
		msg := handler.newDNSReplyMsg()
		msg.SetReply(r)
		msg.Rcode = dns.RcodeServerFailure
		err = (*w).WriteMsg(msg)
		exception.SoftFailWithReason("failed to send secondary DNS answer", err)
	}
}

// ServeDNS TODO: force TCP for 1) clients which requests too fast; 2) non-existent answers
// See: https://labs.apnic.net/?p=382
func (handler *Handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := handler.newDNSReplyMsg()
	msg.SetReply(r)

	// on function return, we send out the current answer
	defer handler.finishAnswer(&w, msg)

	// sanity check
	if len(r.Question) != 1 {
		msg.Rcode = dns.RcodeRefused
		return
	}

	switch r.Question[0].Qclass {
	case dns.ClassINET:
		switch r.Question[0].Qtype {
		case dns.TypeSOA:
			handleSOA(handler, r, msg)
			return

		case dns.TypeNS:
			handleNS(handler, r, msg)
			return

		case dns.TypePTR:
			handlePTR(handler, r, msg)
			return

		default:
			handleDefault(handler, r, msg)
			return
		}
	case dns.ClassCHAOS:
		switch r.Question[0].Qtype {
		case dns.TypeTXT:
			if strings.EqualFold(r.Question[0].Name, "version.bind.") {
				// we need to reply our software version
				// https://serverfault.com/questions/517087/dns-how-to-find-out-which-software-a-remote-dns-server-is-running
				handleTXTVersionRequest(handler, r, msg)
			} else {
				handleDefault(handler, r, msg)
			}
			return

		default:
			handleDefault(handler, r, msg)
			return
		}
	}
}
