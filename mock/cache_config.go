package mock

import (
	"fmt"
	"os"
	"regexp"
	"smlcloudplatform/internal/microservice"

	"github.com/joho/godotenv"
)

type CacherConfig struct{}

func NewCacherConfig() *CacherConfig {
	return &CacherConfig{}
}

func (cfg *CacherConfig) Endpoint() string {
	return os.Getenv("REDIS_CACHE_URI")
}

func (cfg *CacherConfig) TLS() bool {
	return false
}

func (cfg *CacherConfig) UserName() string {
	return ""
}

func (cfg *CacherConfig) Password() string {
	return ""
}

func (cfg *CacherConfig) DB() int {
	return 0
}

func (cfg *CacherConfig) ConnectionSettings() microservice.ICacherConnectionSettings {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/mock/.env`)

	if err != nil {
		fmt.Println("Load Env Failed ")
	}
	return microservice.NewDefaultCacherConnectionSettings()
}
