package customergroup

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/customershop/customergroup/models"
	"smlcloudplatform/pkg/customershop/customergroup/repositories"
	"smlcloudplatform/pkg/customershop/customergroup/services"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
)

type ICustomerGroupHttp interface{}

type CustomerGroupHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc services.ICustomerGroupHttpService
}

func NewCustomerGroupHttp(ms *microservice.Microservice, cfg microservice.IConfig) CustomerGroupHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	repo := repositories.NewCustomerGroupRepository(pst)

	svc := services.NewCustomerGroupHttpService(repo)

	return CustomerGroupHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h CustomerGroupHttp) RouteSetup() {

	h.ms.POST("customershop/customergroup/bulk", h.SaveBulk)

	h.ms.GET("customershop/customergroup", h.SearchCustomerGroupPage)
	h.ms.GET("customershop/customergroup/list", h.SearchCustomerGroupLimit)
	h.ms.POST("customershop/customergroup", h.CreateCustomerGroup)
	h.ms.GET("customershop/customergroup/:id", h.InfoCustomerGroup)
	h.ms.PUT("customershop/customergroup/:id", h.UpdateCustomerGroup)
	h.ms.DELETE("customershop/customergroup/:id", h.DeleteCustomerGroup)
}

// Create CustomerGroup godoc
// @Description Create CustomerGroup
// @Tags		CustomerShop
// @Param		CustomerGroup  body      models.CustomerGroup  true  "CustomerGroup"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /customershop/customergroup [post]
func (h CustomerGroupHttp) CreateCustomerGroup(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.CustomerGroup{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateCustomerGroup(shopID, authUsername, *docReq)

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

// Update CustomerGroup godoc
// @Description Update CustomerGroup
// @Tags		CustomerShop
// @Param		id  path      string  true  "CustomerGroup ID"
// @Param		CustomerGroup  body      models.CustomerGroup  true  "CustomerGroup"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /customershop/customergroup/{id} [put]
func (h CustomerGroupHttp) UpdateCustomerGroup(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.CustomerGroup{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateCustomerGroup(shopID, id, authUsername, *docReq)

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

// Delete CustomerGroup godoc
// @Description Delete CustomerGroup
// @Tags		CustomerShop
// @Param		id  path      string  true  "CustomerGroup ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /customershop/customergroup/{id} [delete]
func (h CustomerGroupHttp) DeleteCustomerGroup(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteCustomerGroup(shopID, id, authUsername)

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

// Get CustomerGroup godoc
// @Description get struct array by ID
// @Tags		CustomerShop
// @Param		id  path      string  true  "CustomerGroup ID"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /customershop/customergroup/{id} [get]
func (h CustomerGroupHttp) InfoCustomerGroup(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get CustomerGroup %v", id)
	doc, err := h.svc.InfoCustomerGroup(shopID, id)

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

// List CustomerGroup godoc
// @Description get struct array by ID
// @Tags		CustomerShop
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Add Category"
// @Param		limit	query	integer		false  "Add Category"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /customershop/customergroup [get]
func (h CustomerGroupHttp) SearchCustomerGroupPage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)
	sort := utils.GetSortParam(ctx.QueryParam)
	docList, pagination, err := h.svc.SearchCustomerGroup(shopID, q, page, limit, sort)

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

// List CustomerGroup godoc
// @Description search limit offset
// @Tags		CustomerShop
// @Param		q		query	string		false  "Search Value"
// @Param		offset	query	integer		false  "offset"
// @Param		limit	query	integer		false  "limit"
// @Param		lang	query	string		false  "lang"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /customershop/customergroup/list [get]
func (h CustomerGroupHttp) SearchCustomerGroupLimit(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	offset, limit := utils.GetParamOffsetLimit(ctx.QueryParam)
	sorts := utils.GetSortParam(ctx.QueryParam)

	lang := ctx.QueryParam("lang")

	docList, total, err := h.svc.SearchCustomerGroupStep(shopID, lang, q, offset, limit, sorts)

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

// Create CustomerGroup Bulk godoc
// @Description Create CustomerGroup
// @Tags		CustomerShop
// @Param		CustomerGroup  body      []models.CustomerGroup  true  "CustomerGroup"
// @Accept 		json
// @Success		201	{object}	common.BulkReponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /customershop/customergroup/bulk [post]
func (h CustomerGroupHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.CustomerGroup{}
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
