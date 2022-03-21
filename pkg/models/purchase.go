package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Purchase struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ShopID    string             `json:"shopID" bson:"shopID"`
	GuidFixed string             `json:"guidFixed,omitempty" bson:"guidFixed"`
	Items     []PurchaseDetail   `json:"items" bson:"items" `
	SumAmount float64            `json:"sumAmount" bson:"sumAmount" `
	Activity
}

func (*Purchase) CollectionName() string {
	return "purchases"
}

type PurchaseDetail struct {
	InventoryID    string  `json:"inventoryID" bson:"inventoryID"`
	ItemSku        string  `json:"itemSku,omitempty" bson:"itemSku,omitempty"`
	CategoryGuid   string  `json:"categoryGuid" bson:"categoryGuid"`
	LineNumber     int     `json:"lineNumber" bson:"lineNumber"`
	Price          float64 `json:"price" bson:"price" `
	Qty            float64 `json:"qty" bson:"qty" `
	DiscountAmount float64 `json:"discountAmount" bson:"discountAmount"`
	DiscountText   string  `json:"discountText" bson:"discountText"`
}

type PurchaseRequest struct {
	ShopID    string           `json:"shopID" `
	GuidFixed string           `json:"guidFixed,omitempty" `
	Items     []PurchaseDetail `json:"items" `
	SumAmount float64          `json:"sumAmount" `
}

func (*PurchaseRequest) IndexName() string {
	return "purchases"
}

func (docReq *PurchaseRequest) MapRequest(doc Purchase) {
	docReq.ShopID = doc.ShopID
	docReq.GuidFixed = doc.GuidFixed
	docReq.Items = doc.Items
	docReq.SumAmount = doc.SumAmount
}
