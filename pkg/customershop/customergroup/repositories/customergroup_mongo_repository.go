package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/customershop/customergroup/models"
	"smlcloudplatform/pkg/repositories"

	"github.com/userplant/mongopagination"
)

type ICustomerGroupRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.CustomerGroupDoc) (string, error)
	CreateInBatch(docList []models.CustomerGroupDoc) error
	Update(shopID string, guid string, doc models.CustomerGroupDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindByGuid(shopID string, guid string) (models.CustomerGroupDoc, error)
	FindByGuids(shopID string, guids []string) ([]models.CustomerGroupDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.CustomerGroupItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.CustomerGroupDoc, error)
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.CustomerGroupInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.CustomerGroupInfo, int, error)
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
