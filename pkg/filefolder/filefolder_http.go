package filefolder

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/filefolder/models"
	"smlcloudplatform/pkg/filefolder/repositories"
	"smlcloudplatform/pkg/filefolder/services"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
)

type IFileFolderHttp interface{}

type FileFolderHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc services.IFileFolderHttpService
}

func NewFileFolderHttp(ms *microservice.Microservice, cfg microservice.IConfig) FileFolderHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	repo := repositories.NewFileFolderRepository(pst)

	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)
	svc := services.NewFileFolderHttpService(repo, masterSyncCacheRepo)

	return FileFolderHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h FileFolderHttp) RouteSetup() {

	h.ms.POST("/file-folder/bulk", h.SaveBulk)

	h.ms.GET("/file-folder", h.SearchFileFolderPage)
	h.ms.GET("/file-folder/list", h.SearchFileFolderLimit)
	h.ms.POST("/file-folder", h.CreateFileFolder)
	h.ms.GET("/file-folder/:id", h.InfoFileFolder)
	h.ms.PUT("/file-folder/:id", h.UpdateFileFolder)
	h.ms.DELETE("/file-folder/:id", h.DeleteFileFolder)
	h.ms.DELETE("/file-folder", h.DeleteFileFolderByGUIDs)
}

// Create FileFolder godoc
// @Description Create FileFolder
// @Tags		FileFolder
// @Param		FileFolder  body      models.FileFolder  true  "FileFolder"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /file-folder [post]
func (h FileFolderHttp) CreateFileFolder(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.FileFolder{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateFileFolder(shopID, authUsername, *docReq)

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

// Update FileFolder godoc
// @Description Update FileFolder
// @Tags		FileFolder
// @Param		id  path      string  true  "FileFolder ID"
// @Param		FileFolder  body      models.FileFolder  true  "FileFolder"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /file-folder/{id} [put]
func (h FileFolderHttp) UpdateFileFolder(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.FileFolder{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateFileFolder(shopID, id, authUsername, *docReq)

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

// Delete FileFolder godoc
// @Description Delete FileFolder
// @Tags		FileFolder
// @Param		id  path      string  true  "FileFolder ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /file-folder/{id} [delete]
func (h FileFolderHttp) DeleteFileFolder(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteFileFolder(shopID, id, authUsername)

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

// Delete FileFolder godoc
// @Description Delete FileFolder
// @Tags		FileFolder
// @Param		FileFolder  body      []string  true  "FileFolder GUIDs"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /file-folder [delete]
func (h FileFolderHttp) DeleteFileFolderByGUIDs(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	input := ctx.ReadInput()

	docReq := []string{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.DeleteFileFolderByGUIDs(shopID, authUsername, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}

// Get FileFolder godoc
// @Description get struct array by ID
// @Tags		FileFolder
// @Param		id  path      string  true  "FileFolder ID"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /file-folder/{id} [get]
func (h FileFolderHttp) InfoFileFolder(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get FileFolder %v", id)
	doc, err := h.svc.InfoFileFolder(shopID, id)

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

// List FileFolder godoc
// @Description get struct array by ID
// @Tags		FileFolder
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Limit"
// @Param		module	query	integer		false  "Module"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /file-folder [get]
func (h FileFolderHttp) SearchFileFolderPage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)

	module := ctx.QueryParam("module")

	docList, pagination, err := h.svc.SearchFileFolder(shopID, module, pageable)

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

// List FileFolder godoc
// @Description search limit offset
// @Tags		FileFolder
// @Param		q		query	string		false  "Search Value"
// @Param		offset	query	integer		false  "offset"
// @Param		limit	query	integer		false  "limit"
// @Param		lang	query	string		false  "lang"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /file-folder/list [get]
func (h FileFolderHttp) SearchFileFolderLimit(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageableStep := utils.GetPageableStep(ctx.QueryParam)

	lang := ctx.QueryParam("lang")

	docList, total, err := h.svc.SearchFileFolderStep(shopID, lang, pageableStep)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    docList,
		Total:   total,
	})
	return nil
}

// Create FileFolder Bulk godoc
// @Description Create FileFolder
// @Tags		FileFolder
// @Param		FileFolder  body      []models.FileFolder  true  "FileFolder"
// @Accept 		json
// @Success		201	{object}	common.BulkReponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /file-folder/bulk [post]
func (h FileFolderHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.FileFolder{}
	err := json.Unmarshal([]byte(input), &dataReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	bulkResponse, err := h.svc.SaveInBatch(shopID, authUsername, dataReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusCreated,
		common.BulkReponse{
			Success:    true,
			BulkImport: bulkResponse,
		},
	)

	return nil
}
