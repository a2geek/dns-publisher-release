package processors

func NewStaticMapper(mappings []MappingConfig) boshDnsMapper {
	return &boshStaticMapper{
		mappings,
	}
}

type boshStaticMapper struct {
	mappings []MappingConfig
}

func (m *boshStaticMapper) GetMappings() ([]MappingConfig, error) {
	return m.mappings, nil
}

func (m *boshStaticMapper) IsReady() (bool, error) {
	return true, nil
}
