package api

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

}

func (svc *ImageService) UploadImage(ctx microservice.IContext) error {

	return nil
}
