package main

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/migration"
	"smlcloudplatform/pkg/api/syncdata"
)

func main() {

	cfg := microservice.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	migration.StartMigrateModel(ms, cfg)

	cacher := ms.Cacher(cfg.CacherConfig())
	authService := microservice.NewAuthService(cacher, 24*3)

	ms.HttpMiddleware(authService.MWFuncWithRedis(cacher))

	syncapi := syncdata.NewSyncDataHttp(ms, cfg)
	syncapi.RouteSetup()

	ms.Start()
}
