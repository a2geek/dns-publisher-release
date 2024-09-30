package publishers

import (
	"bufio"
	"net"
	"regexp"
	"strings"
)

type dhcpDomainPublisher struct {
	openwrtCommon
}

func (p *dhcpDomainPublisher) Current() (map[string][]net.IP, error) {
	// Expected format:
	//   dhcp.cfg08f37d=domain
	//   dhcp.cfg08f37d.name='fake.lan'
	//   dhcp.cfg08f37d.ip='3.4.5.6'
	//   dhcp.cfg09f37d=domain
	//   dhcp.cfg09f37d.name='fake.lan'
	//   dhcp.cfg09f37d.ip='4.5.6.7'
	output, err := p.outputCommand("uci -X show dhcp")
	if err != nil {
		return nil, err
	}

	domain := regexp.MustCompile(`^dhcp\.(cfg\w+)(\.\w+)?='?([^']*)'?$`)

	var section, name string
	var ip net.IP
	scanner := bufio.NewScanner(strings.NewReader(output))
	p.entries = []entry{}
	for scanner.Scan() {
		// 0 = full matching text
		// 1 = section id (ex: cfg09f37d)
		// 2 = option name ('.ip' or '.name')
		// 3 = option value, without quotes
		values := domain.FindStringSubmatch(scanner.Text())
		if values == nil || len(values) != 4 {
			// not the format of our line, just other dhcp entries we can safely ignore
			continue
		}

		switch values[2] {
		case ".ip":
			ip = net.ParseIP(values[3])
			if ip == nil {
				p.logger.Warn("openwrt", "expecting IP but got '%s'", values[3])
			}
		case ".name":
			name = values[3]
		case "": // anything else is a section heading (both for 'domain' and everything else as well)
			if values[3] == "domain" {
				p.appendEntry(section, name, ip)
				section = values[1]
			} else {
				p.logger.Debug("openwrt", "ignoring section '%s'", values[3])
			}
			name = ""
			ip = nil
		}
	}
	p.appendEntry(section, name, ip)
	p.logger.Debug("openwrt", "entries found: %v", p.entries)
	return p.entriesToMap(), nil
}

func (p *dhcpDomainPublisher) Add(host string, ips []net.IP) error {
	for _, ip := range ips {
		section, err := p.outputCommand("uci add dhcp domain")
		if err != nil {
			return err
		}
		section = strings.TrimSpace(section)

		err = p.runCommand("uci set dhcp.%s.name='%s'; uci set dhcp.%s.ip='%s'", section, host, section, ip.String())
		if err != nil {
			return err
		}

		p.appendEntry(section, host, ip)
	}
	return nil
}

func (p *dhcpDomainPublisher) Delete(host string) error {
	keep := []entry{}
	for _, e := range p.entries {
		if e.name == host {
			err := p.runCommand("uci delete dhcp.%s", e.section)
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
