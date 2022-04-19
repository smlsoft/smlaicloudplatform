package main

import (
	"smlcloudplatform/internal/microservice"
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

	h := inventory.NewInventoryImportHttp(ms, cfg)
	h.RouteSetup()
	// categoryHttp := category.NewCategoryHttp(ms, cfg)
	// categoryHttp.RouteSetup()
	ms.Start()
}
