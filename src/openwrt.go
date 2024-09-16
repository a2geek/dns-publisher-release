package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"strings"

	"golang.org/x/crypto/ssh"
)

func NewOpenWrtPublisher(config map[string]string) (Publisher, error) {
	user, ok := config["user"]
	if !ok {
		user = "root"
	}

	privateKey, ok := config["private-key"]
	if !ok {
		return nil, errors.New("private-key not specified for openwrt configuration")
	}

	host, ok := config["host"]
	if !ok {
		return nil, errors.New("ssh host must be specified in openwrt configuration")
	}
	if !strings.Contains(host, ":") {
		host += ":22"
	}

	signer, err := ssh.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		return nil, err
	}

	auth := ssh.PublicKeys(signer)

	clientConfig := ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return openwrtPublisher{
		clientConfig: clientConfig,
		hostAndPort:  host,
	}, nil
}

type openwrtPublisher struct {
	clientConfig ssh.ClientConfig
	hostAndPort  string
}

func (p openwrtPublisher) Current() (map[string][]net.IP, error) {
	conn, err := ssh.Dial("tcp", p.hostAndPort, &p.clientConfig)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return nil, err
	}

	output, err := session.Output("uci get dhcp.@dnsmasq[].address")
	if err != nil {
		return nil, err
	}

	// output should be '/DNS/IP,IP,IP'
	data := make(map[string][]net.IP)
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "/")
		if len(parts) != 3 {
			return nil, fmt.Errorf("unexpected format in output: %s", line)
		}
		host := parts[1]

		ips := strings.Split(parts[2], ",")
		if len(ips) == 0 {
			return nil, fmt.Errorf("no IP addresses found for host %s", host)
		}

		for _, ip := range ips {
			ipAddr := net.ParseIP(ip)
			if ipAddr == nil {
				return nil, fmt.Errorf("unexpected IP format for host %s with '%s'", host, ip)
			}
			list, ok := data[host]
			if ok {
				list = append(list, ipAddr)
			} else {
				list = []net.IP{ipAddr}
			}
			data[host] = list
		}
	}

	return data, nil
}
