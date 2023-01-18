package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/smsreceive/smspaymentsettings/models"

	"github.com/userplant/mongopagination"
)

type ISmsPaymentSettingsRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.SmsPaymentSettingsDoc) (string, error)
	CreateInBatch(docList []models.SmsPaymentSettingsDoc) error
	Update(shopID string, guid string, doc models.SmsPaymentSettingsDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindOne(shopID string, filters interface{}) (models.SmsPaymentSettingsDoc, error)
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.SmsPaymentSettingsInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.SmsPaymentSettingsDoc, error)
}

type SmsPaymentSettingsRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.SmsPaymentSettingsDoc]
	repositories.SearchRepository[models.SmsPaymentSettingsInfo]
	repositories.GuidRepository[models.SmsPaymentSettingsItemGuid]
}

func NewSmsPaymentSettingsRepository(pst microservice.IPersisterMongo) SmsPaymentSettingsRepository {

	insRepo := SmsPaymentSettingsRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.SmsPaymentSettingsDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.SmsPaymentSettingsInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.SmsPaymentSettingsItemGuid](pst)

	return insRepo
}
