package main

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/config"
	"smlcloudplatform/pkg/product/inventory"

	_ "net/http/pprof"
)

func main() {

	cfg := config.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	inventory.StartInventoryComsumeCreated(ms, cfg)

	ms.Start()
}
