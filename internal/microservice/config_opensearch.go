package microservice

type PersisterOpenSearchConfig struct{}

func NewPersisterOpenSearchConfig() *PersisterOpenSearchConfig {
	return &PersisterOpenSearchConfig{}
}

func (c *PersisterOpenSearchConfig) Address() []string {

	return []string{
		getEnv("OPEN_SEARCH_ADDRESS", "http://192.168.2.204:9200"),
	}
}

func (c *PersisterOpenSearchConfig) Username() string {
	return getEnv("OPEN_SEARCH_USERNAME", "elastic")
}

func (c *PersisterOpenSearchConfig) Password() string {
	return getEnv("OPEN_SEARCH_PASSWORD", "smlSoft2021")
}
