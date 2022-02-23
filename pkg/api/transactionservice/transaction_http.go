package transactionservice

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"strconv"
)

type ITransactionHttp interface {
	CreateTransaction(ctx microservice.IServiceContext) error
	UpdateTransaction(ctx microservice.IServiceContext) error
	DeleteTransaction(ctx microservice.IServiceContext) error
	InfoTransaction(ctx microservice.IServiceContext) error
	SearchTransaction(ctx microservice.IServiceContext) error
	SearchTransactionItems(ctx microservice.IServiceContext) error
}

type TransactionHttp struct {
	ms      *microservice.Microservice
	cfg     microservice.IConfig
	service ITransactionService
}

func NewTransactionHttp(ms *microservice.Microservice, cfg microservice.IConfig) ITransactionHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	transRepo := NewTransactionRepository(pst)
	service := NewTransactionService(transRepo)
	return &TransactionHttp{
		ms:      ms,
		cfg:     cfg,
		service: service,
	}
}

func (h *TransactionHttp) CreateTransaction(ctx microservice.IServiceContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	merchantId := userInfo.MerchantId

	input := ctx.ReadInput()

	trans := &models.Transaction{}
	err := json.Unmarshal([]byte(input), &trans)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.service.CreateTransaction(merchantId, authUsername, *trans)

	if err != nil {
		ctx.ResponseError(400, err.Error())
	}

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
		Id:      idx,
	})

	return nil
}

func (h *TransactionHttp) UpdateTransaction(ctx microservice.IServiceContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	merchantId := userInfo.MerchantId

	id := ctx.Param("id")
	input := ctx.ReadInput()

	transReq := &models.Transaction{}
	err := json.Unmarshal([]byte(input), &transReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.service.UpdateTransaction(id, merchantId, authUsername, *transReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})
	return nil
}

func (h *TransactionHttp) DeleteTransaction(ctx microservice.IServiceContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	merchantId := userInfo.MerchantId

	id := ctx.Param("id")

	err := h.service.DeleteTransaction(id, merchantId, authUsername)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})
	return nil
}

func (h *TransactionHttp) InfoTransaction(ctx microservice.IServiceContext) error {

	userInfo := ctx.UserInfo()
	merchantId := userInfo.MerchantId

	id := ctx.Param("id")

	trans, err := h.service.InfoTransaction(id, merchantId)

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

func (h *TransactionHttp) SearchTransaction(ctx microservice.IServiceContext) error {

	userInfo := ctx.UserInfo()
	merchantId := userInfo.MerchantId

	q := ctx.QueryParam("q")
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))

	if err != nil {
		limit = 20
	}

	docList, pagination, err := h.service.SearchTransaction(merchantId, q, page, limit)

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

func (h *TransactionHttp) SearchTransactionItems(ctx microservice.IServiceContext) error {

	userInfo := ctx.UserInfo()
	merchantId := userInfo.MerchantId

	transId := ctx.Param("trans_id")

	q := ctx.QueryParam("q")
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))

	if err != nil {
		limit = 20
	}

	docList, pagination, err := h.service.SearchItemsTransaction(transId, merchantId, q, page, limit)

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
