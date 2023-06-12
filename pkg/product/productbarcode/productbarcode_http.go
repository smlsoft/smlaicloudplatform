package productbarcode

import (
	"encoding/json"
	"net/http"
	"net/url"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/config"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/productbarcode/models"
	"smlcloudplatform/pkg/product/productbarcode/repositories"
	"smlcloudplatform/pkg/product/productbarcode/services"
	"smlcloudplatform/pkg/utils"
	"strings"
)

type IProductBarcodeHttp interface{}

type ProductBarcodeHttp struct {
	ms  *microservice.Microservice
	cfg config.IConfig
	svc services.IProductBarcodeHttpService
}

func NewProductBarcodeHttp(ms *microservice.Microservice, cfg config.IConfig) ProductBarcodeHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	pstClickHouse := ms.ClickHousePersister(cfg.ClickHouseConfig())
	cache := ms.Cacher(cfg.CacherConfig())
	prod := ms.Producer(cfg.MQConfig())

	repo := repositories.NewProductBarcodeRepository(pst, cache)
	clickHouseRepo := repositories.NewProductBarcodeClickhouseRepository(pstClickHouse)
	mqRepo := repositories.NewProductBarcodeMessageQueueRepository(prod)
	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)
	svc := services.NewProductBarcodeHttpService(repo, mqRepo, clickHouseRepo, masterSyncCacheRepo)

	return ProductBarcodeHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h ProductBarcodeHttp) RouteSetup() {

	h.ms.POST("/product/barcode/bulk", h.SaveBulk)

	h.ms.GET("/product/barcode", h.SearchProductBarcodePage)
	h.ms.GET("/product/barcode2", h.SearchProductBarcodePage2)
	h.ms.GET("/product/barcode/list", h.SearchProductBarcodeLimit)
	h.ms.POST("/product/barcode", h.CreateProductBarcode)
	h.ms.GET("/product/barcode/:id", h.InfoProductBarcode)
	h.ms.GET("/product/barcode/ref/:barcode", h.GetroductBarcodeByRef)
	h.ms.GET("/product/barcode/pk/:barcode", h.InfoProductBarcodeByBarcode)
	h.ms.GET("/product/barcode/by-code", h.InfoArray)
	h.ms.GET("/product/barcode/master", h.InfoArrayMaster)
	h.ms.PUT("/product/barcode/xsort", h.UpdateProductBarcodeXSort)
	h.ms.PUT("/product/barcode/:id", h.UpdateProductBarcode)
	h.ms.DELETE("/product/barcode/:id", h.DeleteProductBarcode)
	h.ms.DELETE("/product/barcode", h.DeleteProductBarcodeByGUIDs)

	h.ms.GET("/product/barcode/units", h.GetroductBarcodeByAllUnits)
}

