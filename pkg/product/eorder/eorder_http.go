package eorder

import (
	"encoding/json"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/config"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/eorder/services"
	category_repositories "smlcloudplatform/pkg/product/productcategory/repositories"
	category_services "smlcloudplatform/pkg/product/productcategory/services"
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
	"smlcloudplatform/pkg/restaurant/zone"
	"smlcloudplatform/pkg/shop"
	saleinvoice_repositories "smlcloudplatform/pkg/transaction/saleinvoice/repositories"
	saleinvoice_services "smlcloudplatform/pkg/transaction/saleinvoice/services"
	"smlcloudplatform/pkg/utils"
)

type IEOrderHttp interface{}

type EOrderHttp struct {
	ms             *microservice.Microservice
	cfg            config.IConfig
	svcCategory    category_services.IProductCategoryHttpService
	svcProduct     serviceproduct.IProductBarcodeHttpService
	svcEOrder      services.EOrderService
	svcZone        zone.IZoneService
	svcTable       table.TableService
	svcKitchen     kitchen.IKitchenService
	svcSaleInvoice saleinvoice_services.ISaleInvoiceHttpService
}

func NewEOrderHttp(ms *microservice.Microservice, cfg config.IConfig) EOrderHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	pstClickHouse := ms.ClickHousePersister(cfg.ClickHouseConfig())
	cache := ms.Cacher(cfg.CacherConfig())
	prod := ms.Producer(cfg.MQConfig())

	masterSyncCacheRepo := mastersync.NewMasterSyncCacheRepository(cache)

	repoCategory := category_repositories.NewProductCategoryRepository(pst)
	svcCategory := category_services.NewProductCategoryHttpService(repoCategory, masterSyncCacheRepo)

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

	repoZone := zone.NewZoneRepository(pst)

	svcZone := zone.NewZoneService(repoZone, masterSyncCacheRepo)
	svcEOrder := services.NewEOrderService(repoShop, repoTable, repoOrder, repoMedia, repoKitchen, repoDevice, repoSaleChannel, repoBranch)

	repoSaleInvoice := saleinvoice_repositories.NewSaleInvoiceRepository(pst)
	svcSaleInvoice := saleinvoice_services.NewSaleInvoiceHttpService(repoSaleInvoice, nil, nil, nil)

	svcTable := table.NewTableService(repoTable, masterSyncCacheRepo)
	svcKitchen := kitchen.NewKitchenService(repoKitchen, masterSyncCacheRepo)

	return EOrderHttp{
		ms:             ms,
		cfg:            cfg,
		svcCategory:    svcCategory,
		svcProduct:     svcProduct,
		svcEOrder:      svcEOrder,
		svcZone:        svcZone,
		svcTable:       *svcTable,
		svcKitchen:     svcKitchen,
		svcSaleInvoice: svcSaleInvoice,
	}
}

func (h EOrderHttp) RegisterHttp() {

	h.ms.GET("/e-order/category", h.SearchProductCategoryPage)
	// h.ms.GET("/e-order/product", h.SearchProductBarcodePage)
	h.ms.GET("/e-order/product-barcode", h.SearchProductBarcodePage)
	h.ms.GET("/e-order/shop-info/v1.1", h.ShopInfo)
	h.ms.GET("/e-order/shop-info", h.ShopInfoOld)
	h.ms.GET("/e-order/restaurant/zone", h.SearchZone)
	h.ms.GET("/e-order/restaurant/kitchen", h.SearchKitchen)
	h.ms.GET("/e-order/restaurant/table", h.SearchTable)
	h.ms.GET("/e-order/sale-invoice/last-pos-docno", h.GetLastPOSDocNo)

}

