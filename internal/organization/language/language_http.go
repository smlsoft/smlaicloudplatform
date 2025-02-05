package language

import (
	"encoding/json"
	"net/http"
	"smlaicloudplatform/internal/config"
	mastersync "smlaicloudplatform/internal/mastersync/repositories"
	common "smlaicloudplatform/internal/models"
	"smlaicloudplatform/internal/organization/language/models"
	"smlaicloudplatform/internal/organization/language/repositories"
	"smlaicloudplatform/internal/organization/language/services"
	"smlaicloudplatform/internal/utils"
	"smlaicloudplatform/pkg/microservice"
)

type ILanguageHttp interface{}

type LanguageHttp struct {
	ms  *microservice.Microservice
	cfg config.IConfig
	svc services.ILanguageHttpService
}

func NewLanguageHttp(ms *microservice.Microservice, cfg config.IConfig) LanguageHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	repo := repositories.NewLanguageRepository(pst)

	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)
	svc := services.NewLanguageHttpService(repo, masterSyncCacheRepo)

	return LanguageHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h LanguageHttp) RegisterHttp() {

	h.ms.POST("/organization/language/bulk", h.SaveBulk)

	h.ms.GET("/organization/language", h.SearchLanguagePage)
	h.ms.GET("/organization/language/list", h.SearchLanguageStep)
	h.ms.POST("/organization/language", h.CreateLanguage)
	h.ms.GET("/organization/language/:id", h.InfoLanguage)
	h.ms.GET("/organization/language/:languageCode", h.InfoLanguageByCode)
	h.ms.PUT("/organization/language/:id", h.UpdateLanguage)
	h.ms.DELETE("/organization/language/:id", h.DeleteLanguage)
	h.ms.DELETE("/organization/language", h.DeleteLanguageByGUIDs)
}

// Create Language godoc
// @Description Create Language
// @Tags		Language
// @Param		Language  body      models.Language  true  "Language"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /organization/language [post]
func (h LanguageHttp) CreateLanguage(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.Language{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateLanguage(shopID, authUsername, *docReq)

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

// Update Language godoc
// @Description Update Language
// @Tags		Language
// @Param		id  path      string  true  "Language ID"
// @Param		Language  body      models.Language  true  "Language"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /organization/language/{id} [put]
func (h LanguageHttp) UpdateLanguage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.Language{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateLanguage(shopID, id, authUsername, *docReq)

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

// Delete Language godoc
// @Description Delete Language
// @Tags		Language
// @Param		id  path      string  true  "Language ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /organization/language/{id} [delete]
func (h LanguageHttp) DeleteLanguage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteLanguage(shopID, id, authUsername)

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

// Delete Language godoc
// @Description Delete Language
// @Tags		Language
// @Param		Language  body      []string  true  "Language GUIDs"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /organization/language [delete]
func (h LanguageHttp) DeleteLanguageByGUIDs(ctx microservice.IContext) error {
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

	err = h.svc.DeleteLanguageByGUIDs(shopID, authUsername, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}

// Get Language godoc
// @Description get Language info by guidfixed
// @Tags		Language
// @Param		id  path      string  true  "Language guidfixed"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /organization/language/{id} [get]
func (h LanguageHttp) InfoLanguage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get Language %v", id)
	doc, err := h.svc.InfoLanguage(shopID, id)

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

// Get Language By Code godoc
// @Description get Language info by Code
// @Tags		Language
// @Param		languageCode  path      string  true  "Language Code"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /organization/language/{languageCode} [get]
func (h LanguageHttp) InfoLanguageByCode(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	languageCode := ctx.Param("languageCode")

	doc, err := h.svc.InfoLanguageByCode(shopID, languageCode)

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

// List Language step godoc
// @Description get list step
// @Tags		Language
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Limit"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /organization/language [get]
func (h LanguageHttp) SearchLanguagePage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)

	docList, pagination, err := h.svc.SearchLanguage(shopID, map[string]interface{}{}, pageable)

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

// List Language godoc
// @Description search limit offset
// @Tags		Language
// @Param		q		query	string		false  "Search Value"
// @Param		offset	query	integer		false  "offset"
// @Param		limit	query	integer		false  "limit"
// @Param		lang	query	string		false  "lang"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /organization/language/list [get]
func (h LanguageHttp) SearchLanguageStep(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageableStep := utils.GetPageableStep(ctx.QueryParam)

	lang := ctx.QueryParam("lang")

	docList, total, err := h.svc.SearchLanguageStep(shopID, lang, pageableStep)

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

// Create Language Bulk godoc
// @Description Create Language
// @Tags		Language
// @Param		Language  body      []models.Language  true  "Language"
// @Accept 		json
// @Success		201	{object}	common.BulkResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /organization/language/bulk [post]
func (h LanguageHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.Language{}
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
		common.BulkResponse{
			Success:    true,
			BulkImport: bulkResponse,
		},
	)

	return nil
}
