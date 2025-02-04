package main

import (
	"smlaicloudplatform/internal/config"
	"smlaicloudplatform/internal/migration"
	"smlaicloudplatform/pkg/microservice"
)

func main() {
	cfg := config.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	migration.StartMigrateModel(ms, cfg)
}
