package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockAdjustment struct {
	Id        primitive.ObjectID      `json:"id" bson:"_id,omitempty"`
	ShopId    string                  `json:"shopId" bson:"shopId"`
	GuidFixed string                  `json:"guidFixed,omitempty" bson:"guidFixed"`
	Items     []StockAdjustmentDetail `json:"items" bson:"items" `
	SumAmount float64                 `json:"sumAmount" bson:"sumAmount" `
	Activity
}

func (*StockAdjustment) CollectionName() string {
	return "stockAdjustment"
}

type StockAdjustmentDetail struct {
	InventoryId    string  `json:"inventoryId" bson:"inventoryId"`
	ItemSku        string  `json:"itemSku,omitempty" bson:"itemSku,omitempty"`
	CategoryGuid   string  `json:"categoryGuid" bson:"categoryGuid"`
	LineNumber     int     `json:"lineNumber" bson:"lineNumber"`
	Price          float64 `json:"price" bson:"price" `
	Qty            float64 `json:"qty" bson:"qty" `
	DiscountAmount float64 `json:"discountAmount" bson:"discountAmount"`
	DiscountText   string  `json:"discountText" bson:"discountText"`
}

type StockAdjustmentRequest struct {
	ShopId    string                  `json:"shopId" `
	GuidFixed string                  `json:"guidFixed,omitempty" `
	Items     []StockAdjustmentDetail `json:"items" `
	SumAmount float64                 `json:"sumAmount" `
}

func (*StockAdjustmentRequest) IndexName() string {
	return "stockAdjustment"
}

func (docReq *StockAdjustmentRequest) MapRequest(doc StockAdjustment) {
	docReq.ShopId = doc.ShopId
	docReq.GuidFixed = doc.GuidFixed
	docReq.Items = doc.Items
	docReq.SumAmount = doc.SumAmount
}
