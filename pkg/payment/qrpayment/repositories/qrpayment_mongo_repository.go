package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/payment/qrpayment/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"github.com/userplant/mongopagination"
)

type IQrPaymentRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.QrPaymentDoc) (string, error)
	CreateInBatch(docList []models.QrPaymentDoc) error
	Update(shopID string, guid string, doc models.QrPaymentDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.QrPaymentInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.QrPaymentDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.QrPaymentItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.QrPaymentDoc, error)

	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, selectFields map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.QrPaymentInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.QrPaymentDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.QrPaymentActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.QrPaymentDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.QrPaymentActivity, error)
}

type QrPaymentRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.QrPaymentDoc]
	repositories.SearchRepository[models.QrPaymentInfo]
	repositories.GuidRepository[models.QrPaymentItemGuid]
	repositories.ActivityRepository[models.QrPaymentActivity, models.QrPaymentDeleteActivity]
}

func NewQrPaymentRepository(pst microservice.IPersisterMongo) *QrPaymentRepository {

	insRepo := &QrPaymentRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.QrPaymentDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.QrPaymentInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.QrPaymentItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.QrPaymentActivity, models.QrPaymentDeleteActivity](pst)

	return insRepo
}
