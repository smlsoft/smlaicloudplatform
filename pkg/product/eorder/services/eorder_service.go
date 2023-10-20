package services

import (
	"context"
	order_device_repo "smlcloudplatform/pkg/order/device/repositories"
	order_setting_repo "smlcloudplatform/pkg/order/setting/repositories"
	media_repo "smlcloudplatform/pkg/pos/media/repositories"
	"smlcloudplatform/pkg/product/eorder/models"
	"smlcloudplatform/pkg/restaurant/kitchen"
	"smlcloudplatform/pkg/restaurant/table"
	"smlcloudplatform/pkg/shop"
	"time"
)

type EOrderService struct {
	shopRepo    shop.IShopRepository
	tableRepo   table.ITableRepository
	repoOrder   order_setting_repo.ISettingRepository
	repoMedia   media_repo.IMediaRepository
	repoKitchen kitchen.IKitchenRepository
	repoDevice  order_device_repo.IDeviceRepository

	contextTimeout time.Duration
}

func NewEOrderService(
	shopRepo shop.IShopRepository,
	tableRepo table.ITableRepository,
	repoOrder order_setting_repo.ISettingRepository,
	repoMedia media_repo.IMediaRepository,
	repoKitchen kitchen.IKitchenRepository,
	repoDevice order_device_repo.IDeviceRepository,
) EOrderService {
	contextTimeout := time.Duration(15) * time.Second
	return EOrderService{
		shopRepo:       shopRepo,
		tableRepo:      tableRepo,
		repoOrder:      repoOrder,
		repoMedia:      repoMedia,
		repoKitchen:    repoKitchen,
		repoDevice:     repoDevice,
		contextTimeout: contextTimeout,
	}
}

func (svc EOrderService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc EOrderService) GetShopInfo(shopID string, orderStationCode string) (models.EOrderShop, error) {
	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	result := models.EOrderShop{}

	shopInfo, err := svc.shopRepo.FindByGuid(ctx, shopID)

	if err != nil {
		return models.EOrderShop{}, err
	}

	tableCount, err := svc.tableRepo.Count(ctx, shopID)

	if err != nil {
		return models.EOrderShop{}, err
	}

	tempOrderStation := models.EOrderShopOrderSetting{}
	order, err := svc.repoOrder.FindByDocIndentityGuid(ctx, shopID, "code", orderStationCode)

	if err != nil {
		return models.EOrderShop{}, err
	}

	tempOrderStation.OrderStation = order.OrderSetting

	if order.Code != "" {
		media, err := svc.repoMedia.FindByGuid(ctx, shopID, order.MediaGUID)

		if err != nil {
			return models.EOrderShop{}, err
		}

		result.Media = media.Media

		device, err := svc.repoDevice.FindByGuid(ctx, shopID, order.DeviceNumber)

		if err != nil {
			return models.EOrderShop{}, err
		}

		tempOrderStation.DeviceInfo = device.OrderDevice
	}

	kitchens, err := svc.repoKitchen.All(ctx, shopID)

	if err != nil {
		return models.EOrderShop{}, err
	}

	for _, tempKitchen := range kitchens {
		result.Kitchens = append(result.Kitchens, tempKitchen.Kitchen)
	}

	result.ShopID = shopInfo.ID.Hex()
	result.Name1 = shopInfo.Name1
	result.TotalTable = tableCount
	result.OrderStation = tempOrderStation

	return result, nil
}
