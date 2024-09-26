package sources

import (
	"errors"
	"net"
)

type Source interface {
	Start() <-chan interface{}
	Lookup(query string) ([]net.IP, error)
}

func NewSource(config SourceConfig) (Source, error) {
	if config.Path != "" && len(config.ByWatcher) > 0 {

	}
	if config.Refresh != "" && len(config.ByQuery) > 0 {
		return NewDnsQuerySource(config.QuerySourceConfig)
	}
	return nil, errors.New("unrecognized source configuration")
}
