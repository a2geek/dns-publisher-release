package main

import (
	"errors"
	"fmt"
	"net"
)

type Publisher interface {
	Current() (map[string][]net.IP, error)
}

func NewPublisher(config map[string]string) (Publisher, error) {
	pubType, ok := config["type"]
	if !ok {
		return nil, errors.New("publisher type not specified")
	}

	switch pubType {
	case "openwrt":
		return NewOpenWrtPublisher(config)
	default:
		return nil, fmt.Errorf("unsupported publisher type: %s", pubType)
	}
}
