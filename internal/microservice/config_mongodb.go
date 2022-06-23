package microservice

import "fmt"

type MongoPersisterConfig struct{}

func NewMongoPersisterConfig() *MongoPersisterConfig {
	return &MongoPersisterConfig{}
}
func (cfg *MongoPersisterConfig) MongodbURI() string {

	uri := getEnv("MONGODB_URI", "") // mongodb://root:rootx@localhost:27017/
	if uri != "" {
		return uri
	}

	userNamePassword := fmt.Sprintf("%s:%s@", cfg.MongodbUserName(), cfg.MongodbPassWord())
	if userNamePassword == ":@" {
		userNamePassword = ""
	}

	connectionUri := fmt.Sprintf("%s://%s%s:%s/",
		cfg.MongodbProtocal(),
		userNamePassword, cfg.MongodbServer(),
		cfg.MongodbPort())
	return connectionUri
}

func (cfg *MongoPersisterConfig) MongodbProtocal() string {
	return getEnv("MONGODB_SERVER", "")
}

func (cfg *MongoPersisterConfig) MongodbServer() string {
	return getEnv("MONGODB_SERVER", "")
}

func (cfg *MongoPersisterConfig) MongodbPort() string {
	return getEnv("MONGODB_PORT", "27017")
}

func (cfg *MongoPersisterConfig) DB() string {
	return getEnv("MONGODB_DB", "smldev")
}

func (cfg *MongoPersisterConfig) MongodbUserName() string {
	return getEnv("MONGODB_USERNAME", "")
}

func (cfg *MongoPersisterConfig) MongodbPassWord() string {
	return getEnv("MONGODB_PASSWORD", "")
}
