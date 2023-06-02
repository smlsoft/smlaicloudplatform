package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/transaction/paid/models"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IPaidRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.PaidDoc) (string, error)
	CreateInBatch(docList []models.PaidDoc) error
	Update(shopID string, guid string, doc models.PaidDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.PaidInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.PaidDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.PaidItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.PaidDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.PaidInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.PaidInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.PaidDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.PaidActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.PaidDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.PaidActivity, error)

	FindLastDocNo(shopID string, prefixDocNo string) (models.PaidDoc, error)
}

type PaidRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.PaidDoc]
	repositories.SearchRepository[models.PaidInfo]
	repositories.GuidRepository[models.PaidItemGuid]
	repositories.ActivityRepository[models.PaidActivity, models.PaidDeleteActivity]
}

func NewPaidRepository(pst microservice.IPersisterMongo) *PaidRepository {

	insRepo := &PaidRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.PaidDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.PaidInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.PaidItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.PaidActivity, models.PaidDeleteActivity](pst)

	return insRepo
}

func (repo PaidRepository) FindLastDocNo(shopID string, prefixDocNo string) (models.PaidDoc, error) {
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

	doc := models.PaidDoc{}
	err := repo.pst.FindOne(models.PaidDoc{}, filters, &doc, &optSort)

	if err != nil {
		return doc, err
	}

	return doc, nil
}
