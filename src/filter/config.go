package filter

import (
	"fmt"
	"net"
)

type IPFilters struct {
	Subnets []string
	Type    string

	ipnets []*net.IPNet
}

const (
	AllIPTypes = "ALL"
	IPv4Type   = "IPv4"
	IPv6Type   = "IPv6"
)

var validFilterTypes = map[string]bool{
	AllIPTypes: true,
	IPv4Type:   true,
	IPv6Type:   true,
}

func (f IPFilters) Validate() error {
	if !validFilterTypes[f.Type] {
		return fmt.Errorf("invalid IP filter type '%s'", f.Type)
	}

	for _, subnet := range f.Subnets {
		_, ipnet, err := net.ParseCIDR(subnet)
		if err != nil {
			return err
		}
		f.ipnets = append(f.ipnets, ipnet)
	}
	return nil
}

type IPFilter func(ipaddr net.IP) bool

func (f IPFilters) GetFilters() []IPFilter {
	var filters []IPFilter
	if len(f.ipnets) > 0 {
		filter := func(ipaddr net.IP) bool {
			for _, ipnet := range f.ipnets {
				if ipnet.Contains(ipaddr) {
					return true
				}
			}
			return false
		}
		filters = append(filters, filter)
	}

	if f.Type != "" {
		filter := func(ipaddr net.IP) bool {
			switch f.Type {
			case AllIPTypes:
				return true
			case IPv4Type:
				return ipaddr.To4() != nil
			case IPv6Type:
				return ipaddr.To4() == nil
			}
			return true
		}
		filters = append(filters, filter)
	}
	return filters
}
