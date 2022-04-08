package mock

import (
	"fmt"
	"os"
	"regexp"
	"smlcloudplatform/internal/microservice"

	"github.com/joho/godotenv"
)

type PersisterPostgresqlConfig struct{}

func NewPersisterPostgresqlConfig() microservice.IPersisterConfig {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/mock/.env`)

	if err != nil {
		fmt.Println("Load Env Failed ")
	}

	return &PersisterPostgresqlConfig{}
}

func (cfg *PersisterPostgresqlConfig) Host() string {
	return os.Getenv("POSTGRES_HOST")
}

func (cfg *PersisterPostgresqlConfig) Port() string {
	return os.Getenv("POSTGRES_PORT")
}

func (cfg *PersisterPostgresqlConfig) DB() string {
	return os.Getenv("POSTGRES_DB_NAME")
}

func (cfg *PersisterPostgresqlConfig) Username() string {
	return os.Getenv("POSTGRES_USERNAME")
}

func (cfg *PersisterPostgresqlConfig) Password() string {
	return os.Getenv("POSTGRES_PASSWORD")
}

func (cfg *PersisterPostgresqlConfig) SSLMode() string {
	sslMode := os.Getenv("POSTGRES_SSL_MODE")
	if sslMode != "" {
		sslMode = "disable"
	}
	return sslMode
}

func (cfg *PersisterPostgresqlConfig) TimeZone() string {
	return "Asia/Bangkok"
}

func (cfg *PersisterPostgresqlConfig) LoggerLevel() string {
	return "debug"
}
