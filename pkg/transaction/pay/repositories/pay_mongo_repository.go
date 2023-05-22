package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/transaction/pay/models"
	"time"

	"github.com/userplant/mongopagination"
)

type IPayRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.PayDoc) (string, error)
	CreateInBatch(docList []models.PayDoc) error
	Update(shopID string, guid string, doc models.PayDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.PayInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.PayDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.PayItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.PayDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.PayInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.PayInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.PayDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.PayActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.PayDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.PayActivity, error)
}

type PayRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.PayDoc]
	repositories.SearchRepository[models.PayInfo]
	repositories.GuidRepository[models.PayItemGuid]
	repositories.ActivityRepository[models.PayActivity, models.PayDeleteActivity]
}

func NewPayRepository(pst microservice.IPersisterMongo) *PayRepository {

	insRepo := &PayRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.PayDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.PayInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.PayItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.PayActivity, models.PayDeleteActivity](pst)

	return insRepo
}
