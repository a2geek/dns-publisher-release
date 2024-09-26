package sources

import "time"

type SourceConfig struct {
	QuerySourceConfig
	WatcherSourceConfig
}

type QuerySourceConfig struct {
	Refresh string
	ByQuery map[string][]string

	duration time.Duration
}

type WatcherSourceConfig struct {
	Path      string
	ByWatcher []WatcherConfig
}

type WatcherConfig struct {
	Deployment string
	Instance   string
	Domain     []string
}
