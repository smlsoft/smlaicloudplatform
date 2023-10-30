package models

import (
	salechannel_models "smlcloudplatform/pkg/channel/salechannel/models"
	device_models "smlcloudplatform/pkg/order/device/models"
	order_models "smlcloudplatform/pkg/order/setting/models"
	branch_models "smlcloudplatform/pkg/organization/branch/models"
	"smlcloudplatform/pkg/pos/media/models"
	kitchen_models "smlcloudplatform/pkg/restaurant/kitchen/models"
)

type EOrderShop struct {
	ShopID       string                   `json:"shopid"`
	Name1        string                   `json:"name1"`
	TotalTable   int                      `json:"totaltable"`
	OrderStation EOrderShopOrderStation   `json:"orderstation,omitempty"`
	Kitchens     []kitchen_models.Kitchen `json:"kitchens" bson:"kitchens"`
}

type EOrderShopOrderStation struct {
	device_models.OrderDevice
	Setting EOrderSetting `json:"ordersetting"`
}

type EOrderSetting struct {
	order_models.OrderSetting
	Branch       branch_models.Branch             `json:"branch"`
	Media        models.Media                     `json:"media"`
	SaleChannels []salechannel_models.SaleChannel `json:"salechannels" `
}

// previous version
type EOrderShopOld struct {
	ShopID       string                   `json:"shopid"`
	Name1        string                   `json:"name1"`
	TotalTable   int                      `json:"totaltable"`
	OrderStation EOrderShopOrderOld       `json:"orderstation,omitempty"`
	Kitchens     []kitchen_models.Kitchen `json:"kitchens" bson:"kitchens"`
}

type EOrderShopOrderOld struct {
	EOrderSettingOld
	DeviceInfo device_models.OrderDevice `json:"deviceinfo" bson:"deviceinfo"`
}

type EOrderSettingOld struct {
	order_models.OrderSetting
	Branch       branch_models.Branch             `json:"branch"`
	Media        models.Media                     `json:"media"`
	SaleChannels []salechannel_models.SaleChannel `json:"salechannels" `
}
