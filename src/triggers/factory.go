package triggers

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type Trigger interface {
	Start() (<-chan interface{}, error)
}

const defaultFileWatcher = "/var/vcap/instance/dns/records.json"

func NewTrigger(config TriggerConfig, logger boshlog.Logger) (Trigger, error) {
	if config.FileWatcher != "" {
		return newFileWatcherTrigger(config.FileWatcher, logger)
	}
	if config.Refresh != "" {
		return newRefreshTrigger(config.Refresh)
	}
	return newFileWatcherTrigger(defaultFileWatcher, logger)
}
