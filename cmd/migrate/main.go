package main

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/migration"
)

func main() {
	cfg := microservice.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	migration.StartMigrateModel(ms, cfg)
}
