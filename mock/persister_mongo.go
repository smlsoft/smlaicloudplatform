package mock

import (
	"fmt"
	"os"
	"regexp"
	"smlcloudplatform/internal/microservice"

	"github.com/joho/godotenv"
)

type PersisterMongoConfig struct{}

func (*PersisterMongoConfig) MongodbURI() string {
	return os.Getenv("MONGODB_URI")
}

func (*PersisterMongoConfig) DB() string {
	return os.Getenv("MONGODB_DB")
}

const projectDirName = "smlcloudplatform"

func NewPersisterMongo() microservice.IPersisterMongoConfig {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/mock/.env`)

	if err != nil {
		fmt.Println("Load Env Failed ")
	}

	return &PersisterMongoConfig{}
}
