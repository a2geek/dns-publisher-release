package publishers

import (
	"net"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type openwrtIp struct {
	logger  boshlog.Logger
	entries []ipEntry
}
type ipEntry struct {
	section, name string
	ip            net.IP
}

func (p *openwrtIp) reset() {
	p.entries = []ipEntry{}
}

func (p *openwrtIp) appendEntry(section, name string, ip net.IP) {
	if section != "" && name != "" && ip != nil {
		p.entries = append(p.entries, ipEntry{
			section: section,
			name:    name,
			ip:      ip,
		})
	}
}

func (p *openwrtIp) entriesToMap() map[string][]net.IP {
	data := make(map[string][]net.IP)
	for _, e := range p.entries {
		ips, ok := data[e.name]
		if !ok {
			ips = []net.IP{e.ip}
		} else {
			ips = append(ips, e.ip)
		}
		data[e.name] = ips
	}
	p.logger.Debug("openwrt", "domains found: %v", data)
	return data
}
