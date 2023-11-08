package stockbalanceimport

import (
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/config"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	productbarcode_repo "smlcloudplatform/pkg/product/productbarcode/repositories"
	"smlcloudplatform/pkg/stockbalanceimport/models"
	"smlcloudplatform/pkg/stockbalanceimport/repositories"
	"smlcloudplatform/pkg/stockbalanceimport/services"
	trancache "smlcloudplatform/pkg/transaction/repositories"
	stockbalance_models "smlcloudplatform/pkg/transaction/stockbalance/models"
	stockbalance_repositories "smlcloudplatform/pkg/transaction/stockbalance/repositories"
	stockbalance_serrvices "smlcloudplatform/pkg/transaction/stockbalance/services"
	stockbalancedetail_serrvices "smlcloudplatform/pkg/transaction/stockbalancedetail/services"

	stockbalancedetail_repositories "smlcloudplatform/pkg/transaction/stockbalancedetail/repositories"
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
	pstClickHouse := ms.ClickHousePersister(cfg.ClickHouseConfig())

	repo := stockbalance_repositories.NewStockBalanceRepository(pst)
	repoMq := stockbalance_repositories.NewStockBalanceMessageQueueRepository(producer)

	transRepo := trancache.NewCacheRepository(cache)
	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)
	stockbalanceDetailRepo := stockbalancedetail_repositories.NewStockBalanceDetailRepository(pst)
	stockbalanceDetailMqRepo := stockbalancedetail_repositories.NewStockBalanceDetailMessageQueueRepository(producer)

	stockBalanceDetailSvc := stockbalancedetail_serrvices.NewStockBalanceDetailHttpService(stockbalanceDetailRepo, transRepo, stockbalanceDetailMqRepo, masterSyncCacheRepo)
	stockBalanceSvc := stockbalance_serrvices.NewStockBalanceHttpService(stockBalanceDetailSvc, repo, transRepo, repoMq, masterSyncCacheRepo)

	chRepo := repositories.NewStockBalanceImportClickHouseRepository(pstClickHouse)

	productBarcodeRepo := productbarcode_repo.NewProductBarcodeRepository(pst, cache)

	svc := services.NewStockBalanceImportService(chRepo, productBarcodeRepo, stockBalanceSvc, stockBalanceDetailSvc, utils.RandStringBytesMaskImprSrcUnsafe, utils.NewGUID)

	return StockBalanceImportHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h StockBalanceImportHttp) RegisterHttp() {
	h.ms.POST("/stockbalanceimport/upload", h.UploadExcel)
	h.ms.POST("/stockbalanceimport", h.Create)
	h.ms.GET("/stockbalanceimport/:task-id", h.List)
	h.ms.DELETE("/stockbalanceimport/task/:task-id", h.DeleteByTask)
	h.ms.POST("/stockbalanceimport/task/:task-id", h.SaveTask)
	h.ms.PUT("/stockbalanceimport/:guid", h.Update)
	h.ms.DELETE("/stockbalanceimport/:guid", h.Delete)
}

// Create StockBalanceImport godoc
// @Description Create StockBalanceImport
// @Tags		StockBalanceImport
// @Param		file  formData      file  true  "excel file"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /stockbalanceimport/upload [post]
func (h StockBalanceImportHttp) UploadExcel(ctx microservice.IContext) error {
	shopID := ctx.UserInfo().ShopID
	tempFile, err := ctx.FormFile("file")

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	// Check if the file is an Excel file
	if filepath.Ext(tempFile.Filename) != ".xlsx" {
		ctx.ResponseError(400, "Invalid file xlsx type.")
		return errors.New("invalid file type")
	}

	file, err := tempFile.Open()

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}
	defer file.Close()

	taskID, err := h.svc.ImportFromFile(shopID, file)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		ID:      taskID,
	})

	return nil
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
func (h StockBalanceImportHttp) Create(ctx microservice.IContext) error {
	shopID := ctx.UserInfo().ShopID

	input := ctx.ReadInput()

	docReq := models.StockBalanceImport{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.Create(shopID, &docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
	})
	return nil
}

// List StockBalanceImport godoc
// @Description List StockBalanceImport
// @Tags		StockBalanceImport
// @Param		task-id		path		string		true		"task id"
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Limit"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /stockbalanceimport/{task-id} [get]
func (h StockBalanceImportHttp) List(ctx microservice.IContext) error {
	shopID := ctx.UserInfo().ShopID

	taskID := ctx.Param("task-id")

	if taskID == "" {
		ctx.Response(http.StatusCreated, common.ApiResponse{
			Success: true,
			Data:    []string{},
		})
		return nil
	}

	pageable := utils.GetPageable(ctx.QueryParam)

	results, page, err := h.svc.List(shopID, taskID, pageable)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success:    true,
		Pagination: page,
		Data:       results,
	})
	return nil
}

// Get StockBalanceImport Part godoc
// @Description Get StockBalanceImport Part
// @Tags		StockBalanceImport
// @Param		guid		path		string		true		"guid"
// @Accept 		json
// @Success		201	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /stockbalanceimport/{guid} [put]
func (h StockBalanceImportHttp) Update(ctx microservice.IContext) error {
	shopID := ctx.UserInfo().ShopID

	input := ctx.ReadInput()

	guid := ctx.Param("guid")

	docReq := models.StockBalanceImportRaw{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.Update(shopID, guid, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
	})
	return nil
}

// Delete StockBalanceImport By GUID godoc
// @Description Delete StockBalanceImport By GUID
// @Tags		StockBalanceImport
// @Param		guid		path		string		true		"guid"
// @Accept 		json
// @Success		201	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /stockbalanceimport/{guid} [delete]
func (h StockBalanceImportHttp) Delete(ctx microservice.IContext) error {
	shopID := ctx.UserInfo().ShopID

	guid := ctx.Param("guid")

	err := h.svc.Delete(shopID, guid)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
	})
	return nil
}

// Delete StockBalanceImport By Task ID godoc
// @Description Delete StockBalanceImport By Task ID
// @Tags		StockBalanceImport
// @Param		task-id		path		string		true		"task id"
// @Accept 		json
// @Success		201	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /stockbalanceimport/task/{task-id} [delete]
func (h StockBalanceImportHttp) DeleteByTask(ctx microservice.IContext) error {
	shopID := ctx.UserInfo().ShopID

	taskID := ctx.Param("task-id")

	err := h.svc.DeleteTask(shopID, taskID)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
	})
	return nil
}

// Delete StockBalanceImport By Task ID godoc
// @Description Delete StockBalanceImport By Task ID
// @Tags		StockBalanceImport
// @Param		task-id		path		string		true		"task id"
// @Param		StockBalanceHeader  body      models.StockBalanceHeader  true  "Stock Balance Header"
// @Accept 		json
// @Success		201	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /stockbalanceimport/task/{task-id} [post]
func (h StockBalanceImportHttp) SaveTask(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	taskID := ctx.Param("task-id")

	input := ctx.ReadInput()
	docReq := stockbalance_models.StockBalanceHeader{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	docNo, err := h.svc.SaveTask(shopID, authUsername, taskID, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
		DocNo:   docNo,
	})
	return nil
}
