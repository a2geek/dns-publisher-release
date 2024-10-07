package processors

import (
	"context"
	"dns-publisher/publishers"
	"dns-publisher/sources"
	"dns-publisher/triggers"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cloudfoundry/go-cfclient/v3/client"
	"github.com/cloudfoundry/go-cfclient/v3/config"
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

func NewCloudFoundryProcessor(cfConfig CloudFoundryConfig, publisher publishers.Publisher, logger boshlog.Logger) (Processor, error) {
	trigger, err := triggers.NewTrigger(cfConfig.Trigger, logger)
	if err != nil {
		return nil, err
	}

	options := []config.Option{config.ClientCredentials(cfConfig.ClientId, cfConfig.ClientSecret)}
	if cfConfig.SkipSslValidation {
		options = append(options, config.SkipTLSValidation())
	}
	cfg, err := config.New(cfConfig.URL, options...)
	if err != nil {
		return nil, err
	}

	cf, err := client.New(cfg)
	if err != nil {
		return nil, err
	}

	err = cf.Validate()
	if err != nil {
		return nil, err
	}

	root, err := cf.Root.Get(context.Background())
	if err != nil {
		return nil, err
	}
	logger.Info("cloud-foundry", "Connected. API version is %s", root.Links.CloudControllerV3.Meta.Version)

	return &cloudFoundryProcessor{
		trigger:   trigger,
		cf:        cf,
		publisher: publisher,
		logger:    logger,
	}, nil
}
