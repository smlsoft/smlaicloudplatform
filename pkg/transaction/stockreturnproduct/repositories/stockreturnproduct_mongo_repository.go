package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/transaction/stockreturnproduct/models"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IStockReturnProductRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.StockReturnProductDoc) (string, error)
	CreateInBatch(docList []models.StockReturnProductDoc) error
	Update(shopID string, guid string, doc models.StockReturnProductDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.StockReturnProductInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.StockReturnProductDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.StockReturnProductItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.StockReturnProductDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.StockReturnProductInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.StockReturnProductInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.StockReturnProductDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.StockReturnProductActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.StockReturnProductDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.StockReturnProductActivity, error)

	FindLastDocNo(shopID string, prefixDocNo string) (models.StockReturnProductDoc, error)
}

type StockReturnProductRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.StockReturnProductDoc]
	repositories.SearchRepository[models.StockReturnProductInfo]
	repositories.GuidRepository[models.StockReturnProductItemGuid]
	repositories.ActivityRepository[models.StockReturnProductActivity, models.StockReturnProductDeleteActivity]
}

func NewStockReturnProductRepository(pst microservice.IPersisterMongo) *StockReturnProductRepository {

	insRepo := &StockReturnProductRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.StockReturnProductDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.StockReturnProductInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.StockReturnProductItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.StockReturnProductActivity, models.StockReturnProductDeleteActivity](pst)

	return insRepo
}

func (repo StockReturnProductRepository) FindLastDocNo(shopID string, prefixDocNo string) (models.StockReturnProductDoc, error) {
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

	doc := models.StockReturnProductDoc{}
	err := repo.pst.FindOne(models.StockReturnProductDoc{}, filters, &doc, &optSort)

	if err != nil {
		return doc, err
	}

	return doc, nil
}
