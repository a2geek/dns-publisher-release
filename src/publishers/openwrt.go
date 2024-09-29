package publishers

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"golang.org/x/crypto/ssh"
)

func NewOpenWrtPublisher(config map[string]string, logger boshlog.Logger, dryRun bool) (Publisher, error) {
	user, ok := config["user"]
	if !ok {
		user = "root"
	}

	privateKey, ok := config["private-key"]
	if !ok {
		return nil, errors.New("private-key not specified for openwrt configuration")
	}

	host, ok := config["host"]
	if !ok {
		return nil, errors.New("ssh host must be specified in openwrt configuration")
	}
	if !strings.Contains(host, ":") {
		host += ":22"
	}

	signer, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		return nil, err
	}

	auth := ssh.PublicKeys(signer)

	clientConfig := ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return &openwrtPublisher{
		clientConfig: clientConfig,
		hostAndPort:  host,
		logger:       logger,
		dryRun:       dryRun,
	}, nil
}

type openwrtPublisher struct {
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

func (p *openwrtPublisher) appendEntry(section, name string, ip net.IP) {
	if section != "" && name != "" && ip != nil {
		p.entries = append(p.entries, entry{
			section: section,
			name:    name,
			ip:      ip,
		})
	}
}

func (p *openwrtPublisher) Current() (map[string][]net.IP, error) {
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
	return data, nil
}

func (p *openwrtPublisher) Add(host string, ips []net.IP) error {
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

func (p *openwrtPublisher) Commit() error {
	p.entries = []entry{}
	if p.dryRun {
		return p.runCommand("uci revert dhcp")
	} else {
		return p.runCommand("uci commit dhcp; reload_config")
	}
}

func (p *openwrtPublisher) Delete(host string) error {
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

func (p *openwrtPublisher) runCommand(msg string, args ...interface{}) error {
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

func (p *openwrtPublisher) outputCommand(msg string, args ...interface{}) (string, error) {
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
