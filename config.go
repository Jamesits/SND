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
	Listen                []*string       `toml:"listen"`
	DefaultTTL            uint32          `toml:"default_ttl"`
	DefaultSOARecord      *SOARecord      `toml:"SOA"`
	PerNetConfigs         []*perNetConfig `toml:"net"`
	DefaultNSes           []*string       `toml:"ns"`
	CompressDNSMessages   bool            `toml:"compress_dns_messages"`
	AllowVersionReporting bool            `toml:"allow_version_reporting"`
	OverrideVersionString string          `toml:"version_string"`
}

type SOARecord struct {
	MName   *string `toml:"MNAME"`
	RName   *string `toml:"RNAME"`
	Serial  uint32  `toml:"SERIAL"`
	Refresh uint32  `toml:"REFRESH"`
	Retry   uint32  `toml:"RETRY"`
	Expire  uint32  `toml:"EXPIRE"`
	TTL     uint32  `toml:"TTL"`
}

type perNetConfig struct {
	IPNetString             *string           `toml:"net"`
	IPNet                   *net.IPNet        `toml:""`
	PtrGenerationModeString *string           `toml:"mode"`
	PtrGenerationMode       PtrGenerationMode `toml:""`
	IPv6NotationString      *string           `toml:"ipv6_notation"`
	IPv6NotationMode        IPv6NotationMode  `toml:""`
	Domain                  *string           `toml:"domain"`
	DomainPrefix            *string           `toml:"domain_prefix"`
	TTL                     uint32            `toml:"ttl"`
	SOARecord               *SOARecord        `toml:"SOA"`
}
