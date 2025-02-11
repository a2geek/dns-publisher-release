package processors

type boshDnsMapper interface {
	GetMappings() ([]MappingConfig, error)
	IsReady() (bool, error)
}
