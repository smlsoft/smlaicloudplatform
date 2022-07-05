package stockinout

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/transaction/stockinout/models"
	"smlcloudplatform/pkg/utils"
)

type IStockInOutHttp interface {
	RouteSetup()
	CreateStockInOut(ctx microservice.IContext) error
	UpdateStockInOut(ctx microservice.IContext) error
	DeleteStockInOut(ctx microservice.IContext) error
	InfoStockInOut(ctx microservice.IContext) error
	SearchStockInOut(ctx microservice.IContext) error
	SearchStockInOutItems(ctx microservice.IContext) error
}

type StockInOutHttp struct {
	ms      *microservice.Microservice
	cfg     microservice.IConfig
	service IStockInOutService
}

func NewStockInOutHttp(ms *microservice.Microservice, cfg microservice.IConfig) StockInOutHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	prod := ms.Producer(cfg.MQConfig())

	stockinoutRepo := NewStockInOutRepository(pst)
	stockinoutMQRepo := NewStockInOutMQRepository(prod)

	service := NewStockInOutService(stockinoutRepo, stockinoutMQRepo)
	return StockInOutHttp{
		ms:      ms,
		cfg:     cfg,
		service: service,
	}
}

func (h StockInOutHttp) RouteSetup() {

	h.ms.GET("/stockinout/:id", h.InfoStockInOut)
	h.ms.GET("/stockinout", h.SearchStockInOut)
	h.ms.GET("/stockinout/:id/items", h.SearchStockInOutItems)

	h.ms.POST("/stockinout", h.CreateStockInOut)
	h.ms.PUT("/stockinout/:id", h.UpdateStockInOut)
	h.ms.DELETE("/stockinout/:id", h.DeleteStockInOut)
}

func (h StockInOutHttp) CreateStockInOut(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	doc := models.StockInOut{}
	err := json.Unmarshal([]byte(input), &doc)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.service.CreateStockInOut(shopID, authUsername, doc)

	if err != nil {
		ctx.ResponseError(400, err.Error())
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		ID:      idx,
	})

	return nil
}

func (h StockInOutHttp) UpdateStockInOut(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.StockInOut{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.service.UpdateStockInOut(shopID, id, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})
	return nil
}

func (h StockInOutHttp) DeleteStockInOut(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	err := h.service.DeleteStockInOut(shopID, id, authUsername)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})
	return nil
}

func (h StockInOutHttp) InfoStockInOut(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	doc, err := h.service.InfoStockInOut(shopID, id)

	if err != nil && err.Error() != "mongo: no documents in result" {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    doc,
	})
	return nil
}

func (h StockInOutHttp) SearchStockInOut(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)

	docList, pagination, err := h.service.SearchStockInOut(shopID, q, page, limit)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusOK,
		common.ApiResponse{
			Success:    true,
			Data:       docList,
			Pagination: pagination,
		})

	return nil
}

func (h StockInOutHttp) SearchStockInOutItems(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	docID := ctx.Param("id")

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)

	docList, pagination, err := h.service.SearchItemsStockInOut(shopID, docID, q, page, limit)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusOK,
		common.ApiResponse{
			Success:    true,
			Data:       docList,
			Pagination: pagination,
		})
	return nil
}
