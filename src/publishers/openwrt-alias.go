package publishers

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type openwrtAlias struct {
	logger  boshlog.Logger
	entries []aliasEntry
}
type aliasEntry struct {
	section, cname, target string
}

func (p *openwrtAlias) reset() {
	p.entries = []aliasEntry{}
}

func (p *openwrtAlias) appendAlias(section, cname, target string) {
	if section != "" && cname != "" && target != "" {
		p.entries = append(p.entries, aliasEntry{
			section: section,
			cname:   cname,
			target:  target,
		})
	}
}

func (p *openwrtAlias) aliasesToList() []string {
	data := []string{}
	for _, e := range p.entries {
		data = append(data, e.cname)
	}
	p.logger.Debug("openwrt", "aliases found: %v", data)
	return data
}
