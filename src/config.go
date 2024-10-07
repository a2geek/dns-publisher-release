package main

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

	if c.BoshDns != nil {
		err := c.BoshDns.Validate()
		if err != nil {
			return err
		}
	}

	if c.CloudFoundry != nil {
		err := c.CloudFoundry.Validate()
		if err != nil {
			return err
		}
	}

	return nil
}
