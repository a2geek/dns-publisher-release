package processors

import (
	"context"
	"dns-publisher/publishers"
	"dns-publisher/triggers"
	"os"
	"regexp"
	"slices"
	"sort"
	"time"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cloudfoundry/go-cfclient/v3/client"
	"github.com/cloudfoundry/go-cfclient/v3/resource"
)

type cloudFoundryProcessor struct {
	trigger   triggers.Trigger
	cf        *client.Client
	alias     string
	regexps   []*regexp.Regexp
	publisher publishers.AliasPublisher
	logger    boshlog.Logger
	priorRun  time.Time
}

const (
	ROUTE_CREATE = "audit.route.create"
	ROUTE_DELETE = "audit.route.delete-request"
)

func (p *cloudFoundryProcessor) Run() {
	events, err := p.trigger.Start()
	if err != nil {
		p.logger.Error("cloud-foundry", "Starting event trigger: %s", err.Error())
		os.Exit(1)
	}

	for range events {
		p.checkEvents()
	}
}

func (p *cloudFoundryProcessor) checkEvents() {
	p.logger.Debug("cloud-foundry", "checking at %s", time.Now())
	opts := client.NewAuditEventListOptions()
	if !p.priorRun.IsZero() {
		opts.ListOptions.CreateAts = client.TimestampFilter{
			Timestamp: []time.Time{p.priorRun},
			Operator:  client.FilterModifierGreaterThan,
		}
	}
	opts.Types = client.Filter{
		Values: []string{ROUTE_CREATE, ROUTE_DELETE},
	}

	p.priorRun = time.Now()
	urls := map[string]string{}    // GUID:URL
	actions := map[string]string{} // URL:action

	for {
		events, pager, err := p.cf.AuditEvents.List(context.Background(), opts)
		if err != nil {
			p.logger.Error("cloud-foundry", "error retrieving audit events: %s", err.Error())
			return
		}

		for _, event := range events {
			guid := event.Target.GUID

			url, ok := urls[guid]
			if !ok {
				route, err := p.cf.Routes.Get(context.Background(), guid)
				if err != nil {
					if resource.IsResourceNotFoundError(err) {
						url = ""
						urls[guid] = url
						p.logger.Warn("cloud-foundry", "route not found: %v", err)
						continue
					}
					p.logger.Error("cloud-foundry", "error retrieving route: %v", err)
					return
				}
				p.logger.Info("cloud-foundry", "route is '%s' for GUID '%s'", route.URL, guid)
				url = route.URL
				urls[guid] = url
			}

			if url == "" {
				p.logger.Debug("cloud-foundry", "route is undefined for guid %s", guid)
				continue
			}

			match := false
			for _, re := range p.regexps {
				if re.Match([]byte(url)) {
					match = true
					p.logger.Debug("cloud-foundry", "route match: '%s' with '%s'", url, re.String())
				}
			}
			if !match {
				p.logger.Debug("cloud-foundry", "DID NOT MATCH: '%s'", url)
				continue
			}

			actions[url] = event.Type
		}

		if !pager.HasNextPage() {
			break
		}
		pager.NextPage(opts)
	}

	current, err := p.publisher.Current()
	if err != nil {
		p.logger.Error("cloud-foundry", "cannot retrieve current state: %v", err)
		return
	}
	p.logger.Debug("cloud-foundry", "current alias listing: %v", current)
	p.logger.Debug("cloud-foundry", "actions captures: %v", actions)

	slices.Sort(current)
	changed := false
	for url, action := range actions {
		i := sort.SearchStrings(current, url)
		exists := i < len(current) && current[i] == url
		p.logger.Debug("cloud-foundry", "found '%s' in slice '%s': %t (at %d)", url, current, exists, i)
		switch action {
		case ROUTE_CREATE:
			if !exists {
				p.logger.Info("cloud-foundry", "making alias for '%s' to '%s'", url, p.alias)
				p.publisher.Add(url, p.alias)
				changed = true
			}
		case ROUTE_DELETE:
			if exists {
				p.logger.Info("cloud-foundry", "removing alias for '%s'", url)
				p.publisher.Delete(url)
				changed = true
			}
		}
	}
	if changed {
		p.publisher.Commit()
	}
}
