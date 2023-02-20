package mastersync

import (
	"fmt"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/member"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/shop/employee"
	"smlcloudplatform/pkg/utils"
	"strings"
	"time"

	"smlcloudplatform/pkg/restaurant/kitchen"
	"smlcloudplatform/pkg/restaurant/printer"
	"smlcloudplatform/pkg/restaurant/shoptable"
	"smlcloudplatform/pkg/restaurant/shopzone"

	"smlcloudplatform/pkg/mastersync/services"

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

	"smlcloudplatform/pkg/mastersync/repositories"

	"github.com/userplant/mongopagination"

	micromodels "smlcloudplatform/internal/microservice/models"
)

type MasterSyncHttp struct {
	ms                    *microservice.Microservice
	cfg                   microservice.IConfig
	activityModuleManager *ActivityModuleManager

	svcMasterSync services.IMasterSyncService
	// svcProductBarcode productbarcodeService.ProductBarcodeHttpService
}

func NewMasterSyncHttp(ms *microservice.Microservice, cfg microservice.IConfig) MasterSyncHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	pstPg := ms.Persister(cfg.PersisterConfig())
	// prod := ms.Producer(cfg.MQConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	activityModuleManager := NewActivityModuleManager()

	masterSyncCacheRepo := repositories.NewMasterSyncCacheRepository(cache)

	//############

	// Product
	svcProduct := productService.NewProductHttpService(productRepo.NewProductRepository(pst), nil, masterSyncCacheRepo)
	activityModuleManager.Add(svcProduct)

	// Product Category
	svcProductCategory := productcategoryService.NewProductCategoryHttpService(productcategoryRepo.NewProductCategoryRepository(pst), masterSyncCacheRepo)
	activityModuleManager.Add(svcProductCategory)

	// Product Barcode
	svcProductBarcode := productbarcodeService.NewProductBarcodeHttpService(productbarcodeRepo.NewProductBarcodeRepository(pst), nil, masterSyncCacheRepo)
	activityModuleManager.Add(svcProductBarcode)

	// Product Unit
	svcProductUnit := productunitService.NewUnitHttpService(productunitRepo.NewUnitRepository(pst), masterSyncCacheRepo)
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
	repoEmployee := employee.NewEmployeeRepository(pst)
	svcEmployee := employee.NewEmployeeService(repoEmployee, masterSyncCacheRepo)
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
// @Success		200	{array}		models.ApiResponse
// @Failure		401 {object}	models.AuthResponseFailed
// @Security     AccessToken
// @Router /master-sync [get]
func (h MasterSyncHttp) LastActivitySync(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	layout := "2006-01-02T15:04:05"
	lastUpdateStr := ctx.QueryParam("lastupdate")

	lastUpdateStr = strings.Trim(lastUpdateStr, " ")
	if len(lastUpdateStr) < 1 {
		ctx.ResponseError(400, "lastupdate format invalid.")
		return nil
	}

	lastUpdateStr = strings.Trim(lastUpdateStr, " ")
	if len(lastUpdateStr) < 1 {
		ctx.ResponseError(400, "lastUpdate format invalid.")
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

type ActivityModuleManager struct {
	activityModuleList map[string]ActivityModule
}

func NewActivityModuleManager() *ActivityModuleManager {
	return &ActivityModuleManager{
		activityModuleList: map[string]ActivityModule{},
	}
}

func (m *ActivityModuleManager) Add(activityModule ActivityModule) *ActivityModuleManager {
	m.activityModuleList[activityModule.GetModuleName()] = activityModule
	return m
}

func (m ActivityModuleManager) GetList() map[string]ActivityModule {
	return m.activityModuleList
}

func (m ActivityModuleManager) GetModules() []string {
	modules := []string{}
	for module := range m.activityModuleList {
		modules = append(modules, module)
	}
	return modules
}

func (m ActivityModuleManager) GetPage(moduleSelectList map[string]struct{}, activityParam ActivityParamPage) (map[string]interface{}, mongopagination.PaginationData, error) {
	moduleList := map[string]ActivityModule{}

	for _, activityModule := range m.activityModuleList {
		moduleList[activityModule.GetModuleName()] = activityModule
	}

	return listDataModulePage(moduleList, moduleSelectList, activityParam)
}

type ActivityModule interface {
	LastActivity(string, string, time.Time, micromodels.Pageable) (models.LastActivity, mongopagination.PaginationData, error)
	LastActivityStep(string, string, time.Time, micromodels.PageableStep) (models.LastActivity, error)
	GetModuleName() string
}

type ActivityParamPage struct {
	ShopID     string
	Action     string
	LastUpdate time.Time
	Pageable   micromodels.Pageable
}

type ActivityParamOffset struct {
	ShopID       string
	Action       string
	LastUpdate   time.Time
	PageableStep micromodels.PageableStep
}

func listDataModulePage(appModules map[string]ActivityModule, moduleSelectList map[string]struct{}, param ActivityParamPage) (map[string]interface{}, mongopagination.PaginationData, error) {

	result := map[string]interface{}{}

	resultPagination := mongopagination.PaginationData{}
	for moduleName, appModule := range appModules {
		if len(moduleSelectList) == 0 || isSelectModule(moduleSelectList, moduleName) {
			docList, pagination, err := appModule.LastActivity(param.ShopID, param.Action, param.LastUpdate, param.Pageable)

			if err != nil {
				return map[string]interface{}{}, mongopagination.PaginationData{}, err
			}

			result[moduleName] = docList

			if pagination.Total > resultPagination.Total {
				resultPagination = pagination
			}
		}
	}

	return result, resultPagination, nil
}

func listDataModuleOffset(appModules map[string]ActivityModule, moduleSelectList map[string]struct{}, param ActivityParamOffset) (map[string]interface{}, error) {

	result := map[string]interface{}{}

	for moduleName, appModule := range appModules {
		if len(moduleSelectList) == 0 || isSelectModule(moduleSelectList, moduleName) {

			docList, err := appModule.LastActivityStep(param.ShopID, param.Action, param.LastUpdate, param.PageableStep)

			if err != nil {
				return map[string]interface{}{}, err
			}

			result[moduleName] = docList

		}
	}

	return result, nil
}

func isSelectModule(moduleList map[string]struct{}, moduleKey string) bool {
	_, ok := moduleList[moduleKey]
	return ok
}
