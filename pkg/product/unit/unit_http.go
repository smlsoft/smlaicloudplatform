package unit

import (
	"encoding/json"
	"net/http"
	"net/url"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/config"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/unit/models"
	"smlcloudplatform/pkg/product/unit/repositories"
	"smlcloudplatform/pkg/product/unit/services"
	"smlcloudplatform/pkg/utils"
	"strings"
)

type IUnitHttp interface{}

type UnitHttp struct {
	ms  *microservice.Microservice
	cfg config.IConfig
	svc services.IUnitHttpService
}

func NewUnitHttp(ms *microservice.Microservice, cfg config.IConfig) UnitHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	repo := repositories.NewUnitRepository(pst)

	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)
	svc := services.NewUnitHttpService(repo, masterSyncCacheRepo)

	return UnitHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h UnitHttp) RouteSetup() {

	h.ms.POST("/unit/bulk", h.SaveBulk)

	h.ms.GET("/unit", h.SearchUnit)
	h.ms.GET("/unit/list", h.SearchUnitLimit)
	h.ms.POST("/unit", h.CreateUnit)
	h.ms.GET("/unit/:id", h.InfoUnit)
	h.ms.GET("/unit/by-code", h.InfoArray)
	h.ms.GET("/unit/master", h.InfoArrayMaster)
	h.ms.PUT("/unit/:id", h.UpdateUnit)
	h.ms.PATCH("/unit/:id", h.UpdateFieldUnit)
	h.ms.DELETE("/unit/:id", h.DeleteUnit)
	h.ms.DELETE("/unit", h.DeleteByGUIDs)
}

// Create Unit godoc
// @Description Create Unit
// @Tags		Unit
// @Param		Unit  body      models.Unit  true  "Unit"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /unit [post]
func (h UnitHttp) CreateUnit(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.Unit{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateUnit(shopID, authUsername, *docReq)

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

// Update Unit godoc
// @Description Update Unit
// @Tags		Unit
// @Param		id  path      string  true  "Unit ID"
// @Param		Unit  body      models.Unit  true  "Unit"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /unit/{id} [put]
func (h UnitHttp) UpdateUnit(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	if len(id) < 1 {
		ctx.ResponseError(http.StatusBadRequest, "guid is empty")
		return nil
	}

	input := ctx.ReadInput()

	docReq := &models.Unit{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateUnit(shopID, id, authUsername, *docReq)

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

// Update Field Unit godoc
// @Description Update Unit
// @Tags		Unit
// @Param		id  path      string  true  "Unit ID"
// @Param		Unit  body      models.Unit  true  "Unit"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /unit/{id} [patch]
func (h UnitHttp) UpdateFieldUnit(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	if len(id) < 1 {
		ctx.ResponseError(http.StatusBadRequest, "guid is empty")
		return nil
	}

	input := ctx.ReadInput()

	docReq := &models.Unit{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateFieldUnit(shopID, id, authUsername, *docReq)

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

// Get Unit godoc
// @Description get struct array by ID
// @Tags		Unit
// @Param		id  path      string  true  "Unit ID"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /unit/{id} [get]
func (h UnitHttp) InfoUnit(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	if len(id) < 1 {
		ctx.ResponseError(http.StatusBadRequest, "guid is empty")
		return nil
	}

	h.ms.Logger.Debugf("Get Unit %v", id)
	doc, err := h.svc.InfoUnit(shopID, id)

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

// Get Unit By unit code array godoc
// @Description get unit by unit code array
// @Tags		Unit
// @Param		codes	query	string		false  "Code filter, json array encode "
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /unit/by-code [get]
func (h UnitHttp) InfoArray(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	codesReq, err := url.QueryUnescape(ctx.QueryParam("codes"))

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	docReq := []string{}
	err = json.Unmarshal([]byte(codesReq), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	// where to filter array
	doc, err := h.svc.InfoUnitWTFArray(shopID, docReq)

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

// Get Master Unit By code array godoc
// @Description get master Unit by code array
// @Tags		Unit
// @Param		codes	query	string		false  "Code filter, json array encode "
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /unit/master [get]
func (h UnitHttp) InfoArrayMaster(ctx microservice.IContext) error {
	codesReq, err := url.QueryUnescape(ctx.QueryParam("codes"))

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	docReq := []string{}
	err = json.Unmarshal([]byte(codesReq), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	// where to filter array master
	doc, err := h.svc.InfoWTFArrayMaster(docReq)

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

// List Unit godoc
// @Description get struct array by ID
// @Tags		Unit
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "page "
// @Param		limit	query	integer		false  "liumit "
// @Param		unitcode	query	string		false  "unitcode filter ex. \"u001,u002,u003\""
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /unit [get]
func (h UnitHttp) SearchUnit(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)

	unitCode := ctx.QueryParam("unitcode")

	unitCodeFilters := []string{}
	if len(unitCode) > 0 {
		unitCodeFilters = strings.Split(unitCode, ",")
	}

	docList, pagination, err := h.svc.SearchUnit(shopID, unitCodeFilters, pageable)

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

// List Unit godoc
// @Description search limit offset
// @Tags		Unit
// @Param		q		query	string		false  "Search Value"
// @Param		offset	query	integer		false  "offset"
// @Param		limit	query	integer		false  "limit"
// @Param		lang	query	string		false  "lang ex. en,th"
// @Param		unitcode	query	string		false  "unitcode filter ex. \"u001,u002,u003\" "
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /unit/list [get]
func (h UnitHttp) SearchUnitLimit(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageableStep := utils.GetPageableStep(ctx.QueryParam)

	lang := ctx.QueryParam("lang")

	unitCode := ctx.QueryParam("unitcode")

	unitCodeFilters := []string{}
	if len(unitCode) > 0 {
		unitCodeFilters = strings.Split(unitCode, ",")
	}

	docList, total, err := h.svc.SearchUnitLimit(shopID, lang, unitCodeFilters, pageableStep)

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

// Create Unit Bulk godoc
// @Description Create Unit
// @Tags		Unit
// @Param		Unit  body      []models.Unit  true  "Unit"
// @Accept 		json
// @Success		201	{object}	common.BulkReponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /unit/bulk [post]
func (h UnitHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.Unit{}
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

// Delete Unit godoc
// @Description Delete Unit
// @Tags		Unit
// @Param		id  path      string  true  "Unit ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /unit/{id} [delete]
func (h UnitHttp) DeleteUnit(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	if len(id) < 1 {
		ctx.ResponseError(http.StatusBadRequest, "unit guid is empty")
		return nil
	}

	authHeader := ctx.Header("Authorization")

	err := h.svc.DeleteUnit(shopID, id, authHeader, authUsername)

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

// Delete Unit By GUIDs godoc
// @Description Delete Unit
// @Tags		Unit
// @Param		Unit  body      []string  true  "Unit GUIDs"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /unit [delete]
func (h UnitHttp) DeleteByGUIDs(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	authHeader := ctx.Header("Authorization")

	input := ctx.ReadInput()

	docReq := []string{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.DeleteUnitByGUIDs(shopID, authHeader, authUsername, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}
