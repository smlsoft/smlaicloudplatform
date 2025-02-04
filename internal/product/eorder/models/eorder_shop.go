package models

import (
	salechannel_models "smlaicloudplatform/internal/channel/salechannel/models"
	notify_models "smlaicloudplatform/internal/notify/models"
	device_models "smlaicloudplatform/internal/order/device/models"
	order_models "smlaicloudplatform/internal/order/setting/models"
	branch_models "smlaicloudplatform/internal/organization/branch/models"
	"smlaicloudplatform/internal/pos/media/models"
	kitchen_models "smlaicloudplatform/internal/restaurant/kitchen/models"
)

type EOrderShop struct {
	ShopID         string                     `json:"shopid"`
	Name1          string                     `json:"name1"`
	ProfilePicture string                     `json:"profilepicture"`
	TotalTable     int                        `json:"totaltable"`
	OrderStation   EOrderShopOrderStation     `json:"orderstation,omitempty"`
	Kitchens       []kitchen_models.Kitchen   `json:"kitchens" bson:"kitchens"`
	Notify         []notify_models.NotifyInfo `json:"notify"`
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
	ShopID         string                   `json:"shopid"`
	Name1          string                   `json:"name1"`
	ProfilePicture string                   `json:"profilepicture"`
	TotalTable     int                      `json:"totaltable"`
	OrderStation   EOrderShopOrderOld       `json:"orderstation,omitempty"`
	Kitchens       []kitchen_models.Kitchen `json:"kitchens" bson:"kitchens"`
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
