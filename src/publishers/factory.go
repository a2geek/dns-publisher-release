package publishers

import (
	"errors"
	"fmt"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

func NewPublisher(config map[string]string, logger boshlog.Logger) (Publisher, error) {
	pubType, ok := config["type"]
	if !ok {
		return nil, errors.New("publisher type not specified")
	}

	dryRun := config["dry-run"] == "true"

	switch pubType {
	case "openwrt":
		return NewOpenWrtPublisher(config, logger, dryRun)
	case "fake":
		return NewFakePublisher(config)
	default:
		return nil, fmt.Errorf("unsupported publisher type: %s", pubType)
	}
}
