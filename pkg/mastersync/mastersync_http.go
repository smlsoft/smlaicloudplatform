package mastersync

import (
	"fmt"
	"net/http"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/api/category"
	"smlcloudplatform/pkg/api/member"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"strings"
	"time"

	"smlcloudplatform/pkg/api/inventory"
	"smlcloudplatform/pkg/api/restaurant/kitchen"
	"smlcloudplatform/pkg/api/restaurant/shopprinter"
	"smlcloudplatform/pkg/api/restaurant/shoptable"
	"smlcloudplatform/pkg/api/restaurant/shopzone"
	"smlcloudplatform/pkg/mastersync/services"

	"smlcloudplatform/pkg/mastersync/repositories"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type MasterSyncHttp struct {
	ms             *microservice.Microservice
	cfg            microservice.IConfig
	svcMasterSync  services.IMasterSyncService
	svcCategory    category.ICategoryService
	svcMember      member.IMemberService
	svcInventory   inventory.IInventoryService
	svcKitchen     kitchen.IKitchenService
	svcShopPrinter shopprinter.IShopPrinterService
	svcShopTable   shoptable.ShopTableService
	svcShopZone    shopzone.ShopZoneService
}

func NewMasterSyncHttp(ms *microservice.Microservice, cfg microservice.IConfig) MasterSyncHttp {
	pst := ms.MongoPersister(cfg.MongoPersisterConfig())
	pstPg := ms.Persister(cfg.PersisterConfig())
	prod := ms.Producer(cfg.MQConfig())
	cache := ms.Cacher(cfg.CacherConfig())

	// Category
	repoCategory := category.NewCategoryRepository(pst)
	repoCacheSyncCategory := repositories.NewMasterSyncCacheRepository(cache, "category")
	svcCategory := category.NewCategoryService(repoCategory, repoCacheSyncCategory)

	// Member
	repoMember := member.NewMemberRepository(pst)
	pgRepoMember := member.NewMemberPGRepository(pstPg)
	repoCacheSyncMember := repositories.NewMasterSyncCacheRepository(cache, "member")
	svcMember := member.NewMemberService(repoMember, pgRepoMember, repoCacheSyncMember)

	// Inventory
	repoInv := inventory.NewInventoryRepository(pst)
	mqRepoInv := inventory.NewInventoryMQRepository(prod)
	invCacheSyncRepo := repositories.NewMasterSyncCacheRepository(cache, "inventory")
	svcInventory := inventory.NewInventoryService(repoInv, mqRepoInv, invCacheSyncRepo)

	// Kitchen
	repoKitchen := kitchen.NewKitchenRepository(pst)
	kitchenCacheSyncRepo := repositories.NewMasterSyncCacheRepository(cache, "kitchen")
	svcKitchen := kitchen.NewKitchenService(repoKitchen, kitchenCacheSyncRepo)

	// Shop Printer
	repoShopPrinter := shopprinter.NewShopPrinterRepository(pst)
	shopPrinterCacheSyncRepo := repositories.NewMasterSyncCacheRepository(cache, "shopprinter")
	svcShopPrinter := shopprinter.NewShopPrinterService(repoShopPrinter, shopPrinterCacheSyncRepo)

	// Shop Table
	repoShopTable := shoptable.NewShopTableRepository(pst)
	shopTableCacheSyncRepo := repositories.NewMasterSyncCacheRepository(cache, "shoptable")
	svcShopTable := shoptable.NewShopTableService(repoShopTable, shopTableCacheSyncRepo)

	// Shop Zone
	repoShopZone := shopzone.NewShopZoneRepository(pst)
	shopZoneCacheSyncRepo := repositories.NewMasterSyncCacheRepository(cache, "shopzone")
	svcShopZone := shopzone.NewShopZoneService(repoShopZone, shopZoneCacheSyncRepo)

	masterCacheSyncRepo := repositories.NewMasterSyncCacheRepository(cache, "mastersync")
	svcMasterSync := services.NewMasterSyncService(masterCacheSyncRepo)

	return MasterSyncHttp{
		ms:             ms,
		cfg:            cfg,
		svcMasterSync:  svcMasterSync,
		svcCategory:    svcCategory,
		svcInventory:   svcInventory,
		svcMember:      svcMember,
		svcKitchen:     svcKitchen,
		svcShopPrinter: svcShopPrinter,
		svcShopTable:   svcShopTable,
		svcShopZone:    svcShopZone,
	}
}

func (h MasterSyncHttp) RouteSetup() {
	h.ms.GET("/master-sync", h.LastActivityCategory)
	h.ms.GET("/master-sync/status", h.SyncStatus)
}

func (h MasterSyncHttp) SyncStatus(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID
	status, _ := h.svcMasterSync.GetStatus(shopID)

	ctx.Response(
		http.StatusOK,
		status,
	)

	return nil
}

func (h MasterSyncHttp) LastActivityCategory(ctx microservice.IContext) error {
	userInfo := ctx.UserInfo()
	shopID := userInfo.ShopID

	layout := "2006-01-02T15:04"
	lastUpdateStr := ctx.QueryParam("lastUpdate")

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

	moduleSelectList := []string{}
	keySelectList := map[string]bool{}

	if moduleParam != "" {
		moduleSelectList = strings.Split(moduleParam, ",")
		for _, module := range moduleSelectList {
			module = strings.ToLower(module)
			keySelectList[module] = true
		}
	}

	isSelectAll := false

	if len(moduleSelectList) < 1 {
		isSelectAll = true
	} else if strings.ToLower(moduleSelectList[0]) == "all" {
		isSelectAll = true
	}

	moduleList := map[string]ActivityModule{}

	moduleList["category"] = h.svcCategory
	moduleList["member"] = h.svcMember
	moduleList["inventory"] = h.svcInventory
	moduleList["kitchen"] = h.svcKitchen
	moduleList["shopprinter"] = h.svcShopPrinter
	moduleList["shoptable"] = h.svcShopTable
	moduleList["shopzone"] = h.svcShopZone

	result, pagination, err := runModule(moduleList, isSelectAll, keySelectList, ActivityParam{
		ShopID:     shopID,
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
			Data:       result,
			Pagination: pagination,
		})

	return nil
}

type ActivityModule interface {
	LastActivity(string, time.Time, int, int) (models.LastActivity, mongopagination.PaginationData, error)
}

type ActivityParam struct {
	ShopID     string
	LastUpdate time.Time
	Page       int
	Limit      int
}

func runModule(appModules map[string]ActivityModule, isSelectAll bool, keySelectList map[string]bool, param ActivityParam) (map[string]interface{}, mongopagination.PaginationData, error) {

	result := map[string]interface{}{}

	resultPagination := mongopagination.PaginationData{}
	for moduleName, appModule := range appModules {
		if isSelectAll || isSelect(keySelectList, moduleName) {
			docList, pagination, err := appModule.LastActivity(param.ShopID, param.LastUpdate, param.Page, param.Limit)

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

func isSelect(keyList map[string]bool, key string) bool {
	if _, ok := keyList[key]; ok {
		return true
	}
	return false
}
