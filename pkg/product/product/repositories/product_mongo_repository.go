package repositories

import (
	"context"
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
	Count(ctx context.Context, shopID string) (int, error)
	Create(ctx context.Context, doc models.ProductDoc) (string, error)
	CreateInBatch(ctx context.Context, docList []models.ProductDoc) error
	Update(ctx context.Context, shopID string, guid string, doc models.ProductDoc) error
	DeleteByGuidfixed(ctx context.Context, shopID string, guid string, username string) error
	Delete(ctx context.Context, shopID string, username string, filters map[string]interface{}) error
	FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.ProductInfo, mongopagination.PaginationData, error)
	FindByGuid(ctx context.Context, shopID string, guid string) (models.ProductDoc, error)
	FindByGuids(ctx context.Context, shopID string, guids []string) ([]models.ProductDoc, error)
	FindMasterProductByCode(ctx context.Context, code string) (models.ProductInfo, error)

	FindInItemGuid(ctx context.Context, shopID string, columnName string, itemGuidList []string) ([]models.ProductItemGuid, error)
	FindByDocIndentityGuid(ctx context.Context, shopID string, indentityField string, indentityValue interface{}) (models.ProductDoc, error)

	FindStep(ctx context.Context, shopID string, filters map[string]interface{}, searchInFields []string, selectFields map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ProductInfo, int, error)
	FindPageFilter(ctx context.Context, shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.ProductInfo, mongopagination.PaginationData, error)

	FindDeletedPage(ctx context.Context, shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.ProductDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(ctx context.Context, shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.ProductActivity, mongopagination.PaginationData, error)
	FindDeletedStep(ctx context.Context, shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ProductDeleteActivity, error)
	FindCreatedOrUpdatedStep(ctx context.Context, shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ProductActivity, error)

	FindByBarcode(ctx context.Context, shopID string, barcode string) (models.ProductInfo, error)
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

func (repo *ProductRepository) FindMasterProductByCode(ctx context.Context, code string) (models.ProductInfo, error) {
	masterShopID := os.Getenv("MASTER_SHOP_ID")

	if len(masterShopID) == 0 {
		return models.ProductInfo{}, errors.New("master shop id is empty")
	}

	doc := models.ProductInfo{}

	filters := bson.M{
		"shopid":   masterShopID,
		"itemcode": code,
	}

	err := repo.pst.FindOne(ctx, models.ProductInfo{}, filters, &doc)

	if err != nil {
		return models.ProductInfo{}, err
	}

	return doc, nil
}

func (repo *ProductRepository) FindByBarcode(ctx context.Context, shopID string, barcode string) (models.ProductInfo, error) {

	filters := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"barcodes":  barcode,
	}

	doc := models.ProductInfo{}
	err := repo.pst.FindOne(ctx, models.ProductInfo{}, filters, &doc)

	return doc, err
}
