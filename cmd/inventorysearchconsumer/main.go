package main

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/config"
)

func main() {
	cfg := config.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	ms.Start()
}
