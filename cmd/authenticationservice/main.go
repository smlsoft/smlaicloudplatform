package main

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api"

	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {

	cfg := microservice.NewConfig()
	ms := microservice.NewMicroservice(cfg)

	svc := api.NewAuthenticationService(ms, cfg)
	svc.RouteSetup()

	ms.Echo().GET("/swagger/*", echoSwagger.WrapHandler)

	ms.Start()
	// find by shop_id
}
