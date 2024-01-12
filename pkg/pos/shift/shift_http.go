package shift

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/config"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/pos/shift/models"
	"smlcloudplatform/pkg/pos/shift/repositories"
	"smlcloudplatform/pkg/pos/shift/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/requestfilter"
	"time"
)

type IShiftHttp interface{}

type ShiftHttp struct {
	ms  *microservice.Microservice
	cfg config.IConfig
	svc services.IShiftHttpService
}

func NewShiftHttp(ms *microservice.Microservice, cfg config.IConfig) ShiftHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	repo := repositories.NewShiftRepository(pst)

	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)
	svc := services.NewShiftHttpService(repo, masterSyncCacheRepo, 15*time.Second)

	return ShiftHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h ShiftHttp) RegisterHttp() {

	h.ms.POST("/pos/shift/bulk", h.SaveBulk)

	h.ms.GET("/pos/shift", h.SearchShiftPage)
	h.ms.GET("/pos/shift/list", h.SearchShiftStep)
	h.ms.POST("/pos/shift", h.CreateShift)
	h.ms.GET("/pos/shift/:id", h.InfoShift)
	h.ms.GET("/pos/shift/code/:code", h.InfoShiftByCode)
	h.ms.PUT("/pos/shift/:id", h.UpdateShift)
	h.ms.DELETE("/pos/shift/:id", h.DeleteShift)
	h.ms.DELETE("/pos/shift", h.DeleteShiftByGUIDs)
}

// Create Shift godoc
// @Description Create Shift
// @Tags		Shift
// @Param		Shift  body      models.Shift  true  "Shift"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /pos/shift [post]
func (h ShiftHttp) CreateShift(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.Shift{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateShift(shopID, authUsername, *docReq)

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

// Update Shift godoc
// @Description Update Shift
// @Tags		Shift
// @Param		id  path      string  true  "Shift ID"
// @Param		Shift  body      models.Shift  true  "Shift"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /pos/shift/{id} [put]
func (h ShiftHttp) UpdateShift(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.Shift{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateShift(shopID, id, authUsername, *docReq)

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

// Delete Shift godoc
// @Description Delete Shift
// @Tags		Shift
// @Param		id  path      string  true  "Shift ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /pos/shift/{id} [delete]
func (h ShiftHttp) DeleteShift(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteShift(shopID, id, authUsername)

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

// Delete Shift godoc
// @Description Delete Shift
// @Tags		Shift
// @Param		Shift  body      []string  true  "Shift GUIDs"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /pos/shift [delete]
func (h ShiftHttp) DeleteShiftByGUIDs(ctx microservice.IContext) error {
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

	err = h.svc.DeleteShiftByGUIDs(shopID, authUsername, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}

// Get Shift godoc
// @Description get Shift info by guidfixed
// @Tags		Shift
// @Param		id  path      string  true  "Shift guidfixed"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /pos/shift/{id} [get]
func (h ShiftHttp) InfoShift(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get Shift %v", id)
	doc, err := h.svc.InfoShift(shopID, id)

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

// Get Shift By Code godoc
// @Description get Shift info by Code
// @Tags		Shift
// @Param		code  path      string  true  "Shift Code"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /pos/shift/code/{code} [get]
func (h ShiftHttp) InfoShiftByCode(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	code := ctx.Param("code")

	doc, err := h.svc.InfoShiftByCode(shopID, code)

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

// List Shift step godoc
// @Description get list step
// @Tags		Shift
// @Param		q		query	string		false  "Search Value"
// @Param 		doctype	query 	int8	false	"DocType"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Limit"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /pos/shift [get]
func (h ShiftHttp) SearchShiftPage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)

	filters := requestfilter.GenerateFilters(ctx.QueryParam, []requestfilter.FilterRequest{
		{
			Param: "doctype",
			Type:  requestfilter.FieldTypeInt,
		},
	})

	docList, pagination, err := h.svc.SearchShift(shopID, filters, pageable)

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

// List Shift godoc
// @Description search limit offset
// @Tags		Shift
// @Param		q		query	string		false  "Search Value"
// @Param 		doctype	query 	int8		false	"DocType"
// @Param		offset	query	integer		false  "offset"
// @Param		limit	query	integer		false  "limit"
// @Param		lang	query	string		false  "lang"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /pos/shift/list [get]
func (h ShiftHttp) SearchShiftStep(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageableStep := utils.GetPageableStep(ctx.QueryParam)

	lang := ctx.QueryParam("lang")

	filters := requestfilter.GenerateFilters(ctx.QueryParam, []requestfilter.FilterRequest{
		{
			Param: "doctype",
			Type:  requestfilter.FieldTypeInt,
		},
	})

	docList, total, err := h.svc.SearchShiftStep(shopID, lang, filters, pageableStep)

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

// Create Shift Bulk godoc
// @Description Create Shift
// @Tags		Shift
// @Param		Shift  body      []models.Shift  true  "Shift"
// @Accept 		json
// @Success		201	{object}	common.BulkResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /pos/shift/bulk [post]
func (h ShiftHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.Shift{}
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
