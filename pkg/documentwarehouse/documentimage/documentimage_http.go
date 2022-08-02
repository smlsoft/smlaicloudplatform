package documentimage

import (
	"encoding/json"
	"errors"
	"net/http"

	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/documentwarehouse/documentimage/models"
	"smlcloudplatform/pkg/documentwarehouse/documentimage/repositories"
	"smlcloudplatform/pkg/documentwarehouse/documentimage/services"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
)

type IDocumentImageHttp interface {
	RouteSetup()
	SearchDocumentImage(ctx microservice.IContext) error
	GetDocumentImageInfo(ctx microservice.IContext) error
	UploadDocumentImage(ctx microservice.IContext) error
	UpdateDocumentImage(ctx microservice.IContext) error
	DeleteDocumentImage(ctx microservice.IContext) error
}

type DocumentImageHttp struct {
	ms      *microservice.Microservice
	cfg     microservice.IConfig
	service services.IDocumentImageService
}

func NewDocumentImageHttp(ms *microservice.Microservice, cfg microservice.IConfig) *DocumentImageHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	repo := repositories.NewDocumentImageRepository(pst)

	azureblob := microservice.NewPersisterAzureBlob()
	svc := services.NewDocumentImageService(repo, azureblob)

	return &DocumentImageHttp{
		ms:      ms,
		cfg:     cfg,
		service: svc,
	}
}

func (h DocumentImageHttp) RouteSetup() {
	h.ms.GET("/documentimage", h.SearchDocumentImage)
	h.ms.GET("/documentimage/:id", h.GetDocumentImageInfo)
	h.ms.POST("/documentimage/upload", h.UploadDocumentImage)
	h.ms.POST("/documentimage", h.CreateDocumentImage)
	h.ms.PUT("/documentimage/status/:id", h.UpdateDocumentImageStatus)
	h.ms.PUT("/documentimage/:id", h.UpdateDocumentImage)
	h.ms.DELETE("/documentimage/:id", h.DeleteDocumentImage)
}

// List Document Image
// @Summary		List Document Image
// @Description List Document Image
// @Tags		DocumentImage
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Accept 		json
// @Success		200	{object}	models.DocumentImagePageResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /documentimage [get]
func (h DocumentImageHttp) SearchDocumentImage(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)
	//sort := utils.GetSortParam(ctx.QueryParam)
	docList, pagination, err := h.service.SearchDocumentImage(shopID, q, page, limit)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success:    true,
		Data:       docList,
		Pagination: pagination,
	})

	return nil
}

// Get Document Image Infomation godoc
// @Summary		Get Document Image Infomation
// @Description Get Document Image Infomation
// @Tags		DocumentImage
// @Param		id  path      string  true  "Id"
// @Accept 		json
// @Success		200	{object}	models.DocumentImageInfoResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /documentimage/{id} [get]
func (h DocumentImageHttp) GetDocumentImageInfo(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	doc, err := h.service.InfoDocumentImage(shopID, id)
	if err != nil {
		h.ms.Logger.Errorf("Error getting document %v: %v", id, err)
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    doc,
	})
	return nil
}

// Create Document Image godoc
// @Summary		Create Document Image
// @Description Create Document Image
// @Tags		DocumentImage
// @Param		DocumentImage  body      models.DocumentImage  true  "DocumentImage"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /documentimage [post]
func (h DocumentImageHttp) CreateDocumentImage(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.DocumentImage{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.service.CreateDocumentImage(shopID, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		ID:      idx,
	})
	return nil
}

// Update Document Image godoc
// @Summary		Update Document Image
// @Description Update Document Image
// @Tags		DocumentImage
// @Param		id  path      string  true  "ID"
// @Param		DocumentImage  body      models.DocumentImage  true  "DocumentImage"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /documentimage/{id} [put]
func (h DocumentImageHttp) UpdateDocumentImage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.DocumentImage{}
	err := json.Unmarshal([]byte(input), &docReq)
	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.service.UpdateDocumentImage(shopID, id, authUsername, *docReq)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		ID:      id,
	})
	return nil
}

// Update Document Image Status godoc
// @Summary		Update Document Image Status
// @Description Update Document Image Status
// @Tags		DocumentImageStatus
// @Param		id  path      string  true  "ID"
// @Param		DocumentImageStatus  body      models.DocumentImageStatus  true  "DocumentImageStatus"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /documentimage/status/{id} [put]
func (h DocumentImageHttp) UpdateDocumentImageStatus(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.DocumentImageStatus{}
	err := json.Unmarshal([]byte(input), &docReq)
	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.service.UpdateDocumentImageStatus(shopID, id, authUsername, *docReq)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		ID:      id,
	})
	return nil
}

// Delete Document Image godoc
// @Summary		Delete Document Image
// @Description Delete Document Image
// @Tags		GL
// @Param		id  path      string  true  "ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /documentimage/{id} [delete]
func (h DocumentImageHttp) DeleteDocumentImage(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.service.DeleteDocumentImage(shopID, id, authUsername)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		ID:      id,
	})
	return nil
}

// Upload Document Image
// @Description Upload Document Image
// @Tags		DocumentImage
// @Accept 		json
// @Param		module  query      string  true  "Module"
// @Param		file  formData      file  true  "Image"
// @Success		200	{array}	models.DocumentImageInfoResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Failure		400	{object}	common.AuthResponseFailed
// @Failure		500	{object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /documentimage/upload [post]
func (h DocumentImageHttp) UploadDocumentImage(ctx microservice.IContext) error {

	moduleName := ctx.QueryParam("module")
	if moduleName == "" {
		ctx.ResponseError(400, "No Module Special")
		return errors.New("Upload Image Without Module")
	}

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	idx, err := h.service.UploadDocumentImage(shopID, authUsername, moduleName, fileHeader)

	if err != nil {
		ctx.ResponseError(400, err.Error())
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		Data:    idx,
	})
	return nil
}
