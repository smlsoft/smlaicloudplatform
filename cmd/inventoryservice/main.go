package main

import (

	// "github.com/swaggo/echo-swagger"

	// _ "smlcloudplatform/cmd/productapi/docs"
	// "smlcloudplatform/pkg/models"

	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/category"
	"smlcloudplatform/pkg/api/inventory"
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

	categoryHttp := category.NewCategoryHttp(ms, cfg)
	categoryHttp.RouteSetup()
	ms.Start()

}