// List Product Category
// @Description List Product Category
// @Tags		E-Order
// @Param		shopid		query	string		false  "Shop ID"
// @Param		q		query	string		false  "Search Value"
// @Param		group-number		query	int		false  "group number"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Limit"
// @Accept 		json
// @Success		200	{array}		common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Router /e-order/category [get]
func (h EOrderHttp) SearchProductCategoryPage(ctx microservice.IContext) error {
	shopID := ctx.QueryParam("shopid")

	if len(shopID) == 0 {
		ctx.ResponseError(http.StatusBadRequest, "shopid is empty")
		return nil
	}

	pageable := utils.GetPageable(ctx.QueryParam)
	filters := requestfilter.GenerateFilters(ctx.QueryParam, []requestfilter.FilterRequest{
		{
			Param: "group-number",
			Field: "groupnumber",
			Type:  requestfilter.FieldTypeInt,
		},
	})
	docList, pagination, err := h.svcCategory.SearchProductCategory(shopID, filters, pageable)

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

// List E Order Restaurant Zone godoc
// @Description List Restaurant Zone Category
// @Tags		E-Order
// @Param		group-number	query	integer		false  "Group Number"
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Accept 		json
// @Success		200	{object}	models.ZonePageResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Router /e-order/restaurant/zone [get]
func (h EOrderHttp) SearchZone(ctx microservice.IContext) error {
	shopID := ctx.QueryParam("shopid")

	if len(shopID) == 0 {
		ctx.ResponseError(http.StatusBadRequest, "shopid is empty")
		return nil
	}

	pageable := utils.GetPageable(ctx.QueryParam)
	filters := requestfilter.GenerateFilters(ctx.QueryParam, []requestfilter.FilterRequest{
		{
			Param: "group-number",
			Field: "groupnumber",
			Type:  requestfilter.FieldTypeInt,
		},
	})

	docList, pagination, err := h.svcZone.SearchZone(shopID, filters, pageable)

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

// List E Order Restaurant Kitchen godoc
// @Description List Restaurant Kitchen Category
// @Tags		E-Order
// @Param		group-number	query	integer		false  "Group Number"
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Accept 		json
// @Success		200	{object}	models.KitchenPageResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Router /e-order/restaurant/kitchen [get]
func (h EOrderHttp) SearchKitchen(ctx microservice.IContext) error {
	shopID := ctx.QueryParam("shopid")

	if len(shopID) == 0 {
		ctx.ResponseError(http.StatusBadRequest, "shopid is empty")
		return nil
	}

	pageable := utils.GetPageable(ctx.QueryParam)

	filters := requestfilter.GenerateFilters(ctx.QueryParam, []requestfilter.FilterRequest{
		{
			Param: "group-number",
			Field: "groupnumber",
			Type:  requestfilter.FieldTypeInt,
		},
	})

	docList, pagination, err := h.svcKitchen.SearchKitchen(shopID, filters, pageable)

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

// List E Order Restaurant  Table godoc
// @Description List Restaurant  Table Category
// @Tags		E-Order
// @Param		group-number	query	integer		false  "Group Number"
// @Param		q		query	string		false  "Search Value"
// @Param		page	query	integer		false  "Page"
// @Param		limit	query	integer		false  "Size"
// @Accept 		json
// @Success		200	{object}	models.TablePageResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Router /e-order/restaurant/table [get]
func (h EOrderHttp) SearchTable(ctx microservice.IContext) error {
	shopID := ctx.QueryParam("shopid")

	if len(shopID) == 0 {
		ctx.ResponseError(http.StatusBadRequest, "shopid is empty")
		return nil
	}

	pageable := utils.GetPageable(ctx.QueryParam)

	filters := requestfilter.GenerateFilters(ctx.QueryParam, []requestfilter.FilterRequest{
		{
			Param: "group-number",
			Field: "groupnumber",
			Type:  requestfilter.FieldTypeInt,
		},
	})

	docList, pagination, err := h.svcTable.SearchTable(shopID, filters, pageable)

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

// Get E Order SaleInvoice Last DocNo godoc
// @Description get SaleInvoice Last DocNo
// @Tags		E-Order
// @Param		posid	query	string		false  "POS ID"
// @Param		maxdocno	query	string		false  "Max DocNo"
// @Accept 		json
// @Success		200	{object}	common.ApiResponse
// @Failure		401 {object}	common.AuthResponseFailed
// @Router /e-order/sale-invoice/last-pos-docno [get]
func (h EOrderHttp) GetLastPOSDocNo(ctx microservice.IContext) error {
	shopID := ctx.QueryParam("shopid")
	if len(shopID) == 0 {
		ctx.ResponseError(http.StatusBadRequest, "shopid is empty")
		return nil
	}

	posID := ctx.QueryParam("posid")
	maxDocNo := ctx.QueryParam("maxdocno")

	if posID == "" || maxDocNo == "" {
		ctx.ResponseError(http.StatusBadRequest, "posid and maxdocno is required")
		return nil
	}

	doc, err := h.svcSaleInvoice.GetLastPOSDocNo(shopID, posID, maxDocNo)

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
