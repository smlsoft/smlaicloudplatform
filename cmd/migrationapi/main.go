package main

import (
	"smlaicloudplatform/cmd/migrationapi/api"
	"smlaicloudplatform/internal/config"
	"smlaicloudplatform/pkg/microservice"
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
