package eorder

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/config"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	repositorycategory "smlcloudplatform/pkg/product/productcategory/repositories"
	servicecategory "smlcloudplatform/pkg/product/productcategory/services"

	repositoryproduct "smlcloudplatform/pkg/product/productbarcode/repositories"
	serviceproduct "smlcloudplatform/pkg/product/productbarcode/services"
	"smlcloudplatform/pkg/utils"
)

type IEOrderHttp interface{}

type EOrderHttp struct {
	ms          *microservice.Microservice
	cfg         config.IConfig
	svcCategory servicecategory.IProductCategoryHttpService
	svcProduct  serviceproduct.IProductBarcodeHttpService
}

func NewEOrderHttp(ms *microservice.Microservice, cfg config.IConfig) EOrderHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	pstClickHouse := ms.ClickHousePersister(cfg.ClickHouseConfig())
	cache := ms.Cacher(cfg.CacherConfig())
	prod := ms.Producer(cfg.MQConfig())

	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)

	repoCategory := repositorycategory.NewProductCategoryRepository(pst)
	svcCategory := servicecategory.NewProductCategoryHttpService(repoCategory, masterSyncCacheRepo)

	repo := repositoryproduct.NewProductBarcodeRepository(pst, cache)
	clickHouseRepo := repositoryproduct.NewProductBarcodeClickhouseRepository(pstClickHouse)
	mqRepo := repositoryproduct.NewProductBarcodeMessageQueueRepository(prod)

	svcProduct := serviceproduct.NewProductBarcodeHttpService(repo, mqRepo, clickHouseRepo, masterSyncCacheRepo)

	return EOrderHttp{
		ms:          ms,
		cfg:         cfg,
		svcCategory: svcCategory,
		svcProduct:  svcProduct,
	}
}

func (h EOrderHttp) RouteSetup() {

	h.ms.GET("/e-order/category", h.SearchProductCategoryPage)
	h.ms.GET("/e-order/product", h.SearchProductBarcodePage)
	h.ms.GET("/e-order/product-barcode", h.SearchProductBarcodePage)

}

// List Product Category
// @Description List Product Category
// @Tags		E-Order
// @Param		shopid		query	string		false  "Shop ID"
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Limit"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /e-order/category [get]
func (h EOrderHttp) SearchProductCategoryPage(ctx microservice.IContext) error {
	shopID := ctx.QueryParam("shopid")

	if len(shopID) == 0 {
		ctx.ResponseError(http.StatusBadRequest, "shopid is empty")
		return nil
	}

	pageable := utils.GetPageable(ctx.QueryParam)
	docList, pagination, err := h.svcCategory.SearchProductCategory(shopID, pageable)

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

// List Product
// @Description List Product
// @Tags		E-Order
// @Param		shopid		query	string		false  "Shop ID"
// @Param		isalacarte		query	string		false  "is A La Carte"
// @Param		ordertypes		query	string		false  "order types ex. a01,a02"
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Limit"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /e-order/product [get]
func (h EOrderHttp) SearchProductBarcodePage(ctx microservice.IContext) error {
	shopID := ctx.QueryParam("shopid")

	if len(shopID) == 0 {
		ctx.ResponseError(http.StatusBadRequest, "shopid is empty")
		return nil
	}

	filters := utils.GetFilters(ctx.QueryParam, []utils.FilterRequest{
		{
			Param: "isalacarte",
			Field: "isalacarte",
			Type:  "boolean",
		},
		{
			Param: "ordertypes",
			Field: "ordertypes.code",
			Type:  "array",
		},
	})

	pageable := utils.GetPageable(ctx.QueryParam)
	docList, pagination, err := h.svcProduct.SearchProductBarcode(shopID, filters, pageable)

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

// List Product By Barcodes
// @Description List Product By Barcodes
// @Tags		E-Order
// @Param		shopid		query	string		false  "Shop ID"
// @Param		barcodes		query	string		false  "barcode json array"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /e-order/product-barcode [get]
func (h EOrderHttp) GetProductBarcodeByBarcodes(ctx microservice.IContext) error {
	shopID := ctx.QueryParam("shopid")

	if len(shopID) == 0 {
		ctx.ResponseError(http.StatusBadRequest, "shopid is empty")
		return nil
	}

	rawBarcodes := ctx.QueryParam("barcodes")

	barcodes := []string{}
	if len(rawBarcodes) > 0 {
		json.Unmarshal([]byte(rawBarcodes), &barcodes)
	}

	docList, err := h.svcProduct.GetProductBarcodeByBarcodes(shopID, barcodes)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    docList,
	})
	return nil
}
