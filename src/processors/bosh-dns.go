package processors

import (
	"dns-publisher/publishers"
	"dns-publisher/sources"
	"dns-publisher/triggers"
	"net"
	"os"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type boshDnsProcessor struct {
	source    sources.Source
	trigger   triggers.Trigger
	mappings  func() ([]MappingConfig, error)
	publisher publishers.IPPublisher
	logger    boshlog.Logger
}

func (p *boshDnsProcessor) Run(actionChan chan<- Action) {
	state, err := p.publisher.Current()
	if err != nil {
		p.logger.Error("bosh-dns", "Retrieving current configuration - %s", err.Error())
		os.Exit(1)
	}
	p.logger.Info("bosh-dns", "Startup state includes %d entries: %v", len(state), hostKeysAsString(state))

	triggers, err := p.trigger.Start()
	if err != nil {
		p.logger.Error("bosh-dns", "Starting event trigger: %s", err.Error())
	}

	for range triggers {
		actionChan <- p
	}
}

func (p *boshDnsProcessor) Name() string {
	return "bosh-dns"
}

func (p *boshDnsProcessor) Act() {
	// check and refresh
	p.logger.Info("bosh-dns", "Updating from DNS")

	state, err := p.publisher.Current()
	if err != nil {
		p.logger.Error("bosh-dns", "Retrieving current configuration - %s", err.Error())
		os.Exit(1)
	}
	p.logger.Info("bosh-dns", "Current state includes %d entries: %v\n", len(state), hostKeysAsString(state))
	changes := false

	mappings, err := p.mappings()
	if err != nil {
		p.logger.Error("bosh-dns", "retrieving mappings: %s", err.Error())
		os.Exit(1)
	}
	p.logger.Debug("bosh-dns", "mappings found: %v", mappings)

	for _, mapping := range mappings {
		query := mapping.Query()
		ips, err := p.source.Lookup(query)
		if err != nil {
			p.logger.Warn("bosh-dns", "unable to lookup '%s': %s", query, err.Error())
			continue
		}
		p.logger.Debug("bosh-dns", "found '%s' for %v = %v", query, mapping.FQDNs, ips)
		for _, fqdn := range mapping.FQDNs {
			change, err := p.adjustState(state, fqdn, ips)
			if err != nil {
				p.logger.Error("bosh-dns", "error adjusting state for '%s': %v", fqdn, err)
			}
			changes = changes || change
		}
	}
	if changes {
		err = p.publisher.Commit()
		if err != nil {
			p.logger.Error("bosh-dns", "unable to commit changes: %s", err.Error())
		}
	}
}

func (p *boshDnsProcessor) adjustState(state map[string][]net.IP, host string, newIPs []net.IP) (bool, error) {
	currentIPs, ok := state[host]
	if ok && sameContents(currentIPs, newIPs) {
		return false, nil
	}
	if ok {
		err := p.publisher.Delete(host)
		if err != nil {
			return false, err
		}
		delete(state, host)
	}
	err := p.publisher.Add(host, newIPs)
	if err != nil {
		return true, err
	}
	state[host] = newIPs
	return true, nil
}
