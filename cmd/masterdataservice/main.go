package main

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/config"
	"smlcloudplatform/pkg/member"
	"smlcloudplatform/pkg/product/inventory"
	"smlcloudplatform/pkg/product/productcategory"
)

func main() {

	cfg := config.NewConfig()
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

	memberapi := member.NewMemberHttp(ms, cfg)
	memberapi.RouteSetup()

	ms.Start()

}
