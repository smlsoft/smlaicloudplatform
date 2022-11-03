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
	"smlcloudplatform/pkg/restaurant/shopprinter"
	"smlcloudplatform/pkg/restaurant/shoptable"
	"smlcloudplatform/pkg/restaurant/shopzone"

	"smlcloudplatform/pkg/mastersync/services"
	inventoryRepo "smlcloudplatform/pkg/product/inventory/repositories"
	inventoryService "smlcloudplatform/pkg/product/inventory/services"

	productcategoryRepo "smlcloudplatform/pkg/product/productcategory/repositories"
	productcategoryService "smlcloudplatform/pkg/product/productcategory/services"

	productbarcodeRepo "smlcloudplatform/pkg/product/productbarcode/repositories"
	productbarcodeService "smlcloudplatform/pkg/product/productbarcode/services"

	productunitRepo "smlcloudplatform/pkg/product/unit/repositories"
	productunitService "smlcloudplatform/pkg/product/unit/services"

	"smlcloudplatform/pkg/mastersync/repositories"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type MasterSyncHttp struct {
	ms                    *microservice.Microservice
	cfg                   microservice.IConfig
	activityModuleManager *ActivityModuleManager

	svcMasterSync  services.IMasterSyncService
	svcMember      member.IMemberService
	svcInventory   inventoryService.IInventoryService
	svcKitchen     kitchen.IKitchenService
	svcShopPrinter shopprinter.IShopPrinterService
	svcShopTable   shoptable.ShopTableService
	svcShopZone    shopzone.ShopZoneService
	svcEmployee    employee.EmployeeService
	// svcProductBarcode productbarcodeService.ProductBarcodeHttpService
}

func NewMasterSyncHttp(ms *microservice.Microservice, cfg microservice.IConfig) MasterSyncHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	pstPg := ms.Persister(cfg.PersisterConfig())
	prod := ms.Producer(cfg.MQConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	activityModuleManager := NewActivityModuleManager()

	masterSyncCacheRepo := repositories.NewMasterSyncCacheRepository(cache)

	// Category

	// Member
	repoMember := member.NewMemberRepository(pst)
	pgRepoMember := member.NewMemberPGRepository(pstPg)
	repoCacheSyncMember := repositories.NewMasterSyncCacheRepository(cache)
	svcMember := member.NewMemberService(repoMember, pgRepoMember, repoCacheSyncMember)

	// Inventory
	repoInv := inventoryRepo.NewInventoryRepository(pst)
	mqRepoInv := inventoryRepo.NewInventoryMQRepository(prod)
	invCacheSyncRepo := repositories.NewMasterSyncCacheRepository(cache)
	svcInventory := inventoryService.NewInventoryService(repoInv, mqRepoInv, invCacheSyncRepo)

	// Kitchen
	repoKitchen := kitchen.NewKitchenRepository(pst)
	kitchenCacheSyncRepo := repositories.NewMasterSyncCacheRepository(cache)
	svcKitchen := kitchen.NewKitchenService(repoKitchen, kitchenCacheSyncRepo)

	// Shop Printer
	repoShopPrinter := shopprinter.NewShopPrinterRepository(pst)
	shopPrinterCacheSyncRepo := repositories.NewMasterSyncCacheRepository(cache)
	svcShopPrinter := shopprinter.NewShopPrinterService(repoShopPrinter, shopPrinterCacheSyncRepo)

	// Shop Table
	repoShopTable := shoptable.NewShopTableRepository(pst)
	shopTableCacheSyncRepo := repositories.NewMasterSyncCacheRepository(cache)
	svcShopTable := shoptable.NewShopTableService(repoShopTable, shopTableCacheSyncRepo)

	// Shop Zone
	repoShopZone := shopzone.NewShopZoneRepository(pst)
	shopZoneCacheSyncRepo := repositories.NewMasterSyncCacheRepository(cache)
	svcShopZone := shopzone.NewShopZoneService(repoShopZone, shopZoneCacheSyncRepo)

	// Employee
	repoEmployee := employee.NewEmployeeRepository(pst)
	employeeCacheSyncRepo := repositories.NewMasterSyncCacheRepository(cache)
	svcEmployee := employee.NewEmployeeService(repoEmployee, employeeCacheSyncRepo)

	//############

	// Product Category
	svcProductCategory := productcategoryService.NewProductCategoryHttpService(productcategoryRepo.NewProductCategoryRepository(pst), masterSyncCacheRepo)
	activityModuleManager.Add(svcProductCategory)

	// Product Barcode
	svcProductBarcode := productbarcodeService.NewProductBarcodeHttpService(productbarcodeRepo.NewProductBarcodeRepository(pst), masterSyncCacheRepo)
	activityModuleManager.Add(svcProductBarcode)

	// Product Unit
	svcProductUnit := productunitService.NewUnitHttpService(productunitRepo.NewUnitRepository(pst), masterSyncCacheRepo)
	activityModuleManager.Add(svcProductUnit)

	masterCacheSyncRepo := repositories.NewMasterSyncCacheRepository(cache)
	svcMasterSync := services.NewMasterSyncService(masterCacheSyncRepo)

	return MasterSyncHttp{
		ms:                    ms,
		cfg:                   cfg,
		activityModuleManager: activityModuleManager,

		svcMasterSync:  svcMasterSync,
		svcInventory:   svcInventory,
		svcMember:      svcMember,
		svcKitchen:     svcKitchen,
		svcShopPrinter: svcShopPrinter,
		svcShopTable:   svcShopTable,
		svcShopZone:    svcShopZone,
		svcEmployee:    *svcEmployee,
		// svcProductBarcode: *svcProductBarcode,
	}
}

