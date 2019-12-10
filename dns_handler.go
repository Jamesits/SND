package main

import (
	"github.com/miekg/dns"
	"strings"
)

type handler struct{}

func newDnsReplyMsg() *dns.Msg {
	msg := dns.Msg{}

	msg.Compress = conf.CompressDNSMessages

	// this is an authoritative DNS server
	msg.Authoritative = true
	msg.RecursionAvailable = false

	// DNSSEC disabled for now
	// TODO: fix DNSSEC
	msg.AuthenticatedData = false
	msg.CheckingDisabled = true

	return &msg
}

func finishAnswer(w *dns.ResponseWriter, r *dns.Msg) {
	err := (*w).WriteMsg(r)
	if err != nil {
		softFailIf(err)

		// if answer sanity check (miekg/dns automatically does this) fails, reply with SERVFAIL
		msg := newDnsReplyMsg()
		msg.SetReply(r)
		msg.Rcode = dns.RcodeServerFailure
		err = (*w).WriteMsg(msg)
		softFailIf(err)
	}
}

// TODO: force TCP for 1) clients which requests too fast; 2) non-existent answers
// See: https://labs.apnic.net/?p=382
func (this *handler) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	msg := newDnsReplyMsg()
	msg.SetReply(r)

	// on function return, we send out the current answer
	defer finishAnswer(&w, msg)

	// sanity check
	if len(r.Question) != 1 {
		msg.Rcode = dns.RcodeRefused
		return
	}

	switch r.Question[0].Qtype {
	case dns.TypeSOA:
		handleSOA(this, r, msg)
		return

	case dns.TypeNS:
		handleNS(this, r, msg)
		return

	case dns.TypePTR:
		handlePTR(this, r, msg)
		return

	case dns.TypeTXT:
		if r.Question[0].Qclass == dns.ClassCHAOS && strings.ToLower(r.Question[0].Name) == "version.bind." {
			// we need to reply our software version
			// https://serverfault.com/questions/517087/dns-how-to-find-out-which-software-a-remote-dns-server-is-running
			handleTXTVersionRequest(this, r, msg)
		} else {
			handleDefault(this, r, msg)
		}
		return

	default:
		handleDefault(this, r, msg)
		return
	}
}
