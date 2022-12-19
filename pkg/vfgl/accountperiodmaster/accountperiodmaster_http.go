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
	"strings"
	"time"
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

	h.ms.GET("/gl/accountperiodmaster", h.SearchAccountPeriodMasterPage)
	h.ms.GET("/gl/accountperiodmaster/list", h.SearchAccountPeriodMasterLimit)
	h.ms.POST("/gl/accountperiodmaster", h.CreateAccountPeriodMaster)
	h.ms.POST("/gl/accountperiodmaster/bulk", h.SaveBulkAccountPeriodMaster)
	h.ms.GET("/gl/accountperiodmaster/:id", h.InfoAccountPeriodMaster)
	h.ms.GET("/gl/accountperiodmaster/bydate", h.InfoAccountPeriodMasterByDate)
	h.ms.PUT("/gl/accountperiodmaster/:id", h.UpdateAccountPeriodMaster)
	h.ms.DELETE("/gl/accountperiodmaster/:id", h.DeleteAccountPeriodMaster)
	h.ms.DELETE("/gl/accountperiodmaster", h.DeleteAccountPeriodMasterByGUIDs)
}

// Create AccountPeriodMaster godoc
// @Description Create AccountPeriodMaster
// @Tags		AccountPeriodMaster
// @Param		AccountPeriodMaster  body      models.AccountPeriodMaster true  "AccountPeriodMaster"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /gl/accountperiodmaster [post]
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
// @Param		AccountPeriodMaster  body      models.AccountPeriodMaster true  "AccountPeriodMaster"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /gl/accountperiodmaster/{id} [put]
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
// @Router /gl/accountperiodmaster/{id} [delete]
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
// @Router /gl/accountperiodmaster [delete]
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
// @Router /gl/accountperiodmaster/{id} [get]
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

// Get AccountPeriodMaster godoc
// @Description Get AccountPeriodMaster by date
// @Tags		AccountPeriodMaster
// @Param		date		query	string		false  "date"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /gl/accountperiodmaster/bydate [get]
func (h AccountPeriodMasterHttp) InfoAccountPeriodMasterByDate(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	layout := "2006-01-02"
	dateStr := ctx.QueryParam("date")

	dateStr = strings.Trim(dateStr, " ")
	if len(dateStr) < 1 {
		ctx.ResponseError(400, "date format invalid.")
		return nil
	}

	findDate, err := time.Parse(layout, dateStr)

	h.ms.Logger.Debugf("Get AccountPeriodMaster %v", findDate)
	doc, err := h.svc.InfoAccountPeriodMasterByDate(shopID, findDate)

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
// @Router /gl/accountperiodmaster [get]
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
// @Router /gl/accountperiodmaster/list [get]
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

// Bulk Create AccountPeriodMaster godoc
// @Description Bulk Create AccountPeriodMaster
// @Tags		AccountPeriodMaster
// @Param		AccountPeriodMaster  body      []models.AccountPeriodMaster true  "Array AccountPeriodMaster"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /gl/accountperiodmaster/bulk [post]
func (h AccountPeriodMasterHttp) SaveBulkAccountPeriodMaster(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &[]models.AccountPeriodMasterRequest{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	tempDocListReq := []models.AccountPeriodMaster{}
	for _, doc := range *docReq {
		tempDocListReq = append(tempDocListReq, doc.ToAccountPeriodMaster())
	}

	err = h.svc.SaveInBatch(shopID, authUsername, tempDocListReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
	})
	return nil
}
