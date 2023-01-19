package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/customershop/customer/models"
	"smlcloudplatform/pkg/repositories"

	"github.com/userplant/mongopagination"
)

type ICustomerRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.CustomerDoc) (string, error)
	CreateInBatch(docList []models.CustomerDoc) error
	Update(shopID string, guid string, doc models.CustomerDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindByGuid(shopID string, guid string) (models.CustomerDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.CustomerItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.CustomerDoc, error)
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.CustomerInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.CustomerInfo, int, error)
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
