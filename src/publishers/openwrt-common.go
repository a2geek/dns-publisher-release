package publishers

import (
	"fmt"
	"net"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"golang.org/x/crypto/ssh"
)

type openwrtCommon struct {
	clientConfig ssh.ClientConfig
	hostAndPort  string
	logger       boshlog.Logger
	dryRun       bool
	entries      []entry
}
type entry struct {
	section, name string
	ip            net.IP
}

func (p *openwrtCommon) appendEntry(section, name string, ip net.IP) {
	if section != "" && name != "" && ip != nil {
		p.entries = append(p.entries, entry{
			section: section,
			name:    name,
			ip:      ip,
		})
	}
}

func (p *openwrtCommon) entriesToMap() map[string][]net.IP {
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

func (p *openwrtCommon) Commit() error {
	p.entries = []entry{}
	if p.dryRun {
		return p.runCommand("uci revert dhcp")
	} else {
		return p.runCommand("uci commit dhcp; reload_config")
	}
}

func (p *openwrtCommon) runCommand(msg string, args ...interface{}) error {
	cmd := fmt.Sprintf(msg, args...)

	conn, err := ssh.Dial("tcp", p.hostAndPort, &p.clientConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return err
	}

	if p.dryRun {
		// prevent change commands
		p.logger.Info("openwrt", "dry-run cmd: %s", cmd)
		return nil
	} else {
		p.logger.Debug("openwrt", "cmd: %s", cmd)
		return session.Run(cmd)
	}
}

func (p *openwrtCommon) outputCommand(msg string, args ...interface{}) (string, error) {
	cmd := fmt.Sprintf(msg, args...)

	conn, err := ssh.Dial("tcp", p.hostAndPort, &p.clientConfig)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	p.logger.Debug("openwrt", "cmd: %s", cmd)
	bytes, err := session.CombinedOutput(cmd)
	output := string(bytes)
	p.logger.Debug("openwrt", "output: %s", output)
	return output, err
}
