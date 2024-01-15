package main

import (
	"smlcloudplatform/internal/config"
	"smlcloudplatform/internal/product/productbarcode"
	"smlcloudplatform/pkg/microservice"
)

func main() {
	cfg := config.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	ms.RegisterLivenessProbeEndpoint("/healthz")

	if err != nil {
		panic(err)
	}

	productbarcode.MigrationProductBarcodeTable(ms, cfg)
	productBarcodeConsumer := productbarcode.NewProductBarcodeConsumer(ms, cfg)

	ms.RegisterConsumer(productBarcodeConsumer)

	ms.Start()
}
