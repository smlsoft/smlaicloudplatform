package main

import (
	"smlcloudplatform/cmd/migrationapi/api"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/config"
)

func main() {

	cfg := config.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)

	if err != nil {
		panic(err)
	}

	ms.RegisterHttp(api.NewMigrationAPI(ms, cfg))

	ms.Start()
}
