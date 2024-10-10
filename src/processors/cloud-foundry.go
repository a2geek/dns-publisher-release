package processors

import (
	"context"
	"dns-publisher/publishers"
	"dns-publisher/triggers"
	"os"
	"regexp"
	"time"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cloudfoundry/go-cfclient/v3/client"
	"github.com/cloudfoundry/go-cfclient/v3/resource"
)

type cloudFoundryProcessor struct {
	trigger   triggers.Trigger
	cf        *client.Client
	regexps   []*regexp.Regexp
	publisher publishers.Publisher
	logger    boshlog.Logger
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

	var priorRun time.Time
	for range events {
		p.logger.Debug("cloud-foundry", "checking at %s", time.Now())
		opts := client.NewAuditEventListOptions()
		if !priorRun.IsZero() {
			opts.ListOptions.CreateAts = client.TimestampFilter{
				Timestamp: []time.Time{priorRun},
				Operator:  client.FilterModifierGreaterThan,
			}
		}
		opts.Types = client.Filter{
			Values: []string{ROUTE_CREATE, ROUTE_DELETE},
		}

		priorRun = time.Now()
		urls := map[string]string{}    // GUID:URL
		actions := map[string]string{} // URL:action

		for {
			events, pager, err := p.cf.AuditEvents.List(context.Background(), opts)
			if err != nil {
				p.logger.Error("cloud-foundry", "error retrieving audit events: %s", err.Error())
				continue
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
						break
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

		p.logger.Debug("cloud-foundry", "urls = %v", urls)
		p.logger.Debug("cloud-foundry", "actions = %v", actions)
	}
}
