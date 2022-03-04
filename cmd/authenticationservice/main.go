package main

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/authentication"

	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {

	cfg := microservice.NewConfig()
	ms := microservice.NewMicroservice(cfg)

	svc := authentication.NewAuthenticationHttp(ms, cfg)
	svc.RouteSetup()

	ms.Echo().GET("/swagger/*", echoSwagger.WrapHandler)

	ms.Start()
}
