package main

import (
	"dns-publisher/filter"
	"encoding/json"
	"errors"
	"os"
	"time"
)

type Config struct {
	Refresh string
	Filters filter.IPFilters
	DNS     DNSConfig
	Publish map[string]string

	duration time.Duration
}

type DNSConfig struct {
	ByQuery map[string][]string
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
	if c.Refresh == "" {
		c.duration = 10 * time.Second
	} else {
		duration, err := time.ParseDuration(c.Refresh)
		if err != nil {
			return err
		}
		c.duration = duration
	}

	err := c.Filters.Validate()
	if err != nil {
		return err
	}

	if len(c.DNS.ByQuery) == 0 {
		return errors.New("expecting dns query configuration")
	}
	if len(c.Publish) == 0 {
		return errors.New("expecting publish destination")
	}
	if c.Publish["type"] == "" {
		return errors.New("type of publisher required")
	}
	return nil
}
