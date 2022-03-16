package main

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/shop"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	cfg := microservice.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	svc := shop.NewShopHttp(ms, cfg)

	svc.RouteSetup()

	//ms.Echo().GET("/swagger/*", echoSwagger.WrapHandler)

	ms.Start()
}
