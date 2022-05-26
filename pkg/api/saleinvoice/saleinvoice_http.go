package saleinvoice

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
)

type ISaleinvoiceHttp interface {
	RouteSetup()
	CreateSaleinvoice(ctx microservice.IContext) error
	UpdateSaleinvoice(ctx microservice.IContext) error
	DeleteSaleinvoice(ctx microservice.IContext) error
	InfoSaleinvoice(ctx microservice.IContext) error
	SearchSaleinvoice(ctx microservice.IContext) error
	SearchSaleinvoiceItems(ctx microservice.IContext) error
}

type SaleinvoiceHttp struct {
	ms      *microservice.Microservice
	cfg     microservice.IConfig
	service ISaleinvoiceService
}

func NewSaleinvoiceHttp(ms *microservice.Microservice, cfg microservice.IConfig) SaleinvoiceHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	prod := ms.Producer(cfg.MQConfig())

	transRepo := NewSaleinvoiceRepository(pst)
	mqRepo := NewSaleinvoiceMQRepository(prod)

	service := NewSaleinvoiceService(transRepo, mqRepo)
	return SaleinvoiceHttp{
		ms:      ms,
		cfg:     cfg,
		service: service,
	}
}

func (h SaleinvoiceHttp) RouteSetup() {

	h.ms.GET("/saleinvoice/:id", h.InfoSaleinvoice)
	h.ms.GET("/saleinvoice", h.SearchSaleinvoice)
	h.ms.GET("/saleinvoice/:id/items", h.SearchSaleinvoiceItems)

	h.ms.POST("/saleinvoice", h.CreateSaleinvoice)
	h.ms.PUT("/saleinvoice/:id", h.UpdateSaleinvoice)
	h.ms.DELETE("/saleinvoice/:id", h.DeleteSaleinvoice)
}

// Create Sale Invoice godoc
// @Description Create Inventory
// @Tags		Sale Invoice
// @Param		SaleInvoice  body      models.Saleinvoice  true  "SaleInvoice"
// @Accept 		json
// @Success		201	{object}	models.ResponseSuccessWithID
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /saleinvoice [post]
func (h SaleinvoiceHttp) CreateSaleinvoice(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	trans := models.Saleinvoice{}
	err := json.Unmarshal([]byte(input), &trans)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.service.CreateSaleinvoice(shopID, authUsername, trans)

	if err != nil {
		ctx.ResponseError(400, err.Error())
	}

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
		ID:      idx,
	})

	return nil
}

// Update Sale Invoice godoc
// @Description Update Sale Invoice
// @Tags		Sale Invoice
// @Param		id  path      string  true  "Document ID"
// @Param		Invoice  body      models.Saleinvoice  true  "Body"
// @Accept 		json
// @Success		201	{object}	models.ResponseSuccessWithID
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /saleinvoice/{id} [put]
func (h SaleinvoiceHttp) UpdateSaleinvoice(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	transReq := &models.Saleinvoice{}
	err := json.Unmarshal([]byte(input), &transReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.service.UpdateSaleinvoice(shopID, id, authUsername, *transReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})
	return nil
}

// Delete Sale Invoice Document godoc
// @Description Delete Document
// @Tags		Sale Invoice
// @Param		id  path      string  true  "Document ID"
// @Accept 		json
// @Success		200	{object}	models.ResponseSuccessWithID
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /saleinvoice/{id} [delete]
func (h SaleinvoiceHttp) DeleteSaleinvoice(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	err := h.service.DeleteSaleinvoice(shopID, id, authUsername)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})
	return nil
}

// Get Sale Invoice Document Info godoc
// @Description get struct array by ID
// @Tags		Sale Invoice
// @Param		id  path      string  true  "Inventory ID"
// @Accept 		json
// @Success		200	{object}	models.SaleinvoiceInfo
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /saleinvoice/{id} [get]
func (h SaleinvoiceHttp) InfoSaleinvoice(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	trans, err := h.service.InfoSaleinvoice(shopID, id)

	if err != nil && err.Error() != "mongo: no documents in result" {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
		Data:    trans,
	})
	return nil
}

// List Sale Invoice godoc
// @Description List Sale Invoice Document
// @Tags		Sale Invoice
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Accept 		json
// @Success		200	{array}	models.SaleInvoiceListPageResponse
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /saleinvoice [get]
func (h SaleinvoiceHttp) SearchSaleinvoice(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)

	docList, pagination, err := h.service.SearchSaleinvoice(shopID, q, page, limit)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusOK,
		models.ApiResponse{
			Success:    true,
			Data:       docList,
			Pagination: pagination,
		})

	return nil
}

func (h SaleinvoiceHttp) SearchSaleinvoiceItems(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	transID := ctx.Param("id")

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)

	docList, pagination, err := h.service.SearchItemsSaleinvoice(transID, shopID, q, page, limit)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusOK,
		models.ApiResponse{
			Success:    true,
			Data:       docList,
			Pagination: pagination,
		})
	return nil
}
