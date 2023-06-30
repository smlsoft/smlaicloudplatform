package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/pos/setting/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
)

type ISettingRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.SettingDoc) (string, error)
	CreateInBatch(docList []models.SettingDoc) error
	Update(shopID string, guid string, doc models.SettingDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.SettingInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.SettingDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.SettingItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.SettingDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.SettingInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.SettingInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SettingDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SettingActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SettingDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SettingActivity, error)
}

type SettingRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.SettingDoc]
	repositories.SearchRepository[models.SettingInfo]
	repositories.GuidRepository[models.SettingItemGuid]
	repositories.ActivityRepository[models.SettingActivity, models.SettingDeleteActivity]
}

func NewSettingRepository(pst microservice.IPersisterMongo) *SettingRepository {

	insRepo := &SettingRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.SettingDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.SettingInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.SettingItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.SettingActivity, models.SettingDeleteActivity](pst)

	return insRepo
}
