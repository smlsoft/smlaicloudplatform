package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/channel/salechannel/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
)

type ISaleChannelRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.SaleChannelDoc) (string, error)
	CreateInBatch(docList []models.SaleChannelDoc) error
	Update(shopID string, guid string, doc models.SaleChannelDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.SaleChannelInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.SaleChannelDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.SaleChannelItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.SaleChannelDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.SaleChannelInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.SaleChannelInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SaleChannelDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SaleChannelActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SaleChannelDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SaleChannelActivity, error)
}

type SaleChannelRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.SaleChannelDoc]
	repositories.SearchRepository[models.SaleChannelInfo]
	repositories.GuidRepository[models.SaleChannelItemGuid]
	repositories.ActivityRepository[models.SaleChannelActivity, models.SaleChannelDeleteActivity]
}

func NewSaleChannelRepository(pst microservice.IPersisterMongo) *SaleChannelRepository {

	insRepo := &SaleChannelRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.SaleChannelDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.SaleChannelInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.SaleChannelItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.SaleChannelActivity, models.SaleChannelDeleteActivity](pst)

	return insRepo
}
