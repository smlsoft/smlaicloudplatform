package main

import (
	"log"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/purchase"
	"smlcloudplatform/pkg/api/saleinvoice"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := microservice.NewConfig()
	ms, _ := microservice.NewMicroservice(cfg)

	saleinvoice.StartSaleinvoiceComsumeCreated(ms, cfg)

	purchase.StartPurchaseComsume(ms, cfg)

	ms.Start()
}
