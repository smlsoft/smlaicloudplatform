package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/payment/bookbank/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
)

type IBookBankRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.BookBankDoc) (string, error)
	CreateInBatch(docList []models.BookBankDoc) error
	Update(shopID string, guid string, doc models.BookBankDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.BookBankInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.BookBankDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.BookBankItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.BookBankDoc, error)
	FindPageSort(shopID string, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]models.BookBankInfo, mongopagination.PaginationData, error)
	FindLimit(shopID string, filters map[string]interface{}, colNameSearch []string, q string, skip int, limit int, sorts map[string]int, projects map[string]interface{}) ([]models.BookBankInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.BookBankDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.BookBankActivity, mongopagination.PaginationData, error)
	FindDeletedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.BookBankDeleteActivity, error)
	FindCreatedOrUpdatedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.BookBankActivity, error)
}

type BookBankRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.BookBankDoc]
	repositories.SearchRepository[models.BookBankInfo]
	repositories.GuidRepository[models.BookBankItemGuid]
	repositories.ActivityRepository[models.BookBankActivity, models.BookBankDeleteActivity]
}

func NewBookBankRepository(pst microservice.IPersisterMongo) *BookBankRepository {

	insRepo := &BookBankRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.BookBankDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.BookBankInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.BookBankItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.BookBankActivity, models.BookBankDeleteActivity](pst)

	return insRepo
}
