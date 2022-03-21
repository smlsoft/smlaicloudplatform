package purchase

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"strconv"
)

type IPurchaseHttp interface {
	RouteSetup()
	CreatePurchase(ctx microservice.IContext) error
	UpdatePurchase(ctx microservice.IContext) error
	DeletePurchase(ctx microservice.IContext) error
	InfoPurchase(ctx microservice.IContext) error
	SearchPurchase(ctx microservice.IContext) error
	SearchPurchaseItems(ctx microservice.IContext) error
}

type PurchaseHttp struct {
	ms      *microservice.Microservice
	cfg     microservice.IConfig
	service IPurchaseService
}

func NewPurchaseHttp(ms *microservice.Microservice, cfg microservice.IConfig) IPurchaseHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	prod := ms.Producer(cfg.MQConfig())

	purchaseRepo := NewPurchaseRepository(pst)
	purchaseMQRepo := NewPurchaseMQRepository(prod)

	service := NewPurchaseService(purchaseRepo, purchaseMQRepo)
	return &PurchaseHttp{
		ms:      ms,
		cfg:     cfg,
		service: service,
	}
}

func (h *PurchaseHttp) RouteSetup() {

	h.ms.GET("/purchase/:id", h.InfoPurchase)
	h.ms.GET("/purchase", h.SearchPurchase)
	h.ms.GET("/purchase/:id/items", h.SearchPurchaseItems)

	h.ms.POST("/purchase", h.CreatePurchase)
	h.ms.PUT("/purchase/:id", h.UpdatePurchase)
	h.ms.DELETE("/purchase/:id", h.DeletePurchase)
}

func (h *PurchaseHttp) CreatePurchase(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	doc := &models.Purchase{}
	err := json.Unmarshal([]byte(input), &doc)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.service.CreatePurchase(shopID, authUsername, doc)

	if err != nil {
		ctx.ResponseError(400, err.Error())
	}

	ctx.Response(http.StatusCreated, models.ApiResponse{
		Success: true,
		ID:      idx,
	})

	return nil
}

func (h *PurchaseHttp) UpdatePurchase(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.Purchase{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.service.UpdatePurchase(id, shopID, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})
	return nil
}

func (h *PurchaseHttp) DeletePurchase(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	err := h.service.DeletePurchase(id, shopID, authUsername)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, models.ApiResponse{
		Success: true,
	})
	return nil
}

func (h *PurchaseHttp) InfoPurchase(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	doc, err := h.service.InfoPurchase(id, shopID)

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

func (h *PurchaseHttp) SearchPurchase(ctx microservice.IContext) error {

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

	docList, pagination, err := h.service.SearchPurchase(shopID, q, page, limit)

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

func (h *PurchaseHttp) SearchPurchaseItems(ctx microservice.IContext) error {

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

	docList, pagination, err := h.service.SearchItemsPurchase(docID, shopID, q, page, limit)

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
