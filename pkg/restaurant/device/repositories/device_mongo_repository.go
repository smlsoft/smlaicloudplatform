package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/restaurant/device/models"
	"time"

	"github.com/userplant/mongopagination"
)

type IDeviceRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.DeviceDoc) (string, error)
	CreateInBatch(docList []models.DeviceDoc) error
	Update(shopID string, guid string, doc models.DeviceDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.DeviceInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.DeviceDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.DeviceItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.DeviceDoc, error)
	FindPageSort(shopID string, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]models.DeviceInfo, mongopagination.PaginationData, error)
	FindLimit(shopID string, filters map[string]interface{}, colNameSearch []string, q string, skip int, limit int, sorts map[string]int, projects map[string]interface{}) ([]models.DeviceInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.DeviceDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.DeviceActivity, mongopagination.PaginationData, error)
	FindDeletedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.DeviceDeleteActivity, error)
	FindCreatedOrUpdatedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.DeviceActivity, error)
}

type DeviceRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.DeviceDoc]
	repositories.SearchRepository[models.DeviceInfo]
	repositories.GuidRepository[models.DeviceItemGuid]
	repositories.ActivityRepository[models.DeviceActivity, models.DeviceDeleteActivity]
}

func NewDeviceRepository(pst microservice.IPersisterMongo) *DeviceRepository {

	insRepo := &DeviceRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.DeviceDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.DeviceInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.DeviceItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.DeviceActivity, models.DeviceDeleteActivity](pst)

	return insRepo
}
