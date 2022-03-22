package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const stockInOutCollectionName = "stockInOuts"
const stockInOutIndexName = "stockInOuts"

type StockInOut struct {
	Items     []StockInOutDetail `json:"items" bson:"items" `
	SumAmount float64            `json:"sumAmount" bson:"sumAmount" `
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

type StockInOutInfo struct {
	DocIdentity
	StockInOut
}

func (StockInOutInfo) CollectionName() string {
	return stockInOutCollectionName
}

type StockInOutData struct {
	ShopIdentity
	StockInOutInfo
}

func (StockInOutData) IndexName() string {
	return stockInOutIndexName
}

type StockInOutDoc struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	StockInOutData
	Activity
}

func (StockInOutDoc) CollectionName() string {
	return stockInOutCollectionName
}
