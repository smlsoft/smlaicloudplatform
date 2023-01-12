package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/product/models"
	"smlcloudplatform/pkg/repositories"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type IProductRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.ProductDoc) (string, error)
	CreateInBatch(docList []models.ProductDoc) error
	Update(shopID string, guid string, doc models.ProductDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.ProductInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.ProductDoc, error)
	FindByGuids(shopID string, guids []string) ([]models.ProductDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.ProductItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.ProductDoc, error)
	FindPageSort(shopID string, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]models.ProductInfo, mongopagination.PaginationData, error)
	FindLimit(shopID string, filters map[string]interface{}, colNameSearch []string, q string, skip int, limit int, sorts map[string]int, projects map[string]interface{}) ([]models.ProductInfo, int, error)
	FindPageFilterSort(shopID string, filters map[string]interface{}, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]models.ProductInfo, mongopagination.PaginationData, error)
}

type ProductRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.ProductDoc]
	repositories.SearchRepository[models.ProductInfo]
	repositories.GuidRepository[models.ProductItemGuid]
}

func NewProductRepository(pst microservice.IPersisterMongo) *ProductRepository {

	insRepo := &ProductRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.ProductDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.ProductInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.ProductItemGuid](pst)

	return insRepo
}
