package productimport

import (
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"
	"smlaicloudplatform/internal/config"
	creditorRepo "smlaicloudplatform/internal/debtaccount/creditor/repositories"
	mastersync "smlaicloudplatform/internal/mastersync/repositories"
	common "smlaicloudplatform/internal/models"
	productmaster "smlaicloudplatform/internal/product/product/repositories"
	product_repositories "smlaicloudplatform/internal/product/productbarcode/repositories"
	product_serrvices "smlaicloudplatform/internal/product/productbarcode/services"
	productcategory_repositories "smlaicloudplatform/internal/product/productcategory/repositories"
	productcategory_services "smlaicloudplatform/internal/product/productcategory/services"
	productunit_repo "smlaicloudplatform/internal/product/unit/repositories"
	unitmaster "smlaicloudplatform/internal/product/unit/repositories"
	"smlaicloudplatform/internal/productimport/models"
	"smlaicloudplatform/internal/productimport/repositories"
	"smlaicloudplatform/internal/productimport/services"
	"smlaicloudplatform/internal/utils"
	"smlaicloudplatform/pkg/microservice"
	"time"
)

type IProductImportHttp interface{}

type ProductImportHttp struct {
	ms  *microservice.Microservice
	cfg config.IConfig
	svc services.IProductImportService
}

func NewProductImportHttp(ms *microservice.Microservice, cfg config.IConfig) ProductImportHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())
	producer := ms.Producer(cfg.MQConfig())
	pstClickHouse := ms.ClickHousePersister(cfg.ClickHouseConfig())
	pstPg := ms.Persister(cfg.PersisterConfig())

	repo := product_repositories.NewProductBarcodeRepository(pst, cache)
	unitmaster := unitmaster.NewUnitPGRepository(pstPg)
	repoMq := product_repositories.NewProductBarcodeMessageQueueRepository(producer)
	repoCh := product_repositories.NewProductBarcodeClickhouseRepository(pstClickHouse)
	creditorRepo := creditorRepo.NewCreditorRepository(pst)
	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)

	repoMaster := productmaster.NewProductPGRepository(pstPg)
	productcategoryRepo := productcategory_repositories.NewProductCategoryRepository(pst)
	productcategorySvc := productcategory_services.NewProductCategoryHttpService(productcategoryRepo, masterSyncCacheRepo)

	unitRepo := productunit_repo.NewUnitRepository(pst)

	chRepo := repositories.NewProductImportClickHouseRepository(pstClickHouse)
	stockBalanceSvc := product_serrvices.NewProductBarcodeHttpService(repo, repoMaster, unitmaster, *creditorRepo, repoMq, repoCh, productcategorySvc, masterSyncCacheRepo)

	svc := services.NewProductImportService(chRepo, repo, stockBalanceSvc, unitRepo, utils.RandStringBytesMaskImprSrcUnsafe, utils.NewGUID, time.Now)

	return ProductImportHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h ProductImportHttp) RegisterHttp() {
	h.ms.POST("/productimport/upload", h.UploadExcel)
	h.ms.POST("/productimport", h.Create)
	h.ms.GET("/productimport/:task-id", h.List)
	h.ms.DELETE("/productimport/:task-id", h.DeleteByTask)
	h.ms.POST("/productimport/:task-id", h.SaveTask)
	h.ms.PUT("/productimport/item/:guid", h.Update)
	h.ms.DELETE("/productimport/item/:guid", h.Delete)

	h.ms.POST("/productimport/:task-id/verify", h.Verify)
}

// Create ProductImport godoc
// @Description Create ProductImport
// @Tags		ProductImport
// @Param		file  formData      file  true  "excel file"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /productimport/upload [post]
func (h ProductImportHttp) UploadExcel(ctx microservice.IContext) error {
	shopID := ctx.UserInfo().ShopID
	authUsername := ctx.UserInfo().Username
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

	taskID, err := h.svc.ImportFromFile(shopID, authUsername, file)

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

// Create ProductImport godoc
// @Description Create ProductImport
// @Tags		ProductImport
// @Param		ProductImport  body      models.ProductImport  true  "ProductImport"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /productimport [post]
func (h ProductImportHttp) Create(ctx microservice.IContext) error {
	shopID := ctx.UserInfo().ShopID
	authUsername := ctx.UserInfo().Username

	input := ctx.ReadInput()

	docReq := models.ProductImport{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.Create(shopID, authUsername, &docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
	})
	return nil
}

// List ProductImport godoc
// @Description List ProductImport
// @Tags		ProductImport
// @Param		task-id		path		string		true		"task id"
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Limit"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /productimport/{task-id} [get]
func (h ProductImportHttp) List(ctx microservice.IContext) error {
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

// Get ProductImport Part godoc
// @Description Get ProductImport Part
// @Tags		ProductImport
// @Param		guid		path		string		true		"guid"
// @Accept 		json
// @Success		201	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /productimport/item/{guid} [put]
func (h ProductImportHttp) Update(ctx microservice.IContext) error {
	shopID := ctx.UserInfo().ShopID

	input := ctx.ReadInput()

	guid := ctx.Param("guid")

	docReq := models.ProductImportRaw{}
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

// Delete ProductImport By GUID godoc
// @Description Delete ProductImport By GUID
// @Tags		ProductImport
// @Param		guid		path		string		true		"guid"
// @Accept 		json
// @Success		201	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /productimport/item/{guid} [delete]
func (h ProductImportHttp) Delete(ctx microservice.IContext) error {
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

// Delete ProductImport By Task ID godoc
// @Description Delete ProductImport By Task ID
// @Tags		ProductImport
// @Param		task-id		path		string		true		"task id"
// @Accept 		json
// @Success		201	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /productimport/{task-id} [delete]
func (h ProductImportHttp) DeleteByTask(ctx microservice.IContext) error {
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

// Delete ProductImport By Task ID godoc
// @Description Delete ProductImport By Task ID
// @Tags		ProductImport
// @Param		task-id		path		string		true		"task id"
// @Param		ProductImportHeader  body      models.ProductImportHeader  true  "ProductImportHeader"
// @Accept 		json
// @Success		201	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /productimport/{task-id} [post]
func (h ProductImportHttp) SaveTask(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	taskID := ctx.Param("task-id")

	payload := ctx.ReadInput()

	var docReq models.ProductImportHeader
	err := json.Unmarshal([]byte(payload), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.SaveTask(shopID, authUsername, taskID, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
	})
	return nil
}

// Verify ProductImport By Task ID godoc
// @Description Verify ProductImport By Task ID
// @Tags		ProductImport
// @Param		task-id		path		string		true		"task id"
// @Accept 		json
// @Success		201	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /productimport/{task-id}/verify [post]
func (h ProductImportHttp) Verify(ctx microservice.IContext) error {
	shopID := ctx.UserInfo().ShopID

	guid := ctx.Param("task-id")

	err := h.svc.Verify(shopID, guid)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
	})
	return nil
}
