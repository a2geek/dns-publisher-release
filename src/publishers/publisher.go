package publishers

import (
	"errors"
	"fmt"
	"net"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type Publisher interface {
	Current() (map[string][]net.IP, error)
	Add(host string, ips []net.IP) error
	Delete(host string) error
}

func NewPublisher(config map[string]string, logger boshlog.Logger, dryRun bool) (Publisher, error) {
	pubType, ok := config["type"]
	if !ok {
		return nil, errors.New("publisher type not specified")
	}

	switch pubType {
	case "openwrt":
		return NewOpenWrtPublisher(config, logger, dryRun)
	case "fake":
		return NewFakePublisher(config)
	default:
		return nil, fmt.Errorf("unsupported publisher type: %s", pubType)
	}
}
