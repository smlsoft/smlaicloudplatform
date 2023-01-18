package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/unit/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
)

type IUnitRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.UnitDoc) (string, error)
	CreateInBatch(docList []models.UnitDoc) error
	Update(shopID string, guid string, doc models.UnitDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.UnitInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.UnitDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.UnitItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.UnitDoc, error)
	FindPageSort(shopID string, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]models.UnitInfo, mongopagination.PaginationData, error)
	FindLimit(shopID string, filters map[string]interface{}, colNameSearch []string, q string, skip int, limit int, sorts map[string]int, projects map[string]interface{}) ([]models.UnitInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.UnitDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.UnitActivity, mongopagination.PaginationData, error)
	FindDeletedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.UnitDeleteActivity, error)
	FindCreatedOrUpdatedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.UnitActivity, error)
}

type UnitRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.UnitDoc]
	repositories.SearchRepository[models.UnitInfo]
	repositories.GuidRepository[models.UnitItemGuid]
	repositories.ActivityRepository[models.UnitActivity, models.UnitDeleteActivity]
}

func NewUnitRepository(pst microservice.IPersisterMongo) *UnitRepository {

	insRepo := &UnitRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.UnitDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.UnitInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.UnitItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.UnitActivity, models.UnitDeleteActivity](pst)

	return insRepo
}
