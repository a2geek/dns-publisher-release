package publishers

import (
	"errors"
	"strings"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"golang.org/x/crypto/ssh"
)

const (
	DhcpDnsmasqAddress = "dhcp-dnsmasq-address"
	DhcpDomain         = "dhcp-domain"
)

func NewOpenWrtPublisher(config map[string]string, logger boshlog.Logger, dryRun bool) (Publisher, error) {
	strategy, ok := config["strategy"]
	if !ok {
		strategy = DhcpDnsmasqAddress
	}

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
	shared := openwrtCommon{
		clientConfig: clientConfig,
		hostAndPort:  host,
		logger:       logger,
		dryRun:       dryRun,
	}

	switch strategy {
	case DhcpDomain:
		return &dhcpDomainPublisher{
			openwrtCommon: shared,
		}, nil
	default:
		return &dhcpDnsmasqAddressPublisher{
			openwrtCommon: shared,
		}, nil
	}
}
