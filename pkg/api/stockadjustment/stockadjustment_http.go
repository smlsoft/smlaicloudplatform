package stockadjustment

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"strconv"
)

type IStockAdjustmentHttp interface {
	RouteSetup()
	CreateStockAdjustment(ctx microservice.IContext) error
	UpdateStockAdjustment(ctx microservice.IContext) error
	DeleteStockAdjustment(ctx microservice.IContext) error
	InfoStockAdjustment(ctx microservice.IContext) error
	SearchStockAdjustment(ctx microservice.IContext) error
	SearchStockAdjustmentItems(ctx microservice.IContext) error
}

type StockAdjustmentHttp struct {
	ms      *microservice.Microservice
	cfg     microservice.IConfig
	service IStockAdjustmentService
}

func NewStockAdjustmentHttp(ms *microservice.Microservice, cfg microservice.IConfig) StockAdjustmentHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	prod := ms.Producer(cfg.MQConfig())

	stockadjustmentRepo := NewStockAdjustmentRepository(pst)
	stockadjustmentMQRepo := NewStockAdjustmentMQRepository(prod)

	service := NewStockAdjustmentService(stockadjustmentRepo, stockadjustmentMQRepo)
	return StockAdjustmentHttp{
		ms:      ms,
		cfg:     cfg,
		service: service,
	}
}

func (h StockAdjustmentHttp) RouteSetup() {

	h.ms.GET("/stockadjustment/:id", h.InfoStockAdjustment)
	h.ms.GET("/stockadjustment", h.SearchStockAdjustment)
	h.ms.GET("/stockadjustment/:id/items", h.SearchStockAdjustmentItems)

	h.ms.POST("/stockadjustment", h.CreateStockAdjustment)
	h.ms.PUT("/stockadjustment/:id", h.UpdateStockAdjustment)
	h.ms.DELETE("/stockadjustment/:id", h.DeleteStockAdjustment)
}

func (h StockAdjustmentHttp) CreateStockAdjustment(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	doc := models.StockAdjustment{}
	err := json.Unmarshal([]byte(input), &doc)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.service.CreateStockAdjustment(shopID, authUsername, doc)

	if err != nil {
		ctx.ResponseError(400, err.Error())
	}

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
		ID:      idx,
	})

	return nil
}

func (h StockAdjustmentHttp) UpdateStockAdjustment(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.StockAdjustment{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.service.UpdateStockAdjustment(id, shopID, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})
	return nil
}

func (h StockAdjustmentHttp) DeleteStockAdjustment(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	err := h.service.DeleteStockAdjustment(id, shopID, authUsername)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})
	return nil
}

func (h StockAdjustmentHttp) InfoStockAdjustment(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	doc, err := h.service.InfoStockAdjustment(id, shopID)

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

func (h StockAdjustmentHttp) SearchStockAdjustment(ctx microservice.IContext) error {

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

	docList, pagination, err := h.service.SearchStockAdjustment(shopID, q, page, limit)

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

func (h StockAdjustmentHttp) SearchStockAdjustmentItems(ctx microservice.IContext) error {

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

	docList, pagination, err := h.service.SearchItemsStockAdjustment(docID, shopID, q, page, limit)

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
