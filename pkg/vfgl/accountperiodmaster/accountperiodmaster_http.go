package accountperiodmaster

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/vfgl/accountperiodmaster/models"
	"smlcloudplatform/pkg/vfgl/accountperiodmaster/repositories"
	"smlcloudplatform/pkg/vfgl/accountperiodmaster/services"
)

type IAccountPeriodMasterHttp interface{}

type AccountPeriodMasterHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc services.IAccountPeriodMasterHttpService
}

func NewAccountPeriodMasterHttp(ms *microservice.Microservice, cfg microservice.IConfig) AccountPeriodMasterHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	repo := repositories.NewAccountPeriodMasterRepository(pst)

	svc := services.NewAccountPeriodMasterHttpService(repo)

	return AccountPeriodMasterHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h AccountPeriodMasterHttp) RouteSetup() {

	h.ms.GET("/accountperiodmaster", h.SearchAccountPeriodMasterPage)
	h.ms.GET("/accountperiodmaster/list", h.SearchAccountPeriodMasterLimit)
	h.ms.POST("/accountperiodmaster", h.CreateAccountPeriodMaster)
	h.ms.GET("/accountperiodmaster/:id", h.InfoAccountPeriodMaster)
	h.ms.PUT("/accountperiodmaster/:id", h.UpdateAccountPeriodMaster)
	h.ms.DELETE("/accountperiodmaster/:id", h.DeleteAccountPeriodMaster)
	h.ms.DELETE("/accountperiodmaster", h.DeleteAccountPeriodMasterByGUIDs)
}

// Create AccountPeriodMaster godoc
// @Description Create AccountPeriodMaster
// @Tags		AccountPeriodMaster
// @Param		AccountPeriodMaster  body      models.AccountPeriodMasterRequest  true  "AccountPeriodMaster"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /accountperiodmaster [post]
func (h AccountPeriodMasterHttp) CreateAccountPeriodMaster(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.AccountPeriodMasterRequest{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateAccountPeriodMaster(shopID, authUsername, docReq.ToAccountPeriodMaster())

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

// Update AccountPeriodMaster godoc
// @Description Update AccountPeriodMaster
// @Tags		AccountPeriodMaster
// @Param		id  path      string  true  "AccountPeriodMaster ID"
// @Param		AccountPeriodMaster  body      models.AccountPeriodMasterRequest  true  "AccountPeriodMaster"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /accountperiodmaster/{id} [put]
func (h AccountPeriodMasterHttp) UpdateAccountPeriodMaster(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.AccountPeriodMasterRequest{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateAccountPeriodMaster(shopID, id, authUsername, docReq.ToAccountPeriodMaster())

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

// Delete AccountPeriodMaster godoc
// @Description Delete AccountPeriodMaster
// @Tags		AccountPeriodMaster
// @Param		id  path      string  true  "AccountPeriodMaster ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /accountperiodmaster/{id} [delete]
func (h AccountPeriodMasterHttp) DeleteAccountPeriodMaster(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteAccountPeriodMaster(shopID, id, authUsername)

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

// Delete AccountPeriodMaster godoc
// @Description Delete AccountPeriodMaster
// @Tags		AccountPeriodMaster
// @Param		AccountPeriodMaster  body      []string  true  "AccountPeriodMaster GUIDs"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /accountperiodmaster [delete]
func (h AccountPeriodMasterHttp) DeleteAccountPeriodMasterByGUIDs(ctx microservice.IContext) error {
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

	err = h.svc.DeleteAccountPeriodMasterByGUIDs(shopID, authUsername, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}

// Get AccountPeriodMaster godoc
// @Description get struct array by ID
// @Tags		AccountPeriodMaster
// @Param		id  path      string  true  "AccountPeriodMaster ID"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /accountperiodmaster/{id} [get]
func (h AccountPeriodMasterHttp) InfoAccountPeriodMaster(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get AccountPeriodMaster %v", id)
	doc, err := h.svc.InfoAccountPeriodMaster(shopID, id)

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

// List AccountPeriodMaster godoc
// @Description get struct array by ID
// @Tags		AccountPeriodMaster
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Add Category"
// @Param		limit	query	integer		false  "Add Category"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /accountperiodmaster [get]
func (h AccountPeriodMasterHttp) SearchAccountPeriodMasterPage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)
	sort := utils.GetSortParam(ctx.QueryParam)
	docList, pagination, err := h.svc.SearchAccountPeriodMaster(shopID, q, page, limit, sort)

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

// List AccountPeriodMaster godoc
// @Description search limit offset
// @Tags		AccountPeriodMaster
// @Param		q		query	string		false  "Search Value"
// @Param		offset	query	integer		false  "offset"
// @Param		limit	query	integer		false  "limit"
// @Param		lang	query	string		false  "lang"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /accountperiodmaster/list [get]
func (h AccountPeriodMasterHttp) SearchAccountPeriodMasterLimit(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	offset, limit := utils.GetParamOffsetLimit(ctx.QueryParam)
	sorts := utils.GetSortParam(ctx.QueryParam)

	lang := ctx.QueryParam("lang")

	docList, total, err := h.svc.SearchAccountPeriodMasterStep(shopID, lang, q, offset, limit, sorts)

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
