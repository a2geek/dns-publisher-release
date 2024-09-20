package publishers

import (
	"errors"
	"fmt"
	"net"
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
}

func (p *openwrtPublisher) Current() (map[string][]net.IP, error) {
	entries, err := p.currentEntries()
	if err != nil {
		return nil, err
	}

	// Format:  /.sys.cf.lan/ip-addr /.app.cf.lan/ip-addr
	data := make(map[string][]net.IP)
	for _, entry := range entries {
		parts := strings.Split(entry, "/")
		if len(parts) != 3 {
			return nil, fmt.Errorf("unexpected format in output: %s", entry)
		}
		host := parts[1]

		ips := strings.Split(parts[2], ",")
		if len(ips) == 0 {
			return nil, fmt.Errorf("no IP addresses found for host %s", host)
		}

		for _, ip := range ips {
			ipAddr := net.ParseIP(ip)
			if ipAddr == nil {
				return nil, fmt.Errorf("unexpected IP format for host %s with '%s'", host, ip)
			}
			list, ok := data[host]
			if ok {
				list = append(list, ipAddr)
			} else {
				list = []net.IP{ipAddr}
			}
			data[host] = list
		}
	}

	return data, nil
}

func (p *openwrtPublisher) currentEntries() ([]string, error) {
	output, err := p.outputCommand("uci get dhcp.@dnsmasq[].address")
	if err != nil {
		if !strings.Contains(err.Error(), "uci: Entry not found") {
			return []string{}, err
		}
	}
	line := strings.TrimSpace(string(output))
	p.logger.Debug("openwrt", "output: %s", line)

	// Format:  /.sys.cf.lan/ip-addr /.app.cf.lan/ip-addr
	return strings.Split(line, " "), nil
}

func (p *openwrtPublisher) Add(host string, ips []net.IP) error {
	var builder strings.Builder
	for _, ip := range ips {
		if builder.Len() > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(ip.String())
	}

	err := p.runCommand("uci add_list dhcp.@dnsmasq[].address='/%s/%s'", host, builder.String())
	return err
}

func (p *openwrtPublisher) Delete(host string) error {
	entries, err := p.currentEntries()
	if err != nil {
		return err
	}

	prefix := fmt.Sprintf("/%s/", host)
	for _, entry := range entries {
		if strings.HasPrefix(entry, prefix) {
			err = p.runCommand("uci del_list dhcp.@dnsmasq[].address='%s'", entry)
			return err
		}
	}
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
	}
	p.logger.Debug("openwrt", "cmd: %s", cmd)
	return session.Run(cmd)
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

	p.logger.Debug("openwrt", "cmd: %s", cmd)
	bytes, err := session.Output(cmd)
	return string(bytes), err
}
