package main

import (
	"dns-publisher/triggers"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	Trigger   triggers.TriggerConfig
	Mappings  []MappingConfig
	Publisher map[string]string
}
type MappingConfig struct {
	InstanceGroup string // required
	Network       string // defaults to 'default'
	Deployment    string // required
	TLD           string // defaults to 'bosh'
	FQDN          string // required
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

	if len(c.Mappings) == 0 {
		return errors.New("expecting dns query configuration")
	}
	for _, m := range c.Mappings {
		err := m.Validate()
		if err != nil {
			return err
		}
	}

	err := c.Trigger.Validate()
	return err
}

func (c *MappingConfig) Validate() error {
	if c.InstanceGroup == "" || c.Deployment == "" || c.FQDN == "" {
		return errors.New("mappings require an instance group, deployment, and FQDN")
	}
	if c.Network == "" {
		c.Network = "default"
	}
	if c.TLD == "" {
		c.TLD = "bosh"
	}
	return nil
}

func (c *MappingConfig) Query() string {
	return fmt.Sprintf("q-s0.%s.%s.%s.%s", c.InstanceGroup, c.Network, c.Deployment, c.TLD)
}
