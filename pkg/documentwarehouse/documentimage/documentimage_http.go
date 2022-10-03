package documentimage

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/documentwarehouse/documentimage/models"
	"smlcloudplatform/pkg/documentwarehouse/documentimage/repositories"
	"smlcloudplatform/pkg/documentwarehouse/documentimage/services"

	journalRepo "smlcloudplatform/pkg/vfgl/journal/repositories"
	journalSvc "smlcloudplatform/pkg/vfgl/journal/services"

	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
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

	svcWsJournal journalSvc.IJournalWebsocketService
}

func NewDocumentImageHttp(ms *microservice.Microservice, cfg microservice.IConfig) *DocumentImageHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	repo := repositories.NewDocumentImageRepository(pst)

	azureblob := microservice.NewPersisterAzureBlob()
	svc := services.NewDocumentImageService(repo, azureblob)

	docImageRepo := repositories.NewDocumentImageRepository(pst)

	cacheRepo := journalRepo.NewJournalCacheRepository(cache)
	svcWsJournal := journalSvc.NewJournalWebsocketService(docImageRepo, cacheRepo, time.Duration(30)*time.Minute)

	return &DocumentImageHttp{
		ms:           ms,
		cfg:          cfg,
		service:      svc,
		svcWsJournal: svcWsJournal,
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
	h.ms.PUT("/documentimage/documentref/status/:docref", h.UpdateDocumentImageStatusByDocumentRef)

	h.ms.GET("/documentimagegroup", h.ListDocumentImageGroup)
	h.ms.GET("/documentimagegroup/:docref", h.GetDocumentImageGroup)
	h.ms.POST("/documentimagegroup", h.SaveDocumentImageGroup)
}

// List Document Image
// @Summary		List Document Image
// @Description List Document Image
// @Tags		DocumentImage
// @Param		q		query	string		false  "Search Value"
// @Param		status		query	integer		false  "Status Value"
// @Param		module		query	string		false  "Module Value"
// @Param		docguidref		query	string		false  "Doc GUID Ref Value"
// @Param		documentref		query	string		false  "Document Ref Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Param		sort	query	string		false  "Sort"
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
	sorts := utils.GetSortParam(ctx.QueryParam)

	matchFilters := map[string]interface{}{}

	statusRaw := strings.TrimSpace(ctx.QueryParam("status"))

	statusFilter := []int{}
	if len(statusRaw) > 0 {

		statusRawArray := strings.Split(statusRaw, ",")

		for _, status := range statusRawArray {
			tempStatus, err := strconv.Atoi(status)
			if err == nil {
				statusFilter = append(statusFilter, tempStatus)
			}
		}
	}

	lenStatus := len(statusFilter)
	if lenStatus > 0 {

		if lenStatus == 1 {
			matchFilters["status"] = statusFilter[0]
		} else {
			matchFilters["status"] = bson.M{"$in": statusFilter}
		}
	}

	module := strings.TrimSpace(ctx.QueryParam("module"))

	if len(module) > 0 {
		matchFilters["module"] = module
	}

	docGuidRef := strings.TrimSpace(ctx.QueryParam("docguidref"))

	if len(docGuidRef) > 0 {
		matchFilters["docguidref"] = docGuidRef
	}

	documentRef := strings.TrimSpace(ctx.QueryParam("documentref"))

	docRefReserve := strings.TrimSpace(ctx.QueryParam("docref-reserve"))

	if len(docRefReserve) > 0 && docRefReserve != "0" {
		docRefPoolList, err := h.svcWsJournal.GetAllDocRefPool(shopID)

		if err == nil {
			docRefList := []string{}
			for docRef := range docRefPoolList {
				docRefList = append(docRefList, docRef)
			}
			matchFilters["documentref"] = bson.M{"$eq": documentRef, "$nin": docRefList}
		}

	} else if len(documentRef) > 0 {
		matchFilters["documentref"] = documentRef
	}

	docList, pagination, err := h.service.SearchDocumentImage(shopID, matchFilters, q, page, limit, sorts)
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
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.DocumentImageStatus{}
	err := json.Unmarshal([]byte(input), &docReq)
	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.service.UpdateDocumentImageStatus(shopID, id, docReq.DocGUIDRef, docReq.Status)
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

// Update Document Image Status By Document Ref godoc
// @Summary		Update Document Image Status By Document Ref
// @Description Update Document Image Status By Document Ref
// @Tags		DocumentImageStatusByDocumentRef
// @Param		docref  path      string  true  "Document Ref"
// @Param		DocumentImageStatus  body      models.DocumentImageStatus  true  "DocumentImageStatus"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /documentimage/documentref/status/{docref} [put]
func (h DocumentImageHttp) UpdateDocumentImageStatusByDocumentRef(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	docref := ctx.Param("docref")
	input := ctx.ReadInput()

	docReq := &models.DocumentImageStatus{}
	err := json.Unmarshal([]byte(input), &docReq)
	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.service.UpdateDocumentImageStatusByDocumentRef(shopID, docref, docReq.DocGUIDRef, docReq.Status)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		ID:      docref,
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

// List Document Image Group
// @Description Get Document Image Group
// @Tags		DocumentImageGroup
// @Accept 		json
// @Tags		Restaurant
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /documentimagegroup [get]
func (h DocumentImageHttp) ListDocumentImageGroup(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)

	matchFilters := map[string]interface{}{}

	docRefReserve := strings.TrimSpace(ctx.QueryParam("docref-reserve"))

	status := strings.TrimSpace(ctx.QueryParam("status"))

	if len(status) > 0 {
		tempStatus, err := strconv.Atoi(status)
		if err == nil {
			matchFilters["status"] = tempStatus
		}
	}

	if len(docRefReserve) > 0 && docRefReserve != "0" {
		docRefPoolList, err := h.svcWsJournal.GetAllDocRefPool(shopID)

		if err == nil {
			docRefList := []string{}
			for docRef := range docRefPoolList {
				docRefList = append(docRefList, docRef)
			}
			matchFilters["documentref"] = bson.M{"$nin": docRefList}
		}

	}

	docList, pagination, err := h.service.ListDocumentImageDocRefGroup(shopID, matchFilters, q, page, limit)
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

// Get Document Image Group
// @Description Get Document Image Group
// @Tags		DocumentImageGroup
// @Accept 		json
// @Tags		Restaurant
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /documentimagegroup/{docref} [get]
func (h DocumentImageHttp) GetDocumentImageGroup(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	docRef := ctx.Param("docref")

	doc, err := h.service.GetDocumentImageDocRefGroup(shopID, docRef)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    doc,
	})

	return nil
}

// Save Document Image Group
// @Description Save Document Image Group
// @Tags		DocumentImageGroup
// @Param		DocumentImageGroup  body      models.DocumentImageGroup  true  "Document Image Group"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /documentimagegroup [post]
func (h DocumentImageHttp) SaveDocumentImageGroup(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	docImages := &models.DocumentImageGroupRequest{}

	err := json.Unmarshal([]byte(input), &docImages)
	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.service.SaveDocumentImageDocRefGroup(shopID, docImages.DocumentRef, docImages.DocumentImages)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}
