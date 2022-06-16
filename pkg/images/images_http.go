package images

import (
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/inventory"
	common "smlcloudplatform/pkg/models"
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
	svc.ms.POST("/upload/productimage", svc.UploadImageToProduct)

	svc.ms.GET("/productimage/:shopid/:itemguid", svc.GetProductImage)
	svc.ms.GET("/productimage/:shopid/:itemguid/:index", svc.GetProductImage)

	svc.ms.Echo().Static("/images", storageConfig.StorageDataPath())
	// check config storage

}

func (svc ImagesHttp) GetProductImage(ctx microservice.IContext) error {

	// get image format {shopid}-{itemguid}-{index} ex xxx-xxx-1
	// queryParams := strings.Split(ctx.Param("id"), "-")

	// if len(queryParams) < 2 {
	// 	ctx.Response(http.StatusBadRequest, &common.ApiResponse{
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
		ctx.Response(http.StatusBadRequest, &common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return nil
	}

	if index == 0 {
		index = 1
	}

	fileName, buffer, err := svc.service.GetImageByProductCode(shopId, itemguid, index)

	if err != nil {
		ctx.Response(http.StatusBadRequest, &common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return nil
	}

	// if strings.HasPrefix(fileName, "http") {
	// resp, err := http.Get(fileName)
	// if err != nil {
	// 	ctx.Response(http.StatusBadRequest, &common.ApiResponse{
	// 		Success: false,
	// 		Message: err.Error(),
	// 	})
	// 	return nil
	// }
	// defer resp.Body.Close()

	// downloadedData := &bytes.Buffer{}
	// downloadedData.ReadFrom(resp.Body)

	// return ctx.EchoContext().Blob(http.StatusOK, "", downloadedData.Bytes())
	// }
	if buffer != nil {
		return ctx.EchoContext().Blob(http.StatusOK, "", buffer.Bytes())
	}

	return ctx.EchoContext().File(fileName)
}

// Upload Image
// @Description Update Image
// @Tags		Common
// @Accept 		json
// @Param		file  formData      file  true  "Image"
// @Success		200	{array}	models.Image
// @Failure		401 {object}	common.AuthResponseFailed
// @Failure		400	{object}	common.AuthResponseFailed
// @Failure		500	{object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /upload/images [post]
func (svc ImagesHttp) UploadImage(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopId := userInfo.ShopID
	fileHeader, _ := ctx.FormFile("file")
	image, err := svc.service.UploadImage(shopId, fileHeader)

	if err != nil {
		ctx.Response(http.StatusBadRequest, &common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return nil
	}

	ctx.Response(http.StatusOK, &common.ApiResponse{
		Success: true,
		Data:    image,
	})
	return nil
}

// Upload Image
// @Description Update Image
// @Tags		Common
// @Accept 		json
// @Param		file  formData      file  true  "Image"
// @Success		200	{array}	models.Image
// @Failure		401 {object}	common.AuthResponseFailed
// @Failure		400	{object}	common.AuthResponseFailed
// @Failure		500	{object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /upload/productimage [post]
func (svc ImagesHttp) UploadImageToProduct(ctx microservice.IContext) error {

	fileHeader, _ := ctx.FormFile("file")
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	// find
	err := svc.service.UploadImageToProduct(shopID, fileHeader)

	if err != nil {
		ctx.Response(http.StatusBadRequest, &common.ApiResponse{
			Success: false,
			Message: err.Error(),
		})
		return nil
	}

	ctx.Response(http.StatusOK, &common.ApiResponse{
		Success: true,
	})
	return nil
}
