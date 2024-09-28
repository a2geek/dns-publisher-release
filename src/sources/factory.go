package sources

import (
	"net"
)

type Source interface {
	Lookup(query string) ([]net.IP, error)
}

func NewSource() (Source, error) {
	return NewDnsQuerySource()
}
