package mastersync

import (
	"fmt"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/config"
	"smlcloudplatform/pkg/member"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"strings"
	"time"

	"smlcloudplatform/pkg/restaurant/kitchen"
	"smlcloudplatform/pkg/restaurant/printer"
	"smlcloudplatform/pkg/restaurant/shoptable"
	"smlcloudplatform/pkg/restaurant/shopzone"

	"smlcloudplatform/pkg/mastersync/services"

	employeeRepo "smlcloudplatform/pkg/shop/employee/repositories"
	employeeService "smlcloudplatform/pkg/shop/employee/services"

	productRepo "smlcloudplatform/pkg/product/product/repositories"
	productService "smlcloudplatform/pkg/product/product/services"

	productcategoryRepo "smlcloudplatform/pkg/product/productcategory/repositories"
	productcategoryService "smlcloudplatform/pkg/product/productcategory/services"

	productbarcodeRepo "smlcloudplatform/pkg/product/productbarcode/repositories"
	productbarcodeService "smlcloudplatform/pkg/product/productbarcode/services"

	productunitRepo "smlcloudplatform/pkg/product/unit/repositories"
	productunitService "smlcloudplatform/pkg/product/unit/services"

	bankmasterRepo "smlcloudplatform/pkg/payment/bankmaster/repositories"
	bankmasterService "smlcloudplatform/pkg/payment/bankmaster/services"

	bookbankRepo "smlcloudplatform/pkg/payment/bookbank/repositories"
	bookbankService "smlcloudplatform/pkg/payment/bookbank/services"

	qrpaymentRepo "smlcloudplatform/pkg/payment/qrpayment/repositories"
	qrpaymentService "smlcloudplatform/pkg/payment/qrpayment/services"

	restaurantDeviceRepo "smlcloudplatform/pkg/restaurant/device/repositories"
	restaurantDeviceService "smlcloudplatform/pkg/restaurant/device/services"

	restaurantStaffRepo "smlcloudplatform/pkg/restaurant/staff/repositories"
	restaurantStaffService "smlcloudplatform/pkg/restaurant/staff/services"

	ordertype_repo "smlcloudplatform/pkg/product/ordertype/repositories"
	ordertype_service "smlcloudplatform/pkg/product/ordertype/services"

	"smlcloudplatform/pkg/mastersync/repositories"
)

type MasterSyncHttp struct {
	ms                    *microservice.Microservice
	cfg                   config.IConfig
	activityModuleManager *ActivityModuleManager

	svcMasterSync services.IMasterSyncService
	// svcProductBarcode productbarcodeService.ProductBarcodeHttpService
}

