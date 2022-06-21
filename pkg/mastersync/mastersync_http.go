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

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type MasterSyncHttp struct {
	ms             *microservice.Microservice
	cfg            microservice.IConfig
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

	// Category
	repoCategory := category.NewCategoryRepository(pst)
	svcCategory := category.NewCategoryService(repoCategory)

	// Member
	repoMember := member.NewMemberRepository(pst)
	pgRepoMember := member.NewMemberPGRepository(pstPg)
	svcMember := member.NewMemberService(repoMember, pgRepoMember)

	// Inventory
	repoInv := inventory.NewInventoryRepository(pst)
	mqRepoInv := inventory.NewInventoryMQRepository(prod)
	svcInventory := inventory.NewInventoryService(repoInv, mqRepoInv)

	// Kitchen
	repoKitchen := kitchen.NewKitchenRepository(pst)
	svcKitchen := kitchen.NewKitchenService(repoKitchen)

	// Shop Printer
	repoShopPrinter := shopprinter.NewShopPrinterRepository(pst)
	svcShopPrinter := shopprinter.NewShopPrinterService(repoShopPrinter)

	// Shop Table
	repoShopTable := shoptable.NewShopTableRepository(pst)
	svcShopTable := shoptable.NewShopTableService(repoShopTable)

	// Shop Zone
	repoShopZone := shopzone.NewShopZoneRepository(pst)
	svcShopZone := shopzone.NewShopZoneService(repoShopZone)

	return MasterSyncHttp{
		ms:             ms,
		cfg:            cfg,
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
