package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/customershop/customergroup/models"
	"smlcloudplatform/pkg/repositories"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type ICustomerGroupRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.CustomerGroupDoc) (string, error)
	CreateInBatch(docList []models.CustomerGroupDoc) error
	Update(shopID string, guid string, doc models.CustomerGroupDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.CustomerGroupInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.CustomerGroupDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.CustomerGroupItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.CustomerGroupDoc, error)
	FindPageSort(shopID string, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]models.CustomerGroupInfo, mongopagination.PaginationData, error)
	FindLimit(shopID string, filters map[string]interface{}, colNameSearch []string, q string, skip int, limit int, sorts map[string]int, projects map[string]interface{}) ([]models.CustomerGroupInfo, int, error)
}

type CustomerGroupRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.CustomerGroupDoc]
	repositories.SearchRepository[models.CustomerGroupInfo]
	repositories.GuidRepository[models.CustomerGroupItemGuid]
}

func NewCustomerGroupRepository(pst microservice.IPersisterMongo) *CustomerGroupRepository {

	insRepo := &CustomerGroupRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.CustomerGroupDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.CustomerGroupInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.CustomerGroupItemGuid](pst)

	return insRepo
}
