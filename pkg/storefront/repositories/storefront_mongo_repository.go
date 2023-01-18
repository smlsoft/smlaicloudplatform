package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/storefront/models"

	"github.com/userplant/mongopagination"
)

type IStorefrontRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.StorefrontDoc) (string, error)
	CreateInBatch(docList []models.StorefrontDoc) error
	Update(shopID string, guid string, doc models.StorefrontDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.StorefrontInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.StorefrontDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.StorefrontItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.StorefrontDoc, error)
	FindPageSort(shopID string, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]models.StorefrontInfo, mongopagination.PaginationData, error)
}

type StorefrontRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.StorefrontDoc]
	repositories.SearchRepository[models.StorefrontInfo]
	repositories.GuidRepository[models.StorefrontItemGuid]
}

func NewStorefrontRepository(pst microservice.IPersisterMongo) *StorefrontRepository {

	insRepo := &StorefrontRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.StorefrontDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.StorefrontInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.StorefrontItemGuid](pst)

	return insRepo
}
