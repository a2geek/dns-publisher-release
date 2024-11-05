package publishers

import (
	"fmt"
	"net"
	"strings"
)

func NewFakeIPPublisher(config map[string]string) (*FakeIPPublisher, error) {
	data := make(map[string][]net.IP)
	for k, v := range config {
		if k == "type" {
			continue
		}
		var ips []net.IP
		for _, s := range strings.Split(v, ",") {
			ip := net.ParseIP(s)
			if ip == nil {
				return &FakeIPPublisher{}, fmt.Errorf("invalid IP addr: %s", s)
			}
			ips = append(ips, ip)
		}
		if len(ips) > 0 {
			data[k] = ips
		}
	}
	return &FakeIPPublisher{
		ips: data,
	}, nil
}

type FakeIPPublisher struct {
	IPPublisher
	ips map[string][]net.IP
}

func (p *FakeIPPublisher) Current() (map[string][]net.IP, error) {
	return p.ips, nil
}

func (p *FakeIPPublisher) Add(host string, ips []net.IP) error {
	p.ips[host] = ips
	return nil
}

func (p *FakeIPPublisher) Delete(host string) error {
	delete(p.ips, host)
	return nil
}
