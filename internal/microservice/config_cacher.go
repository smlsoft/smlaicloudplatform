package microservice

import "time"

type CacherConfig struct{}

// ICacherConfig is cacher configuration interface
type ICacherConfig interface {
	Endpoint() string
	Password() string
	UserName() string
	TLS() bool
	DB() int
	ConnectionSettings() ICacherConnectionSettings
}

// ICacherConnectionSettings is connection settings for cacher
type ICacherConnectionSettings interface {
	PoolSize() int
	MinIdleConns() int
	MaxRetries() int
	MinRetryBackoff() time.Duration
	MaxRetryBackoff() time.Duration
	IdleTimeout() time.Duration
	IdleCheckFrequency() time.Duration
	PoolTimeout() time.Duration
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
}

func NewCacherConfig() *CacherConfig {
	return &CacherConfig{}
}

func (cfg *CacherConfig) Endpoint() string {
	return getEnv("REDIS_CACHE_URI", "")
}

func (cfg *CacherConfig) Password() string {
	return getEnv("REDIS_CACHE_PASSWORD", "")
}

func (cfg *CacherConfig) UserName() string {
	return getEnv("REDIS_CACHE_USERNAME", "")
}

func (cfg *CacherConfig) TLS() bool {
	tlsEnable := getEnv("REDIS_CACHE_TLS_ENABLE", "")

	if tlsEnable != "" && tlsEnable == "true" {
		return true
	}
	return false
}

func (cfg *CacherConfig) DB() int {
	return 0
}

func (cfg *CacherConfig) ConnectionSettings() ICacherConnectionSettings {
	return NewDefaultCacherConnectionSettings()
}
