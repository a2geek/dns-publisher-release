package publishers

import (
	"net"
)

type Publisher interface {
	Current() (map[string][]net.IP, error)
	Add(host string, ips []net.IP) error
	Delete(host string) error
	Commit() error
}
