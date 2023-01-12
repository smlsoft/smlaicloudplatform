package main

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/inventory"

	_ "net/http/pprof"
)

func main() {

	cfg := microservice.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	inventory.StartInventoryComsumeCreated(ms, cfg)

	ms.Start()
}
