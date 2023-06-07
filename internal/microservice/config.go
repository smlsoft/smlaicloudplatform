package microservice

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type IConfig interface {
	PathPrefix() string
	PersisterConfig() IPersisterConfig
	MongoPersisterConfig() IPersisterMongoConfig
	ClickHouseConfig() IPersisterClickHouseConfig
	ElkPersisterConfig() IPersisterElkConfig
	OpenSearchPersisterConfig() IPersisterOpenSearchConfig
	CacherConfig() ICacherConfig
	MQConfig() IMQConfig
	TopicName() string
	HttpCORS() []string

	// SignKeyPath() string
	// VerifyKeyPath() string
	JwtSecretKey() string
	ApplicationName() string
}

func getEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

type Config struct {
	Mode string
}

func NewConfig() IConfig {
	config := &Config{}
	config.LoadConfig()
	return config
}

func (cfg *Config) LoadConfig() {
	env := os.Getenv("MODE")
	if env == "" {
		os.Setenv("MODE", "development")
		env = "development"
	}

	cfg.Mode = env

	godotenv.Load(".env." + env + ".local")
	if env != "test" {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load() //
}

func (cfg *Config) PathPrefix() string {
	return getEnv("PATH_PREFIX", "")
}

func (*Config) PersisterConfig() IPersisterConfig {
	return NewPersisterConfig()
}

func (cfg *Config) MongoPersisterConfig() IPersisterMongoConfig {
	return NewMongoPersisterConfig()
}

func (cfg *Config) ClickHouseConfig() IPersisterClickHouseConfig {
	return NewPersisterClickHouseConfig()
}

func (*Config) TopicName() string {
	return os.Getenv("TOPIC_NAME")
}

func (*Config) ApplicationName() string {
	return getEnv("SERVICE_NAME", "microservice")
}

func (*Config) HttpCORS() []string {
	rawCORS := getEnv("HTTP_CORS", "*")

	return strings.Split(rawCORS, " ")
}

// func (*Config) SignKeyPath() string {
// 	return getEnv("PUBLIC_KEY_PATH", "./../../private.key")
// }

// func (*Config) VerifyKeyPath() string {
// 	return getEnv("PRIVATE_KEY_PATH", "./../../public.key")
// }

func (*Config) JwtSecretKey() string {
	return getEnv("JWT_SECRET_KEY", "54cfcbf5437a029d48a9f67552eeb04b48a65703")
}

func (cfg *Config) ElkPersisterConfig() IPersisterElkConfig {
	return NewPersisterElkConfig()
}

func (cfg *Config) OpenSearchPersisterConfig() IPersisterOpenSearchConfig {
	return NewPersisterOpenSearchConfig()
}

///

func (cfg *Config) CacherConfig() ICacherConfig {
	return NewCacherConfig()
}

// imprement iloggerconfig
type ILoggerConfig interface {
	LogLevel() string
	DevMode() bool
	Encoder() string
}

type LoggerConfig struct{}

func NewLoggerConfig() ILoggerConfig {
	return &LoggerConfig{}
}

func (*LoggerConfig) LogLevel() string {
	return getEnv("LOG_LEVEL", "info")
}

func (*LoggerConfig) Encoder() string {
	return getEnv("LOG_ENCODER", "")
}

func (*LoggerConfig) DevMode() bool {
	env := getEnv("MODE", "development")
	return env == "development"
}
