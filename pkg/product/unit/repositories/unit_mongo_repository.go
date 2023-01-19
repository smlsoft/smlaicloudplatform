package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
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
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.UnitInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.UnitDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.UnitItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.UnitDoc, error)

	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, selectFields map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.UnitInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, pageable micromodels.Pageable) ([]models.UnitDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, pageable micromodels.Pageable) ([]models.UnitActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, pageableStep micromodels.PageableStep) ([]models.UnitDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, pageableStep micromodels.PageableStep) ([]models.UnitActivity, error)
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
