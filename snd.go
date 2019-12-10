package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/miekg/dns"
	"log"
	"strings"
	"time"
)

var conf *config
var configFilePath *string
var showVersionOnly *bool

// Listen on a specific endpoint
func listen(proto, endpoint string) {
	log.Printf("Listening on %s %s", proto, endpoint)
	srv := &dns.Server{Addr: endpoint, Net: proto}
	srv.Handler = &handler{}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Failed to set udp listener %s\n", err.Error())
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
		go listen(r[0], r[1])
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
