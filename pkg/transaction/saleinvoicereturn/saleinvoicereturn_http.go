package saleinvoicereturn

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/transaction/saleinvoicereturn/models"
	"smlcloudplatform/pkg/transaction/saleinvoicereturn/repositories"
	"smlcloudplatform/pkg/transaction/saleinvoicereturn/services"
	"smlcloudplatform/pkg/utils"
)

type ISaleInvoiceReturnHttp interface{}

type SaleInvoiceReturnHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc services.ISaleInvoiceReturnHttpService
}

func NewSaleInvoiceReturnHttp(ms *microservice.Microservice, cfg microservice.IConfig) SaleInvoiceReturnHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	repo := repositories.NewSaleInvoiceReturnRepository(pst)

	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)
	svc := services.NewSaleInvoiceReturnHttpService(repo, masterSyncCacheRepo)

	return SaleInvoiceReturnHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h SaleInvoiceReturnHttp) RouteSetup() {

	h.ms.POST("/transaction/sale-invoice-return/bulk", h.SaveBulk)

	h.ms.GET("/transaction/sale-invoice-return", h.SearchSaleInvoiceReturnPage)
	h.ms.GET("/transaction/sale-invoice-return/list", h.SearchSaleInvoiceReturnStep)
	h.ms.POST("/transaction/sale-invoice-return", h.CreateSaleInvoiceReturn)
	h.ms.GET("/transaction/sale-invoice-return/:id", h.InfoSaleInvoiceReturn)
	h.ms.GET("/transaction/sale-invoice-return/code/:code", h.InfoSaleInvoiceReturnByCode)
	h.ms.PUT("/transaction/sale-invoice-return/:id", h.UpdateSaleInvoiceReturn)
	h.ms.DELETE("/transaction/sale-invoice-return/:id", h.DeleteSaleInvoiceReturn)
	h.ms.DELETE("/transaction/sale-invoice-return", h.DeleteSaleInvoiceReturnByGUIDs)
}

// Create SaleInvoiceReturn godoc
// @Description Create SaleInvoiceReturn
// @Tags		SaleInvoiceReturn
// @Param		SaleInvoiceReturn  body      models.SaleInvoiceReturn  true  "SaleInvoiceReturn"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/sale-invoice-return [post]
func (h SaleInvoiceReturnHttp) CreateSaleInvoiceReturn(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.SaleInvoiceReturn{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateSaleInvoiceReturn(shopID, authUsername, *docReq)

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

// Update SaleInvoiceReturn godoc
// @Description Update SaleInvoiceReturn
// @Tags		SaleInvoiceReturn
// @Param		id  path      string  true  "SaleInvoiceReturn ID"
// @Param		SaleInvoiceReturn  body      models.SaleInvoiceReturn  true  "SaleInvoiceReturn"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/sale-invoice-return/{id} [put]
func (h SaleInvoiceReturnHttp) UpdateSaleInvoiceReturn(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.SaleInvoiceReturn{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateSaleInvoiceReturn(shopID, id, authUsername, *docReq)

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

// Delete SaleInvoiceReturn godoc
// @Description Delete SaleInvoiceReturn
// @Tags		SaleInvoiceReturn
// @Param		id  path      string  true  "SaleInvoiceReturn ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/sale-invoice-return/{id} [delete]
func (h SaleInvoiceReturnHttp) DeleteSaleInvoiceReturn(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteSaleInvoiceReturn(shopID, id, authUsername)

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

// Delete SaleInvoiceReturn godoc
// @Description Delete SaleInvoiceReturn
// @Tags		SaleInvoiceReturn
// @Param		SaleInvoiceReturn  body      []string  true  "SaleInvoiceReturn GUIDs"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/sale-invoice-return [delete]
func (h SaleInvoiceReturnHttp) DeleteSaleInvoiceReturnByGUIDs(ctx microservice.IContext) error {
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

	err = h.svc.DeleteSaleInvoiceReturnByGUIDs(shopID, authUsername, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}

// Get SaleInvoiceReturn godoc
// @Description get SaleInvoiceReturn info by guidfixed
// @Tags		SaleInvoiceReturn
// @Param		id  path      string  true  "SaleInvoiceReturn guidfixed"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/sale-invoice-return/{id} [get]
func (h SaleInvoiceReturnHttp) InfoSaleInvoiceReturn(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get SaleInvoiceReturn %v", id)
	doc, err := h.svc.InfoSaleInvoiceReturn(shopID, id)

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

// Get SaleInvoiceReturn By Code godoc
// @Description get SaleInvoiceReturn info by Code
// @Tags		SaleInvoiceReturn
// @Param		code  path      string  true  "SaleInvoiceReturn Code"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/sale-invoice-return/code/{code} [get]
func (h SaleInvoiceReturnHttp) InfoSaleInvoiceReturnByCode(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	code := ctx.Param("code")

	doc, err := h.svc.InfoSaleInvoiceReturnByCode(shopID, code)

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

// List SaleInvoiceReturn step godoc
// @Description get list step
// @Tags		SaleInvoiceReturn
// @Param		custcode	query	string		false  "customer code"
// @Param		q		query	string		false  "Search Value"
// @Param		fromdate	query	string		false  "from date"
// @Param		todate	query	string		false  "to date"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Limit"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/sale-invoice-return [get]
func (h SaleInvoiceReturnHttp) SearchSaleInvoiceReturnPage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)

	filters := utils.GetFilters(ctx.QueryParam, []utils.FilterRequest{
		{
			Param: "custcode",
			Type:  "string",
		},
		{
			Param: "-",
			Field: "docdatetime",
			Type:  "rangeDate",
		},
	})

	docList, pagination, err := h.svc.SearchSaleInvoiceReturn(shopID, filters, pageable)

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

// List SaleInvoiceReturn godoc
// @Description search limit offset
// @Tags		SaleInvoiceReturn
// @Param		q		query	string		false  "Search Value"
// @Param		fromdate	query	string		false  "from date"
// @Param		todate	query	string		false  "to date"
// @Param		offset	query	integer		false  "offset"
// @Param		limit	query	integer		false  "limit"
// @Param		lang	query	string		false  "lang"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/sale-invoice-return/list [get]
func (h SaleInvoiceReturnHttp) SearchSaleInvoiceReturnStep(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageableStep := utils.GetPageableStep(ctx.QueryParam)

	lang := ctx.QueryParam("lang")

	filters := utils.GetFilters(ctx.QueryParam, []utils.FilterRequest{
		{
			Param: "custcode",
			Type:  "string",
		},
		{
			Param: "-",
			Field: "docdatetime",
			Type:  "rangeDate",
		},
	})

	docList, total, err := h.svc.SearchSaleInvoiceReturnStep(shopID, lang, filters, pageableStep)

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

// Create SaleInvoiceReturn Bulk godoc
// @Description Create SaleInvoiceReturn
// @Tags		SaleInvoiceReturn
// @Param		SaleInvoiceReturn  body      []models.SaleInvoiceReturn  true  "SaleInvoiceReturn"
// @Accept 		json
// @Success		201	{object}	common.BulkReponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/sale-invoice-return/bulk [post]
func (h SaleInvoiceReturnHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.SaleInvoiceReturn{}
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
