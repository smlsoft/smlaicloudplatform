package repositories

import (
	"smlaicloudplatform/internal/product/productbarcode/config"
	"smlaicloudplatform/internal/product/productbarcode/models"
	"smlaicloudplatform/internal/repositories"
	"smlaicloudplatform/pkg/microservice"
)

type IProductBarcodeMessageQueueRepository interface {
	Create(doc models.ProductBarcodeDoc) error
	Update(doc models.ProductBarcodeDoc) error
	Delete(doc models.ProductBarcodeDoc) error
	CreateInBatch(docList []models.ProductBarcodeDoc) error
	UpdateInBatch(docList []models.ProductBarcodeDoc) error
	DeleteInBatch(docList []models.ProductBarcodeDoc) error
}

type ProductBarcodeMessageQueueRepository struct {
	prod  microservice.IProducer
	mqKey string
	repositories.KafkaRepository[models.ProductBarcodeDoc]
}

func NewProductBarcodeMessageQueueRepository(prod microservice.IProducer) ProductBarcodeMessageQueueRepository {
	mqKey := ""

	insRepo := ProductBarcodeMessageQueueRepository{
		prod:  prod,
		mqKey: mqKey,
	}
	insRepo.KafkaRepository = repositories.NewKafkaRepository[models.ProductBarcodeDoc](prod, config.ProductMessageQueueConfig{}, "")
	return insRepo
}
