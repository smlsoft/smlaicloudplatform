package mastersync

import (
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

	// Category
	categoryDocList, categoryPagination, err := h.svcCategory.LastActivity(shopID, lastUpdate, page, limit)
	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	// Member
	memberDocList, memberPagination, err := h.svcMember.LastActivityCategory(shopID, lastUpdate, page, limit)
	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	// Inventory
	invDocList, invPagination, err := h.svcInventory.LastActivityInventory(shopID, lastUpdate, page, limit)
	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	// Kitchen
	kitchenDocList, kitchenPagination, err := h.svcKitchen.LastActivity(shopID, lastUpdate, page, limit)
	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	// Shop Printer
	shopPrinterDocList, shopPrinterPagination, err := h.svcShopPrinter.LastActivity(shopID, lastUpdate, page, limit)
	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	// Shop Table
	shopTableDocList, shopTablePagination, err := h.svcShopTable.LastActivity(shopID, lastUpdate, page, limit)
	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	// Shop Zone
	shopZoneDocList, shopZonePagination, err := h.svcShopZone.LastActivity(shopID, lastUpdate, page, limit)
	if err != nil {
		ctx.ResponseError(400, err.Error())
		return err
	}

	maxPagination := mongopagination.PaginationData{}
	paginateList := []mongopagination.PaginationData{}

	paginateList = append(paginateList, categoryPagination)
	paginateList = append(paginateList, memberPagination)
	paginateList = append(paginateList, invPagination)
	paginateList = append(paginateList, kitchenPagination)
	paginateList = append(paginateList, shopPrinterPagination)
	paginateList = append(paginateList, shopTablePagination)
	paginateList = append(paginateList, shopZonePagination)

	maxPagination = paginateList[0]
	for idx, tempPagination := range paginateList {
		if idx != 0 {
			if tempPagination.Total > maxPagination.Total {
				maxPagination = tempPagination
			}
		}
	}

	ctx.Response(
		http.StatusOK,
		models.ApiResponse{
			Success: true,
			Data: map[string]interface{}{
				"category":    categoryDocList,
				"member":      memberDocList,
				"inventory":   invDocList,
				"kitchen":     kitchenDocList,
				"shopprinter": shopPrinterDocList,
				"shoptable":   shopTableDocList,
				"shopzone":    shopZoneDocList,
			},
			Pagination: maxPagination,
		})

	return nil
}
