package repositories

import (
	"smlaicloudplatform/internal/repositories"
	"smlaicloudplatform/internal/transaction/saleinvoice/config"
	"smlaicloudplatform/internal/transaction/saleinvoice/models"
	"smlaicloudplatform/pkg/microservice"
)

type ISaleInvoiceMessageQueueRepository interface {
	Create(doc models.SaleInvoiceDoc) error
	Update(doc models.SaleInvoiceDoc) error
	Delete(doc models.SaleInvoiceDoc) error
	CreateInBatch(docList []models.SaleInvoiceDoc) error
	UpdateInBatch(docList []models.SaleInvoiceDoc) error
	DeleteInBatch(docList []models.SaleInvoiceDoc) error
}

type SaleInvoiceMessageQueueRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[models.SaleInvoiceDoc]
}

func NewSaleInvoiceMessageQueueRepository(prod microservice.IProducer) SaleInvoiceMessageQueueRepository {
	mqKey := ""

	insRepo := SaleInvoiceMessageQueueRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[models.SaleInvoiceDoc](prod, config.SaleInvoiceMessageQueueConfig{}, "")
	return insRepo
}
