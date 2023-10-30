package eorder

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/config"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/eorder/services"
	repositorycategory "smlcloudplatform/pkg/product/productcategory/repositories"
	servicecategory "smlcloudplatform/pkg/product/productcategory/services"
	"smlcloudplatform/pkg/utils/requestfilter"

	salechannel_repo "smlcloudplatform/pkg/channel/salechannel/repositories"
	repo_order_device "smlcloudplatform/pkg/order/device/repositories"
	repo_order_setting "smlcloudplatform/pkg/order/setting/repositories"
	branch_repo "smlcloudplatform/pkg/organization/branch/repositories"
	repo_media "smlcloudplatform/pkg/pos/media/repositories"
	repo_product "smlcloudplatform/pkg/product/productbarcode/repositories"
	serviceproduct "smlcloudplatform/pkg/product/productbarcode/services"
	"smlcloudplatform/pkg/restaurant/kitchen"
	"smlcloudplatform/pkg/restaurant/table"
	"smlcloudplatform/pkg/shop"
	"smlcloudplatform/pkg/utils"
)

type IEOrderHttp interface{}

type EOrderHttp struct {
	ms          *microservice.Microservice
	cfg         config.IConfig
	svcCategory servicecategory.IProductCategoryHttpService
	svcProduct  serviceproduct.IProductBarcodeHttpService
	svcEOrder   services.EOrderService
}

func NewEOrderHttp(ms *microservice.Microservice, cfg config.IConfig) EOrderHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	pstClickHouse := ms.ClickHousePersister(cfg.ClickHouseConfig())
	cache := ms.Cacher(cfg.CacherConfig())
	prod := ms.Producer(cfg.MQConfig())

	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)

	repoCategory := repositorycategory.NewProductCategoryRepository(pst)
	svcCategory := servicecategory.NewProductCategoryHttpService(repoCategory, masterSyncCacheRepo)

	repo := repo_product.NewProductBarcodeRepository(pst, cache)
	clickHouseRepo := repo_product.NewProductBarcodeClickhouseRepository(pstClickHouse)
	mqRepo := repo_product.NewProductBarcodeMessageQueueRepository(prod)

	svcProduct := serviceproduct.NewProductBarcodeHttpService(repo, mqRepo, clickHouseRepo, nil, masterSyncCacheRepo)

	repoShop := shop.NewShopRepository(pst)
	repoTable := table.NewTableRepository(pst)
	repoOrder := repo_order_setting.NewSettingRepository(pst)
	repoDevice := repo_order_device.NewDeviceRepository(pst)
	repoMedia := repo_media.NewMediaRepository(pst)
	repoKitchen := kitchen.NewKitchenRepository(pst)
	repoSaleChannel := salechannel_repo.NewSaleChannelRepository(pst)
	repoBranch := branch_repo.NewBranchRepository(pst)

	svcEOrder := services.NewEOrderService(repoShop, repoTable, repoOrder, repoMedia, repoKitchen, repoDevice, repoSaleChannel, repoBranch)

	return EOrderHttp{
		ms:          ms,
		cfg:         cfg,
		svcCategory: svcCategory,
		svcProduct:  svcProduct,
		svcEOrder:   svcEOrder,
	}
}

func (h EOrderHttp) RegisterHttp() {

	h.ms.GET("/e-order/category", h.SearchProductCategoryPage)
	h.ms.GET("/e-order/product", h.SearchProductBarcodePage)
	h.ms.GET("/e-order/product-barcode", h.SearchProductBarcodePage)
	h.ms.GET("/e-order/shop-info/v1.1", h.ShopInfo)
	h.ms.GET("/e-order/shop-info", h.ShopInfoOld)

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

	filters := requestfilter.GenerateFilters(ctx.QueryParam, []requestfilter.FilterRequest{
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

// Get Shop Info
// @Description Get Shop Info
// @Tags		E-Order
// @Param		shopid		query	string		false  "Shop ID"
// @Param		order-station		query	string		false  "Order station code"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /e-order/shop-info [get]
func (h EOrderHttp) ShopInfoOld(ctx microservice.IContext) error {
	shopID := ctx.QueryParam("shopid")
	orderStationCode := ctx.QueryParam("order-station")

	if len(shopID) == 0 {
		ctx.ResponseError(http.StatusBadRequest, "shopid is empty")
		return nil
	}

	data, err := h.svcEOrder.GetShopInfoOld(shopID, orderStationCode)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    data,
	})
	return nil
}

// Get Shop Info v1.1
// @Description Get Shop Info v1.1
// @Tags		E-Order
// @Param		shopid		query	string		false  "Shop ID"
// @Param		order-station		query	string		false  "Order station code"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Security     AccessToken
// @Router /e-order/shop-info/v1.1 [get]
func (h EOrderHttp) ShopInfo(ctx microservice.IContext) error {
	shopID := ctx.QueryParam("shopid")
	orderStationCode := ctx.QueryParam("order-station")

	if len(shopID) == 0 {
		ctx.ResponseError(http.StatusBadRequest, "shopid is empty")
		return nil
	}

	data, err := h.svcEOrder.GetShopInfo(shopID, orderStationCode)

	if err != nil {
		ctx.ResponseError(http.StatusBadRequest, err.Error())
		return err
	}

	ctx.Response(http.StatusOK, common.ApiResponse{
		Success: true,
		Data:    data,
	})
	return nil
}
