package main

import (
	"smlaicloudplatform/internal/config"
	"smlaicloudplatform/pkg/microservice"
)

func main() {
	cfg := config.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	ms.Start()
}
