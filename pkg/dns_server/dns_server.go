package dns_server

import (
	"github.com/jamesits/libiferr/exception"
	"github.com/jamesits/snd/pkg/config"
	"github.com/miekg/dns"
	"log"
)

// ListenSync listens on a specific endpoint
// proto can be "udp" or "tcp"
// endpoint is "ip.address:port"
func ListenSync(config *config.Config, proto, endpoint string) {
	log.Printf("Listening on %s %s", proto, endpoint)
	srv := &dns.Server{Addr: endpoint, Net: proto}
	srv.Handler = NewHandler(config)
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("Failed to set %s listener %s\n", proto, endpoint)
		exception.HardFailWithReason("failed to enable listener", err)
	}
}