func NewMasterSyncHttp(ms *microservice.Microservice, cfg config.IConfig) MasterSyncHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	pstPg := ms.Persister(cfg.PersisterConfig())
	// prod := ms.Producer(cfg.MQConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	activityModuleManager := NewActivityModuleManager(pst)

	masterSyncCacheRepo := repositories.NewMasterSyncCacheRepository(cache)

	//############

	// pdt1 := productRepo.ProductRepository{}
	// pdt1.InitialActivityRepository(pst)

	// pdt2 := productService.ProductHttpService{}
	// pdt2.InitialActivityService(pst, &productRepo.ProductRepository{})

	// Product
	svcProduct := productService.NewProductHttpService(productRepo.NewProductRepository(pst), nil, masterSyncCacheRepo)
	activityModuleManager.Add(svcProduct)

	// Product Category
	svcProductCategory := productcategoryService.NewProductCategoryHttpService(productcategoryRepo.NewProductCategoryRepository(pst), masterSyncCacheRepo)
	activityModuleManager.Add(svcProductCategory)

	// Product Barcode
	svcProductBarcode := productbarcodeService.NewProductBarcodeHttpService(productbarcodeRepo.NewProductBarcodeRepository(pst, cache), nil, nil, masterSyncCacheRepo)
	activityModuleManager.Add(svcProductBarcode)

	// Product Unit
	svcProductUnit := productunitService.NewUnitHttpService(productunitRepo.NewUnitRepository(pst), productbarcodeRepo.NewProductBarcodeRepository(pst, cache), cfg.UnitServiceConfig(), masterSyncCacheRepo)
	activityModuleManager.Add(svcProductUnit)

	// Kitchen
	repoKitchen := kitchen.NewKitchenRepository(pst)
	svcKitchen := kitchen.NewKitchenService(repoKitchen, masterSyncCacheRepo)
	activityModuleManager.Add(svcKitchen)

	// Shop Printer
	repoShopPrinter := printer.NewPrinterRepository(pst)
	svcShopPrinter := printer.NewPrinterService(repoShopPrinter, masterSyncCacheRepo)
	activityModuleManager.Add(svcShopPrinter)

	// Shop Table
	repoShopTable := shoptable.NewShopTableRepository(pst)
	svcShopTable := shoptable.NewShopTableService(repoShopTable, masterSyncCacheRepo)
	activityModuleManager.Add(svcShopTable)

	// Shop Zone
	repoShopZone := shopzone.NewShopZoneRepository(pst)
	svcShopZone := shopzone.NewShopZoneService(repoShopZone, masterSyncCacheRepo)
	activityModuleManager.Add(svcShopZone)

	// device
	repoRestaurantDevice := restaurantDeviceRepo.NewDeviceRepository(pst)
	svcRestaurantDevice := restaurantDeviceService.NewDeviceHttpService(repoRestaurantDevice, masterSyncCacheRepo)
	activityModuleManager.Add(svcRestaurantDevice)

	// staff
	repoRestaurantStaff := restaurantStaffRepo.NewStaffRepository(pst)
	svcRestaurantStaff := restaurantStaffService.NewStaffHttpService(repoRestaurantStaff, masterSyncCacheRepo)
	activityModuleManager.Add(svcRestaurantStaff)

	// Member
	repoMember := member.NewMemberRepository(pst)
	pgRepoMember := member.NewMemberPGRepository(pstPg)
	svcMember := member.NewMemberService(repoMember, pgRepoMember, masterSyncCacheRepo)
	activityModuleManager.Add(svcMember)

	// Employee
	repoEmployee := employeeRepo.NewEmployeeRepository(pst)
	svcEmployee := employeeService.NewEmployeeHttpService(repoEmployee, masterSyncCacheRepo, utils.HashPassword)
	activityModuleManager.Add(svcEmployee)

	// Bank Master
	repoBankMaster := bankmasterRepo.NewBankMasterRepository(pst)
	svcBankMaster := bankmasterService.NewBankMasterHttpService(repoBankMaster, masterSyncCacheRepo)
	activityModuleManager.Add(svcBankMaster)

	// Book Bank
	repoBookBank := bookbankRepo.NewBookBankRepository(pst)
	svcBookBank := bookbankService.NewBookBankHttpService(repoBookBank, masterSyncCacheRepo)
	activityModuleManager.Add(svcBookBank)

	// Qr Payment
	qrpaymentRepo := qrpaymentRepo.NewQrPaymentRepository(pst)
	svcQrPayment := qrpaymentService.NewQrPaymentHttpService(qrpaymentRepo, masterSyncCacheRepo)
	activityModuleManager.Add(svcQrPayment)

	// Order type
	repoOrdertype := ordertype_repo.NewOrderTypeRepository(pst)
	svcOrdertype := ordertype_service.NewOrderTypeHttpService(repoOrdertype, masterSyncCacheRepo)
	activityModuleManager.Add(svcOrdertype)

	masterCacheSyncRepo := repositories.NewMasterSyncCacheRepository(cache)
	svcMasterSync := services.NewMasterSyncService(masterCacheSyncRepo)

	return MasterSyncHttp{
		ms:                    ms,
		cfg:                   cfg,
		activityModuleManager: activityModuleManager,

		svcMasterSync: svcMasterSync,
		// svcProductBarcode: *svcProductBarcode,
	}
}

func (h MasterSyncHttp) RouteSetup() {
	h.ms.GET("/master-sync", h.LastActivitySync)
	h.ms.GET("/master-sync/status", h.SyncStatus)
	h.ms.GET("/master-sync/list", h.LastActivitySyncOffset)
}

// List Master Sync Status godoc
// @Description  Master Sync Status
// @Tags		MasterSync
// @Success		200	{array}		interface{}
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /master-sync/status [get]
func (h MasterSyncHttp) SyncStatus(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	status, _ := h.svcMasterSync.GetStatus(shopID, h.activityModuleManager.GetModules())

	ctx.Response(
		http.StatusOK,
		status,
	)

	return nil
}

