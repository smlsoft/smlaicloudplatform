package shop

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/shop/models"
)

type IShopUserAccessLogRepository interface {
	Create(shopUserAccessLog models.ShopUserAccessLog) error
}

type ShopUserAccessLogRepository struct {
	pst microservice.IPersisterMongo
}

func NewShopUserAccessLogRepository(pst microservice.IPersisterMongo) ShopUserAccessLogRepository {
	return ShopUserAccessLogRepository{
		pst: pst,
	}
}

func (svc ShopUserAccessLogRepository) Create(shopUserAccessLog models.ShopUserAccessLog) error {

	_, err := svc.pst.Create(&models.ShopUserAccessLog{}, shopUserAccessLog)

	if err != nil {
		return err
	}

	return nil
}
