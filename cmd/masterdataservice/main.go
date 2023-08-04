package main

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/config"
	"smlcloudplatform/pkg/member"
	"smlcloudplatform/pkg/product/inventory"
	"smlcloudplatform/pkg/product/productcategory"
	"time"
)

func main() {

	cfg := config.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	cacher := ms.Cacher(cfg.CacherConfig())
	authService := microservice.NewAuthService(cacher, 24*3*time.Hour, 24*30*time.Hour)
	ms.HttpMiddleware(authService.MWFuncWithRedis(cacher))

	inventoryapi := inventory.NewInventoryHttp(ms, cfg)
	inventoryapi.RegisterHttp()

	categoryHttp := productcategory.NewProductCategoryHttp(ms, cfg)
	categoryHttp.RegisterHttp()

	memberapi := member.NewMemberHttp(ms, cfg)
	memberapi.RegisterHttp()

	ms.Start()

}
