package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/channel/transportchannel/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
)

type ITransportChannelRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.TransportChannelDoc) (string, error)
	CreateInBatch(docList []models.TransportChannelDoc) error
	Update(shopID string, guid string, doc models.TransportChannelDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.TransportChannelInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.TransportChannelDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.TransportChannelItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.TransportChannelDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.TransportChannelInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.TransportChannelInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.TransportChannelDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.TransportChannelActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.TransportChannelDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.TransportChannelActivity, error)
}

type TransportChannelRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.TransportChannelDoc]
	repositories.SearchRepository[models.TransportChannelInfo]
	repositories.GuidRepository[models.TransportChannelItemGuid]
	repositories.ActivityRepository[models.TransportChannelActivity, models.TransportChannelDeleteActivity]
}

func NewTransportChannelRepository(pst microservice.IPersisterMongo) *TransportChannelRepository {

	insRepo := &TransportChannelRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.TransportChannelDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.TransportChannelInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.TransportChannelItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.TransportChannelActivity, models.TransportChannelDeleteActivity](pst)

	return insRepo
}
