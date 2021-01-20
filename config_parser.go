package main

import (
	"log"
	"net"
	"strings"
)

func ensureDotAtRight(s *string) *string {
	if !strings.HasSuffix(*s, ".") {
		ret := *s + "."
		return &ret
	} else {
		return s
	}
}

func ensureNoDotAtLeft(s *string) *string {
	if !strings.HasPrefix(*s, ".") {
		ret := strings.TrimLeft(*s, ".")
		return &ret
	} else {
		return s
	}
}

func SOARecordFillDefault(r *SOARecord, useDefaultRecord bool) {
	if useDefaultRecord {
		if len(*r.MName) == 0 {
			r.MName = conf.DefaultSOARecord.MName
		}

		if len(*r.RName) == 0 {
			r.RName = conf.DefaultSOARecord.RName
		}

		if r.Serial == 0 {
			r.Serial = conf.DefaultSOARecord.Serial
		}

		if r.Refresh == 0 {
			r.Refresh = conf.DefaultSOARecord.Refresh
		}

		if r.Retry == 0 {
			r.Retry = conf.DefaultSOARecord.Retry
		}

		if r.Expire == 0 {
			r.Expire = conf.DefaultSOARecord.Expire
		}

		if r.TTL == 0 {
			r.TTL = conf.DefaultSOARecord.TTL
		}
	} else {
		if r.Serial == 0 {
			r.Serial = 114514
		}

		if r.Refresh == 0 {
			r.Refresh = 86400
		}

		if r.Retry == 0 {
			r.Retry = 7200
		}

		if r.Expire == 0 {
			r.Expire = 3600000
		}

		if r.TTL == 0 {
			r.TTL = 172800
		}
	}

	if len(*r.MName) == 0 {
		log.Fatalf("a SOA record is missing MName")
	}

	if len(*r.RName) == 0 {
		log.Fatalf("a SOA record is missing RName")
	}

	r.RName = ensureDotAtRight(r.RName)
	r.MName = ensureDotAtRight(r.MName)
}

// fix config and fill in defaults
func fixConfig() {
	var err error

	if conf.DefaultTTL == 0 {
		conf.DefaultTTL = 114514
	}

	for index, ns := range conf.DefaultNSes {
		conf.DefaultNSes[index] = ensureDotAtRight(ns)
	}

	if conf.DefaultSOARecord == nil {
		conf.DefaultSOARecord = new(SOARecord)
	}
	SOARecordFillDefault(conf.DefaultSOARecord, false)

	// note that range is byVal so we use index here
	for _, currentConfig := range conf.PerNetConfigs {
		// fill IPNet
		_, currentConfig.IPNet, err = net.ParseCIDR(*currentConfig.IPNetString)
		hardFailIf(err)

		log.Printf("Loading network %s\n", currentConfig.IPNet.String())

		// fill Mode
		if currentConfig.PtrGenerationModeString == nil {
			log.Fatalf("Missing PTR generation method")
		}
		switch strings.ToLower(*currentConfig.PtrGenerationModeString) {
		case "fixed":
			currentConfig.PtrGenerationMode = FIXED
		case "prefix_ltr":
			currentConfig.PtrGenerationMode = PREPEND_LEFT_TO_RIGHT
		case "prefix_rtl":
			currentConfig.PtrGenerationMode = PREPEND_RIGHT_TO_LEFT
		case "prefix_ltr_dash":
			currentConfig.PtrGenerationMode = PREPEND_LEFT_TO_RIGHT_DASH
		case "prefix_rtl_dash":
			currentConfig.PtrGenerationMode = PREPEND_RIGHT_TO_LEFT_DASH
		default:
			log.Fatalf("Unknown mode \"%s\"", *currentConfig.PtrGenerationModeString)
		}

		// fill IPv6Notation
		if currentConfig.IPv6NotationString == nil {
			currentConfig.IPv6NotationMode = ARPA_NOTATION
		} else {
			switch strings.ToLower(*currentConfig.IPv6NotationString) {
			case "arpa":
				currentConfig.IPv6NotationMode = ARPA_NOTATION
			case "four_hexs":
				currentConfig.IPv6NotationMode = FOUR_HEXS_NOTATION
			default:
				log.Fatalf("Unknown ipv6_notation \"%s\"", *currentConfig.PtrGenerationModeString)
			}
		}

		// check domain
		currentConfig.Domain = ensureDotAtRight(currentConfig.Domain)
		currentConfig.Domain = ensureNoDotAtLeft(currentConfig.Domain)

		// fill TTL
		if currentConfig.TTL == 0 {
			currentConfig.TTL = conf.DefaultTTL
		}

		// fill SOA
		if currentConfig.SOARecord == nil {
			currentConfig.SOARecord = conf.DefaultSOARecord
		} else {
			SOARecordFillDefault(currentConfig.SOARecord, true)
		}
	}
}
