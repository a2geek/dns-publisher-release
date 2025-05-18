package config

import (
	"dns-publisher/processors"
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	BoshDns      *processors.BoshDnsConfig
	CloudFoundry *processors.CloudFoundryConfig
	Publisher    map[string]string
	Web          struct {
		HTTP int
	}
}

func NewConfigFromPath(path string) (Config, error) {
	config := Config{}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(bytes, &config)
	if err != nil {
		return config, err
	}

	return config, config.Validate()
}

func (c *Config) Validate() error {
	if len(c.Publisher) == 0 {
		return errors.New("expecting publish destination")
	}
	if c.Publisher["type"] == "" {
		return errors.New("type of publisher required")
	}

	count := 0
	if c.BoshDns != nil {
		err := c.BoshDns.Validate()
		if err != nil {
			return err
		}
		count++
	}

	if c.CloudFoundry != nil {
		err := c.CloudFoundry.Validate()
		if err != nil {
			return err
		}
		count++
	}

	if count == 0 {
		return errors.New("at least one of bosh-dns or cloud-foundry must be configured")
	}

	return nil
}
