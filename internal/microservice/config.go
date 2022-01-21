package microservice

import (
	"os"

	"github.com/joho/godotenv"
)

type IConfig interface {
	PersisterConfig() IPersisterConfig
	MongoPersisterConfig() IPersisterConfig
	MQServer() string
	TopicName() string
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
	if "" == env {
		os.Setenv("MODE", "development")
		env = "development"
	}

	cfg.Mode = env

	godotenv.Load(".env." + env + ".local")
	if "test" != env {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load() //
}

func (cfg *Config) PersisterConfig() IPersisterConfig {
	return NewPersisterConfig()
}

func (cfg *Config) MongoPersisterConfig() IPersisterConfig {
	return NewMongoPersisterConfig()
}

//kafka server
func (cfg *Config) MQServer() string {
	return os.Getenv("KAFKA_SERVER_URL")
}

func (cfg *Config) TopicName() string {
	return os.Getenv("TOPIC_NAME")
}

type PersisterConfig struct{}

func NewPersisterConfig() *PersisterConfig {
	return &PersisterConfig{}
}

func (cfg *PersisterConfig) Host() string {
	return os.Getenv("POSTGRES_HOST")
}

func (cfg *PersisterConfig) Port() string {
	return os.Getenv("POSTGRES_PORT")
}

func (cfg *PersisterConfig) DB() string {
	return os.Getenv("POSTGRES_DB_NAME")
}

func (cfg *PersisterConfig) Username() string {
	return os.Getenv("POSTGRES_USERNAME")
}

func (cfg *PersisterConfig) Password() string {
	return os.Getenv("POSTGRES_PASSWORD")
}

func (cfg *PersisterConfig) SSLMode() string {
	sslMode := os.Getenv("POSTGRES_SSL_MODE")
	if sslMode != "" {
		sslMode = "disable"
	}
	return sslMode
}

func (cfg *PersisterConfig) TimeZone() string {

	timezoneEnvironment := os.Getenv("timezone")
	if timezoneEnvironment != "" {
		timezoneEnvironment = "Asia/Bangkok"
	}

	return timezoneEnvironment
}

type MongoPersisterConfig struct{}

func NewMongoPersisterConfig() *MongoPersisterConfig {
	return &MongoPersisterConfig{}
}

func (cfg *MongoPersisterConfig) Host() string {
	return "localhost"
}

func (cfg *MongoPersisterConfig) Port() string {
	return "27017"
}

func (cfg *MongoPersisterConfig) DB() string {
	return "smldev"
}

func (cfg *MongoPersisterConfig) Username() string {
	return "root"
}

func (cfg *MongoPersisterConfig) Password() string {
	return "rootx"
}

func (cfg *MongoPersisterConfig) SSLMode() string {
	return ""
}

func (cfg *MongoPersisterConfig) TimeZone() string {
	return ""
}
