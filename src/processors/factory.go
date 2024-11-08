package processors

import (
	"context"
	"dns-publisher/publishers"
	"dns-publisher/sources"
	"dns-publisher/triggers"
	"fmt"
	"regexp"
	"strings"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cloudfoundry/go-cfclient/v3/client"
	"github.com/cloudfoundry/go-cfclient/v3/config"
)

func NewBoshDnsProcessor(config BoshDnsConfig, publisher publishers.IPPublisher, logger boshlog.Logger) (Processor, error) {
	source, err := sources.NewSource()
	if err != nil {
		return nil, err
	}

	trigger, err := triggers.NewTrigger(config.Trigger, logger)
	if err != nil {
		return nil, err
	}

	processor := &boshDnsProcessor{
		source:    source,
		trigger:   trigger,
		publisher: publisher,
		logger:    logger,
	}

	if config.Type == "manual" {
		processor.mappings = func() ([]MappingConfig, error) {
			return config.Mappings, nil
		}
		return processor, nil
	} else if config.Type == "manifest" {
		client, err := NewBoshConnection(config.Director, logger)
		if err != nil {
			return nil, err
		}
		processor.mappings = client.GetMappings
		return processor, nil
	} else {
		return nil, fmt.Errorf("unknown dns processor type: %s", config.Type)
	}
}

func NewCloudFoundryProcessor(cfConfig CloudFoundryConfig, publisher publishers.AliasPublisher, logger boshlog.Logger) (Processor, error) {
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

	res := []*regexp.Regexp{}
	for _, match := range cfConfig.Mappings {
		str := strings.ReplaceAll(match, "*", "[-0-9a-zA-Z]+")
		re, err := regexp.Compile("^" + str + "$")
		if err != nil {
			return nil, err
		}
		res = append(res, re)
	}

	return &cloudFoundryProcessor{
		trigger:   trigger,
		cf:        cf,
		alias:     cfConfig.Alias,
		regexps:   res,
		publisher: publisher,
		logger:    logger,
	}, nil
}
