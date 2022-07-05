package main

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/inventorysearchconsumer"
)

func main() {
	cfg := microservice.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	inventorysearchconsumer.StartInventorySearchComsumerOnProductCreated(ms, cfg)
	inventorysearchconsumer.StartInventorySearchComsumerOnProductUpdated(ms, cfg)
	inventorysearchconsumer.StartInventorySearchComsumerOnProductDeleted(ms, cfg)

	ms.Start()
}
