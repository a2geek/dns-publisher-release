package processors

import (
	"crypto/x509"
	"dns-publisher/triggers"
	"encoding/pem"
	"errors"
	"fmt"
)

type BoshDnsConfig struct {
	Trigger  triggers.TriggerConfig
	Type     string
	Mappings []MappingConfig
	Director DirectorConfig
}
type MappingConfig struct {
	InstanceGroup string   // required
	Network       string   // defaults to 'default'
	Deployment    string   // required
	TLD           string   // defaults to 'bosh'
	FQDNs         []string // required
}
type DirectorConfig struct {
	URL               string // required
	Certificate       string // some valid combination of a certificate and/or skip ssl validation
	SkipSslValidation bool
	ClientId          string // required
	ClientSecret      string // required
	FQDNAllowed       []string
}

func (c *BoshDnsConfig) Validate() error {
	if c.Type == "manual" {
		if len(c.Mappings) == 0 {
			return errors.New("expecting dns query configuration")
		}
		for _, m := range c.Mappings {
			err := m.Validate()
			if err != nil {
				return err
			}
		}
	} else if c.Type == "manifest" {
		err := c.Director.Validate()
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("unknown BOSH DNS configuration type: %s", c.Type)
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

func (c *DirectorConfig) Validate() error {
	if c.URL == "" || c.ClientId == "" || c.ClientSecret == "" {
		return fmt.Errorf("manifest configuration requires url, client id, and client secret")
	}

	if c.Certificate != "" {
		block, _ := pem.Decode([]byte(c.Certificate))
		if block == nil {
			return fmt.Errorf("failed to parse certificate PEM")
		}

		_, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return err
		}
	}

	return nil
}
