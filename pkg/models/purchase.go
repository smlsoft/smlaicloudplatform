package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const purchaseCollectionName = "purchases"
const purchaseIndexName = "purchases"

type Purchase struct {
	Items     *[]PurchaseDetail `json:"items" bson:"items" `
	SumAmount float64           `json:"sumAmount" bson:"sumAmount" `
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

type PurchaseInfo struct {
	DocIdentity `bson:"inline"`
	Purchase    `bson:"inline"`
}

func (PurchaseInfo) CollectionName() string {
	return purchaseCollectionName
}

type PurchaseData struct {
	ShopIdentity `bson:"inline"`
	PurchaseInfo `bson:"inline"`
}

func (PurchaseData) IndexName() string {
	return purchaseIndexName
}

type PurchaseDoc struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PurchaseData `bson:"inline"`
	Activity     `bson:"inline"`
}

func (PurchaseDoc) CollectionName() string {
	return purchaseCollectionName
}
