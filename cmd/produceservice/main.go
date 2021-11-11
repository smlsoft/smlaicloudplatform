package main

import (
	"fmt"
	"os"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/productservice"
)

func main() {

	fmt.Println("Start Product Service")
	cfg := microservice.NewConfig()

	ms := microservice.NewMicroservice(cfg)

	serviceID := os.Getenv("SERVICE_ID")

	fmt.Printf("Service ID :: %s \n", serviceID)

	productservice.NewProductService(ms, cfg)
	ms.Start()
	// find by shop_id

}
