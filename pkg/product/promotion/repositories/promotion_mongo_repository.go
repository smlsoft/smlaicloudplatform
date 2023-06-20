package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/product/promotion/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
)

type IPromotionRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.PromotionDoc) (string, error)
	CreateInBatch(docList []models.PromotionDoc) error
	Update(shopID string, guid string, doc models.PromotionDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.PromotionInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.PromotionDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.PromotionItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.PromotionDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.PromotionInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.PromotionInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.PromotionDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.PromotionActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.PromotionDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.PromotionActivity, error)
}

type PromotionRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.PromotionDoc]
	repositories.SearchRepository[models.PromotionInfo]
	repositories.GuidRepository[models.PromotionItemGuid]
	repositories.ActivityRepository[models.PromotionActivity, models.PromotionDeleteActivity]
}

func NewPromotionRepository(pst microservice.IPersisterMongo) *PromotionRepository {

	insRepo := &PromotionRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.PromotionDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.PromotionInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.PromotionItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.PromotionActivity, models.PromotionDeleteActivity](pst)

	return insRepo
}
