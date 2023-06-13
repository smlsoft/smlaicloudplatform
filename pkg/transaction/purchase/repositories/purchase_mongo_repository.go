package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/transaction/purchase/models"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IPurchaseRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.PurchaseDoc) (string, error)
	CreateInBatch(docList []models.PurchaseDoc) error
	Update(shopID string, guid string, doc models.PurchaseDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.PurchaseInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.PurchaseDoc, error)
	FindByGuids(shopID string, guids []string) ([]models.PurchaseDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.PurchaseItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.PurchaseDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.PurchaseInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.PurchaseInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.PurchaseDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.PurchaseActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.PurchaseDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.PurchaseActivity, error)

	FindLastDocNo(shopID string, prefixDocNo string) (models.PurchaseDoc, error)
}

type PurchaseRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.PurchaseDoc]
	repositories.SearchRepository[models.PurchaseInfo]
	repositories.GuidRepository[models.PurchaseItemGuid]
	repositories.ActivityRepository[models.PurchaseActivity, models.PurchaseDeleteActivity]
}

func NewPurchaseRepository(pst microservice.IPersisterMongo) *PurchaseRepository {

	insRepo := &PurchaseRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.PurchaseDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.PurchaseInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.PurchaseItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.PurchaseActivity, models.PurchaseDeleteActivity](pst)

	return insRepo
}
func (repo PurchaseRepository) FindLastDocNo(shopID string, prefixDocNo string) (models.PurchaseDoc, error) {
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

	doc := models.PurchaseDoc{}
	err := repo.pst.FindOne(models.PurchaseDoc{}, filters, &doc, &optSort)

	if err != nil {
		return doc, err
	}

	return doc, nil
}
