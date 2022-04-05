package config

import "strings"

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
