package productbarcode

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/productbarcode/models"
	"smlcloudplatform/pkg/product/productbarcode/repositories"
	"smlcloudplatform/pkg/product/productbarcode/services"
	"smlcloudplatform/pkg/utils"
)

type IProductBarcodeHttp interface{}

type ProductBarcodeHttp struct {
	ms  *microservice.Microservice
	cfg microservice.IConfig
	svc services.IProductBarcodeHttpService
}

func NewProductBarcodeHttp(ms *microservice.Microservice, cfg microservice.IConfig) ProductBarcodeHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	repo := repositories.NewProductBarcodeRepository(pst)

	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)
	svc := services.NewProductBarcodeHttpService(repo, masterSyncCacheRepo)

	return ProductBarcodeHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h ProductBarcodeHttp) RouteSetup() {

	h.ms.POST("/product/barcode/bulk", h.SaveBulk)

	h.ms.GET("/product/barcode", h.SearchProductBarcodePage)
	h.ms.GET("/product/barcode/list", h.SearchProductBarcodeLimit)
	h.ms.POST("/product/barcode", h.CreateProductBarcode)
	h.ms.GET("/product/barcode/:id", h.InfoProductBarcode)
	h.ms.PUT("/product/barcode/xsort", h.UpdateProductBarcodeXSort)
	h.ms.PUT("/product/barcode/:id", h.UpdateProductBarcode)
	h.ms.DELETE("/product/barcode/:id", h.DeleteProductBarcode)
}

// Create ProductBarcode godoc
// @Description Create ProductBarcode
// @Tags		ProductBarcode
// @Param		ProductBarcode  body      models.ProductBarcode  true  "ProductBarcode"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/barcode [post]
func (h ProductBarcodeHttp) CreateProductBarcode(ctx microservice.IContext) error {
	authUsername := ctx.UserInfo().Username
	shopID := ctx.UserInfo().ShopID
	input := ctx.ReadInput()

	docReq := &models.ProductBarcode{}
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
// @Param		ProductBarcode  body      models.ProductBarcode  true  "ProductBarcode"
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

	docReq := &models.ProductBarcode{}
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

	err = h.svc.XSortsSave(shopID, *req)

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

// List ProductBarcode godoc
// @Description get struct array by ID
// @Tags		ProductBarcode
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Add Category"
// @Param		limit	query	integer		false  "Add Category"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/barcode [get]
func (h ProductBarcodeHttp) SearchProductBarcodePage(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	q := ctx.QueryParam("q")
	page, limit := utils.GetPaginationParam(ctx.QueryParam)
	sort := utils.GetSortParam(ctx.QueryParam)
	docList, pagination, err := h.svc.SearchProductBarcode(shopID, q, page, limit, sort)

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

	q := ctx.QueryParam("q")
	offset, limit := utils.GetParamOffsetLimit(ctx.QueryParam)
	sorts := utils.GetSortParam(ctx.QueryParam)

	lang := ctx.QueryParam("lang")

	docList, total, err := h.svc.SearchProductBarcodeStep(shopID, lang, q, offset, limit, sorts)

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
