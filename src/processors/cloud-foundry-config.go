package processors

import (
	"dns-publisher/triggers"
	"errors"
)

type CloudFoundryConfig struct {
	Trigger           triggers.TriggerConfig
	URL               string
	SkipSslValidation bool
	ClientId          string
	ClientSecret      string
	Alias             string
	Mappings          []string
}

func (c *CloudFoundryConfig) Validate() error {
	if c.URL == "" {
		return errors.New("url required")
	}
	if c.ClientId == "" || c.ClientSecret == "" {
		return errors.New("client id and secret required")
	}
	if c.Alias == "" || len(c.Mappings) == 0 {
		return errors.New("alias required and url mappings expected")
	}

	return c.Trigger.Validate()
}
