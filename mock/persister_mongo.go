package mock

import "smlcloudplatform/internal/microservice"

type PersisterMongoConfig struct{}

func (*PersisterMongoConfig) MongodbURI() string {
	return "mongodb://root:rootx@localhost:27017/"
}

func (*PersisterMongoConfig) DB() string {
	return "micro_test"
}

func NewPersisterMongo() microservice.IPersisterMongoConfig {
	return &PersisterMongoConfig{}
}
