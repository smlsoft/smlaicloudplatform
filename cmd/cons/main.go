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
	ms := microservice.NewMicroservice(cfg)

	transaction.StartTransactionComsume(ms, cfg)

	ms.Start()
}
