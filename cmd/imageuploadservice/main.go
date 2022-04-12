package main

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/images"
)

func main() {

	cfg := microservice.NewConfig()
	ms, err := microservice.NewMicroservice(cfg)
	if err != nil {
		panic(err)
	}

	publicPath := []string{
		"/images/*",
		"/productimage/",
	}

	cacher := ms.Cacher(cfg.CacherConfig())
	authService := microservice.NewAuthService(cacher, 24*3)
	ms.HttpMiddleware(authService.MWFuncWithRedis(cacher, publicPath...))

	filePersister := microservice.NewPersisterFile(microservice.NewStorageFileConfig())
	imagePersister := microservice.NewPersisterImage(filePersister)

	ms.Logger.Debugf("Store File Path %v", filePersister.StoreFilePath)
	imageHttp := images.NewImagesHttp(ms, cfg, imagePersister)
	imageHttp.RouteSetup()

	ms.Start()
}
