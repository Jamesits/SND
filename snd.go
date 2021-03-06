package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/miekg/dns"
	"log"
	"strings"
	"sync"
)

var conf *config
var configFilePath *string
var showVersionOnly *bool
var mainThreadWaitGroup = &sync.WaitGroup{}

// Listen on a specific endpoint
// proto can be "udp" or "tcp"
// endpoint is "ip.address:port"
func listen(proto, endpoint string) {
	defer mainThreadWaitGroup.Done()
	log.Printf("Listening on %s %s", proto, endpoint)
	srv := &dns.Server{Addr: endpoint, Net: proto}
	srv.Handler = &handler{}
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("Failed to set %s listener %s\n", proto, endpoint)
		hardFailIf(err)
	}
}

func main() {
	// parse flags
	var err error
	configFilePath = flag.String("config", "/etc/snd/config.toml", "config directory")
	showVersionOnly = flag.Bool("version", false, "show version and quit")
	flag.Parse()

	if *showVersionOnly {
		fmt.Println(getVersionFullString())
		return
	} else {
		log.Println(getVersionFullString())
	}

	// parse config file
	conf = &config{}
	metaData, err := toml.DecodeFile(*configFilePath, conf)
	hardFailIf(err)

	// print unknown configs
	for _, key := range metaData.Undecoded() {
		log.Printf("Unknown key %q in the config file, maybe a typo?", key.String())
	}

	// fix config and fill in defaults
	fixConfig()

	// listen them
	for _, elem := range conf.Listen {
		r := strings.SplitN(*elem, ":", 2)
		mainThreadWaitGroup.Add(1)
		go listen(r[0], r[1])
	}

	mainThreadWaitGroup.Wait()
}
