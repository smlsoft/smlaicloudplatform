package repositories

import (
	"smlcloudplatform/internal/microservice"
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
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.QrPaymentInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.QrPaymentDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.QrPaymentItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.QrPaymentDoc, error)
	FindPageSort(shopID string, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]models.QrPaymentInfo, mongopagination.PaginationData, error)
	FindLimit(shopID string, filters map[string]interface{}, colNameSearch []string, q string, skip int, limit int, sorts map[string]int, projects map[string]interface{}) ([]models.QrPaymentInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.QrPaymentDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.QrPaymentActivity, mongopagination.PaginationData, error)
	FindDeletedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.QrPaymentDeleteActivity, error)
	FindCreatedOrUpdatedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.QrPaymentActivity, error)
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
