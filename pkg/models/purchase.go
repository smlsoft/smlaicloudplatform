package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const purchaseCollectionName = "purchases"
const purchaseIndexName = "purchases"

type Purchase struct {
	Items     *[]PurchaseDetail `json:"items" bson:"items" `
	SumAmount float64           `json:"sumamount" bson:"sumamount" `
}

type PurchaseDetail struct {
	InventoryID    string  `json:"inventoryid" bson:"inventoryid"`
	ItemSku        string  `json:"itemsku,omitempty" bson:"itemsku,omitempty"`
	CategoryGuid   string  `json:"categoryguid" bson:"categoryguid"`
	LineNumber     int     `json:"linenumber" bson:"linenumber"`
	Price          float64 `json:"price" bson:"price" `
	Qty            float64 `json:"qty" bson:"qty" `
	DiscountAmount float64 `json:"discountamount" bson:"discountamount"`
	DiscountText   string  `json:"discounttext" bson:"discounttext"`
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

type PurchaseListPageResponse struct {
	Success    bool                   `json:"success"`
	Data       []PurchaseInfo         `json:"data,omitempty"`
	Pagination PaginationDataResponse `json:"pagination,omitempty"`
}
