package documentimage

import (
	"encoding/json"
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
	Module       string
	ms           *microservice.Microservice
	cfg          microservice.IConfig
	service      services.IDocumentImageService
	svcWsJournal journalSvc.IJournalWebsocketService
}

func NewDocumentImageHttp(ms *microservice.Microservice, cfg microservice.IConfig) *DocumentImageHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	repo := repositories.NewDocumentImageRepository(pst)
	repoImageGroup := repositories.NewDocumentImageGroupRepository(pst)

	azureblob := microservice.NewPersisterAzureBlob()
	svc := services.NewDocumentImageService(repo, repoImageGroup, azureblob)

	docImageRepo := repositories.NewDocumentImageRepository(pst)

	cacheRepo := journalRepo.NewJournalCacheRepository(cache)
	svcWsJournal := journalSvc.NewJournalWebsocketService(docImageRepo, cacheRepo, time.Duration(30)*time.Minute)

	return &DocumentImageHttp{
		Module:       "GL",
		ms:           ms,
		cfg:          cfg,
		service:      svc,
		svcWsJournal: svcWsJournal,
	}
}

func (h DocumentImageHttp) RouteSetup() {
	h.ms.GET("/documentimage", h.SearchDocumentImage)
	// h.ms.GET("/documentimage/special", h.DocumentImageSpecial)
	h.ms.GET("/documentimage/:guid", h.GetDocumentImageInfo)
	h.ms.POST("/documentimage/upload", h.UploadDocumentImage)
	h.ms.POST("/documentimage", h.CreateDocumentImage)
	h.ms.POST("/documentimage/bulk", h.BulkCreateDocumentImage)
	// h.ms.PUT("/documentimage/status/:id", h.UpdateDocumentImageStatus)
	h.ms.PUT("/documentimage/:guid", h.UpdateDocumentImage)
	h.ms.PUT("/documentimage/:guid/reject", h.UpdateDocumentImageReject)
	h.ms.DELETE("/documentimage/:guid", h.DeleteDocumentImage)
	// h.ms.PUT("/documentimage/documentref/status/:docref", h.UpdateDocumentImageStatusByDocumentRef)

	h.ms.GET("/documentimagegroup", h.ListDocumentImageGroup)
	h.ms.GET("/documentimagegroup/:guid", h.GetDocumentImageGroup)
	h.ms.POST("/documentimagegroup", h.CreateDocumentImageGroup)
	h.ms.PUT("/documentimagegroup/:guid", h.UpdateDocumentImageGroup)
	h.ms.PUT("/documentimagegroup/:guid/documentimages", h.UpdateImageReferenceByDocumentImageGroup)
	h.ms.PUT("/documentimagegroup/:guid/reference", h.UpdateReferenceByDocumentImageGroup)
	h.ms.PUT("/documentimagegroup/:guid/ungroup", h.UngroupDocumentImageGroup)
	h.ms.PUT("/documentimagegroup/:guid/images", h.UpdateDocumentImageGroup)

	h.ms.GET("/documentimagegroup/docref/:docref", h.GetDocumentImageGroupByDocRefInfo)
}

func (h DocumentImageHttp) DocumentImageSpecial(ctx microservice.IContext) error {

	err := h.service.UpdateDocumentImageRederenceGroup()
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})
	return nil
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
// @Param		guid  path      string  true  "document image guid"
// @Accept 		json
// @Success		200	{object}	models.DocumentImageInfoResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /documentimage/{guid} [get]
func (h DocumentImageHttp) GetDocumentImageInfo(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("guid")
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

	idx, imageGroupGUID, err := h.service.CreateDocumentImage(shopID, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, map[string]interface{}{
		"success": true,
		"id":      idx,
		"groupid": imageGroupGUID,
	})
	return nil
}

// Bulk Create Document Image godoc
// @Summary		Create Document Image
// @Description Create Document Image
// @Tags		DocumentImage
// @Param		DocumentImage  body      []models.DocumentImage  true  "DocumentImage"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /documentimage/bulk [post]
func (h DocumentImageHttp) BulkCreateDocumentImage(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &[]models.DocumentImage{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.service.BulkCreateDocumentImage(shopID, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
	})
	return nil
}

