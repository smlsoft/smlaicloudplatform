package main

import (
	"log"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/transaction"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := microservice.NewConfig()
	ms, _ := microservice.NewMicroservice(cfg)

	transaction.StartTransactionComsumeCreated(ms, cfg)

	ms.Start()
}
