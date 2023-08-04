package main

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/config"
	"smlcloudplatform/pkg/product/productbarcode"
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
