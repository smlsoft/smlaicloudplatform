package mock

import (
	"fmt"
	"os"
	"regexp"
	"smlcloudplatform/internal/config"

	"github.com/joho/godotenv"
)

type MqConfig struct{}

func (MqConfig) URI() string {
	return os.Getenv("KAFKA_SERVER_URL")
}

func NewMqConfig() config.IMQConfig {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/mock/.env`)

	if err != nil {
		fmt.Println("Load Env Failed ")
	}

	return &MqConfig{}
}