func (h MasterSyncHttp) RouteSetup() {
	h.ms.GET("/master-sync", h.LastActivitySync)
	h.ms.GET("/master-sync/status", h.SyncStatus)
	h.ms.GET("/master-sync/list", h.LastActivitySyncOffset)
}

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

func (h MasterSyncHttp) LastActivitySync(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	layout := "2006-01-02T15:04"
	lastUpdateStr := ctx.QueryParam("lastupdate")

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

	page, limit := utils.GetPaginationParam(ctx.QueryParam)

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
		Page:       page,
		Limit:      limit,
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

func (h MasterSyncHttp) LastActivitySyncOffset(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	layout := "2006-01-02T15:04"
	lastUpdateStr := ctx.QueryParam("lastupdate")

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

	offset, limit := utils.GetParamOffsetLimit(ctx.QueryParam)

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
		ShopID:     shopID,
		Action:     action,
		LastUpdate: lastUpdate,
		Offset:     offset,
		Limit:      limit,
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
	LastActivity(string, string, time.Time, int, int) (models.LastActivity, mongopagination.PaginationData, error)
	LastActivityOffset(string, string, time.Time, int, int) (models.LastActivity, error)
	GetModuleName() string
}

type ActivityParamPage struct {
	ShopID     string
	Action     string
	LastUpdate time.Time
	Page       int
	Limit      int
}

type ActivityParamOffset struct {
	ShopID     string
	Action     string
	LastUpdate time.Time
	Offset     int
	Limit      int
}

func listDataModulePage(appModules map[string]ActivityModule, moduleSelectList map[string]struct{}, param ActivityParamPage) (map[string]interface{}, mongopagination.PaginationData, error) {

	result := map[string]interface{}{}

	resultPagination := mongopagination.PaginationData{}
	for moduleName, appModule := range appModules {
		if len(moduleSelectList) == 0 || isSelectModule(moduleSelectList, moduleName) {
			docList, pagination, err := appModule.LastActivity(param.ShopID, param.Action, param.LastUpdate, param.Page, param.Limit)

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
			docList, err := appModule.LastActivityOffset(param.ShopID, param.Action, param.LastUpdate, param.Offset, param.Limit)

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
