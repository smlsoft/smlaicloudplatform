package main

import (

	// "github.com/swaggo/echo-swagger"

	// _ "smlcloudplatform/cmd/productapi/docs"
	// "smlcloudplatform/pkg/models"

	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/inventory"
	"smlcloudplatform/pkg/product/productcategory"
)

func main() {

	cfg := microservice.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	cacher := ms.Cacher(cfg.CacherConfig())
	authService := microservice.NewAuthService(cacher, 24*3)
	ms.HttpMiddleware(authService.MWFuncWithRedis(cacher))

	inventoryapi := inventory.NewInventoryHttp(ms, cfg)
	inventoryapi.RouteSetup()

	categoryHttp := productcategory.NewProductCategoryHttp(ms, cfg)
	categoryHttp.RouteSetup()
	ms.Start()

}
