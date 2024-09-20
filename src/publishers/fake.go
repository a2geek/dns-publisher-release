package publishers

import (
	"fmt"
	"net"
	"strings"
)

func NewFakePublisher(config map[string]string) (*FakePublisher, error) {
	data := make(map[string][]net.IP)
	for k, v := range config {
		if k == "type" {
			continue
		}
		var ips []net.IP
		for _, s := range strings.Split(v, ",") {
			ip := net.ParseIP(s)
			if ip == nil {
				return &FakePublisher{}, fmt.Errorf("invalid IP addr: %s", s)
			}
			ips = append(ips, ip)
		}
		if len(ips) > 0 {
			data[k] = ips
		}
	}
	return &FakePublisher{
		data: data,
	}, nil
}

type FakePublisher struct {
	Publisher
	data map[string][]net.IP
}

func (p *FakePublisher) Current() (map[string][]net.IP, error) {
	return p.data, nil
}

func (p *FakePublisher) Add(host string, ips []net.IP) error {
	p.data[host] = ips
	return nil
}

func (p *FakePublisher) Delete(host string) error {
	delete(p.data, host)
	return nil
}
