package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/transaction/stocktransfer/models"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IStockTransferRepository interface {
	FindDocOne(shopID, docno string, transFlag int) (models.StockTransferDoc, error)
	Count(shopID string) (int, error)
	Create(doc models.StockTransferDoc) (string, error)
	CreateInBatch(docList []models.StockTransferDoc) error
	Update(shopID string, guid string, doc models.StockTransferDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.StockTransferInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.StockTransferDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.StockTransferItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.StockTransferDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.StockTransferInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.StockTransferInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.StockTransferDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.StockTransferActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.StockTransferDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.StockTransferActivity, error)

	FindLastDocNo(shopID string, prefixDocNo string) (models.StockTransferDoc, error)
}

type StockTransferRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.StockTransferDoc]
	repositories.SearchRepository[models.StockTransferInfo]
	repositories.GuidRepository[models.StockTransferItemGuid]
	repositories.ActivityRepository[models.StockTransferActivity, models.StockTransferDeleteActivity]
}

func NewStockTransferRepository(pst microservice.IPersisterMongo) *StockTransferRepository {

	insRepo := &StockTransferRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.StockTransferDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.StockTransferInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.StockTransferItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.StockTransferActivity, models.StockTransferDeleteActivity](pst)

	return insRepo
}

func (repo StockTransferRepository) FindDocOne(shopID, docno string, transFlag int) (models.StockTransferDoc, error) {
	doc := models.StockTransferDoc{}
	err := repo.pst.FindOne(models.StockTransferDoc{}, bson.M{"shopid": shopID, "docno": docno, "transflag": transFlag}, &doc)

	if err != nil {
		return doc, err
	}

	return doc, nil
}

func (repo StockTransferRepository) FindLastDocNo(shopID string, prefixDocNo string) (models.StockTransferDoc, error) {
	filters := bson.M{
		"shopid": shopID,
		"deletedat": bson.M{
			"$exists": false,
		},
		"docno": bson.M{
			"$regex": "^" + prefixDocNo + ".*$",
		},
	}

	doc := models.StockTransferDoc{}
	err := repo.pst.FindOne(models.StockTransferDoc{}, filters, &doc)

	if err != nil {
		return doc, err
	}

	return doc, nil
}
