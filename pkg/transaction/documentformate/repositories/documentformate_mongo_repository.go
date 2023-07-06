package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/transaction/documentformate/models"
	"time"

	"github.com/userplant/mongopagination"
)

type IDocumentFormateRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.DocumentFormateDoc) (string, error)
	CreateInBatch(docList []models.DocumentFormateDoc) error
	Update(shopID string, guid string, doc models.DocumentFormateDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.DocumentFormateInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.DocumentFormateDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.DocumentFormateItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.DocumentFormateDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.DocumentFormateInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.DocumentFormateInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.DocumentFormateDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.DocumentFormateActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.DocumentFormateDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.DocumentFormateActivity, error)
}

type DocumentFormateRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.DocumentFormateDoc]
	repositories.SearchRepository[models.DocumentFormateInfo]
	repositories.GuidRepository[models.DocumentFormateItemGuid]
	repositories.ActivityRepository[models.DocumentFormateActivity, models.DocumentFormateDeleteActivity]
}

func NewDocumentFormateRepository(pst microservice.IPersisterMongo) *DocumentFormateRepository {

	insRepo := &DocumentFormateRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.DocumentFormateDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.DocumentFormateInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.DocumentFormateItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.DocumentFormateActivity, models.DocumentFormateDeleteActivity](pst)

	return insRepo
}
