package sources

import (
	"net"
	"time"
)

func NewDnsQuerySource() (Source, error) {
	return &dnsQuerySource{}, nil
}

type dnsQuerySource struct {
}
type dnsTick struct {
	data time.Time
}

func (s *dnsQuerySource) Lookup(query string) ([]net.IP, error) {
	return net.LookupIP(query)
}
