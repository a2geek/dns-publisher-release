package processors

import (
	"net"
	"strings"
)

func sameContents(a, b []net.IP) bool {
	if len(a) != len(b) {
		return false
	}
	m := make(map[string]bool)
	for _, ip := range a {
		m[ip.String()] = true
	}
	for _, ip := range b {
		if _, ok := m[ip.String()]; !ok {
			return false
		}
	}
	return true
}

func hostKeysAsString(state map[string][]net.IP) string {
	var builder strings.Builder
	for key := range state {
		if builder.Len() > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(key)
	}
	return builder.String()
}
