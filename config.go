package main

import "net"

type PtrGenerationMode int

const (
	FIXED PtrGenerationMode = iota
	PREPEND_LEFT_TO_RIGHT
	PREPEND_RIGHT_TO_LEFT
)

type IPv6NotationMode int

const (
	ARPA_NOTATION IPv6NotationMode = iota
	FOUR_HEXS_NOTATION
)

// config file
type config struct {
	Listen           []string       `toml:"listen"`
	DefaultTTL       uint32         `toml:"default_ttl"`
	DefaultSOARecord *SOARecord     `toml:"SOA"`
	PerNetConfigs    []perNetConfig `toml:"net"`
}

type SOARecord struct {
	MName   string `toml:"MNAME"`
	RName   string `toml:"RNAME"`
	Serial  uint32 `toml:"SERIAL"`
	Refresh uint32 `toml:"REFRESH"`
	Retry   uint32 `toml:"RETRY"`
	Expire  uint32 `toml:"EXPIRE"`
	TTL     uint32 `toml:"TTL"`
}

type perNetConfig struct {
	IPNetString             string            `toml:"net"`
	IPNet                   *net.IPNet        `toml:""`
	PtrGenerationModeString string            `toml:"mode"`
	PtrGenerationMode       PtrGenerationMode `toml:""`
	IPv6NotationString      string            `toml:"ipv6_notation"`
	IPv6NotationMode        IPv6NotationMode  `toml:""`
	Domain                  string            `toml:"domain"`
	DomainPrefix            string            `toml:"domain_prefix"`
	TTL                     uint32            `toml:"ttl"`
	SOARecord               *SOARecord        `toml:"SOA"`
}

func SOARecordFillDefault(r *SOARecord, useDefaultRecord bool) {
	// TODO: fill in RName
	// TODO: check MName and RName format (dot at the end)
	if useDefaultRecord {
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

}
