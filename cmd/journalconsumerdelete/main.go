package main

import (
	"os"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/vfgl/journal"
)

func main() {
	cfg := microservice.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	consumerGroupName := os.Getenv("CONSUMER_GROUP_NAME")
	if consumerGroupName == "" {
		consumerGroupName = "00"
	}
	journal.StartJournalComsumeDeleted(ms, cfg, consumerGroupName)
	// purchase.StartPurchaseComsume(ms, cfg)

	ms.Start()
}
