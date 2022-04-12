package images

import (
	"net/http"
	"path/filepath"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/inventory"
	"smlcloudplatform/pkg/models"
	"strconv"
)

type IImageHttp interface {
	RouteSetup()
	GetProductImage(ctx microservice.IContext) error
	UploadImage(ctx microservice.IContext) error
}

type ImagesHttp struct {
	ms      *microservice.Microservice
	cfg     microservice.IConfig
	service IImagesService
}

func NewImagesHttp(
	ms *microservice.Microservice, cfg microservice.IConfig,
	persisterImage *microservice.PersisterImage,
) ImagesHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	inventoryRepo := inventory.NewInventoryRepository(pst)
	imgSrv := NewImageService(persisterImage, inventoryRepo)
	return ImagesHttp{
		ms:      ms,
		cfg:     cfg,
		service: imgSrv,
	}
}

func (svc ImagesHttp) RouteSetup() {

	storageConfig := microservice.NewStorageFileConfig()

	svc.ms.POST("/upload/images", svc.UploadImage)

	svc.ms.GET("/productimage/:shopid/:itemguid", svc.GetProductImage)
	svc.ms.GET("/productimage/:shopid/:itemguid/:index", svc.GetProductImage)

	svc.ms.Echo().Static("/images", storageConfig.StorageDataPath())
	// check config storage

}

func (svc ImagesHttp) GetProductImage(ctx microservice.IContext) error {

	// get image format {shopid}-{itemguid}-{index} ex xxx-xxx-1
	// queryParams := strings.Split(ctx.Param("id"), "-")

	// if len(queryParams) < 2 {
	// 	ctx.Response(http.StatusBadRequest, &models.ApiResponse{
	// 		Success: false,
	// 		Message: "Invalid Payload",
	// 	})
	// 	return nil
	// }

	shopId := ctx.Param("shopid")
	itemguid := ctx.Param("itemguid")
	imageIndex := ctx.Param("index")

	if imageIndex == "" {
		imageIndex = "1"
	}

	// index, err := strconv.Atoi(queryParams[2])
	index, err := strconv.Atoi(imageIndex)
	if err != nil {
		ctx.Response(http.StatusBadRequest, &models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return nil
	}

	if index == 0 {
		index = 1
	}

	fileName, err := svc.service.GetImageByProductCode(shopId, itemguid, index)
	if err != nil {
		ctx.Response(http.StatusBadRequest, &models.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return nil
	}

	storateFileName := filepath.Join(svc.service.GetStoragePath(), fileName)

	return ctx.EchoContext().File(storateFileName)
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
