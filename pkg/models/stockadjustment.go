package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const stockAdjustmentCollectionName = "stockAdjustments"
const stockAdjustmentIndexName = "stockAdjustments"

type StockAdjustment struct {
	Items     []StockAdjustmentDetail `json:"items" bson:"items" `
	SumAmount float64                 `json:"sumAmount" bson:"sumAmount" `
}

type StockAdjustmentDetail struct {
	InventoryID    string  `json:"inventoryID" bson:"inventoryID"`
	ItemSku        string  `json:"itemSku,omitempty" bson:"itemSku,omitempty"`
	CategoryGuid   string  `json:"categoryGuid" bson:"categoryGuid"`
	LineNumber     int     `json:"lineNumber" bson:"lineNumber"`
	Price          float64 `json:"price" bson:"price" `
	Qty            float64 `json:"qty" bson:"qty" `
	DiscountAmount float64 `json:"discountAmount" bson:"discountAmount"`
	DiscountText   string  `json:"discountText" bson:"discountText"`
}

type StockAdjustmentInfo struct {
	DocIdentity
	StockAdjustment
}

func (StockAdjustmentInfo) CollectionName() string {
	return stockAdjustmentCollectionName
}

type StockAdjustmentData struct {
	ShopIdentity
	StockAdjustmentInfo
}

func (StockAdjustmentData) IndexName() string {
	return stockAdjustmentIndexName
}

type StockAdjustmentDoc struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	StockAdjustmentData
	Activity
}

func (StockAdjustmentDoc) CollectionName() string {
	return stockAdjustmentCollectionName
}
