package main

import (
	"log"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/authentication"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := microservice.NewConfig()
	ms := microservice.NewMicroservice(cfg)

	svc := authentication.NewAuthenticationHttp(ms, cfg)
	svc.RouteSetup()

	//ms.Echo().GET("/swagger/*", echoSwagger.WrapHandler)

	ms.Start()
}
