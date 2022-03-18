package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Purchase struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ShopId    string             `json:"shopId" bson:"shopId"`
	GuidFixed string             `json:"guidFixed,omitempty" bson:"guidFixed"`
	Items     []PurchaseDetail   `json:"items" bson:"items" `
	SumAmount float64            `json:"sumAmount" bson:"sumAmount" `
	Activity
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
	ShopId    string           `json:"shopId" `
	GuidFixed string           `json:"guidFixed,omitempty" `
	Items     []PurchaseDetail `json:"items" `
	SumAmount float64          `json:"sumAmount" `
}

func (*PurchaseRequest) IndexName() string {
	return "purchases"
}

func (docReq *PurchaseRequest) MapRequest(doc Purchase) {
	docReq.ShopId = doc.ShopId
	docReq.GuidFixed = doc.GuidFixed
	docReq.Items = doc.Items
	docReq.SumAmount = doc.SumAmount
}
