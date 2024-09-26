package main

import (
	"net"
	"testing"
)

func TestSameContents(t *testing.T) {
	a := []net.IP{
		net.ParseIP("192.168.1.2"),
		net.ParseIP("192.168.1.1"),
	}
	b := []net.IP{
		net.ParseIP("192.168.1.1"),
		net.ParseIP("192.168.1.2"),
	}

	if !sameContents(a, b) {
		t.Fail()
	}
}

func TestDifferentContents(t *testing.T) {
	a := []net.IP{
		net.ParseIP("192.168.1.1"),
	}
	b := []net.IP{
		net.ParseIP("192.168.1.2"),
	}

	if sameContents(a, b) {
		t.Fail()
	}
}

func TestDifferentLengths(t *testing.T) {
	a := []net.IP{
		net.ParseIP("192.168.1.2"),
		net.ParseIP("192.168.1.1"),
	}
	b := []net.IP{
		net.ParseIP("192.168.1.2"),
	}

	if sameContents(a, b) {
		t.Fail()
	}
}
