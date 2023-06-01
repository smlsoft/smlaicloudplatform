package stockreceiveproduct

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/transaction/stockreceiveproduct/models"
	"smlcloudplatform/pkg/transaction/stockreceiveproduct/repositories"
	"smlcloudplatform/pkg/transaction/stockreceiveproduct/services"
	"smlcloudplatform/pkg/utils"
)

type IStockReceiveProductHttp interface{}

type StockReceiveProductHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc services.IStockReceiveProductHttpService
}

func NewStockReceiveProductHttp(ms *microservice.Microservice, cfg microservice.IConfig) StockReceiveProductHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	repo := repositories.NewStockReceiveProductRepository(pst)

	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)
	svc := services.NewStockReceiveProductHttpService(repo, masterSyncCacheRepo)

	return StockReceiveProductHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h StockReceiveProductHttp) RouteSetup() {

	h.ms.POST("/transaction/stock-receive-product/bulk", h.SaveBulk)

	h.ms.GET("/transaction/stock-receive-product", h.SearchStockReceiveProductPage)
	h.ms.GET("/transaction/stock-receive-product/list", h.SearchStockReceiveProductStep)
	h.ms.POST("/transaction/stock-receive-product", h.CreateStockReceiveProduct)
	h.ms.GET("/transaction/stock-receive-product/:id", h.InfoStockReceiveProduct)
	h.ms.GET("/transaction/stock-receive-product/code/:code", h.InfoStockReceiveProductByCode)
	h.ms.PUT("/transaction/stock-receive-product/:id", h.UpdateStockReceiveProduct)
	h.ms.DELETE("/transaction/stock-receive-product/:id", h.DeleteStockReceiveProduct)
	h.ms.DELETE("/transaction/stock-receive-product", h.DeleteStockReceiveProductByGUIDs)
}

// Create StockReceiveProduct godoc
// @Description Create StockReceiveProduct
// @Tags		StockReceiveProduct
// @Param		StockReceiveProduct  body      models.StockReceiveProduct  true  "StockReceiveProduct"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/stock-receive-product [post]
func (h StockReceiveProductHttp) CreateStockReceiveProduct(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.StockReceiveProduct{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateStockReceiveProduct(shopID, authUsername, *docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		ID:      idx,
	})
	return nil
}

// Update StockReceiveProduct godoc
// @Description Update StockReceiveProduct
// @Tags		StockReceiveProduct
// @Param		id  path      string  true  "StockReceiveProduct ID"
// @Param		StockReceiveProduct  body      models.StockReceiveProduct  true  "StockReceiveProduct"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/stock-receive-product/{id} [put]
func (h StockReceiveProductHttp) UpdateStockReceiveProduct(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.StockReceiveProduct{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateStockReceiveProduct(shopID, id, authUsername, *docReq)

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

// Delete StockReceiveProduct godoc
// @Description Delete StockReceiveProduct
// @Tags		StockReceiveProduct
// @Param		id  path      string  true  "StockReceiveProduct ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/stock-receive-product/{id} [delete]
func (h StockReceiveProductHttp) DeleteStockReceiveProduct(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteStockReceiveProduct(shopID, id, authUsername)

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

// Delete StockReceiveProduct godoc
// @Description Delete StockReceiveProduct
// @Tags		StockReceiveProduct
// @Param		StockReceiveProduct  body      []string  true  "StockReceiveProduct GUIDs"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/stock-receive-product [delete]
func (h StockReceiveProductHttp) DeleteStockReceiveProductByGUIDs(ctx microservice.IContext) error {
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

	err = h.svc.DeleteStockReceiveProductByGUIDs(shopID, authUsername, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}

// Get StockReceiveProduct godoc
// @Description get StockReceiveProduct info by guidfixed
// @Tags		StockReceiveProduct
// @Param		id  path      string  true  "StockReceiveProduct guidfixed"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/stock-receive-product/{id} [get]
func (h StockReceiveProductHttp) InfoStockReceiveProduct(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get StockReceiveProduct %v", id)
	doc, err := h.svc.InfoStockReceiveProduct(shopID, id)

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

// Get StockReceiveProduct By Code godoc
// @Description get StockReceiveProduct info by Code
// @Tags		StockReceiveProduct
// @Param		code  path      string  true  "StockReceiveProduct Code"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/stock-receive-product/code/{code} [get]
func (h StockReceiveProductHttp) InfoStockReceiveProductByCode(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	code := ctx.Param("code")

	doc, err := h.svc.InfoStockReceiveProductByCode(shopID, code)

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

// List StockReceiveProduct step godoc
// @Description get list step
// @Tags		StockReceiveProduct
// @Param		custcode	query	string		false  "customer code"
// @Param		q		query	string		false  "Search Value"
// @Param		fromdate	query	string		false  "from date"
// @Param		todate	query	string		false  "to date"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Limit"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/stock-receive-product [get]
func (h StockReceiveProductHttp) SearchStockReceiveProductPage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)

	filters := utils.GetFilters(ctx.QueryParam, []utils.FilterRequest{
		{
			Param: "custcode",
			Type:  "string",
		},
		{
			Param: "-",
			Field: "docdatetime",
			Type:  "rangeDate",
		},
	})

	docList, pagination, err := h.svc.SearchStockReceiveProduct(shopID, filters, pageable)

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

// List StockReceiveProduct godoc
// @Description search limit offset
// @Tags		StockReceiveProduct
// @Param		q		query	string		false  "Search Value"
// @Param		fromdate	query	string		false  "from date"
// @Param		todate	query	string		false  "to date"
// @Param		offset	query	integer		false  "offset"
// @Param		limit	query	integer		false  "limit"
// @Param		lang	query	string		false  "lang"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/stock-receive-product/list [get]
func (h StockReceiveProductHttp) SearchStockReceiveProductStep(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageableStep := utils.GetPageableStep(ctx.QueryParam)

	lang := ctx.QueryParam("lang")

	filters := utils.GetFilters(ctx.QueryParam, []utils.FilterRequest{
		{
			Param: "custcode",
			Type:  "string",
		},
		{
			Param: "-",
			Field: "docdatetime",
			Type:  "rangeDate",
		},
	})

	docList, total, err := h.svc.SearchStockReceiveProductStep(shopID, lang, filters, pageableStep)

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

// Create StockReceiveProduct Bulk godoc
// @Description Create StockReceiveProduct
// @Tags		StockReceiveProduct
// @Param		StockReceiveProduct  body      []models.StockReceiveProduct  true  "StockReceiveProduct"
// @Accept 		json
// @Success		201	{object}	common.BulkReponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/stock-receive-product/bulk [post]
func (h StockReceiveProductHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.StockReceiveProduct{}
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
		common.BulkReponse{
			Success:    true,
			BulkImport: bulkResponse,
		},
	)

	return nil
}
