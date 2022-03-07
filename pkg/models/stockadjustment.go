package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StockAdjustment struct {
	Id         primitive.ObjectID      `json:"id" bson:"_id,omitempty"`
	MerchantId string                  `json:"merchantId" bson:"merchantId"`
	GuidFixed  string                  `json:"guidFixed,omitempty" bson:"guidFixed"`
	Items      []StockAdjustmentDetail `json:"items" bson:"items" `
	SumAmount  float64                 `json:"sumAmount" bson:"sumAmount" `
	CreatedBy  string                  `json:"createdBy" bson:"createdBy"`
	CreatedAt  time.Time               `json:"createdAt" bson:"createdAt"`
	UpdatedBy  string                  `json:"updatedBy,omitempty" bson:"updatedBy,omitempty"`
	UpdatedAt  time.Time               `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	Deleted    bool                    `json:"-" bson:"deleted"`
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
	MerchantId string                  `json:"merchantId" `
	GuidFixed  string                  `json:"guidFixed,omitempty" `
	Items      []StockAdjustmentDetail `json:"items" `
	SumAmount  float64                 `json:"sumAmount" `
}

func (*StockAdjustmentRequest) IndexName() string {
	return "stockAdjustment"
}

func (docReq *StockAdjustmentRequest) MapRequest(doc StockAdjustment) {
	docReq.MerchantId = doc.MerchantId
	docReq.GuidFixed = doc.GuidFixed
	docReq.Items = doc.Items
	docReq.SumAmount = doc.SumAmount
}
