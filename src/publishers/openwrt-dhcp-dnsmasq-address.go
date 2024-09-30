package publishers

import (
	"net"
	"strings"
)

type dhcpDnsmasqAddressPublisher struct {
	openwrtCommon
}

func (p *dhcpDnsmasqAddressPublisher) Current() (map[string][]net.IP, error) {
	p.entries = []entry{}
	output, err := p.outputCommand("uci get dhcp.@dnsmasq[].address")
	if err != nil {
		if strings.Contains(output, "uci: Entry not found") {
			return p.entriesToMap(), nil
		}
		return nil, err
	}
	line := strings.TrimSpace(string(output))

	// Format:  /.sys.cf.lan/ip-addr /.app.cf.lan/ip-addr
	for _, record := range strings.Split(line, " ") {
		parts := strings.Split(record, "/")
		if len(parts) != 3 {
			p.logger.Warn("openwrt", "unexpected format in output: %s", record)
			continue
		}

		host := parts[1]
		ips := strings.Split(parts[2], ",")
		if len(ips) == 0 {
			p.logger.Warn("no IP addresses found for host %s", host)
			continue
		}

		for _, ip := range ips {
			ipAddr := net.ParseIP(ip)
			if ipAddr == nil {
				p.logger.Warn("unexpected IP format for host %s with '%s'", host, ip)
				continue
			}
			p.appendEntry("", host, ipAddr)
		}
	}
	p.logger.Debug("openwrt", "entries found: %v", p.entries)
	return p.entriesToMap(), nil
}

func (p *dhcpDnsmasqAddressPublisher) Add(host string, ips []net.IP) error {
	for _, ip := range ips {
		err := p.runCommand("uci add_list dhcp.@dnsmasq[].address='/%s/%s'", host, ip.String())
		if err != nil {
			return err
		}

		p.appendEntry("", host, ip)
	}
	return nil
}

func (p *dhcpDnsmasqAddressPublisher) Delete(host string) error {
	keep := []entry{}
	for _, e := range p.entries {
		if e.name == host {
			err := p.runCommand("uci del_list dhcp.@dnsmasq[].address='/%s/%s'", e.name, e.ip.String())
			if err != nil {
				return err
			}
		} else {
			keep = append(keep, e)
		}
	}
	p.entries = keep
	return nil
}
