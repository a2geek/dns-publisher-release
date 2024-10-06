package sources

import (
	"context"
	"net"
	"time"
)

func NewDnsQuerySource() (Source, error) {
	dialFunc := func(c context.Context, network string, _ string) (net.Conn, error) {
		d := net.Dialer{}
		return d.DialContext(c, network, "169.254.0.2:53")
	}
	return &dnsQuerySource{
		resolver: net.Resolver{
			PreferGo: true,
			Dial:     dialFunc,
		},
	}, nil
}

type dnsQuerySource struct {
	resolver net.Resolver
}
type dnsTick struct {
	data time.Time
}

func (s *dnsQuerySource) Lookup(query string) ([]net.IP, error) {
	return s.resolver.LookupIP(context.Background(), "ip4", query)
}
