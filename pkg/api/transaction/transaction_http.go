package transaction

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"strconv"
)

type ITransactionHttp interface {
	RouteSetup()
	CreateTransaction(ctx microservice.IContext) error
	UpdateTransaction(ctx microservice.IContext) error
	DeleteTransaction(ctx microservice.IContext) error
	InfoTransaction(ctx microservice.IContext) error
	SearchTransaction(ctx microservice.IContext) error
	SearchTransactionItems(ctx microservice.IContext) error
}

type TransactionHttp struct {
	ms      *microservice.Microservice
	cfg     microservice.IConfig
	service ITransactionService
}

func NewTransactionHttp(ms *microservice.Microservice, cfg microservice.IConfig) ITransactionHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	prod := ms.Producer(cfg.MQConfig())

	transRepo := NewTransactionRepository(pst)
	mqRepo := NewTransactionMQRepository(prod)

	service := NewTransactionService(transRepo, mqRepo)
	return &TransactionHttp{
		ms:      ms,
		cfg:     cfg,
		service: service,
	}
}

func (h *TransactionHttp) RouteSetup() {

	h.ms.GET("/transaction/:id", h.InfoTransaction)
	h.ms.GET("/transaction", h.SearchTransaction)
	h.ms.GET("/transaction/:id/items", h.SearchTransactionItems)

	h.ms.POST("/transaction", h.CreateTransaction)
	h.ms.PUT("/transaction/:id", h.UpdateTransaction)
	h.ms.DELETE("/transaction/:id", h.DeleteTransaction)
}

func (h *TransactionHttp) CreateTransaction(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	trans := &models.Transaction{}
	err := json.Unmarshal([]byte(input), &trans)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.service.CreateTransaction(shopID, authUsername, trans)

	if err != nil {
		ctx.ResponseError(400, err.Error())
	}

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
		ID:      idx,
	})

	return nil
}

func (h *TransactionHttp) UpdateTransaction(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	transReq := &models.Transaction{}
	err := json.Unmarshal([]byte(input), &transReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.service.UpdateTransaction(id, shopID, authUsername, *transReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})
	return nil
}

func (h *TransactionHttp) DeleteTransaction(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	err := h.service.DeleteTransaction(id, shopID, authUsername)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})
	return nil
}

func (h *TransactionHttp) InfoTransaction(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	trans, err := h.service.InfoTransaction(id, shopID)

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

func (h *TransactionHttp) SearchTransaction(ctx microservice.IContext) error {

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

	docList, pagination, err := h.service.SearchTransaction(shopID, q, page, limit)

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

func (h *TransactionHttp) SearchTransactionItems(ctx microservice.IContext) error {

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

	docList, pagination, err := h.service.SearchItemsTransaction(transID, shopID, q, page, limit)

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
