package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/organization/businesstype/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
)

type IBusinessTypeRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.BusinessTypeDoc) (string, error)
	CreateInBatch(docList []models.BusinessTypeDoc) error
	Update(shopID string, guid string, doc models.BusinessTypeDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.BusinessTypeInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.BusinessTypeDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.BusinessTypeItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.BusinessTypeDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.BusinessTypeInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.BusinessTypeInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.BusinessTypeDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.BusinessTypeActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.BusinessTypeDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.BusinessTypeActivity, error)
}

type BusinessTypeRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.BusinessTypeDoc]
	repositories.SearchRepository[models.BusinessTypeInfo]
	repositories.GuidRepository[models.BusinessTypeItemGuid]
	repositories.ActivityRepository[models.BusinessTypeActivity, models.BusinessTypeDeleteActivity]
}

func NewBusinessTypeRepository(pst microservice.IPersisterMongo) *BusinessTypeRepository {

	insRepo := &BusinessTypeRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.BusinessTypeDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.BusinessTypeInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.BusinessTypeItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.BusinessTypeActivity, models.BusinessTypeDeleteActivity](pst)

	return insRepo
}
