package stockinout

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"strconv"
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

	ctx.Response(http.StatusCreated, models.ApiResponse{
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

	err = h.service.UpdateStockInOut(id, shopID, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})
	return nil
}

func (h StockInOutHttp) DeleteStockInOut(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	err := h.service.DeleteStockInOut(id, shopID, authUsername)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})
	return nil
}

func (h StockInOutHttp) InfoStockInOut(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	doc, err := h.service.InfoStockInOut(id, shopID)

	if err != nil && err.Error() != "mongo: no documents in result" {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
		Data:    doc,
	})
	return nil
}

func (h StockInOutHttp) SearchStockInOut(ctx microservice.IContext) error {

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

	docList, pagination, err := h.service.SearchStockInOut(shopID, q, page, limit)

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

func (h StockInOutHttp) SearchStockInOutItems(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	docID := ctx.Param("id")

	q := ctx.QueryParam("q")
	page, err := strconv.Atoi(ctx.QueryParam("page"))
	if err != nil {
		page = 1
	}

	limit, err := strconv.Atoi(ctx.QueryParam("limit"))

	if err != nil {
		limit = 20
	}

	docList, pagination, err := h.service.SearchItemsStockInOut(docID, shopID, q, page, limit)

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
