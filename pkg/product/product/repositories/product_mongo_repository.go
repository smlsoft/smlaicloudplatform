package repositories

import (
	"errors"
	"os"
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/product/product/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IProductRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.ProductDoc) (string, error)
	CreateInBatch(docList []models.ProductDoc) error
	Update(shopID string, guid string, doc models.ProductDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.ProductInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.ProductDoc, error)
	FindByGuids(shopID string, guids []string) ([]models.ProductDoc, error)
	FindMasterProductByCode(code string) (models.ProductInfo, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.ProductItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.ProductDoc, error)

	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, selectFields map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ProductInfo, int, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.ProductInfo, mongopagination.PaginationData, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.ProductDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.ProductActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ProductDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ProductActivity, error)
}

type ProductRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.ProductDoc]
	repositories.SearchRepository[models.ProductInfo]
	repositories.GuidRepository[models.ProductItemGuid]
	repositories.ActivityRepository[models.ProductActivity, models.ProductDeleteActivity]
}

func NewProductRepository(pst microservice.IPersisterMongo) *ProductRepository {

	insRepo := &ProductRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.ProductDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.ProductInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.ProductItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.ProductActivity, models.ProductDeleteActivity](pst)

	return insRepo
}

func (repo *ProductRepository) FindMasterProductByCode(code string) (models.ProductInfo, error) {
	masterShopID := os.Getenv("MASTER_SHOP_ID")

	if len(masterShopID) == 0 {
		return models.ProductInfo{}, errors.New("master shop id is empty")
	}

	doc := models.ProductInfo{}

	filters := bson.M{
		"shopid":   masterShopID,
		"itemcode": code,
	}

	err := repo.pst.FindOne(models.ProductInfo{}, filters, &doc)

	if err != nil {
		return models.ProductInfo{}, err
	}

	return doc, nil
}
