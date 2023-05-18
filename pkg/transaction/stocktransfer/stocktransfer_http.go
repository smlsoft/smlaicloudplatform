package stocktransfer

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/transaction/stocktransfer/models"
	"smlcloudplatform/pkg/transaction/stocktransfer/repositories"
	"smlcloudplatform/pkg/transaction/stocktransfer/services"
	"smlcloudplatform/pkg/utils"
)

type IStockTransferHttp interface{}

type StockTransferHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc services.IStockTransferHttpService
}

func NewStockTransferHttp(ms *microservice.Microservice, cfg microservice.IConfig) StockTransferHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	repo := repositories.NewStockTransferRepository(pst)

	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)
	svc := services.NewStockTransferHttpService(repo, masterSyncCacheRepo)

	return StockTransferHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h StockTransferHttp) RouteSetup() {

	h.ms.GET("/transaction/stock-transfer", h.SearchStockTransferPage)
	h.ms.GET("/transaction/stock-transfer/list", h.SearchStockTransferStep)
	h.ms.POST("/transaction/stock-transfer", h.CreateStockTransfer)
	h.ms.GET("/transaction/stock-transfer/:id", h.InfoStockTransfer)
	h.ms.GET("/transaction/stock-transfer/code/:code", h.InfoStockTransferByCode)
	h.ms.PUT("/transaction/stock-transfer/:id", h.UpdateStockTransfer)
	h.ms.DELETE("/transaction/stock-transfer/:id", h.DeleteStockTransfer)
	h.ms.DELETE("/transaction/stock-transfer", h.DeleteStockTransferByGUIDs)
}

// Create StockTransfer godoc
// @Description Create StockTransfer
// @Tags		StockTransfer
// @Param		StockTransfer  body      models.StockTransfer  true  "StockTransfer"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/stock-transfer [post]
func (h StockTransferHttp) CreateStockTransfer(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.StockTransfer{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateStockTransfer(shopID, authUsername, *docReq)

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

// Update StockTransfer godoc
// @Description Update StockTransfer
// @Tags		StockTransfer
// @Param		id  path      string  true  "StockTransfer ID"
// @Param		StockTransfer  body      models.StockTransfer  true  "StockTransfer"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/stock-transfer/{id} [put]
func (h StockTransferHttp) UpdateStockTransfer(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.StockTransfer{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateStockTransfer(shopID, id, authUsername, *docReq)

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

// Delete StockTransfer godoc
// @Description Delete StockTransfer
// @Tags		StockTransfer
// @Param		id  path      string  true  "StockTransfer ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/stock-transfer/{id} [delete]
func (h StockTransferHttp) DeleteStockTransfer(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteStockTransfer(shopID, id, authUsername)

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

// Delete StockTransfer godoc
// @Description Delete StockTransfer
// @Tags		StockTransfer
// @Param		StockTransfer  body      []string  true  "StockTransfer GUIDs"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/stock-transfer [delete]
func (h StockTransferHttp) DeleteStockTransferByGUIDs(ctx microservice.IContext) error {
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

	err = h.svc.DeleteStockTransferByGUIDs(shopID, authUsername, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}

// Get StockTransfer godoc
// @Description get StockTransfer info by guidfixed
// @Tags		StockTransfer
// @Param		id  path      string  true  "StockTransfer guidfixed"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/stock-transfer/{id} [get]
func (h StockTransferHttp) InfoStockTransfer(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get StockTransfer %v", id)
	doc, err := h.svc.InfoStockTransfer(shopID, id)

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

// Get StockTransfer By Code godoc
// @Description get StockTransfer info by Code
// @Tags		StockTransfer
// @Param		code  path      string  true  "StockTransfer Code"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/stock-transfer/code/{code} [get]
func (h StockTransferHttp) InfoStockTransferByCode(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	code := ctx.Param("code")

	doc, err := h.svc.InfoStockTransferByCode(shopID, code)

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

// List StockTransfer step godoc
// @Description get list step
// @Tags		StockTransfer
// @Param		custcode	query	string		false  "customer code"
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Limit"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/stock-transfer [get]
func (h StockTransferHttp) SearchStockTransferPage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)

	filters := utils.GetFilters(ctx.QueryParam, []utils.FilterRequest{
		{
			Param: "custcode",
			Type:  "string",
		},
	})

	docList, pagination, err := h.svc.SearchStockTransfer(shopID, filters, pageable)

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

// List StockTransfer godoc
// @Description search limit offset
// @Tags		StockTransfer
// @Param		q		query	string		false  "Search Value"
// @Param		offset	query	integer		false  "offset"
// @Param		limit	query	integer		false  "limit"
// @Param		lang	query	string		false  "lang"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /transaction/stock-transfer/list [get]
func (h StockTransferHttp) SearchStockTransferStep(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageableStep := utils.GetPageableStep(ctx.QueryParam)

	lang := ctx.QueryParam("lang")

	docList, total, err := h.svc.SearchStockTransferStep(shopID, lang, pageableStep)

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
