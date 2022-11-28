package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/customershop/customer/models"
	"smlcloudplatform/pkg/repositories"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type ICustomerRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.CustomerDoc) (string, error)
	CreateInBatch(docList []models.CustomerDoc) error
	Update(shopID string, guid string, doc models.CustomerDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.CustomerInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.CustomerDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.CustomerItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.CustomerDoc, error)
	FindPageSort(shopID string, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]models.CustomerInfo, mongopagination.PaginationData, error)
	FindLimit(shopID string, colNameSearch []string, q string, skip int, limit int, sorts map[string]int, projects map[string]interface{}) ([]models.CustomerInfo, int, error)
}

type CustomerRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.CustomerDoc]
	repositories.SearchRepository[models.CustomerInfo]
	repositories.GuidRepository[models.CustomerItemGuid]
}

func NewCustomerRepository(pst microservice.IPersisterMongo) *CustomerRepository {

	insRepo := &CustomerRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.CustomerDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.CustomerInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.CustomerItemGuid](pst)

	return insRepo
}
