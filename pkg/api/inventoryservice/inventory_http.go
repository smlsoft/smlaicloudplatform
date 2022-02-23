package inventoryservice

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"strconv"
)

type IInventoryHttp interface {
	RouteSetup()
	CreateInventory(ctx microservice.IServiceContext) error
	UpdateInventory(ctx microservice.IServiceContext) error
	DeleteInventory(ctx microservice.IServiceContext) error
	InfoInventory(ctx microservice.IServiceContext) error
	SearchInventory(ctx microservice.IServiceContext) error
}

type InventoryHttp struct {
	invService IInventoryService
}

func NewInventoryHttp(ms *microservice.Microservice, cfg microservice.IConfig) IInventoryHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	invRepo := NewInventoryRepository(pst)
	invService := NewInventoryService(invRepo)
	return &InventoryHttp{
		invService: invService,
	}
}

func (h *InventoryHttp) RouteSetup() {
}

func (h *InventoryHttp) CreateInventory(ctx microservice.IServiceContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	merchantId := userInfo.MerchantId

	input := ctx.ReadInput()

	inventoryReq := &models.Inventory{}
	err := json.Unmarshal([]byte(input), &inventoryReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.invService.CreateInventory(merchantId, authUsername, *inventoryReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusCreated,
		models.ApiResponse{
			Success: true,
			Id:      idx,
		})

	return nil

}

func (h *InventoryHttp) UpdateInventory(ctx microservice.IServiceContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	merchantId := userInfo.MerchantId

	id := ctx.Param("id")
	input := ctx.ReadInput()

	inventoryReq := &models.Inventory{}
	err := json.Unmarshal([]byte(input), &inventoryReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.invService.UpdateInventory(id, merchantId, authUsername, *inventoryReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusOK,
		models.ApiResponse{
			Success: true,
			Id:      id,
		})

	return nil
}

func (h *InventoryHttp) DeleteInventory(ctx microservice.IServiceContext) error {
	userInfo := ctx.UserInfo()
	merchantId := userInfo.MerchantId

	id := ctx.Param("id")

	err := h.invService.DeleteInventory(id, merchantId)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusOK,
		models.ApiResponse{
			Success: true,
			Id:      id,
		},
	)
	return nil
}

func (h *InventoryHttp) InfoInventory(ctx microservice.IServiceContext) error {

	userInfo := ctx.UserInfo()
	merchantId := userInfo.MerchantId

	id := ctx.Param("id")

	doc, err := h.invService.InfoInventory(id, merchantId)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusOK,
		models.ApiResponse{
			Success: true,
			Data:    doc,
		},
	)

	return nil
}

func (h *InventoryHttp) SearchInventory(ctx microservice.IServiceContext) error {
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

	docList, pagination, err := h.invService.SearchInventory(merchantId, q, page, limit)

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
