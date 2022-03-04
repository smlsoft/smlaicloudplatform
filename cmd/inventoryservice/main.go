package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"

	// "github.com/swaggo/echo-swagger"

	// _ "smlcloudplatform/cmd/productapi/docs"
	// "smlcloudplatform/pkg/models"
	"smlcloudplatform/api/swagger"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/inventory"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := microservice.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	inventoryapi := inventory.NewInventoryHttp(ms, cfg)
	inventoryapi.RouteSetup()

	//ms.Echo().GET("/swagger/*", echoSwagger.WrapHandler)

	ms.Echo().GET("/swagger/*", swagger.EchoWrapHandler)
	fmt.Print("Start Product Service")
	ms.Start()
	// find by shop_id

}
