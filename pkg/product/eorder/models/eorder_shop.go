package models

import (
	order_models "smlcloudplatform/pkg/order/setting/models"
	"smlcloudplatform/pkg/pos/media/models"
)

type EOrderShop struct {
	ShopID       string                    `json:"shopid" bson:"shopid"`
	Name1        string                    `json:"name1" bson:"name1"`
	TotalTable   int                       `json:"totaltable" bson:"totaltable"`
	OrderStation order_models.OrderSetting `json:"orderstation" bson:"orderstation"`
	Media        models.Media              `json:"media" bson:"media"`
}
