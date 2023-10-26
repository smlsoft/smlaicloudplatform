package services

import (
	"context"
	salechannel_repo "smlcloudplatform/pkg/channel/salechannel/repositories"
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
	shopRepo        shop.IShopRepository
	tableRepo       table.ITableRepository
	repoOrder       order_setting_repo.ISettingRepository
	repoMedia       media_repo.IMediaRepository
	repoKitchen     kitchen.IKitchenRepository
	repoDevice      order_device_repo.IDeviceRepository
	repoSaleChannel salechannel_repo.ISaleChannelRepository
	contextTimeout  time.Duration
}

func NewEOrderService(
	shopRepo shop.IShopRepository,
	tableRepo table.ITableRepository,
	repoOrder order_setting_repo.ISettingRepository,
	repoMedia media_repo.IMediaRepository,
	repoKitchen kitchen.IKitchenRepository,
	repoDevice order_device_repo.IDeviceRepository,
	repoSaleChannel salechannel_repo.ISaleChannelRepository,
) EOrderService {
	contextTimeout := time.Duration(15) * time.Second
	return EOrderService{
		shopRepo:        shopRepo,
		tableRepo:       tableRepo,
		repoOrder:       repoOrder,
		repoMedia:       repoMedia,
		repoKitchen:     repoKitchen,
		repoDevice:      repoDevice,
		repoSaleChannel: repoSaleChannel,
		contextTimeout:  contextTimeout,
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

	if orderStationCode != "" {
		orderDevice, err := svc.repoDevice.FindByDocIndentityGuid(ctx, shopID, "code", orderStationCode)

		if err != nil {
			return models.EOrderShop{}, err
		}

		tempOrderStation := models.EOrderShopOrder{}
		if orderDevice.Code != "" {
			order, err := svc.repoOrder.FindByDocIndentityGuid(ctx, shopID, "guidfixed", orderDevice.SettingCode)

			if err != nil {
				return models.EOrderShop{}, err
			}

			tempOrderStation.OrderSetting = order.OrderSetting

			if order.Code != "" {
				// Media
				media, err := svc.repoMedia.FindByGuid(ctx, shopID, order.MediaGUID)

				if err != nil {
					return models.EOrderShop{}, err
				}

				tempOrderStation.Media = media.Media

				// Device info
				tempOrderStation.DeviceInfo = orderDevice.OrderDevice

				// Sale channel
				saleChannels, err := svc.repoSaleChannel.FindByGuids(ctx, shopID, *order.SaleChannels)

				if err != nil {
					return models.EOrderShop{}, err
				}

				for _, tempSaleChannel := range saleChannels {
					tempOrderStation.SaleChannels = append(tempOrderStation.SaleChannels, tempSaleChannel.SaleChannel)
				}

			}

			result.OrderStation = tempOrderStation
		}

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

	return result, nil
}
