package config

import "net"

type PtrGenerationMode int

const (
	FIXED PtrGenerationMode = iota
	PREPEND_LEFT_TO_RIGHT
	PREPEND_RIGHT_TO_LEFT
	PREPEND_LEFT_TO_RIGHT_DASH
	PREPEND_RIGHT_TO_LEFT_DASH
	PREPEND_RIGHT_TO_LEFT_ONLYIP
	PREPEND_LEFT_TO_RIGHT_ONLYIP
)

type IPv6NotationMode int

const (
	ARPA_NOTATION IPv6NotationMode = iota
	FOUR_HEXS_NOTATION
)

// Config file
type Config struct {
	Listen                []*string         `toml:"listen"`
	PerNetConfigs         []*perNetConfig   `toml:"net"`
	PerHostConfigs        map[string]string `toml:"host"`
	DefaultNSes           []*string         `toml:"ns"`
	OverrideVersionString string            `toml:"version_string"`
	DefaultSOARecord      *SOARecord        `toml:"SOA"`
	DefaultTTL            uint32            `toml:"default_ttl"`
	CompressDNSMessages   bool              `toml:"compress_dns_messages"`
	AllowVersionReporting bool              `toml:"allow_version_reporting"`
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
