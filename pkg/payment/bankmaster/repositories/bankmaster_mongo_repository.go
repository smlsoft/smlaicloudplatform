package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/payment/bankmaster/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IBankMasterRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.BankMasterDoc) (string, error)
	CreateInBatch(docList []models.BankMasterDoc) error
	Update(shopID string, guid string, doc models.BankMasterDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.BankMasterInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.BankMasterDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.BankMasterItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.BankMasterDoc, error)
	FindPageSort(shopID string, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]models.BankMasterInfo, mongopagination.PaginationData, error)
	FindLimit(shopID string, filters map[string]interface{}, colNameSearch []string, q string, skip int, limit int, sorts map[string]int, projects map[string]interface{}) ([]models.BankMasterInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.BankMasterDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.BankMasterActivity, mongopagination.PaginationData, error)
	FindDeletedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.BankMasterDeleteActivity, error)
	FindCreatedOrUpdatedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.BankMasterActivity, error)
}

type BankMasterRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.BankMasterDoc]
	repositories.SearchRepository[models.BankMasterInfo]
	repositories.GuidRepository[models.BankMasterItemGuid]
	repositories.ActivityRepository[models.BankMasterActivity, models.BankMasterDeleteActivity]
}

func NewBankMasterRepository(pst microservice.IPersisterMongo) *BankMasterRepository {

	insRepo := &BankMasterRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.BankMasterDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.BankMasterInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.BankMasterItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.BankMasterActivity, models.BankMasterDeleteActivity](pst)

	return insRepo
}
