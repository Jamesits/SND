package config

import "net"

type PtrGenerationMode int

const (
	Fixed PtrGenerationMode = iota
	PrependLeftToRight
	PrependRightToLeft
	PrependLeftToRightDash
	PrependRightToLeftDash
	PrependRightToLeftOnlyip
	PrependLeftToRightOnlyip
)

type IPv6NotationMode int

const (
	ArpaNotation IPv6NotationMode = iota
	FourHexsNotation
)

// Config file
type Config struct {
	Debug                 bool              `toml:"debug"`
	Listen                []*string         `toml:"listen"`
	PerNetConfigs         []*PerNetConfig   `toml:"net"`
	PerIPv4NetConfigs     []*PerNetConfig   `toml:""`
	PerIPv6NetConfigs     []*PerNetConfig   `toml:""`
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

type PerNetConfig struct {
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
