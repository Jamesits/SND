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
	Listen []string `toml:"listen"`
	PerNetConfigs []perNetConfig `toml:"net"`
}

type perNetConfig struct {
	IPNetString    string   `toml:"net"`
	IPNet *net.IPNet `toml:""`
	PtrGenerationModeString string `toml:"mode"`
	PtrGenerationMode PtrGenerationMode `toml:""`
	IPv6NotationString string `toml:"ipv6_notation"`
	IPv6NotationMode IPv6NotationMode `toml:""`
	Domain string `toml:"domain"`
	DomainPrefix string `toml:"domain_prefix"`
}