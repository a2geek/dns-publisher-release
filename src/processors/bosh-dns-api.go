package processors

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudfoundry-community/gogobosh"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"gopkg.in/yaml.v3"
)

func NewBoshConnection(config DirectorConfig, logger boshlog.Logger) (*boshConnection, error) {
	clientConfig := &gogobosh.Config{
		BOSHAddress:       config.URL,
		ClientID:          config.ClientId,
		ClientSecret:      config.ClientSecret,
		SkipSslValidation: config.SkipSslValidation,
	}

	if config.Certificate != "" {
		block, _ := pem.Decode([]byte(config.Certificate))
		if block == nil {
			return nil, fmt.Errorf("failed to parse certificate PEM")
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return nil, err
		}

		tlsConf := &tls.Config{}
		tlsConf.RootCAs = x509.NewCertPool()
		tlsConf.RootCAs.AddCert(cert)

		clientConfig.HttpClient = &http.Client{
			Transport: &http.Transport{TLSClientConfig: tlsConf},
		}
	}

	client, err := gogobosh.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}

	return &boshConnection{
		logger: logger,
		client: client,
	}, nil
}

type boshConnection struct {
	logger boshlog.Logger
	client *gogobosh.Client
}

// boshManifest is a very light structure of a manifest - only the pieces we need
type boshManifest struct {
	InstanceGroups []struct {
		Name     string
		Networks []struct {
			Name string
		}
		Tags map[string]string
	} `yaml:"instance_groups"`
}

func (b *boshConnection) GetMappings() ([]MappingConfig, error) {
	b.logger.Info("bosh-manifest", "refreshing from bosh manifests: start")
	deps, err := b.client.GetDeployments()
	if err != nil {
		return nil, err
	}

	mappings := []MappingConfig{}
	for _, dep := range deps {
		b.logger.Debug("bosh-manifest", "pulling manifest for '%s' deployment", dep.Name)
		data, err := b.client.GetDeployment(dep.Name)
		if err != nil {
			return nil, err
		}

		manifest := boshManifest{}
		err = yaml.Unmarshal([]byte(data.Manifest), &manifest)
		if err != nil {
			return nil, err
		}
		b.logger.Debug("bosh-manifest", "simplified manifest for '%s': %v", dep.Name, manifest)

		for _, ig := range manifest.InstanceGroups {
			if fqdns, ok := ig.Tags["fqdns"]; ok {
				mapping := MappingConfig{
					InstanceGroup: ig.Name,
					Network:       ig.Networks[0].Name,
					Deployment:    dep.Name,
					TLD:           "bosh",
					FQDNs:         strings.Split(fqdns, ","),
				}
				mappings = append(mappings, mapping)
			}
		}
	}
	b.logger.Info("bosh-manifest", "refreshing from bosh manifests: end")
	return mappings, nil
}
