package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/transaction/saleinvoice/models"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ISaleInvoiceRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.SaleInvoiceDoc) (string, error)
	CreateInBatch(docList []models.SaleInvoiceDoc) error
	Update(shopID string, guid string, doc models.SaleInvoiceDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.SaleInvoiceInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.SaleInvoiceDoc, error)
	FindByGuids(shopID string, guids []string) ([]models.SaleInvoiceDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.SaleInvoiceItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.SaleInvoiceDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.SaleInvoiceInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.SaleInvoiceInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SaleInvoiceDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SaleInvoiceActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SaleInvoiceDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SaleInvoiceActivity, error)

	FindLastDocNo(shopID string, prefixDocNo string) (models.SaleInvoiceDoc, error)
}

type SaleInvoiceRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.SaleInvoiceDoc]
	repositories.SearchRepository[models.SaleInvoiceInfo]
	repositories.GuidRepository[models.SaleInvoiceItemGuid]
	repositories.ActivityRepository[models.SaleInvoiceActivity, models.SaleInvoiceDeleteActivity]
}

func NewSaleInvoiceRepository(pst microservice.IPersisterMongo) *SaleInvoiceRepository {

	insRepo := &SaleInvoiceRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.SaleInvoiceDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.SaleInvoiceInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.SaleInvoiceItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.SaleInvoiceActivity, models.SaleInvoiceDeleteActivity](pst)

	return insRepo
}

func (repo SaleInvoiceRepository) FindLastDocNo(shopID string, prefixDocNo string) (models.SaleInvoiceDoc, error) {
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

	doc := models.SaleInvoiceDoc{}
	err := repo.pst.FindOne(models.SaleInvoiceDoc{}, filters, &doc, &optSort)

	if err != nil {
		return doc, err
	}

	return doc, nil
}
