package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/transaction/stockpickupproduct/models"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IStockPickupProductRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.StockPickupProductDoc) (string, error)
	CreateInBatch(docList []models.StockPickupProductDoc) error
	Update(shopID string, guid string, doc models.StockPickupProductDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.StockPickupProductInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.StockPickupProductDoc, error)
	FindByGuids(shopID string, guids []string) ([]models.StockPickupProductDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.StockPickupProductItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.StockPickupProductDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.StockPickupProductInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.StockPickupProductInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.StockPickupProductDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.StockPickupProductActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.StockPickupProductDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.StockPickupProductActivity, error)

	FindLastDocNo(shopID string, prefixDocNo string) (models.StockPickupProductDoc, error)
}

type StockPickupProductRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.StockPickupProductDoc]
	repositories.SearchRepository[models.StockPickupProductInfo]
	repositories.GuidRepository[models.StockPickupProductItemGuid]
	repositories.ActivityRepository[models.StockPickupProductActivity, models.StockPickupProductDeleteActivity]
}

func NewStockPickupProductRepository(pst microservice.IPersisterMongo) *StockPickupProductRepository {

	insRepo := &StockPickupProductRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.StockPickupProductDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.StockPickupProductInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.StockPickupProductItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.StockPickupProductActivity, models.StockPickupProductDeleteActivity](pst)

	return insRepo
}
func (repo StockPickupProductRepository) FindLastDocNo(shopID string, prefixDocNo string) (models.StockPickupProductDoc, error) {
	filters := bson.M{
		"shopid": shopID,
		"deletedat": bson.M{
			"$exists": false,
		},
		"docno": bson.M{
			"$regex": "^" + prefixDocNo + ".*$",
		},
	}

	optSort := options.FindOneOptions{}
	optSort.SetSort(bson.M{
		"docno": -1,
	})

	doc := models.StockPickupProductDoc{}
	err := repo.pst.FindOne(models.StockPickupProductDoc{}, filters, &doc, &optSort)

	if err != nil {
		return doc, err
	}

	return doc, nil
}
