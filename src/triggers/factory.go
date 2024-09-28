package triggers

import (
	"fmt"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type Trigger interface {
	Start() (<-chan interface{}, error)
}

func NewTrigger(config TriggerConfig, logger boshlog.Logger) (Trigger, error) {
	switch config.Type {
	case "file-watcher":
		return newFileWatcherTrigger(config.FileWatcher, logger)
	case "timer":
		return newRefreshTrigger(config.Refresh)
	}
	return nil, fmt.Errorf("unknown trigger type: '%s'", config.Type)
}
