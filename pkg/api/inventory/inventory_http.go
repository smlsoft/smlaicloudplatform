package inventory

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"strconv"
)

type IInventoryHttp interface {
	RouteSetup()
	CreateInventory(ctx microservice.IContext) error
	UpdateInventory(ctx microservice.IContext) error
	DeleteInventory(ctx microservice.IContext) error
	InfoInventory(ctx microservice.IContext) error
	SearchInventory(ctx microservice.IContext) error
}

type InventoryHttp struct {
	ms              *microservice.Microservice
	cfg             microservice.IConfig
	invService      IInventoryService
	cateService     ICategoryService
	invOptService   IInventoryOptionService
	optGroupService IOptionGroupService
}

func NewInventoryHttp(ms *microservice.Microservice, cfg microservice.IConfig) IInventoryHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())

	invRepo := NewInventoryRepository(pst)
	invService := NewInventoryService(invRepo)

	cateRepo := NewCategoryRepository(pst)
	cateService := NewCategoryService(cateRepo)

	invOptRepo := NewInventoryOptionRepository(pst)
	invOptService := NewInventoryOptionService(invOptRepo)

	optGroupRepo := NewOptionGroupRepository(pst)
	optGroupService := NewOptionGroupService(optGroupRepo)

	return &InventoryHttp{
		ms:              ms,
		cfg:             cfg,
		invService:      invService,
		cateService:     cateService,
		invOptService:   invOptService,
		optGroupService: optGroupService,
	}
}

func (h *InventoryHttp) RouteSetup() {
	h.ms.GET("/inventory/:id", h.InfoInventory)
	h.ms.GET("/inventory", h.SearchInventory)
	h.ms.POST("/inventory", h.CreateInventory)
	h.ms.PUT("/inventory/:id", h.UpdateInventory)
	h.ms.DELETE("/inventory/:id", h.DeleteInventory)

	h.ms.GET("/category/:id", h.InfoCategory)
	h.ms.GET("/category", h.SearchCategory)
	h.ms.POST("/category", h.CreateCategory)
	h.ms.PUT("/category/:id", h.UpdateCategory)
	h.ms.DELETE("/category/:id", h.DeleteCategory)

	h.ms.GET("/option/:id", h.InfoInventoryOption)
	h.ms.GET("/option", h.SearchInventoryOption)
	h.ms.POST("/option", h.CreateInventoryOption)
	h.ms.PUT("/option/:id", h.UpdateInventoryOption)
	h.ms.DELETE("/option/:id", h.DeleteInventoryOption)

	h.ms.GET("/optgroup/:id", h.InfoOptionGroup)
	h.ms.GET("/optgroup", h.SearchOptionGroup)
	h.ms.POST("/optgroup", h.CreateOptionGroup)
	h.ms.PUT("/optgroup/:id", h.UpdateOptionGroup)
	h.ms.DELETE("/optgroup/:id", h.DeleteOptionGroup)
}

func (h *InventoryHttp) CreateInventory(ctx microservice.IContext) error {
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

func (h *InventoryHttp) UpdateInventory(ctx microservice.IContext) error {
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

func (h *InventoryHttp) DeleteInventory(ctx microservice.IContext) error {
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

func (h *InventoryHttp) InfoInventory(ctx microservice.IContext) error {

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

func (h *InventoryHttp) SearchInventory(ctx microservice.IContext) error {
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
