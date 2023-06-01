package printer

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/restaurant/printer/models"
	"time"

	"github.com/userplant/mongopagination"
)

type IPrinterRepository interface {
	Count(shopID string) (int, error)
	Create(category models.PrinterDoc) (string, error)
	CreateInBatch(docList []models.PrinterDoc) error
	Update(shopID string, guid string, category models.PrinterDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.PrinterInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.PrinterDoc, error)
	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.PrinterItemGuid, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, pageable micromodels.Pageable) ([]models.PrinterDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, pageable micromodels.Pageable) ([]models.PrinterActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, pageableStep micromodels.PageableStep) ([]models.PrinterDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, pageableStep micromodels.PageableStep) ([]models.PrinterActivity, error)
	FindStep(shopID string, searchInFields []string, q string, skip int, limit int, sorts map[string]int, projects map[string]interface{}) ([]models.PrinterInfo, int, error)
}

type PrinterRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.PrinterDoc]
	repositories.SearchRepository[models.PrinterInfo]
	repositories.GuidRepository[models.PrinterItemGuid]
	repositories.ActivityRepository[models.PrinterActivity, models.PrinterDeleteActivity]
}

func NewPrinterRepository(pst microservice.IPersisterMongo) PrinterRepository {

	insRepo := PrinterRepository{
		pst: pst,
	}
	insRepo.CrudRepository = repositories.NewCrudRepository[models.PrinterDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.PrinterInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.PrinterItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.PrinterActivity, models.PrinterDeleteActivity](pst)

	return insRepo
}
