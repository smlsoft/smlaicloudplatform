package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/transaction/purchasereturn/models"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IPurchaseReturnRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.PurchaseReturnDoc) (string, error)
	CreateInBatch(docList []models.PurchaseReturnDoc) error
	Update(shopID string, guid string, doc models.PurchaseReturnDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.PurchaseReturnInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.PurchaseReturnDoc, error)
	FindByGuids(shopID string, guids []string) ([]models.PurchaseReturnDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.PurchaseReturnItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.PurchaseReturnDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.PurchaseReturnInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.PurchaseReturnInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.PurchaseReturnDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.PurchaseReturnActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.PurchaseReturnDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.PurchaseReturnActivity, error)

	FindLastDocNo(shopID string, prefixDocNo string) (models.PurchaseReturnDoc, error)
}

type PurchaseReturnRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.PurchaseReturnDoc]
	repositories.SearchRepository[models.PurchaseReturnInfo]
	repositories.GuidRepository[models.PurchaseReturnItemGuid]
	repositories.ActivityRepository[models.PurchaseReturnActivity, models.PurchaseReturnDeleteActivity]
}

func NewPurchaseReturnRepository(pst microservice.IPersisterMongo) *PurchaseReturnRepository {

	insRepo := &PurchaseReturnRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.PurchaseReturnDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.PurchaseReturnInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.PurchaseReturnItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.PurchaseReturnActivity, models.PurchaseReturnDeleteActivity](pst)

	return insRepo
}

func (repo PurchaseReturnRepository) FindLastDocNo(shopID string, prefixDocNo string) (models.PurchaseReturnDoc, error) {
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

	doc := models.PurchaseReturnDoc{}
	err := repo.pst.FindOne(models.PurchaseReturnDoc{}, filters, &doc, &optSort)

	if err != nil {
		return doc, err
	}

	return doc, nil
}
