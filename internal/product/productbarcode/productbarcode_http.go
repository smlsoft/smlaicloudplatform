package productbarcode

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"smlcloudplatform/internal/config"
	mastersync "smlcloudplatform/internal/mastersync/repositories"
	common "smlcloudplatform/internal/models"
	"smlcloudplatform/internal/product/productbarcode/models"
	"smlcloudplatform/internal/product/productbarcode/repositories"
	"smlcloudplatform/internal/product/productbarcode/services"
	productcategory_repositories "smlcloudplatform/internal/product/productcategory/repositories"
	productcategory_services "smlcloudplatform/internal/product/productcategory/services"
	"smlcloudplatform/internal/utils"
	"smlcloudplatform/internal/utils/requestfilter"
	"smlcloudplatform/pkg/microservice"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
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

	productcategoryRepo := productcategory_repositories.NewProductCategoryRepository(pst)
	productcategorySvc := productcategory_services.NewProductCategoryHttpService(productcategoryRepo, masterSyncCacheRepo)

	svc := services.NewProductBarcodeHttpService(repo, mqRepo, clickHouseRepo, productcategorySvc, masterSyncCacheRepo)

	return ProductBarcodeHttp{
		ms:  ms,
		cfg: cfg,
		svc: svc,
	}
}

func (h ProductBarcodeHttp) RegisterHttp() {

	h.ms.POST("/product/barcode/bulk", h.SaveBulk)
	h.ms.POST("/product/barcode/import", h.Import)

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
	h.ms.PUT("/product/barcode/branch", h.UpdateProductBarcodeBranch)
	h.ms.PUT("/product/barcode/business-type", h.UpdateProductBarcodeBusinessType)

	h.ms.DELETE("/product/barcode/:id", h.DeleteProductBarcode)
	h.ms.DELETE("/product/barcode", h.DeleteProductBarcodeByGUIDs)

	h.ms.GET("/product/barcode/units", h.GetroductBarcodeByAllUnits)
	h.ms.GET("/product/barcode/groups", h.GetroductBarcodeByGroups)

	h.ms.GET("/product/barcode/export", h.Export)

	h.ms.GET("/product/barcode/bom/:barcode", h.InfoBOMView)
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

// Update ProductBarcode Branch godoc
// @Description Update ProductBarcode Branch
// @Tags		ProductBarcode
// @Param		ProductBarcodeBranchRequest  body      models.ProductBarcodeBranchRequest  true  "Product BarcodeBranch Request"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/barcode/branch [put]
func (h ProductBarcodeHttp) UpdateProductBarcodeBranch(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	docReq := &models.ProductBarcodeBranchRequest{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateProductBarcodeBranch(shopID, authUsername, docReq.Branch, docReq.Products)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
	})

	return nil
}

