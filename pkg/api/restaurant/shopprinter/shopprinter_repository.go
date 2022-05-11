package shopprinter

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models/restaurant"
	"smlcloudplatform/pkg/repositories"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IShopPrinterRepository interface {
	Count(shopID string) (int, error)
	Create(category restaurant.PrinterTerminalDoc) (string, error)
	CreateInBatch(inventories []restaurant.PrinterTerminalDoc) error
	Update(guid string, category restaurant.PrinterTerminalDoc) error
	Delete(shopID string, guid string, username string) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]restaurant.PrinterTerminalInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (restaurant.PrinterTerminalDoc, error)
	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]restaurant.PrinterTerminalItemGuid, error)
	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]restaurant.PrinterTerminalDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]restaurant.PrinterTerminalActivity, mongopagination.PaginationData, error)
}

type ShopPrinterRepository struct {
	pst          microservice.IPersisterMongo
	crudRepo     repositories.CrudRepository[restaurant.PrinterTerminalDoc]
	searchRepo   repositories.SearchRepository[restaurant.PrinterTerminalInfo]
	guidRepo     repositories.GuidRepository[restaurant.PrinterTerminalItemGuid]
	activityRepo repositories.ActivityRepository[restaurant.PrinterTerminalActivity, restaurant.PrinterTerminalDeleteActivity]
}

func NewShopPrinterRepository(pst microservice.IPersisterMongo) ShopPrinterRepository {
	crudRepo := repositories.NewCrudRepository[restaurant.PrinterTerminalDoc](pst)
	searchRepo := repositories.NewSearchRepository[restaurant.PrinterTerminalInfo](pst)
	guidRepo := repositories.NewGuidRepository[restaurant.PrinterTerminalItemGuid](pst)
	activityRepo := repositories.NewActivityRepository[restaurant.PrinterTerminalActivity, restaurant.PrinterTerminalDeleteActivity](pst)

	return ShopPrinterRepository{
		pst:          pst,
		crudRepo:     crudRepo,
		searchRepo:   searchRepo,
		guidRepo:     guidRepo,
		activityRepo: activityRepo,
	}
}
