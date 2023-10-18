package models

import (
	salechannel_models "smlcloudplatform/pkg/channel/salechannel/models"
	order_models "smlcloudplatform/pkg/order/setting/models"
	"smlcloudplatform/pkg/pos/media/models"
	kitchen_models "smlcloudplatform/pkg/restaurant/kitchen/models"
)

type EOrderShop struct {
	ShopID       string                         `json:"shopid" bson:"shopid"`
	Name1        string                         `json:"name1" bson:"name1"`
	TotalTable   int                            `json:"totaltable" bson:"totaltable"`
	OrderStation order_models.OrderSetting      `json:"orderstation" bson:"orderstation"`
	Media        models.Media                   `json:"media" bson:"media"`
	Kitchens     []kitchen_models.Kitchen       `json:"kitchens" bson:"kitchens"`
	SaleChannels salechannel_models.SaleChannel `json:"salechannels" bson:"salechannels"`
}
