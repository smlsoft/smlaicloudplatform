package main

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/category"
	"smlcloudplatform/pkg/api/inventory"
	"smlcloudplatform/pkg/api/member"
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

	memberapi := member.NewMemberHttp(ms, cfg)
	memberapi.RouteSetup()

	ms.Start()

}
