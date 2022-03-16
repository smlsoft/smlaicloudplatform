package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockInOut struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ShopId    string             `json:"shopId" bson:"shopId"`
	GuidFixed string             `json:"guidFixed,omitempty" bson:"guidFixed"`
	Items     []StockInOutDetail `json:"items" bson:"items" `
	SumAmount float64            `json:"sumAmount" bson:"sumAmount" `
	CreatedBy string             `json:"createdBy" bson:"createdBy"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedBy string             `json:"updatedBy,omitempty" bson:"updatedBy,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	Deleted   bool               `json:"-" bson:"deleted"`
}

func (*StockInOut) CollectionName() string {
	return "stockInOut"
}

type StockInOutDetail struct {
	InventoryId    string  `json:"inventoryId" bson:"inventoryId"`
	ItemSku        string  `json:"itemSku,omitempty" bson:"itemSku,omitempty"`
	CategoryGuid   string  `json:"categoryGuid" bson:"categoryGuid"`
	LineNumber     int     `json:"lineNumber" bson:"lineNumber"`
	Price          float64 `json:"price" bson:"price" `
	Qty            float64 `json:"qty" bson:"qty" `
	DiscountAmount float64 `json:"discountAmount" bson:"discountAmount"`
	DiscountText   string  `json:"discountText" bson:"discountText"`
}

type StockInOutRequest struct {
	ShopId    string             `json:"shopId" `
	GuidFixed string             `json:"guidFixed,omitempty" `
	Items     []StockInOutDetail `json:"items" `
	SumAmount float64            `json:"sumAmount" `
}

func (*StockInOutRequest) IndexName() string {
	return "stockInOut"
}

func (docReq *StockInOutRequest) MapRequest(doc StockInOut) {
	docReq.ShopId = doc.ShopId
	docReq.GuidFixed = doc.GuidFixed
	docReq.Items = doc.Items
	docReq.SumAmount = doc.SumAmount
}
