package imageservice

import "smlcloudplatform/internal/microservice"

type ImageService struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
}

func NewImageService(ms *microservice.Microservice, cfg microservice.IConfig) *ImageService {

	return &ImageService{
		ms:  ms,
		cfg: cfg,
	}
}

func (svc *ImageService) RouteSetup() {

	svc.ms.GET("/image/:id", svc.getImageHTTP)
	svc.ms.POST("/image/upload", svc.UploadImageHTTP)
}

func (svc *ImageService) UploadImage(ctx microservice.IContext) error {

	return nil
}
