package sources

import (
	"net"
	"time"
)

func NewDnsQuerySource(config QuerySourceConfig) (Source, error) {
	if config.Refresh == "" {
		config.duration = 10 * time.Second
	} else {
		duration, err := time.ParseDuration(config.Refresh)
		if err != nil {
			return nil, err
		}
		config.duration = duration
	}
	return &dnsQuerySource{
		config: config,
	}, nil
}

type dnsQuerySource struct {
	config QuerySourceConfig
}
type dnsTick struct {
	data time.Time
}

func (s *dnsQuerySource) Start() <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		for t := range time.Tick(s.config.duration) {
			ch <- dnsTick{
				data: t,
			}
		}
	}()
	return ch
}

func (s *dnsQuerySource) Lookup(query string) ([]net.IP, error) {
	return net.LookupIP(query)
}
