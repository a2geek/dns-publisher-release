package processors

import (
	"dns-publisher/publishers"
	"dns-publisher/sources"
	"dns-publisher/triggers"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

func NewBoshDnsProcessor(config BoshDnsConfig, publisher publishers.Publisher, logger boshlog.Logger) (Processor, error) {
	source, err := sources.NewSource()
	if err != nil {
		return nil, err
	}

	trigger, err := triggers.NewTrigger(config.Trigger, logger)
	if err != nil {
		return nil, err
	}

	return &boshDnsProcessor{
		source:    source,
		trigger:   trigger,
		mappings:  config.Mappings,
		publisher: publisher,
		logger:    logger,
	}, nil
}
