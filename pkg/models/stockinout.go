package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockInOut struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ShopID    string             `json:"shopID" bson:"shopID"`
	GuidFixed string             `json:"guidFixed,omitempty" bson:"guidFixed"`
	Items     []StockInOutDetail `json:"items" bson:"items" `
	SumAmount float64            `json:"sumAmount" bson:"sumAmount" `
	Activity
}

func (*StockInOut) CollectionName() string {
	return "stockInOut"
}

type StockInOutDetail struct {
	InventoryID    string  `json:"inventoryID" bson:"inventoryID"`
	ItemSku        string  `json:"itemSku,omitempty" bson:"itemSku,omitempty"`
	CategoryGuid   string  `json:"categoryGuid" bson:"categoryGuid"`
	LineNumber     int     `json:"lineNumber" bson:"lineNumber"`
	Price          float64 `json:"price" bson:"price" `
	Qty            float64 `json:"qty" bson:"qty" `
	DiscountAmount float64 `json:"discountAmount" bson:"discountAmount"`
	DiscountText   string  `json:"discountText" bson:"discountText"`
}

type StockInOutRequest struct {
	ShopID    string             `json:"shopID" `
	GuidFixed string             `json:"guidFixed,omitempty" `
	Items     []StockInOutDetail `json:"items" `
	SumAmount float64            `json:"sumAmount" `
}

func (*StockInOutRequest) IndexName() string {
	return "stockInOut"
}

func (docReq *StockInOutRequest) MapRequest(doc StockInOut) {
	docReq.ShopID = doc.ShopID
	docReq.GuidFixed = doc.GuidFixed
	docReq.Items = doc.Items
	docReq.SumAmount = doc.SumAmount
}
