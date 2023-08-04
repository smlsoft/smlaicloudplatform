package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/transaction/saleinvoicereturn/config"
	"smlcloudplatform/pkg/transaction/saleinvoicereturn/models"
)

type ISaleInvoiceReturnMessageQueueRepository interface {
	Create(doc models.SaleInvoiceReturnDoc) error
	Update(doc models.SaleInvoiceReturnDoc) error
	Delete(doc models.SaleInvoiceReturnDoc) error
	CreateInBatch(docList []models.SaleInvoiceReturnDoc) error
	UpdateInBatch(docList []models.SaleInvoiceReturnDoc) error
	DeleteInBatch(docList []models.SaleInvoiceReturnDoc) error
}

type SaleInvoiceReturnMessageQueueRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[models.SaleInvoiceReturnDoc]
}

func NewSaleInvoiceReturnMessageQueueRepository(prod microservice.IProducer) SaleInvoiceReturnMessageQueueRepository {
	mqKey := ""

	insRepo := SaleInvoiceReturnMessageQueueRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[models.SaleInvoiceReturnDoc](prod, config.SaleInvoiceReturnMessageQueueConfig{}, "")
	return insRepo
}
