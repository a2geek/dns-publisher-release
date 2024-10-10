package publishers

import (
	"errors"
	"fmt"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

func NewIPPublisher(config map[string]string, logger boshlog.Logger) (IPPublisher, error) {
	pubType, ok := config["type"]
	if !ok {
		return nil, errors.New("publisher type not specified")
	}

	dryRun := config["dry-run"] == "true"

	switch pubType {
	case "openwrt":
		return NewOpenWrtIPPublisher(config, logger, dryRun)
	case "fake":
		return NewFakeIPPublisher(config)
	default:
		return nil, fmt.Errorf("unsupported publisher type: %s", pubType)
	}
}

func NewAliasPublisher(config map[string]string, logger boshlog.Logger) (AliasPublisher, error) {
	pubType, ok := config["type"]
	if !ok {
		return nil, errors.New("publisher type not specified")
	}

	dryRun := config["dry-run"] == "true"

	switch pubType {
	case "openwrt":
		return NewOpenWrtAliasPublisher(config, logger, dryRun)
	case "fake":
		return NewFakeAliasPublisher(config)
	default:
		return nil, fmt.Errorf("unsupported publisher type: %s", pubType)
	}
}