// Update ProductBarcode Business Type  godoc
// @Description Update ProductBarcode Business Type
// @Tags		ProductBarcode
// @Param		ProductBarcodeBusinessTypeRequest  body      models.ProductBarcodeBusinessTypeRequest  true  "Product Barcode Business Type Request"
// @Accept 		json
// @Success		201	{object}	common.ResponseSuccessWithID
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/barcode/business-type [put]
func (h ProductBarcodeHttp) UpdateProductBarcodeBusinessType(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	authUsername := userInfo.Username
	shopID := userInfo.ShopID

	input := ctx.ReadInput()

	docReq := &models.ProductBarcodeBusinessTypeRequest{}
	err := json.Unmarshal([]byte(input), &docReq)

	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	if err = ctx.Validate(docReq); err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	err = h.svc.UpdateProductBarcodeBusinessType(shopID, authUsername, docReq.BusinessType, docReq.Products)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusCreated, common.ApiResponse{
		Success: true,
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
// @Param		businesstypecode		query	string		false  "business type code ex. bt1,bt2"
// @Param		branchcode		query	string		false  "branch code ex. b1,b2"
// @Param		isalacarte		query	boolean		false  "is A La Carte"
// @Param		isusesubbarcodes		query	boolean		false  "is use sub barcodes"
// @Param		isbom		query	boolean		false  "is use BOM"
// @Param		ordertypes		query	string		false  "order types ex. a01,a02"
// @Param		itemtype		query	int8		false  "item type"
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

	filters := h.searchFilter(ctx.QueryParam)

	docList, pagination, err := h.svc.SearchProductBarcode(shopID, filters, pageable)

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
// @Param		businesstypecode		query	string		false  "business type code ex. bt1,bt2"
// @Param		branchcode		query	string		false  "branch code ex. b1,b2"
// @Param		isalacarte		query	boolean		false  "is A La Carte"
// @Param		isusesubbarcodes		query	boolean		false  "is use sub barcodes"
// @Param		isbom		query	boolean		false  "is use BOM"
// @Param		ordertypes		query	string		false  "order types ex. a01,a02"
// @Param		itemtype		query	int8		false  "item type"
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

	filters := h.searchFilter(ctx.QueryParam)

	docList, total, err := h.svc.SearchProductBarcodeStep(shopID, lang, filters, pageableStep)

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
// @Success		201	{object}	common.BulkResponse
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
		common.BulkResponse{
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

// Get ProductBarcode By Groups
// @Description get by group codes
// @Tags		ProductBarcode
// @Accept 		json
// @Param		codes	query	string		false  "array of group"
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/barcode/groups [get]
func (h ProductBarcodeHttp) GetroductBarcodeByGroups(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	pageable := utils.GetPageable(ctx.QueryParam)

	reqUnitCodes := ctx.QueryParam("codes")
	groupCodes := []string{}

	tempUnitCodes := strings.Split(reqUnitCodes, ",")
	groupCodes = append(groupCodes, tempUnitCodes...)

	docs, pagination, err := h.svc.GetProductBarcodeByGroups(shopID, groupCodes, pageable)

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

// Get  Export
// @Description ProductBarcode Export
// @Tags		ProductBarcode
// @Param		lang	query	string		false  "language code"
// @Param		barcode	query	string		false  "Label Barcode"
// @Param		productname	query	string		false  "Label Product Name"
// @Param		unitcode	query	string		false  "Label Unit Code"
// @Param		unitname	query	string		false  "Label Unit Name"
// @Param		price	query	string		false  "Label Price"
// @Param		itemtype	query	string		false  "Label Item Type"
// @Param		groupcode	query	string		false  "Label Group Code"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/barcode/export [get]
func (h ProductBarcodeHttp) Export(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	languageCode := ctx.QueryParam("lang")

	if languageCode == "" {
		languageCode = "en"
	}

	keyCols := []string{
		"barcode",     //บาร์โค้ด",
		"productname", //"ชื่อสินค้า",
		"unitcode",    //"หน่วยนับ",
		"unitname",    //"ชื่อหน่วยนับ",
		"price",       //ราคาขาย",
		"itemtype",    //ประเภทสินค้า",
		"groupcode",   //กลุ่มสินค้า",
	}

	languageHeader := map[string]string{}
	for _, key := range keyCols {
		languageHeader[key] = ctx.QueryParam(key)
	}

	results, err := h.svc.Export(shopID, languageCode, languageHeader)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	fileName := fmt.Sprintf("%s_productbarcode_%s.csv", shopID, time.Now().Format("20060102150405"))

	ctx.EchoContext().Response().Header().Set(echo.HeaderContentType, "application/octet-stream; charset=UTF-8")
	ctx.EchoContext().Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename=\""+fileName+"\"")
	ctx.EchoContext().Response().WriteHeader(http.StatusOK)

	t := transform.NewWriter(ctx.EchoContext().Response(), unicode.UTF8BOM.NewEncoder())

	csvWriter := csv.NewWriter(t)
	defer csvWriter.Flush()

	for _, value := range results {

		err := csvWriter.Write(value)
		if err != nil {
			log.Fatal("Error writing record to CSV:", err)
			return err
		}
	}

	return nil
}

// Get ProductBarcode BOM godoc
// @Description get product barcode bom view information
// @Tags		ProductBarcode
// @Param		barcode  path      string  true  "Barcode"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/barcode/bom/{barcode} [get]
func (h ProductBarcodeHttp) InfoBOMView(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	barcode := ctx.Param("barcode")

	doc, err := h.svc.InfoBomView(shopID, barcode)

	if err != nil {
		h.ms.Logger.Errorf("Error getting document %v: %v", barcode, err)
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    doc,
	})
	return nil
}

func (h ProductBarcodeHttp) searchFilter(queryParam func(string) string) map[string]interface{} {
	filters := requestfilter.GenerateFilters(queryParam, []requestfilter.FilterRequest{
		{
			Param: "isalacarte",
			Field: "isalacarte",
			Type:  requestfilter.FieldTypeBoolean,
		},
		{
			Param: "ordertypes",
			Field: "ordertypes.code",
			Type:  requestfilter.FieldTypeString,
		},
		{
			Param: "itemtype",
			Field: "itemtype",
			Type:  requestfilter.FieldTypeInt,
		},
		{
			Param: "businesstypecode",
			Field: "businesstypes.code",
			Type:  requestfilter.FieldTypeString,
		},
		{
			Param: "branchcode",
			Field: "branches.code",
			Type:  requestfilter.FieldTypeString,
		},
		{
			Param: "isusesubbarcodes",
			Field: "isusesubbarcodes",
			Type:  requestfilter.FieldTypeBoolean,
		},
	})

	if temp, ok := filters["branches.code"]; ok {
		if tempBson, ok := temp.(bson.M); ok {
			if tempIn, ok := tempBson["$in"]; ok {
				filters["ignorebranches.code"] = bson.M{
					"$nin": tempIn,
				}
			}
		} else {
			filters["ignorebranches.code"] = bson.M{
				"$ne": temp,
			}
		}
		delete(filters, "branches.code")
	}

	if queryParam("isbom") != "" {
		if queryParam("isbom") == "true" {
			filters["bom"] = bson.M{
				"$exists": true,
				"$ne":     []string{},
			}
		} else {
			filters["bom"] = bson.M{
				"$eq": []string{},
			}
		}
	}

	return filters
}

// Create ProductBarcode import godoc
// @Description Create ProductBarcode
// @Tags		ProductBarcode
// @Param		ProductBarcode  body      []models.ProductBarcode  true  "ProductBarcode"
// @Accept 		json
// @Success		201	{object}	common.BulkResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /product/barcode/import [post]
func (h ProductBarcodeHttp) Import(ctx microservice.IContext) error {

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

	bulkResponse, err := h.svc.Import(shopID, authUsername, dataReq)

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
