package purchaseorder

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/config"
	mastersync "smlcloudplatform/internal/mastersync/repositories"
	common "smlcloudplatform/internal/models"
	"smlcloudplatform/internal/transaction/purchaseorder/models"
	"smlcloudplatform/internal/transaction/purchaseorder/repositories"
	"smlcloudplatform/internal/transaction/purchaseorder/services"
	trancache "smlcloudplatform/internal/transaction/repositories"
	"smlcloudplatform/internal/utils"
	"smlcloudplatform/internal/utils/requestfilter"
	"smlcloudplatform/pkg/microservice"
)

type IPurchaseOrderHttp interface{}

type PurchaseOrderHttp struct {
	ms  *microservice.Microservice
	cfg config.IConfig
	svc services.IPurchaseOrderHttpService
}

func NewPurchaseOrderHttp(ms *microservice.Microservice, cfg config.IConfig) PurchaseOrderHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())
	producer := ms.Producer(cfg.MQConfig())

	repo := repositories.NewPurchaseOrderRepository(pst)
	repoMq := repositories.NewPurchaseOrderMessageQueueRepository(producer)

	transRepo := trancache.NewCacheRepository(cache)
	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)
	svc := services.NewPurchaseOrderHttpService(repo, transRepo, repoMq, masterSyncCacheRepo)

	return PurchaseOrderHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h PurchaseOrderHttp) RegisterHttp() {

	h.ms.POST("/transaction/purchase-order/bulk", h.SaveBulk)

	h.ms.GET("/transaction/purchase-order", h.SearchPurchaseOrderPage)
	h.ms.GET("/transaction/purchase-order/list", h.SearchPurchaseOrderStep)
	h.ms.POST("/transaction/purchase-order", h.CreatePurchaseOrder)
	h.ms.GET("/transaction/purchase-order/:id", h.InfoPurchaseOrder)
	h.ms.GET("/transaction/purchase-order/code/:code", h.InfoPurchaseOrderByCode)
	h.ms.PUT("/transaction/purchase-order/:id", h.UpdatePurchaseOrder)
	h.ms.DELETE("/transaction/purchase-order/:id", h.DeletePurchaseOrder)
	h.ms.DELETE("/transaction/purchase-order", h.DeletePurchaseOrderByGUIDs)
}

// Create PurchaseOrder godoc
// @Description Create PurchaseOrder
// @Tags		PurchaseOrder
// @Param		PurchaseOrder  body      models.PurchaseOrder  true  "PurchaseOrder"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/purchase-order [post]
func (h PurchaseOrderHttp) CreatePurchaseOrder(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.PurchaseOrder{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, docNo, err := h.svc.CreatePurchaseOrder(shopID, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		ID:      idx,
		Data:    docNo,
	})
	return nil
}

// Update PurchaseOrder godoc
// @Description Update PurchaseOrder
// @Tags		PurchaseOrder
// @Param		id  path      string  true  "PurchaseOrder ID"
// @Param		PurchaseOrder  body      models.PurchaseOrder  true  "PurchaseOrder"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/purchase-order/{id} [put]
func (h PurchaseOrderHttp) UpdatePurchaseOrder(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.PurchaseOrder{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdatePurchaseOrder(shopID, id, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		ID:      id,
	})

	return nil
}

// Delete PurchaseOrder godoc
// @Description Delete PurchaseOrder
// @Tags		PurchaseOrder
// @Param		id  path      string  true  "PurchaseOrder ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/purchase-order/{id} [delete]
func (h PurchaseOrderHttp) DeletePurchaseOrder(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeletePurchaseOrder(shopID, id, authUsername)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		ID:      id,
	})

	return nil
}

