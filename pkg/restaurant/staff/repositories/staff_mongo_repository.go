package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/restaurant/staff/models"
	"time"

	"github.com/userplant/mongopagination"
)

type IStaffRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.StaffDoc) (string, error)
	CreateInBatch(docList []models.StaffDoc) error
	Update(shopID string, guid string, doc models.StaffDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.StaffInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.StaffDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.StaffItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.StaffDoc, error)
	FindPageSort(shopID string, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]models.StaffInfo, mongopagination.PaginationData, error)
	FindLimit(shopID string, filters map[string]interface{}, colNameSearch []string, q string, skip int, limit int, sorts map[string]int, projects map[string]interface{}) ([]models.StaffInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.StaffDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.StaffActivity, mongopagination.PaginationData, error)
	FindDeletedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.StaffDeleteActivity, error)
	FindCreatedOrUpdatedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.StaffActivity, error)
}

type StaffRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.StaffDoc]
	repositories.SearchRepository[models.StaffInfo]
	repositories.GuidRepository[models.StaffItemGuid]
	repositories.ActivityRepository[models.StaffActivity, models.StaffDeleteActivity]
}

func NewStaffRepository(pst microservice.IPersisterMongo) *StaffRepository {

	insRepo := &StaffRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.StaffDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.StaffInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.StaffItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.StaffActivity, models.StaffDeleteActivity](pst)

	return insRepo
}