// Update Document Image godoc
// @Summary		Update Document Image
// @Description Update Document Image
// @Tags		DocumentImage
// @Param		guid  path      string  true  "document image guid"
// @Param		DocumentImage  body      models.DocumentImage  true  "DocumentImage"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /documentimage/{guid} [put]
func (h DocumentImageHttp) UpdateDocumentImage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("guid")
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

// Update Reject Status Document Image godoc
// @Summary		Update Reject Document Image
// @Description Update Reject Document Image
// @Tags		DocumentImage
// @Param		guid  path      string  true  "document image guid"
// @Param		RequestDocumentImageReject  body      models.RequestDocumentImageReject  true  "Document Image Reject"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /documentimage/{guid}/reject [put]
func (h DocumentImageHttp) UpdateDocumentImageReject(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	guid := ctx.Param("guid")
	input := ctx.ReadInput()

	docReq := &models.RequestDocumentImageReject{}
	err := json.Unmarshal([]byte(input), &docReq)
	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.service.UpdateDocumentImageReject(shopID, guid, authUsername, docReq.IsReject)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		ID:      guid,
	})
	return nil
}

/*
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
*/

/*
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
*/

// Delete Document Image godoc
// @Summary		Delete Document Image
// @Description Delete Document Image
// @Tags		DocumentImage
// @Param		guid  path      string  true  "document image guid"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /documentimage/{guid} [delete]
func (h DocumentImageHttp) DeleteDocumentImage(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("guid")

	err := h.service.DeleteDocumentImage(shopID, authUsername, id)
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

	// moduleName := ctx.QueryParam("module")
	// if moduleName == "" {
	// 	ctx.ResponseError(400, "No Module Special")
	// 	return errors.New("Upload Image Without Module")
	// }

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	idx, err := h.service.UploadDocumentImage(shopID, authUsername, fileHeader)

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
// @Param		reserve	query	integer		false  "เอกสารที่มีการจอง,0 not filter, 1 filter"
// @Param		reject	query	integer		false  "empty not filter, 0 not reject, 1 reject"
// @Param		ref	query	integer		false  "document reference: empty not filter, 1 not reference, 2 referenced"
// @Param		fromdate		query	string		false  "From Date (YYYY-MM-DD)"
// @Param		todate		query	string		false  "To Date (YYYY-MM-DD)"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /documentimagegroup [get]
func (h DocumentImageHttp) ListDocumentImageGroup(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetSearchPageable(ctx.QueryParam)

	matchFilters := map[string]interface{}{}

	docRefReserve := strings.TrimSpace(ctx.QueryParam("reserve"))
	isreject := strings.TrimSpace(ctx.QueryParam("reject"))
	isref := strings.TrimSpace(ctx.QueryParam("ref"))
	path := strings.TrimSpace(ctx.QueryParam("path"))

	fromDateStr := strings.TrimSpace(ctx.QueryParam("fromdate"))
	toDateStr := strings.TrimSpace(ctx.QueryParam("todate"))

	if len(isreject) > 0 {
		tempStatus, err := strconv.Atoi(isreject)
		if err == nil {
			matchFilters["isreject"] = tempStatus != 0
		}
	}

	if len(isref) > 0 && isref != "0" {
		if isref == "1" {
			matchFilters["references.module"] = bson.M{"$ne": h.Module}
		} else if isref == "2" {
			matchFilters["references.module"] = h.Module
		}
	}

	documentImageGUID := strings.TrimSpace(ctx.QueryParam("documentimageguid"))

	if len(documentImageGUID) > 0 {
		matchFilters["imagereferences.documentimageguid"] = documentImageGUID
	}

	if len(docRefReserve) > 0 && docRefReserve != "0" {
		docRefPoolList, err := h.svcWsJournal.GetAllDocRefPool(shopID)

		if err == nil {
			docRefList := []string{}
			for docRef := range docRefPoolList {
				docRefList = append(docRefList, docRef)
			}
			matchFilters["guidfixed"] = bson.M{"$nin": docRefList}
		}

	}

	if len(fromDateStr) > 0 && len(toDateStr) > 0 {
		fromDate, err1 := time.Parse("2006-01-02", fromDateStr)
		toDate, err2 := time.Parse("2006-01-02", toDateStr)

		if err1 == nil && err2 == nil {
			matchFilters["uploadedat"] = bson.M{
				"$gte": fromDate,
				"$lt":  toDate.AddDate(0, 0, 1),
			}
		}
	}

	if len(path) > 0 {

		matchFilters["path"] = path
	}

	docList, pagination, err := h.service.ListDocumentImageGroup(shopID, matchFilters, pageable)
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
// @Param		guid  path      string  true  "document image group guid"
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /documentimagegroup/{guid} [get]
func (h DocumentImageHttp) GetDocumentImageGroup(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	docImageGroupGUID := ctx.Param("guid")

	doc, err := h.service.GetDocumentImageDocRefGroup(shopID, docImageGroupGUID)
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
func (h DocumentImageHttp) CreateDocumentImageGroup(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	input := ctx.ReadInput()

	docImageGroup := &models.DocumentImageGroup{}

	err := json.Unmarshal([]byte(input), &docImageGroup)
	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.service.CreateDocumentImageGroup(shopID, authUsername, *docImageGroup)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		ID:      idx,
	})

	return nil
}

// Update Document Image Group
// @Description Update Document Image Group
// @Tags		DocumentImageGroup
// @Accept 		json
// @Param		guid  path      string  true  "document image group guid"
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /documentimagegroup/{guid} [put]
func (h DocumentImageHttp) UpdateDocumentImageGroup(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()

	docImageGroupGUID := ctx.Param("guid")

	input := ctx.ReadInput()

	docImageGroup := &models.DocumentImageGroup{}

	err := json.Unmarshal([]byte(input), &docImageGroup)
	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.service.UpdateDocumentImageGroup(userInfo.ShopID, userInfo.Username, docImageGroupGUID, *docImageGroup)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}

// Update Document Image Group
// @Description Update Document Image Group
// @Tags		DocumentImageGroup
// @Accept 		json
// @Param		guid  path      string  true  "document image group guid"
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /documentimagegroup/{guid}/documentimages [put]
func (h DocumentImageHttp) UpdateImageReferenceByDocumentImageGroup(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()

	docImageGroupGUID := ctx.Param("guid")

	input := ctx.ReadInput()

	docImages := &[]models.ImageReferenceBody{}

	err := json.Unmarshal([]byte(input), &docImages)
	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.service.UpdateImageReferenceByDocumentImageGroup(userInfo.ShopID, userInfo.Username, docImageGroupGUID, *docImages)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}

// Update Reference In Document Image Group
// @Description Update Reference In Document Image Group
// @Tags		DocumentImageGroup
// @Accept 		json
// @Param		guid  path      string  true  "document image group guid"
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /documentimagegroup/{guid}/reference [put]
func (h DocumentImageHttp) UpdateReferenceByDocumentImageGroup(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()

	docImageGroupGUID := ctx.Param("guid")

	input := ctx.ReadInput()

	docImages := &models.Reference{}

	err := json.Unmarshal([]byte(input), &docImages)
	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.service.UpdateReferenceByDocumentImageGroup(userInfo.ShopID, userInfo.Username, docImageGroupGUID, *docImages)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}

// Ungroup Document Image Group
// @Description Ungroup Document Image Group
// @Tags		DocumentImageGroup
// @Accept 		json
// @Param		guid  path      string  true  "document image group guid"
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /documentimagegroup/{guid}/ungroup [put]
func (h DocumentImageHttp) UngroupDocumentImageGroup(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()

	docImageGroupGUID := ctx.Param("guid")

	err := h.service.UnGroupDocumentImageGroup(userInfo.ShopID, userInfo.Username, docImageGroupGUID)
	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}

// Get Document Image Group By document reference Infomation godoc
// @Summary		Get Document Image Group By document reference Infomation
// @Description Get Document Image Group By document reference Infomation
// @Tags		DocumentImageGroup
// @Param		docref  path      string  true  "document reference"
// @Accept 		json
// @Success		200	{object}	models.DocumentImageInfoResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /documentimagegroup/docref/:docref [get]
func (h DocumentImageHttp) GetDocumentImageGroupByDocRefInfo(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	docRef := ctx.Param("docref")
	doc, err := h.service.GetDocumentImageGroupByDocRef(shopID, docRef)
	if err != nil {
		h.ms.Logger.Errorf("Error getting document %v: %v", docRef, err)
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    doc,
	})
	return nil
}
