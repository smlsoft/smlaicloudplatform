package main

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/authentication"
)

func main() {

	cfg := microservice.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	svc := authentication.NewAuthenticationHttp(ms, cfg)
	svc.RouteSetup()

	// ms.Echo().GET("/swagger/*", echoSwagger.WrapHandler)

	ms.Start()
}
