package microservice

type CacherConfig struct{}

func NewCacherConfig() *CacherConfig {
	return &CacherConfig{}
}

func (cfg *CacherConfig) Endpoint() string {
	return getEnv("REDIS_CACHE_URI", "")
}

func (cfg *CacherConfig) Password() string {
	return getEnv("REDIS_CACHE_PASSWORD", "")
}

func (cfg *CacherConfig) DB() int {
	return 0
}

func (cfg *CacherConfig) ConnectionSettings() ICacherConnectionSettings {
	return NewDefaultCacherConnectionSettings()
}
