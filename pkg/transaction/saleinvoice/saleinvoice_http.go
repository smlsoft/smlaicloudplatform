package saleinvoice

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/transaction/saleinvoice/models"
	"smlcloudplatform/pkg/transaction/saleinvoice/repositories"
	"smlcloudplatform/pkg/transaction/saleinvoice/services"
	"smlcloudplatform/pkg/utils"
)

type ISaleInvoiceHttp interface{}

type SaleInvoiceHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc services.ISaleInvoiceHttpService
}

func NewSaleInvoiceHttp(ms *microservice.Microservice, cfg microservice.IConfig) SaleInvoiceHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	repo := repositories.NewSaleInvoiceRepository(pst)

	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)
	svc := services.NewSaleInvoiceHttpService(repo, masterSyncCacheRepo)

	return SaleInvoiceHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h SaleInvoiceHttp) RouteSetup() {

	h.ms.POST("/transaction/sale-invoice/bulk", h.SaveBulk)

	h.ms.GET("/transaction/sale-invoice", h.SearchSaleInvoicePage)
	h.ms.GET("/transaction/sale-invoice/list", h.SearchSaleInvoiceStep)
	h.ms.POST("/transaction/sale-invoice", h.CreateSaleInvoice)
	h.ms.GET("/transaction/sale-invoice/:id", h.InfoSaleInvoice)
	h.ms.GET("/transaction/sale-invoice/code/:code", h.InfoSaleInvoiceByCode)
	h.ms.PUT("/transaction/sale-invoice/:id", h.UpdateSaleInvoice)
	h.ms.DELETE("/transaction/sale-invoice/:id", h.DeleteSaleInvoice)
	h.ms.DELETE("/transaction/sale-invoice", h.DeleteSaleInvoiceByGUIDs)
}

// Create SaleInvoice godoc
// @Description Create SaleInvoice
// @Tags		SaleInvoice
// @Param		SaleInvoice  body      models.SaleInvoice  true  "SaleInvoice"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/sale-invoice [post]
func (h SaleInvoiceHttp) CreateSaleInvoice(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.SaleInvoice{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateSaleInvoice(shopID, authUsername, *docReq)

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

// Update SaleInvoice godoc
// @Description Update SaleInvoice
// @Tags		SaleInvoice
// @Param		id  path      string  true  "SaleInvoice ID"
// @Param		SaleInvoice  body      models.SaleInvoice  true  "SaleInvoice"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/sale-invoice/{id} [put]
func (h SaleInvoiceHttp) UpdateSaleInvoice(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.SaleInvoice{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateSaleInvoice(shopID, id, authUsername, *docReq)

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

// Delete SaleInvoice godoc
// @Description Delete SaleInvoice
// @Tags		SaleInvoice
// @Param		id  path      string  true  "SaleInvoice ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/sale-invoice/{id} [delete]
func (h SaleInvoiceHttp) DeleteSaleInvoice(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteSaleInvoice(shopID, id, authUsername)

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

// Delete SaleInvoice godoc
// @Description Delete SaleInvoice
// @Tags		SaleInvoice
// @Param		SaleInvoice  body      []string  true  "SaleInvoice GUIDs"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/sale-invoice [delete]
func (h SaleInvoiceHttp) DeleteSaleInvoiceByGUIDs(ctx microservice.IContext) error {
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

	err = h.svc.DeleteSaleInvoiceByGUIDs(shopID, authUsername, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}

// Get SaleInvoice godoc
// @Description get SaleInvoice info by guidfixed
// @Tags		SaleInvoice
// @Param		id  path      string  true  "SaleInvoice guidfixed"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/sale-invoice/{id} [get]
func (h SaleInvoiceHttp) InfoSaleInvoice(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get SaleInvoice %v", id)
	doc, err := h.svc.InfoSaleInvoice(shopID, id)

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

// Get SaleInvoice By Code godoc
// @Description get SaleInvoice info by Code
// @Tags		SaleInvoice
// @Param		code  path      string  true  "SaleInvoice Code"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/sale-invoice/code/{code} [get]
func (h SaleInvoiceHttp) InfoSaleInvoiceByCode(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	code := ctx.Param("code")

	doc, err := h.svc.InfoSaleInvoiceByCode(shopID, code)

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

// List SaleInvoice step godoc
// @Description get list step
// @Tags		SaleInvoice
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Limit"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/sale-invoice [get]
func (h SaleInvoiceHttp) SearchSaleInvoicePage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)

	filters := utils.GetFilters(ctx.QueryParam, []utils.FilterRequest{
		{
			Param: "custcode",
			Type:  "string",
		},
	})

	docList, pagination, err := h.svc.SearchSaleInvoice(shopID, filters, pageable)

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

// List SaleInvoice godoc
// @Description search limit offset
// @Tags		SaleInvoice
// @Param		q		query	string		false  "Search Value"
// @Param		offset	query	integer		false  "offset"
// @Param		limit	query	integer		false  "limit"
// @Param		lang	query	string		false  "lang"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/sale-invoice/list [get]
func (h SaleInvoiceHttp) SearchSaleInvoiceStep(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageableStep := utils.GetPageableStep(ctx.QueryParam)

	lang := ctx.QueryParam("lang")

	docList, total, err := h.svc.SearchSaleInvoiceStep(shopID, lang, pageableStep)

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

// Create SaleInvoice Bulk godoc
// @Description Create SaleInvoice
// @Tags		SaleInvoice
// @Param		SaleInvoice  body      []models.SaleInvoice  true  "SaleInvoice"
// @Accept 		json
// @Success		201	{object}	common.BulkReponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/sale-invoice/bulk [post]
func (h SaleInvoiceHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.SaleInvoice{}
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
