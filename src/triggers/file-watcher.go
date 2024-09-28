package triggers

import (
	"path/filepath"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/fsnotify/fsnotify"
)

func newFileWatcherTrigger(fullPath string, logger boshlog.Logger) (Trigger, error) {
	return &fileWatchTrigger{
		fullPath: fullPath,
		logger:   logger,
	}, nil
}

type fileWatchTrigger struct {
	fullPath string
	logger   boshlog.Logger
}
type fileWatchTick struct {
	fullPath string
}

func (t *fileWatchTrigger) Start() (<-chan interface{}, error) {
	ch := make(chan interface{})

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				t.logger.Debug("file-watcher", "fswatcher event: %v", event)
				if event.Has(fsnotify.Write) {
					t.logger.Debug("file-watcher", "modified file:", event.Name)
					if event.Name == t.fullPath {
						ch <- fileWatchTick{
							fullPath: event.Name,
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				t.logger.Error("file-watcher", "fsnotify error: %v", err)
			}
		}
	}()

	err = watcher.Add(filepath.Dir(t.fullPath))
	if err != nil {
		return nil, err
	}

	return ch, nil
}
