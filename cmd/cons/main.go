package main

import (
	"log"
	"smlaicloudplatform/internal/config"
	"smlaicloudplatform/pkg/microservice"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.NewConfig()
	ms, _ := microservice.NewMicroservice(cfg)

	// saleinvoice.StartSaleinvoiceComsumeCreated(ms, cfg, "")

	// purchase.StartPurchaseComsume(ms, cfg)

	ms.Start()
}
