package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/filefolder/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IFileFolderRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.FileFolderDoc) (string, error)
	CreateInBatch(docList []models.FileFolderDoc) error
	Update(shopID string, guid string, doc models.FileFolderDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.FileFolderInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.FileFolderDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.FileFolderItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.FileFolderDoc, error)
	FindPageSort(shopID string, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]models.FileFolderInfo, mongopagination.PaginationData, error)
	FindLimit(shopID string, filters map[string]interface{}, colNameSearch []string, q string, skip int, limit int, sorts map[string]int, projects map[string]interface{}) ([]models.FileFolderInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.FileFolderDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.FileFolderActivity, mongopagination.PaginationData, error)
	FindDeletedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.FileFolderDeleteActivity, error)
	FindCreatedOrUpdatedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.FileFolderActivity, error)
}

type FileFolderRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.FileFolderDoc]
	repositories.SearchRepository[models.FileFolderInfo]
	repositories.GuidRepository[models.FileFolderItemGuid]
	repositories.ActivityRepository[models.FileFolderActivity, models.FileFolderDeleteActivity]
}

func NewFileFolderRepository(pst microservice.IPersisterMongo) *FileFolderRepository {

	insRepo := &FileFolderRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.FileFolderDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.FileFolderInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.FileFolderItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.FileFolderActivity, models.FileFolderDeleteActivity](pst)

	return insRepo
}
