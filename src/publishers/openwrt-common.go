package publishers

import (
	"fmt"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"golang.org/x/crypto/ssh"
)

type openwrtCommon struct {
	clientConfig ssh.ClientConfig
	hostAndPort  string
	logger       boshlog.Logger
	dryRun       bool
}

func (p *openwrtCommon) runCommand(msg string, args ...interface{}) error {
	cmd := fmt.Sprintf(msg, args...)

	conn, err := ssh.Dial("tcp", p.hostAndPort, &p.clientConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return err
	}

	if p.dryRun {
		// prevent change commands
		p.logger.Info("openwrt", "dry-run cmd: %s", cmd)
		return nil
	} else {
		p.logger.Debug("openwrt", "cmd: %s", cmd)
		return session.Run(cmd)
	}
}

func (p *openwrtCommon) outputCommand(msg string, args ...interface{}) (string, error) {
	cmd := fmt.Sprintf(msg, args...)

	conn, err := ssh.Dial("tcp", p.hostAndPort, &p.clientConfig)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	p.logger.Debug("openwrt", "cmd: %s", cmd)
	bytes, err := session.CombinedOutput(cmd)
	output := string(bytes)
	p.logger.Debug("openwrt", "output: %s", output)
	return output, err
}
