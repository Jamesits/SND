package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/jamesits/libiferr/exception"
	"github.com/jamesits/snd/pkg/config"
	"github.com/jamesits/snd/pkg/dns_server"
	"github.com/jamesits/snd/pkg/version"
	"log"
	"strings"
	"sync"
)

var conf *config.Config
var configFilePath *string
var showVersionOnly *bool
var mainThreadWaitGroup = &sync.WaitGroup{}

func main() {
	// parse flags
	var err error
	configFilePath = flag.String("config", "/etc/snd/config.toml", "config file")
	showVersionOnly = flag.Bool("version", false, "show version and quit")
	flag.Parse()

	if *showVersionOnly {
		fmt.Println(version.GetVersionFullString())
		return
	} else {
		log.Println(version.GetVersionFullString())
	}

	// parse config file
	conf = &config.Config{}
	metaData, err := toml.DecodeFile(*configFilePath, conf)
	exception.HardFailWithReason("failed to read the config file", err)

	// print unknown configs
	for _, key := range metaData.Undecoded() {
		log.Printf("Unknown key %q in the Config file, maybe a typo?", key.String())
	}

	// fix config and fill in defaults
	conf.FixConfig()

	// listen on all the configured listeners
	for _, elem := range conf.Listen {
		r := strings.SplitN(*elem, ":", 2)
		go func() {
			mainThreadWaitGroup.Add(1)
			defer mainThreadWaitGroup.Done()
			go dns_server.ListenSync(conf, r[0], r[1])
		}()
	}

	mainThreadWaitGroup.Wait()
}
