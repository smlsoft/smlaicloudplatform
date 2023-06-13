package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/transaction/stockreceiveproduct/models"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IStockReceiveProductRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.StockReceiveProductDoc) (string, error)
	CreateInBatch(docList []models.StockReceiveProductDoc) error
	Update(shopID string, guid string, doc models.StockReceiveProductDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.StockReceiveProductInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.StockReceiveProductDoc, error)
	FindByGuids(shopID string, guids []string) ([]models.StockReceiveProductDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.StockReceiveProductItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.StockReceiveProductDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.StockReceiveProductInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.StockReceiveProductInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.StockReceiveProductDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.StockReceiveProductActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.StockReceiveProductDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.StockReceiveProductActivity, error)

	FindLastDocNo(shopID string, prefixDocNo string) (models.StockReceiveProductDoc, error)
}

type StockReceiveProductRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.StockReceiveProductDoc]
	repositories.SearchRepository[models.StockReceiveProductInfo]
	repositories.GuidRepository[models.StockReceiveProductItemGuid]
	repositories.ActivityRepository[models.StockReceiveProductActivity, models.StockReceiveProductDeleteActivity]
}

func NewStockReceiveProductRepository(pst microservice.IPersisterMongo) *StockReceiveProductRepository {

	insRepo := &StockReceiveProductRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.StockReceiveProductDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.StockReceiveProductInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.StockReceiveProductItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.StockReceiveProductActivity, models.StockReceiveProductDeleteActivity](pst)

	return insRepo
}

func (repo StockReceiveProductRepository) FindLastDocNo(shopID string, prefixDocNo string) (models.StockReceiveProductDoc, error) {
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

	doc := models.StockReceiveProductDoc{}
	err := repo.pst.FindOne(models.StockReceiveProductDoc{}, filters, &doc, &optSort)

	if err != nil {
		return doc, err
	}

	return doc, nil
}
