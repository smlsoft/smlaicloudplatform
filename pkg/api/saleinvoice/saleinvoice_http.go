package saleinvoice

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"strconv"
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

	err = h.service.UpdateSaleinvoice(id, shopID, authUsername, *transReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})
	return nil
}

func (h SaleinvoiceHttp) DeleteSaleinvoice(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	err := h.service.DeleteSaleinvoice(id, shopID, authUsername)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})
	return nil
}

func (h SaleinvoiceHttp) InfoSaleinvoice(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	trans, err := h.service.InfoSaleinvoice(id, shopID)

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

func (h SaleinvoiceHttp) SearchSaleinvoice(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))

	if err != nil {
		limit = 20
	}

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
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))

	if err != nil {
		limit = 20
	}

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
