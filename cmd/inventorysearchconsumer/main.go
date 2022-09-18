package main

import (
	"smlcloudplatform/internal/microservice"
)

func main() {
	cfg := microservice.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	ms.Start()
}