// Create ProductBarcode godoc
// @Description Create ProductBarcode
// @Tags		ProductBarcode
// @Param		ProductBarcode  body      models.ProductBarcodeRequest  true  "ProductBarcode"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/barcode [post]
func (h ProductBarcodeHttp) CreateProductBarcode(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.ProductBarcodeRequest{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if docReq.XSorts == nil {
		docReq.XSorts = &[]common.XSort{}
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	idx, err := h.svc.CreateProductBarcode(shopID, authUsername, *docReq)

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

// Update ProductBarcode godoc
// @Description Update ProductBarcode
// @Tags		ProductBarcode
// @Param		id  path      string  true  "ProductBarcode ID"
// @Param		ProductBarcode  body      models.ProductBarcodeRequest  true  "ProductBarcode"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/barcode/{id} [put]
func (h ProductBarcodeHttp) UpdateProductBarcode(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	id := ctx.Param("id")
	input := ctx.ReadInput()

	docReq := &models.ProductBarcodeRequest{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if docReq.XSorts == nil {
		docReq.XSorts = &[]common.XSort{}
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateProductBarcode(shopID, id, authUsername, *docReq)

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

// Update XSort	 ProductBarcode godoc
// @Description Update XSort ProductBarcode
// @Tags		ProductBarcode
// @Param		XSort  body      []common.XSortModifyReqesut  true  "XSort"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/barcode/xsort [put]
func (h ProductBarcodeHttp) UpdateProductBarcodeXSort(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	req := &[]common.XSortModifyReqesut{}
	err := json.Unmarshal([]byte(input), &req)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(req); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.XSortsSave(shopID, authUsername, *req)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
	})

	return nil
}

// Delete ProductBarcode godoc
// @Description Delete ProductBarcode
// @Tags		ProductBarcode
// @Param		id  path      string  true  "ProductBarcode ID"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/barcode/{id} [delete]
func (h ProductBarcodeHttp) DeleteProductBarcode(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	authUsername := userInfo.Username

	id := ctx.Param("id")

	err := h.svc.DeleteProductBarcode(shopID, id, authUsername)

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

// Get ProductBarcode godoc
// @Description get struct array by ID
// @Tags		ProductBarcode
// @Param		id  path      string  true  "ProductBarcode ID"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/barcode/{id} [get]
func (h ProductBarcodeHttp) InfoProductBarcode(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	id := ctx.Param("id")

	h.ms.Logger.Debugf("Get ProductBarcode %v", id)
	doc, err := h.svc.InfoProductBarcode(shopID, id)

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

// Get ProductBarcode By Reference Barcode godoc
// @Description get by reference barcode
// @Tags		ProductBarcode
// @Param		barcode  path      string  true  "Reference Barcode"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/barcode/ref/{barcode} [get]
func (h ProductBarcodeHttp) GetroductBarcodeByRef(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	refBarcode := ctx.Param("barcode")

	docs, err := h.svc.GetProductBarcodeByBarcodeRef(shopID, refBarcode)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    docs,
	})
	return nil
}

// Get ProductBarcode By Barcode godoc
// @Description get data by barcode
// @Tags		ProductBarcode
// @Param		barcode  path      string  true  "Barcode"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/barcode/pk/{barcode} [get]
func (h ProductBarcodeHttp) InfoProductBarcodeByBarcode(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	barcode := ctx.Param("barcode")

	doc, err := h.svc.InfoProductBarcodeByBarcode(shopID, barcode)

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

// Get ProductBarcode By code array godoc
// @Description get ProductBarcode by code array
// @Tags		ProductBarcode
// @Param		codes	query	string		false  "Code filter, json array encode "
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/barcode/by-code [get]
func (h ProductBarcodeHttp) InfoArray(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	codesReq, err := url.QueryUnescape(ctx.QueryParam("codes"))

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	docReq := []string{}
	err = json.Unmarshal([]byte(codesReq), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	// where to filter array
	doc, err := h.svc.InfoWTFArray(shopID, docReq)

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

// Get Master ProductBarcode By code array godoc
// @Description get master ProductBarcode by code array
// @Tags		ProductBarcode
// @Param		codes	query	string		false  "Code filter, json array encode "
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/barcode/master [get]
func (h ProductBarcodeHttp) InfoArrayMaster(ctx microservice.IContext) error {
	codesReq, err := url.QueryUnescape(ctx.QueryParam("codes"))

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	docReq := []string{}
	err = json.Unmarshal([]byte(codesReq), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	// where to filter array master
	doc, err := h.svc.InfoWTFArrayMaster(docReq)

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

// List ProductBarcode godoc
// @Description get struct array by ID
// @Tags		ProductBarcode
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Limit"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/barcode [get]
func (h ProductBarcodeHttp) SearchProductBarcodePage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)
	docList, pagination, err := h.svc.SearchProductBarcode(shopID, pageable)

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

// List ProductBarcode2 godoc
// @Description get struct array by ID
// @Tags		ProductBarcode2
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Limit"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/barcode2 [get]
func (h ProductBarcodeHttp) SearchProductBarcodePage2(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)
	docList, pagination, err := h.svc.SearchProductBarcode2(shopID, pageable)

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

// List ProductBarcode godoc
// @Description search limit offset
// @Tags		ProductBarcode
// @Param		q		query	string		false  "Search Value"
// @Param		offset	query	integer		false  "offset"
// @Param		limit	query	integer		false  "limit"
// @Param		lang	query	string		false  "lang"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/barcode/list [get]
func (h ProductBarcodeHttp) SearchProductBarcodeLimit(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageableStep := utils.GetPageableStep(ctx.QueryParam)

	lang := ctx.QueryParam("lang")

	docList, total, err := h.svc.SearchProductBarcodeStep(shopID, lang, pageableStep)

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

// Create ProductBarcode Bulk godoc
// @Description Create ProductBarcode
// @Tags		ProductBarcode
// @Param		ProductBarcode  body      []models.ProductBarcode  true  "ProductBarcode"
// @Accept 		json
// @Success		201	{object}	common.BulkReponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/barcode/bulk [post]
func (h ProductBarcodeHttp) SaveBulk(ctx microservice.IContext) error {

	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	dataReq := []models.ProductBarcode{}
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

// Delete ProductBarcode By GUIDs godoc
// @Description Delete ProductBarcode
// @Tags		ProductBarcode
// @Param		ProductBarcode  body      []string  true  "ProductBarcode GUIDs"
// @Accept 		json
// @Success		200	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/barcode [delete]
func (h ProductBarcodeHttp) DeleteProductBarcodeByGUIDs(ctx microservice.IContext) error {
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

	err = h.svc.DeleteProductBarcodeByGUIDs(shopID, authUsername, docReq)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
	})

	return nil
}

// Get ProductBarcode By Reference Barcode godoc
// @Description get by reference barcode
// @Tags		ProductBarcode
// @Accept 		json
// @Param		codes	query	string		false  "array of units"
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/barcode/units [get]
func (h ProductBarcodeHttp) GetroductBarcodeByAllUnits(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)

	// inputBody := ctx.ReadInput()

	// unitCodes := []string{}
	// err := json.Unmarshal([]byte(inputBody), &unitCodes)

	// if err != nil {
	// 	ctx.ResponseError(400, err.Error())
	// 	return err
	// }

	reqUnitCodes := ctx.QueryParam("codes")
	unitCodes := []string{}

	tempUnitCodes := strings.Split(reqUnitCodes, ",")
	unitCodes = append(unitCodes, tempUnitCodes...)

	docs, pagination, err := h.svc.GetProductBarcodeByUnits(shopID, unitCodes, pageable)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success:    true,
		Pagination: pagination,
		Data:       docs,
	})
	return nil
}
