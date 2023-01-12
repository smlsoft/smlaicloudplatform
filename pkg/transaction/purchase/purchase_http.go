package purchase

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/transaction/purchase/models"
	"smlcloudplatform/pkg/utils"
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

func NewPurchaseHttp(ms *microservice.Microservice, cfg microservice.IConfig) PurchaseHttp {

	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	prod := ms.Producer(cfg.MQConfig())

	purchaseRepo := NewPurchaseRepository(pst)
	purchaseMQRepo := NewPurchaseMQRepository(prod)

	service := NewPurchaseService(purchaseRepo, purchaseMQRepo)
	return PurchaseHttp{
		ms:      ms,
		cfg:     cfg,
		service: service,
	}
}

func (h PurchaseHttp) RouteSetup() {

	h.ms.GET("/purchase/:id", h.InfoPurchase)
	h.ms.GET("/purchase", h.SearchPurchase)
	h.ms.GET("/purchase/:id/items", h.SearchPurchaseItems)

	h.ms.POST("/purchase", h.CreatePurchase)
	h.ms.PUT("/purchase/:id", h.UpdatePurchase)
	h.ms.DELETE("/purchase/:id", h.DeletePurchase)
}

// Create Purchase Transaction godoc
// @Description Create Purchase Transaction
// @Tags		Purchase
// @Param		Purchase  body      models.Purchase  true  "payload"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /purchase [post]
func (h PurchaseHttp) CreatePurchase(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	doc := models.Purchase{}
	err := json.Unmarshal([]byte(input), &doc)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.service.CreatePurchase(shopID, authUsername, doc)

	if err != nil {
		ctx.ResponseError(400, err.Error())
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		ID:      idx,
	})

	return nil
}

// Update Inventory godoc
// @Description Update Inventory
// @Tags		Purchase
// @Param		id  path      string  true  "Purchase Document ID"
// @Param		Purchase  body      models.Purchase  true  "payload"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /purchase/{id} [put]
func (h PurchaseHttp) UpdatePurchase(ctx microservice.IContext) error {
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

	err = h.service.UpdatePurchase(shopID, id, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})
	return nil
}

// Delete Purchase Transaction godoc
// @Description Tempolaty Delete  Purchase Transaction
// @Tags		Purchase
// @Param		id  path      string  true  "Purchase Doc ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /purchase/{id} [delete]
func (h PurchaseHttp) DeletePurchase(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	err := h.service.DeletePurchase(shopID, id, authUsername)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})
	return nil
}

// Get Purchase By Doc guid godoc
// @Description After Purchare Doc Create system will be return Id of Purchase Document. this Id can get Purchase Document in this service.
// @Tags		Purchase
// @Param		id  path      string  true  "Document ID"
// @Accept 		json
// @Success		200	{object}	models.PurchaseInfo
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /purchase/{id} [get]
func (h PurchaseHttp) InfoPurchase(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	doc, err := h.service.InfoPurchase(shopID, id)

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

// List Purchase Transaction godoc
// @Tags		Purchase
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Add Category"
// @Param		limit	query	integer		false  "Add Category"
// @Accept 		json
// @Success		200	{array}	models.PurchaseListPageResponse
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /purchase [get]
func (h PurchaseHttp) SearchPurchase(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)

	docList, pagination, err := h.service.SearchPurchase(shopID, q, page, limit)

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

func (h PurchaseHttp) SearchPurchaseItems(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	docID := ctx.Param("id")

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)

	docList, pagination, err := h.service.SearchItemsPurchase(docID, shopID, q, page, limit)

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
