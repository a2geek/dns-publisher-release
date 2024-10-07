package processors

import (
	"dns-publisher/triggers"
	"errors"
	"fmt"
)

type BoshDnsConfig struct {
	Trigger  triggers.TriggerConfig
	Mappings []MappingConfig
}
type MappingConfig struct {
	InstanceGroup string   // required
	Network       string   // defaults to 'default'
	Deployment    string   // required
	TLD           string   // defaults to 'bosh'
	FQDNs         []string // required
}

func (c *BoshDnsConfig) Validate() error {
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
	if c.InstanceGroup == "" || c.Deployment == "" || len(c.FQDNs) == 0 {
		return errors.New("mappings require an instance group, deployment, and FQDNs")
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
