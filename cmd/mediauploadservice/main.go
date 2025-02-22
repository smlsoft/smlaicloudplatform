package main

import (
	"smlaicloudplatform/internal/config"
	"smlaicloudplatform/internal/media"
	"smlaicloudplatform/pkg/microservice"
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

	ms.RegisterHttp(media.InitMediaUploadHttp(ms, cfg))
	ms.Start()
}
