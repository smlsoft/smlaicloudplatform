package stockbalanceimport

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/config"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/stockbalanceimport/models"
	"smlcloudplatform/pkg/stockbalanceimport/repositories"
	"smlcloudplatform/pkg/stockbalanceimport/services"
	trancache "smlcloudplatform/pkg/transaction/repositories"
	stockbalance_models "smlcloudplatform/pkg/transaction/stockbalance/models"
	stockbalance_repositories "smlcloudplatform/pkg/transaction/stockbalance/repositories"
	stockbalance_serrvices "smlcloudplatform/pkg/transaction/stockbalance/services"
	"smlcloudplatform/pkg/utils"
)

type IStockBalanceImportHttp interface{}

type StockBalanceImportHttp struct {
	ms  *microservice.Microservice
	cfg config.IConfig
	svc services.IStockBalanceImportService
}

func NewStockBalanceImportHttp(ms *microservice.Microservice, cfg config.IConfig) StockBalanceImportHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())
	producer := ms.Producer(cfg.MQConfig())

	repo := stockbalance_repositories.NewStockBalanceRepository(pst)
	repoMq := stockbalance_repositories.NewStockBalanceMessageQueueRepository(producer)

	transRepo := trancache.NewCacheRepository(cache)
	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)

	stockBalanceSvc := stockbalance_serrvices.NewStockBalanceHttpService(repo, transRepo, repoMq, masterSyncCacheRepo)
	cacheRepo := repositories.NewStockBalanceImportCacheRepository(cache)

	svc := services.NewStockBalanceImportService(cacheRepo, stockBalanceSvc, utils.RandStringBytesMaskImprSrcUnsafe)

	return StockBalanceImportHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h StockBalanceImportHttp) RegisterHttp() {

	h.ms.POST("/stockbalanceimport/task", h.CreateStockBalanceImport)
	h.ms.GET("/stockbalanceimport/task/:id", h.GetStockBalanceImportMeta)
	h.ms.POST("/stockbalanceimport/task/:id", h.SaveTaskComplete)
	h.ms.POST("/stockbalanceimport/task/part/:id", h.SaveStockBalanceImportPart)
	h.ms.GET("/stockbalanceimport/task/part/:id", h.GetStockBalanceImportPart)
}

// Create StockBalanceImport godoc
// @Description Create StockBalanceImport
// @Tags		StockBalanceImport
// @Param		StockBalanceImport  body      models.StockBalanceImport  true  "StockBalanceImport"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /stockbalanceimport [post]
func (h StockBalanceImportHttp) CreateStockBalanceImport(ctx microservice.IContext) error {
	shopID := ctx.UserInfo().ShopID

	input := ctx.ReadInput()

	docReq := models.StockBalanceImportTaskRequest{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateTask(shopID, docReq)

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

// Get StockBalanceImport Part godoc
// @Description Get StockBalanceImport Part
// @Tags		StockBalanceImport
// @Param		id		path		string		true		"StockBalanceImport ID"
// @Accept 		json
// @Success		201	{object}	StockBalanceImportPartResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /stockbalanceimport/task/part/{id} [get]
func (h StockBalanceImportHttp) GetStockBalanceImportPart(ctx microservice.IContext) error {
	shopID := ctx.UserInfo().ShopID

	cacheID := ctx.Param("id")

	result, err := h.svc.GetTaskPart(shopID, cacheID)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		Data:    result,
	})
	return nil
}

// Get StockBalanceImport Meta godoc
// @Description Get StockBalanceImport Meta
// @Tags		StockBalanceImport
// @Param		id		path		string		true		"StockBalanceImport ID"
// @Accept 		json
// @Success		201	{object}	StockBalanceImportMeta
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /stockbalanceimport/task/{id} [get]
func (h StockBalanceImportHttp) GetStockBalanceImportMeta(ctx microservice.IContext) error {
	shopID := ctx.UserInfo().ShopID

	cacheID := ctx.Param("id")

	result, err := h.svc.GetTaskMeta(shopID, cacheID)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		Data:    result,
	})
	return nil
}

// Create StockBalanceImport godoc
// @Description Create StockBalanceImport
// @Tags		StockBalanceImport
// @Param		StockBalanceImport  body      []models.StockBalanceImport  true  "StockBalanceImport"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /stockbalanceimport/task/part/{id} [post]
func (h StockBalanceImportHttp) SaveStockBalanceImportPart(ctx microservice.IContext) error {
	shopID := ctx.UserInfo().ShopID

	cacheID := ctx.Param("id")

	input := ctx.ReadInput()

	docReq := []stockbalance_models.StockBalanceDetail{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.SaveTaskPart(shopID, cacheID, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
	})
	return nil
}

func (h StockBalanceImportHttp) SaveTaskComplete(ctx microservice.IContext) error {
	shopID := ctx.UserInfo().ShopID

	authUsername := ctx.UserInfo().Username

	taskID := ctx.Param("id")

	result, err := h.svc.SaveTaskComplete(shopID, authUsername, taskID)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		Data:    result,
	})
	return nil
}
