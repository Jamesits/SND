package config

import (
	"github.com/jamesits/libiferr/exception"
	"log"
	"net"
	"slices"
	"strings"
)

func (config *Config) SOARecordFillDefault(r *SOARecord, useDefaultRecord bool) {
	if useDefaultRecord {
		if len(*r.MName) == 0 {
			r.MName = config.DefaultSOARecord.MName
		}

		if len(*r.RName) == 0 {
			r.RName = config.DefaultSOARecord.RName
		}

		if r.Serial == 0 {
			r.Serial = config.DefaultSOARecord.Serial
		}

		if r.Refresh == 0 {
			r.Refresh = config.DefaultSOARecord.Refresh
		}

		if r.Retry == 0 {
			r.Retry = config.DefaultSOARecord.Retry
		}

		if r.Expire == 0 {
			r.Expire = config.DefaultSOARecord.Expire
		}

		if r.TTL == 0 {
			r.TTL = config.DefaultSOARecord.TTL
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

// FixConfig fixes Config and fill in defaults
func (config *Config) FixConfig() {
	var err error

	if config.DefaultTTL == 0 {
		config.DefaultTTL = 114514
	}

	for index, ns := range config.DefaultNSes {
		config.DefaultNSes[index] = ensureDotAtRight(ns)
	}

	if config.DefaultSOARecord == nil {
		config.DefaultSOARecord = new(SOARecord)
	}
	config.SOARecordFillDefault(config.DefaultSOARecord, false)

	fixedHosts := make([]*PerNetConfig, 0)
	for network, domain := range config.PerHostConfigs {
		netCIDR := ""
		for i := 0; i < len(network); i++ {
			switch network[i] {
			case '.':
				netCIDR = network + "/32"
				break
			case ':':
				netCIDR = network + "/128"
				break
			}
		}
		if netCIDR == "" {
			break
		}
		log.Printf("Loading host %s -> %s\n", netCIDR, domain)
		mode := "fixed"
		for _, d := range strings.Split(domain, ",") {
			thisHost := &PerNetConfig{
				IPNetString:             &netCIDR,
				PtrGenerationModeString: &mode,
				Domain:                  &d,
			}
			fixedHosts = append(fixedHosts, thisHost)
		}
	}
	config.PerNetConfigs = append(fixedHosts, config.PerNetConfigs...)

	// note that range is byVal so we use index here
	for _, currentConfig := range config.PerNetConfigs {
		// fill IPNet
		_, currentConfig.IPNet, err = net.ParseCIDR(*currentConfig.IPNetString)
		exception.HardFailWithReason("failed to parse CIDR", err)

		log.Printf("Loading network %s\n", currentConfig.IPNet.String())

		// fill Mode
		if currentConfig.PtrGenerationModeString == nil {
			log.Fatalf("Missing PTR generation method")
		}
		switch strings.ToLower(*currentConfig.PtrGenerationModeString) {
		case "fixed":
			currentConfig.PtrGenerationMode = Fixed
		case "prefix_ltr":
			currentConfig.PtrGenerationMode = PrependLeftToRight
		case "prefix_rtl":
			currentConfig.PtrGenerationMode = PrependRightToLeft
		case "prefix_ltr_dash":
			currentConfig.PtrGenerationMode = PrependLeftToRightDash
		case "prefix_rtl_dash":
			currentConfig.PtrGenerationMode = PrependRightToLeftDash
		case "prefix_ltr_onlyip":
			currentConfig.PtrGenerationMode = PrependLeftToRightOnlyip
		case "prefix_rtl_onlyip":
			currentConfig.PtrGenerationMode = PrependRightToLeftOnlyip
		default:
			log.Fatalf("Unknown mode \"%s\"", *currentConfig.PtrGenerationModeString)
		}

		// fill IPv6Notation
		if currentConfig.IPv6NotationString == nil {
			currentConfig.IPv6NotationMode = ArpaNotation
		} else {
			switch strings.ToLower(*currentConfig.IPv6NotationString) {
			case "arpa":
				currentConfig.IPv6NotationMode = ArpaNotation
			case "four_hexs":
				currentConfig.IPv6NotationMode = FourHexsNotation
			default:
				log.Fatalf("Unknown ipv6_notation \"%s\"", *currentConfig.PtrGenerationModeString)
			}
		}

		// check domain
		currentConfig.Domain = ensureDotAtRight(currentConfig.Domain)
		currentConfig.Domain = ensureNoDotAtLeft(currentConfig.Domain)

		// fill TTL
		if currentConfig.TTL == 0 {
			currentConfig.TTL = config.DefaultTTL
		}

		// fill SOA
		if currentConfig.SOARecord == nil {
			currentConfig.SOARecord = config.DefaultSOARecord
		} else {
			config.SOARecordFillDefault(currentConfig.SOARecord, true)
		}

		// Add configuration to dedicated list
		if strings.Contains(*currentConfig.IPNetString, ":") {
			config.PerIPv6NetConfigs = append(config.PerIPv6NetConfigs, currentConfig)
		} else {
			config.PerIPv4NetConfigs = append(config.PerIPv4NetConfigs, currentConfig)
		}
	}

	slices.SortFunc(config.PerNetConfigs, subnetSortingFunc)
	slices.SortFunc(config.PerIPv4NetConfigs, subnetSortingFunc)
	slices.SortFunc(config.PerIPv6NetConfigs, subnetSortingFunc)
}

func subnetSortingFunc(a, b *PerNetConfig) int {
	var aOnes, _ = a.IPNet.Mask.Size()
	var bOnes, _ = b.IPNet.Mask.Size()

	if aOnes > bOnes {
		return -1
	}
	if aOnes < bOnes {
		return 1
	}
	return 0
}
