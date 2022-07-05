package shopprinter

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/restaurant/shopprinter/models"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IShopPrinterRepository interface {
	Count(shopID string) (int, error)
	Create(category models.PrinterTerminalDoc) (string, error)
	CreateInBatch(inventories []models.PrinterTerminalDoc) error
	Update(shopID string, guid string, category models.PrinterTerminalDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.PrinterTerminalInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.PrinterTerminalDoc, error)
	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.PrinterTerminalItemGuid, error)
	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.PrinterTerminalDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.PrinterTerminalActivity, mongopagination.PaginationData, error)
}

type ShopPrinterRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.PrinterTerminalDoc]
	repositories.SearchRepository[models.PrinterTerminalInfo]
	repositories.GuidRepository[models.PrinterTerminalItemGuid]
	repositories.ActivityRepository[models.PrinterTerminalActivity, models.PrinterTerminalDeleteActivity]
}

func NewShopPrinterRepository(pst microservice.IPersisterMongo) ShopPrinterRepository {

	insRepo := ShopPrinterRepository{
		pst: pst,
	}
	insRepo.CrudRepository = repositories.NewCrudRepository[models.PrinterTerminalDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.PrinterTerminalInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.PrinterTerminalItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.PrinterTerminalActivity, models.PrinterTerminalDeleteActivity](pst)

	return insRepo
}
