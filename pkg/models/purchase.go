package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Purchase struct {
	Id         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	MerchantId string             `json:"merchantId" bson:"merchantId"`
	GuidFixed  string             `json:"guidFixed,omitempty" bson:"guidFixed"`
	Items      []PurchaseDetail   `json:"items" bson:"items" `
	SumAmount  float64            `json:"sumAmount" bson:"sumAmount" `
	CreatedBy  string             `json:"createdBy" bson:"createdBy"`
	CreatedAt  time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedBy  string             `json:"updatedBy,omitempty" bson:"updatedBy,omitempty"`
	UpdatedAt  time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	Deleted    bool               `json:"-" bson:"deleted"`
}

func (*Purchase) CollectionName() string {
	return "purchases"
}

type PurchaseDetail struct {
	InventoryId    string  `json:"inventoryId" bson:"inventoryId"`
	ItemSku        string  `json:"itemSku,omitempty" bson:"itemSku,omitempty"`
	CategoryGuid   string  `json:"categoryGuid" bson:"categoryGuid"`
	LineNumber     int     `json:"lineNumber" bson:"lineNumber"`
	Price          float64 `json:"price" bson:"price" `
	Qty            float64 `json:"qty" bson:"qty" `
	DiscountAmount float64 `json:"discountAmount" bson:"discountAmount"`
	DiscountText   string  `json:"discountText" bson:"discountText"`
}

type PurchaseRequest struct {
	MerchantId string           `json:"merchantId" `
	GuidFixed  string           `json:"guidFixed,omitempty" `
	Items      []PurchaseDetail `json:"items" `
	SumAmount  float64          `json:"sumAmount" `
}

func (*PurchaseRequest) IndexName() string {
	return "purchases"
}

func (docReq *PurchaseRequest) MapRequest(doc Purchase) {
	docReq.MerchantId = doc.MerchantId
	docReq.GuidFixed = doc.GuidFixed
	docReq.Items = doc.Items
	docReq.SumAmount = doc.SumAmount
}
