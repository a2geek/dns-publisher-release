package publishers

import (
	"net"
)

type IPPublisher interface {
	Current() (map[string][]net.IP, error)
	Add(host string, ips []net.IP) error
	Delete(host string) error
	Commit() error
}

type AliasPublisher interface {
	Current() ([]string, error)
	Add(url string, alias string) error
	Delete(url string) error
	Commit() error
}
