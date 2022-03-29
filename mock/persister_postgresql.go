package mock

import (
	"fmt"
	"os"
	"regexp"
	"smlcloudplatform/internal/microservice"

	"github.com/joho/godotenv"
)

type PersisterConfig struct{}

func NewPersister() microservice.IPersisterConfig {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/mock/.env`)

	if err != nil {
		fmt.Println("Load Env Failed ")
	}

	return &PersisterConfig{}
}

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
	return "Asia/Bangkok"
}