// Delete PurchaseOrder godoc
// @Description Delete PurchaseOrder
// @Tags		PurchaseOrder
// @Param		PurchaseOrder  body      []string  true  "PurchaseOrder GUIDs"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/purchase-order [delete]
func (h PurchaseOrderHttp) DeletePurchaseOrderByGUIDs(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	input := ctx.ReadInput()

	docReq := []string{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.DeletePurchaseOrderByGUIDs(shopID, authUsername, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}

// Get PurchaseOrder godoc
// @Description get PurchaseOrder info by guidfixed
// @Tags		PurchaseOrder
// @Param		id  path      string  true  "PurchaseOrder guidfixed"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/purchase-order/{id} [get]
func (h PurchaseOrderHttp) InfoPurchaseOrder(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get PurchaseOrder %v", id)
	doc, err := h.svc.InfoPurchaseOrder(shopID, id)

	if err != nil {
		h.ms.Logger.Errorf("Error getting document %v: %v", id, err)
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    doc,
	})
	return nil
}

// Get PurchaseOrder By Code godoc
// @Description get PurchaseOrder info by Code
// @Tags		PurchaseOrder
// @Param		code  path      string  true  "PurchaseOrder Code"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/purchase-order/code/{code} [get]
func (h PurchaseOrderHttp) InfoPurchaseOrderByCode(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	code := ctx.Param("code")

	doc, err := h.svc.InfoPurchaseOrderByCode(shopID, code)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    doc,
	})
	return nil
}

// List PurchaseOrder step godoc
// @Description get list step
// @Tags		PurchaseOrder
// @Param		q		query	string		false  "Search Value"
// @Param		custcode	query	string		false  "cust code"
// @Param		branchcode	query	string		false  "branch code"
// @Param		fromdate	query	string		false  "from date"
// @Param		todate	query	string		false  "to date"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Limit"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/purchase-order [get]
func (h PurchaseOrderHttp) SearchPurchaseOrderPage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)

	filters := requestfilter.GenerateFilters(ctx.QueryParam, []requestfilter.FilterRequest{
		{
			Param: "custcode",
			Type:  requestfilter.FieldTypeString,
		},
		{
			Param: "-",
			Field: "docdatetime",
			Type:  requestfilter.FieldTypeRangeDate,
		},
		{
			Param: "branchcode",
			Field: "branch.code",
			Type:  requestfilter.FieldTypeString,
		},
	})

	docList, pagination, err := h.svc.SearchPurchaseOrder(shopID, filters, pageable)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success:    true,
		Data:       docList,
		Pagination: pagination,
	})
	return nil
}

// List PurchaseOrder godoc
// @Description search limit offset
// @Tags		PurchaseOrder
// @Param		q		query	string		false  "Search Value"
// @Param		custcode	query	string		false  "cust code"
// @Param		branchcode	query	string		false  "branch code"
// @Param		fromdate	query	string		false  "from date"
// @Param		todate	query	string		false  "to date"
// @Param		offset	query	integer		false  "offset"
// @Param		limit	query	integer		false  "limit"
// @Param		lang	query	string		false  "lang"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/purchase-order/list [get]
func (h PurchaseOrderHttp) SearchPurchaseOrderStep(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageableStep := utils.GetPageableStep(ctx.QueryParam)

	lang := ctx.QueryParam("lang")

	filters := requestfilter.GenerateFilters(ctx.QueryParam, []requestfilter.FilterRequest{
		{
			Param: "custcode",
			Type:  requestfilter.FieldTypeString,
		},
		{
			Param: "-",
			Field: "docdatetime",
			Type:  requestfilter.FieldTypeRangeDate,
		},
		{
			Param: "branchcode",
			Field: "branch.code",
			Type:  requestfilter.FieldTypeString,
		},
	})

	docList, total, err := h.svc.SearchPurchaseOrderStep(shopID, lang, filters, pageableStep)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    docList,
		Total:   total,
	})
	return nil
}

// Create PurchaseOrder Bulk godoc
// @Description Create PurchaseOrder
// @Tags		PurchaseOrder
// @Param		PurchaseOrder  body      []models.PurchaseOrder  true  "PurchaseOrder"
// @Accept 		json
// @Success		201	{object}	common.BulkResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/purchase-order/bulk [post]
func (h PurchaseOrderHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.PurchaseOrder{}
	err := json.Unmarshal([]byte(input), &dataReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	bulkResponse, err := h.svc.SaveInBatch(shopID, authUsername, dataReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	ctx.Response(
		http.StatusCreated,
		common.BulkResponse{
			Success:    true,
			BulkImport: bulkResponse,
		},
	)

	return nil
}
