package main

import (

	// "github.com/swaggo/echo-swagger"

	// _ "smlcloudplatform/cmd/productapi/docs"
	// "smlcloudplatform/pkg/models"

	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/inventory"
)

func main() {

	cfg := microservice.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	inventoryapi := inventory.NewInventoryHttp(ms, cfg)
	inventoryapi.RouteSetup()
	ms.Start()

}
