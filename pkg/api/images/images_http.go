package images

import (
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
)

type IImageHttp interface {
	RouteSetup()
	GetImage(ctx microservice.IContext) error
	UploadImage(ctx microservice.IContext) error
}

type ImagesHttp struct {
	ms      *microservice.Microservice
	cfg     microservice.IConfig
	service IImagesService
}

func NewImagesHttp(ms *microservice.Microservice, cfg microservice.IConfig, persisterImage *microservice.PersisterImage) ImagesHttp {

	imgSrv := NewImageService(persisterImage)
	return ImagesHttp{
		ms:      ms,
		cfg:     cfg,
		service: imgSrv,
	}
}

func (svc ImagesHttp) RouteSetup() {

	storageConfig := microservice.NewStorageFileConfig()

	svc.ms.POST("/upload/images", svc.UploadImage)
	svc.ms.Echo().Static("/images", storageConfig.StorageDataPath())
	// check config storage

}

// Upload Image
// @Description Update Image
// @Tags		Common
// @Accept 		json
// @Param		file  formData      file  true  "Image"
// @Success		200	{array}	models.Image
// @Failure		401 {object}	models.AuthResponseFailed
// @Failure		400	{object}	models.AuthResponseFailed
// @Failure		500	{object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /upload/images [post]
func (svc ImagesHttp) UploadImage(ctx microservice.IContext) error {

	fileHeader, _ := ctx.FormFile("file")
	image, err := svc.service.UploadImage(fileHeader)

	if err != nil {
		ctx.Response(http.StatusBadRequest, &models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return nil
	}

	ctx.Response(http.StatusOK, &models.ApiResponse{
		Success: true,
		Data:    image,
	})
	return nil
}
