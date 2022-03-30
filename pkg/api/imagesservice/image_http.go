package imageservice

import "smlcloudplatform/internal/microservice"

type IImageService interface {
	RouteSetup()
	GetImage(ctx microservice.IContext) error
	UploadImage(ctx microservice.IContext) error
}

func NewImageHttp(ms *microservice.Microservice, cfg microservice.IConfig) {

}
