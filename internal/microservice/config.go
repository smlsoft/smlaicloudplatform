package microservice

import (
	"os"

	"github.com/joho/godotenv"
)

type IConfig interface {
	PathPrefix() string
	PersisterConfig() IPersisterConfig
	MongoPersisterConfig() IPersisterMongoConfig
	ElkPersisterConfig() IPersisterElkConfig
	CacherConfig() ICacherConfig
	MQConfig() IMQConfig
	TopicName() string

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

func (*Config) TopicName() string {
	return os.Getenv("TOPIC_NAME")
}

func (*Config) ApplicationName() string {
	return getEnv("SERVICE_NAME", "microservice")
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
func (cfg *MongoPersisterConfig) MongodbURI() string {
	return getEnv("MONGODB_URI", "") // mongodb://root:rootx@localhost:27017/
}

func (cfg *MongoPersisterConfig) DB() string {
	return getEnv("MONGODB_DB", "smldev")
}

func (cfg *Config) ElkPersisterConfig() IPersisterElkConfig {
	return NewPersisterElkConfig()
}

type PersisterElkConfig struct{}

func NewPersisterElkConfig() *PersisterElkConfig {
	return &PersisterElkConfig{}
}

func (c *PersisterElkConfig) ElkAddress() []string {

	return []string{
		getEnv("ELK_ADDRESS", "http://192.168.2.204:9200"),
	}
}

func (c *PersisterElkConfig) Username() string {
	return getEnv("ELK_USERNAME", "elastic")
}

func (c *PersisterElkConfig) Password() string {
	return getEnv("ELK_PASSWORD", "smlSoft2021")
}

func (cfg *Config) CacherConfig() ICacherConfig {
	return NewCacherConfig()
}

type CacherConfig struct{}

func NewCacherConfig() *CacherConfig {
	return &CacherConfig{}
}

func (cfg *CacherConfig) Endpoint() string {
	return getEnv("REDIS_CACHE_URI", "127.0.0.1:6379")
}

func (cfg *CacherConfig) Password() string {
	return ""
}

func (cfg *CacherConfig) DB() int {
	return 0
}

func (cfg *CacherConfig) ConnectionSettings() ICacherConnectionSettings {
	return NewDefaultCacherConnectionSettings()
}

type StorageFileConfig struct{}

func NewStorageFileConfig() *StorageFileConfig {
	return &StorageFileConfig{}
}

func (cfg *StorageFileConfig) StorageDataPath() string {
	return getEnv("STORAGE_DATA_PATH", "")
}

func (cfg *StorageFileConfig) StorageUriAtlas() string {
	return getEnv("STORAGE_DATA_URI", "")
}
