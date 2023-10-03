package services

import (
	"context"
	"smlcloudplatform/pkg/product/eorder/models"
	"smlcloudplatform/pkg/restaurant/table"
	"smlcloudplatform/pkg/shop"
	"time"
)

type EOrderService struct {
	shopRepo       shop.IShopRepository
	tableRepo      table.ITableRepository
	contextTimeout time.Duration
}

func NewEOrderService(
	shopRepo shop.IShopRepository,
	tableRepo table.ITableRepository,
) EOrderService {
	contextTimeout := time.Duration(15) * time.Second
	return EOrderService{
		shopRepo:       shopRepo,
		tableRepo:      tableRepo,
		contextTimeout: contextTimeout,
	}
}

func (svc EOrderService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc EOrderService) GetShopInfo(shopID string) (models.EOrderShop, error) {
	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	shopInfo, err := svc.shopRepo.FindByGuid(ctx, shopID)

	if err != nil {
		return models.EOrderShop{}, err
	}

	tableCount, err := svc.tableRepo.Count(ctx, shopID)

	if err != nil {
		return models.EOrderShop{}, err
	}

	return models.EOrderShop{
		ShopID:     shopInfo.ID.Hex(),
		Name1:      shopInfo.Name1,
		TotalTable: tableCount,
	}, nil
}
