package main

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/inventoryimport"
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

	invImp := inventoryimport.NewInventoryImportHttp(ms, cfg)
	invImp.RouteSetup()

	invOptionImp := inventoryimport.NewInventoryImporOptionMaintHttp(ms, cfg)
	invOptionImp.RouteSetup()

	catImp := inventoryimport.NewCategoryImportHttp(ms, cfg)
	catImp.RouteSetup()

	// categoryHttp := category.NewCategoryHttp(ms, cfg)
	// categoryHttp.RouteSetup()
	ms.Start()
}
