package publishers

import (
	"bufio"
	"regexp"
	"strings"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type dhcpCnamePublisher struct {
	logger boshlog.Logger

	*openwrtCommon
	*openwrtAlias
}

func (p *dhcpCnamePublisher) Current() ([]string, error) {
	// Expected format:
	//   dhcp.cfg0876c9=cname
	//   dhcp.cfg0876c9.cname='www.cf.lan'
	//   dhcp.cfg0876c9.target='alias.sys.cf.lan'
	output, err := p.outputCommand("uci -X show dhcp")
	if err != nil {
		return nil, err
	}

	domain := regexp.MustCompile(`^dhcp\.(cfg\w+)(\.\w+)?='?([^']*)'?$`)

	var section, cname, target string
	scanner := bufio.NewScanner(strings.NewReader(output))
	p.reset()
	for scanner.Scan() {
		// 0 = full matching text
		// 1 = section id (ex: cfg09f37d)
		// 2 = option name ('.cname' or '.target')
		// 3 = option value, without quotes
		values := domain.FindStringSubmatch(scanner.Text())
		if values == nil || len(values) != 4 {
			// not the format of our line, just other dhcp entries we can safely ignore
			continue
		}

		switch values[2] {
		case ".cname":
			cname = values[3]
		case ".alias":
			target = values[3]
		case "": // anything else is a section heading (both for 'domain' and everything else as well)
			if values[3] == "cname" {
				p.appendAlias(section, cname, target)
				section = values[1]
			} else {
				p.logger.Debug("openwrt", "ignoring section '%s'", values[3])
			}
			cname = ""
			target = ""
		}
	}
	p.appendAlias(section, cname, target)
	p.logger.Debug("openwrt", "entries found: %v", p.entries)
	return p.aliasesToList(), nil
}

func (p *dhcpCnamePublisher) Add(url, alias string) error {
	section, err := p.outputCommand("uci add dhcp cname")
	if err != nil {
		return err
	}
	section = strings.TrimSpace(section)

	err = p.runCommand("uci set dhcp.%s.cname='%s'; uci set dhcp.%s.target='%s'", section, url, section, alias)
	if err != nil {
		return err
	}

	p.appendAlias(section, url, alias)
	return nil
}

func (p *dhcpCnamePublisher) Delete(host string) error {
	keep := []aliasEntry{}
	for _, e := range p.entries {
		if e.cname == host {
			err := p.runCommand("uci delete dhcp.%s", e.section)
			if err != nil {
				return err
			}
		} else {
			keep = append(keep, e)
		}
	}
	p.entries = keep
	return nil
}
