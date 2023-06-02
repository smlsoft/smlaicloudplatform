package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/transaction/saleinvoicereturn/models"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type ISaleInvoiceReturnRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.SaleInvoiceReturnDoc) (string, error)
	CreateInBatch(docList []models.SaleInvoiceReturnDoc) error
	Update(shopID string, guid string, doc models.SaleInvoiceReturnDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.SaleInvoiceReturnInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.SaleInvoiceReturnDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.SaleInvoiceReturnItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.SaleInvoiceReturnDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.SaleInvoiceReturnInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.SaleInvoiceReturnInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SaleInvoiceReturnDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SaleInvoiceReturnActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SaleInvoiceReturnDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SaleInvoiceReturnActivity, error)

	FindLastDocNo(shopID string, prefixDocNo string) (models.SaleInvoiceReturnDoc, error)
}

type SaleInvoiceReturnRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.SaleInvoiceReturnDoc]
	repositories.SearchRepository[models.SaleInvoiceReturnInfo]
	repositories.GuidRepository[models.SaleInvoiceReturnItemGuid]
	repositories.ActivityRepository[models.SaleInvoiceReturnActivity, models.SaleInvoiceReturnDeleteActivity]
}

func NewSaleInvoiceReturnRepository(pst microservice.IPersisterMongo) *SaleInvoiceReturnRepository {

	insRepo := &SaleInvoiceReturnRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.SaleInvoiceReturnDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.SaleInvoiceReturnInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.SaleInvoiceReturnItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.SaleInvoiceReturnActivity, models.SaleInvoiceReturnDeleteActivity](pst)

	return insRepo
}

func (repo SaleInvoiceReturnRepository) FindLastDocNo(shopID string, prefixDocNo string) (models.SaleInvoiceReturnDoc, error) {
	filters := bson.M{
		"shopid": shopID,
		"deletedat": bson.M{
			"$exists": false,
		},
		"docno": bson.M{
			"$regex": "^" + prefixDocNo + ".*$",
		},
	}

	doc := models.SaleInvoiceReturnDoc{}
	err := repo.pst.FindOne(models.SaleInvoiceReturnDoc{}, filters, &doc)

	if err != nil {
		return doc, err
	}

	return doc, nil
}
