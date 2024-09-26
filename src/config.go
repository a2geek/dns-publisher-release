package main

import (
	"dns-publisher/sources"
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	Source    sources.SourceConfig
	Publisher map[string]string
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
	if len(c.Source.ByQuery) == 0 {
		return errors.New("expecting dns query configuration")
	}
	if len(c.Publisher) == 0 {
		return errors.New("expecting publish destination")
	}
	if c.Publisher["type"] == "" {
		return errors.New("type of publisher required")
	}
	return nil
}
