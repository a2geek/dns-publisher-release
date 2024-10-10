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

func NewOpenWrtIPPublisher(config map[string]string, logger boshlog.Logger, dryRun bool) (IPPublisher, error) {
	strategy, ok := config["strategy"]
	if !ok {
		strategy = DhcpDnsmasqAddress
	}

	openwrtCommon, err := newOpenWrtCommon(config, logger, dryRun)
	if err != nil {
		return nil, err
	}

	ipMgr := &openwrtIp{
		logger: logger,
	}

	switch strategy {
	case DhcpDomain:
		return &dhcpDomainPublisher{
			logger:        logger,
			openwrtCommon: openwrtCommon,
			openwrtIp:     ipMgr,
		}, nil
	default:
		return &dhcpDnsmasqAddressPublisher{
			logger:        logger,
			openwrtCommon: openwrtCommon,
			openwrtIp:     ipMgr,
		}, nil
	}
}

func NewOpenWrtAliasPublisher(config map[string]string, logger boshlog.Logger, dryRun bool) (AliasPublisher, error) {
	openwrtCommon, err := newOpenWrtCommon(config, logger, dryRun)
	if err != nil {
		return nil, err
	}

	aliasMgr := &openwrtAlias{
		logger: logger,
	}

	return &dhcpCnamePublisher{
		logger:        logger,
		openwrtCommon: openwrtCommon,
		openwrtAlias:  aliasMgr,
	}, nil
}

func newOpenWrtCommon(config map[string]string, logger boshlog.Logger, dryRun bool) (*openwrtCommon, error) {
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

	return &openwrtCommon{
		clientConfig: ssh.ClientConfig{
			User:            user,
			Auth:            []ssh.AuthMethod{auth},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
		hostAndPort: host,
		logger:      logger,
		dryRun:      dryRun,
	}, nil
}
