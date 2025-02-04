package main

import (
	"smlaicloudplatform/internal/config"
	"smlaicloudplatform/internal/product/productbarcode"
	"smlaicloudplatform/pkg/microservice"
)

func main() {
	cfg := config.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	ms.RegisterLivenessProbeEndpoint("/healthz")

	if err != nil {
		panic(err)
	}

	productbarcode.MigrationDatabase(ms, cfg)
	productBarcodeConsumer := productbarcode.NewProductBarcodeConsumer(ms, cfg)
	productBarcodeConsumer.RegisterConsumer()

	ms.Start()
}
