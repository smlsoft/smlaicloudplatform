package main

import (
	"smlaicloudplatform/internal/config"
	"smlaicloudplatform/internal/images"
	"smlaicloudplatform/pkg/microservice"
	"time"
)

func main() {

	cfg := config.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	publicPath := []string{
		"/images*",
		"/productimage",
	}

	cacher := ms.Cacher(cfg.CacherConfig())
	authService := microservice.NewAuthService(cacher, 24*3*time.Hour, 24*30*time.Hour)
	ms.HttpMiddleware(authService.MWFuncWithRedis(cacher, publicPath...))

	filePersister := microservice.NewPersisterFile(config.NewStorageFileConfig())
	imagePersister := microservice.NewPersisterImage(filePersister)

	ms.Logger.Debugf("Store File Path %v", filePersister.StoreFilePath)
	imageHttp := images.NewImagesHttp(ms, cfg, imagePersister)
	imageHttp.RegisterHttp()

	ms.Start()
}
