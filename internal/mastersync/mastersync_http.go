package mastersync

import (
	"fmt"
	"net/http"
	"smlaicloudplatform/internal/config"
	"smlaicloudplatform/internal/models"
	"smlaicloudplatform/internal/utils"
	"smlaicloudplatform/pkg/microservice"
	"strings"
	"time"

	"smlaicloudplatform/internal/restaurant/kitchen"
	"smlaicloudplatform/internal/restaurant/printer"
	"smlaicloudplatform/internal/restaurant/table"
	"smlaicloudplatform/internal/restaurant/zone"

	"smlaicloudplatform/internal/mastersync/services"

	employeeRepo "smlaicloudplatform/internal/shop/employee/repositories"
	employeeService "smlaicloudplatform/internal/shop/employee/services"

	productcategoryRepo "smlaicloudplatform/internal/product/productcategory/repositories"
	productcategoryService "smlaicloudplatform/internal/product/productcategory/services"

	productbarcodeRepo "smlaicloudplatform/internal/product/productbarcode/repositories"
	productbarcodeService "smlaicloudplatform/internal/product/productbarcode/services"

	bankmasterRepo "smlaicloudplatform/internal/payment/bankmaster/repositories"
	bankmasterService "smlaicloudplatform/internal/payment/bankmaster/services"

	bookbankRepo "smlaicloudplatform/internal/payment/bookbank/repositories"
	bookbankService "smlaicloudplatform/internal/payment/bookbank/services"

	qrpaymentRepo "smlaicloudplatform/internal/payment/qrpayment/repositories"
	qrpaymentService "smlaicloudplatform/internal/payment/qrpayment/services"

	restaurantDeviceRepo "smlaicloudplatform/internal/restaurant/device/repositories"
	restaurantDeviceService "smlaicloudplatform/internal/restaurant/device/services"

	restaurantStaffRepo "smlaicloudplatform/internal/restaurant/staff/repositories"
	restaurantStaffService "smlaicloudplatform/internal/restaurant/staff/services"

	ordertype_repo "smlaicloudplatform/internal/product/ordertype/repositories"
	ordertype_service "smlaicloudplatform/internal/product/ordertype/services"

	"smlaicloudplatform/internal/mastersync/repositories"
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
	// pstPg := ms.Persister(cfg.PersisterConfig())
	// prod := ms.Producer(cfg.MQConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	activityModuleManager := NewActivityModuleManager(pst)

	masterSyncCacheRepo := repositories.NewMasterSyncCacheRepository(cache)

	//############

	// pdt1 := productRepo.ProductRepository{}
	// pdt1.InitialActivityRepository(pst)

	// pdt2 := productService.ProductHttpService{}
	// pdt2.InitialActivityService(pst, &productRepo.ProductRepository{})

	// Product Category
	svcProductCategory := productcategoryService.NewProductCategoryHttpService(productcategoryRepo.NewProductCategoryRepository(pst), masterSyncCacheRepo)
	activityModuleManager.Add(svcProductCategory)

	// Product Barcode
	repoProductBarcode := productbarcodeRepo.NewProductBarcodeRepository(pst, cache)
	svcProductBarcode := productbarcodeService.NewProductBarcodeHttpService(repoProductBarcode, nil, nil, nil, masterSyncCacheRepo)
	activityModuleManager.Add(svcProductBarcode)

	// Product Unit
	// svcProductUnit := productunitService.NewUnitHttpService(productunitRepo.NewUnitPGRepository(pstPg))
	// activityModuleManager.Add(svcProductUnit)

	// Kitchen
	repoKitchen := kitchen.NewKitchenRepository(pst)
	svcKitchen := kitchen.NewKitchenService(repoKitchen, masterSyncCacheRepo)
	activityModuleManager.Add(svcKitchen)

	// Shop Printer
	repoShopPrinter := printer.NewPrinterRepository(pst)
	svcShopPrinter := printer.NewPrinterService(repoShopPrinter, masterSyncCacheRepo)
	activityModuleManager.Add(svcShopPrinter)

	// Shop Table
	repoShopTable := table.NewTableRepository(pst)
	svcShopTable := table.NewTableService(repoShopTable, masterSyncCacheRepo)
	activityModuleManager.Add(svcShopTable)

	// Shop Zone
	repoShopZone := zone.NewZoneRepository(pst)
	svcShopZone := zone.NewZoneService(repoShopZone, masterSyncCacheRepo)
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
	// repoMember := member.NewMemberRepository(pst)
	// pgRepoMember := member.NewMemberPGRepository(pstPg)
	// svcMember := member.NewMemberService(repoMember, pgRepoMember, nil, nil, masterSyncCacheRepo)
	// activityModuleManager.Add(svcMember)

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
	svcOrdertype := ordertype_service.NewOrderTypeHttpService(repoOrdertype, nil, repoProductBarcode, masterSyncCacheRepo)
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

func (h MasterSyncHttp) RegisterHttp() {
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
