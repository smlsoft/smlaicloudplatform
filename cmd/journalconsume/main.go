package main

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/purchase"
)

func main() {
	cfg := microservice.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	// journal.StartJournalComsumeCreated(ms, cfg, "")
	purchase.StartPurchaseComsume(ms, cfg)
}
