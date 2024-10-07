package processors

import (
	"context"
	"dns-publisher/publishers"
	"dns-publisher/triggers"
	"os"
	"time"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cloudfoundry/go-cfclient/v3/client"
)

type cloudFoundryProcessor struct {
	trigger   triggers.Trigger
	cf        *client.Client
	publisher publishers.Publisher
	logger    boshlog.Logger
}

func (p *cloudFoundryProcessor) Run() {
	events, err := p.trigger.Start()
	if err != nil {
		p.logger.Error("cloud-foundry", "Starting event trigger: %s", err.Error())
		os.Exit(1)
	}

	priorRun := time.Now().Add(-1 * time.Hour)
	for range events {
		p.logger.Debug("cloud-foundry", "checking at %s", time.Now())
		opts := client.NewAuditEventListOptions()
		opts.ListOptions.CreateAts = client.TimestampFilter{
			Timestamp: []time.Time{priorRun},
			Operator:  client.FilterModifierGreaterThan,
		}
		opts.Types = client.Filter{
			Values: []string{"audit.app.map-route", "audit.app.unmap-route"},
		}
		priorRun = time.Now()
		events, err := p.cf.AuditEvents.ListAll(context.Background(), opts)
		if err != nil {
			p.logger.Error("cloud-foundry", "error retrieving audit events: %s", err.Error())
			continue
		}
		for event := range events {
			p.logger.Info("cloud-foundry", "event: %v", event)
		}
	}
}
