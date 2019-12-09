package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/miekg/dns"
	"log"
	"net"
	"strings"
	"time"
)

var conf *config
var configFilePath *string
var showVersionOnly *bool

func listen(proto string, endpoint string) {
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
		fmt.Printf("SND %d.%d.%d (Commit %s at %s)\n", versionMajor, versionMinor, versionRevision, versionGitCommitHash, versionCompileTime)
		return
	}

	// parse config file
	conf = &config{}
	metaData, err := toml.DecodeFile(*configFilePath, conf)
	softFailIf(err)

	// print unknown configs
	for _, key := range metaData.Undecoded() {
		log.Printf("Warning: unknown option %q", key.String())
	}

	// fix config
	if conf.DefaultTTL == 0 {
		conf.DefaultTTL = 114
	}

	// TODO: check NS records format (dot at the end)
	// TODO: check if DefaultSOARecord exists
	SOARecordFillDefault(conf.DefaultSOARecord, false)

	// note that range is byVal so we use index here
	for index := range conf.PerNetConfigs {
		// fill IPNet
		_, conf.PerNetConfigs[index].IPNet, err = net.ParseCIDR(conf.PerNetConfigs[index].IPNetString)
		hardFailIf(err)

		// fill Mode
		switch strings.ToLower(conf.PerNetConfigs[index].PtrGenerationModeString) {
		case "fixed":
			conf.PerNetConfigs[index].PtrGenerationMode = FIXED
		case "prefix_ltr":
			conf.PerNetConfigs[index].PtrGenerationMode = PREPEND_LEFT_TO_RIGHT
		case "prefix_rtl":
			conf.PerNetConfigs[index].PtrGenerationMode = PREPEND_RIGHT_TO_LEFT
		default:
			log.Printf("Unknown PTR generation \"%s\" for net \"%s\"", conf.PerNetConfigs[index].PtrGenerationModeString, conf.PerNetConfigs[index].IPNetString)
			conf.PerNetConfigs[index].PtrGenerationMode = FIXED
		}

		// fill IPv6Notation
		switch strings.ToLower(conf.PerNetConfigs[index].IPv6NotationString) {
		case "arpa":
			conf.PerNetConfigs[index].IPv6NotationMode = ARPA_NOTATION
		case "four_hexs":
			conf.PerNetConfigs[index].IPv6NotationMode = FOUR_HEXS_NOTATION
		default:
			conf.PerNetConfigs[index].IPv6NotationMode = ARPA_NOTATION
		}

		// check domain
		l := len(conf.PerNetConfigs[index].Domain)
		if conf.PerNetConfigs[index].Domain[l-1] != '.' {
			conf.PerNetConfigs[index].Domain += "."
		}
		conf.PerNetConfigs[index].Domain = strings.TrimLeft(conf.PerNetConfigs[index].Domain, ".")

		// fill TTL
		if conf.PerNetConfigs[index].TTL == 0 {
			conf.PerNetConfigs[index].TTL = conf.DefaultTTL
		}

		// fill SOA
		if conf.PerNetConfigs[index].SOARecord == nil {
			conf.PerNetConfigs[index].SOARecord = conf.DefaultSOARecord
		} else {
			SOARecordFillDefault(conf.PerNetConfigs[index].SOARecord, true)
		}
	}

	// listen them
	for _, elem := range conf.Listen {
		r := strings.SplitN(elem, ":", 2)
		go listen(r[0], r[1])
	}

	for {
		time.Sleep(1)
	}
}
