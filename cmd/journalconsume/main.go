package main

import (
	"smlaicloudplatform/internal/config"
	"smlaicloudplatform/pkg/microservice"
)

func main() {
	cfg := config.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	// journal.StartJournalComsumeCreated(ms, cfg, "journal-bulk-created-group-1")
	// purchase.StartPurchaseComsume(ms, cfg)

	ms.Start()
}