// List Master Sync godoc
// @Description  Master Sync
// @Tags		MasterSync
// @Param		lastupdate		query	string		false  "last update date ex: 2020-01-01T00:00:00"
// @Param		module		query	string		false  "module code ex: product,productcategory,productbarcode"
// @Param		action		query	string		false  "action code (all, new, remove)"
// @Param		filter		query	string		false  "filter data ex. filter=branch:1,department:x01"
// @Success		200	{array}		models.ApiResponse
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /master-sync [get]
func (h MasterSyncHttp) LastActivitySync(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	layout := "2006-01-02T15:04:05"
	lastUpdateStr := ctx.QueryParam("lastupdate")
	if len(lastUpdateStr) < 1 {
		lastUpdateStr = ctx.QueryParam("lastUpdate")
	}

	lastUpdateStr = strings.Trim(lastUpdateStr, " ")
	if len(lastUpdateStr) < 1 {
		ctx.ResponseError(400, "lastupdate format invalid.")
		return nil
	}

	lastUpdate, err := time.Parse(layout, lastUpdateStr)

	if err != nil {
		ctx.ResponseError(400, "lastUpdate format invalid.")
		return err
	}

	pageable := utils.GetPageable(ctx.QueryParam)

	moduleParam := strings.Trim(ctx.QueryParam("module"), " ")
	action := strings.Trim(ctx.QueryParam("action"), " ")
	filterParam := strings.Trim(ctx.QueryParam("filter"), " ")

	requestModuleSelectList := []string{}
	moduleSelectList := map[string]struct{}{}

	if moduleParam != "" {
		requestModuleSelectList = strings.Split(moduleParam, ",")
		for _, module := range requestModuleSelectList {
			module = strings.ToLower(module)
			moduleSelectList[module] = struct{}{}
		}
	}

	if len(requestModuleSelectList) > 0 && strings.ToLower(requestModuleSelectList[0]) == "all" {
		moduleSelectList = map[string]struct{}{}
	}

	results, pagination, err := listDataModulePage(h.activityModuleManager.GetList(), moduleSelectList, ActivityParamPage{
		ShopID:     shopID,
		Action:     action,
		LastUpdate: lastUpdate,
		Filters:    filterParam,
		Pageable:   pageable,
	})

	if err != nil {
		fmt.Println(err)
	}

	ctx.Response(
		http.StatusOK,
		models.ApiResponse{
			Success:    true,
			Data:       results,
			Pagination: pagination,
		})

	return nil
}

// List Master Sync Offset godoc
// @Description  Master Sync Offset
// @Tags		MasterSync
// @Param		lastupdate		query	string		false  "last update date ex: 2020-01-01T00:00:00"
// @Param		module		query	string		false  "module code ex: product,productcategory,productbarcode"
// @Param		action		query	string		false  "action code (all, new, remove)"
// @Success		200	{array}		models.ApiResponse
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /master-sync/list [get]
func (h MasterSyncHttp) LastActivitySyncOffset(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	layout := "2006-01-02T15:04:05"
	lastUpdateStr := ctx.QueryParam("lastupdate")
	if len(lastUpdateStr) < 1 {
		lastUpdateStr = ctx.QueryParam("lastUpdate")
	}

	lastUpdateStr = strings.Trim(lastUpdateStr, " ")
	if len(lastUpdateStr) < 1 {
		ctx.ResponseError(400, "lastupdate format invalid.")
		return nil
	}

	lastUpdate, err := time.Parse(layout, lastUpdateStr)

	if err != nil {
		ctx.ResponseError(400, "lastUpdate format invalid.")
		return err
	}

	pageableStep := utils.GetPageableStep(ctx.QueryParam)

	moduleParam := strings.Trim(ctx.QueryParam("module"), " ")
	action := strings.Trim(ctx.QueryParam("action"), " ")
	filterParam := strings.Trim(ctx.QueryParam("filter"), " ")

	requestModuleSelectList := []string{}
	moduleSelectList := map[string]struct{}{}

	if moduleParam != "" {
		requestModuleSelectList = strings.Split(moduleParam, ",")
		for _, module := range requestModuleSelectList {
			module = strings.ToLower(module)
			moduleSelectList[module] = struct{}{}
		}
	}

	if len(requestModuleSelectList) > 0 && strings.ToLower(requestModuleSelectList[0]) == "all" {
		moduleSelectList = map[string]struct{}{}
	}

	results, err := listDataModuleOffset(h.activityModuleManager.GetList(), moduleSelectList, ActivityParamOffset{
		ShopID:       shopID,
		Action:       action,
		LastUpdate:   lastUpdate,
		Filters:      filterParam,
		PageableStep: pageableStep,
	})

	if err != nil {
		fmt.Println(err)
	}

	ctx.Response(
		http.StatusOK,
		models.ApiResponse{
			Success: true,
			Data:    results,
		})

	return nil
}
