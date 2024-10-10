package publishers

func NewFakeAliasPublisher(config map[string]string) (*FakeAliasPublisher, error) {
	return &FakeAliasPublisher{}, nil
}

type FakeAliasPublisher struct {
	AliasPublisher
	aliases []string
}

func (p *FakeAliasPublisher) Current() ([]string, error) {
	return p.aliases, nil
}

func (p *FakeAliasPublisher) Add(url, alias string) error {
	p.aliases = append(p.aliases, url)
	return nil
}

func (p *FakeAliasPublisher) Delete(url string) error {
	i := -1
	for n, alias := range p.aliases {
		if alias == url {
			i = n
			break
		}
	}
	if i > -1 {
		p.aliases[i] = p.aliases[len(p.aliases)-1]
		p.aliases = p.aliases[:len(p.aliases)-1]
	}
	return nil
}
