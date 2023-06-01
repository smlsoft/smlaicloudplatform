package main

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/transaction/purchase"
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

	purchase := purchase.NewPurchaseHttp(ms, cfg)
	purchase.RouteSetup()

	ms.Start()
}
