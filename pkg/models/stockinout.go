package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const stockInOutCollectionName = "stockInOuts"
const stockInOutIndexName = "stockInOuts"

type StockInOut struct {
	Items     *[]StockInOutDetail `json:"items" bson:"items" `
	SumAmount float64             `json:"sumamount" bson:"sumamount" `
}

type StockInOutDetail struct {
	InventoryID    string  `json:"inventoryid" bson:"inventoryid"`
	ItemSku        string  `json:"itemsku,omitempty" bson:"itemsku,omitempty"`
	CategoryGuid   string  `json:"categoryguid" bson:"categoryguid"`
	LineNumber     int     `json:"linenumber" bson:"linenumber"`
	Price          float64 `json:"price" bson:"price" `
	Qty            float64 `json:"qty" bson:"qty" `
	DiscountAmount float64 `json:"discountamount" bson:"discountamount"`
	DiscountText   string  `json:"discounttext" bson:"discounttext"`
}

type StockInOutInfo struct {
	DocIdentity `bson:"inline"`
	StockInOut  `bson:"inline"`
}

func (StockInOutInfo) CollectionName() string {
	return stockInOutCollectionName
}

type StockInOutData struct {
	ShopIdentity   `bson:"inline"`
	StockInOutInfo `bson:"inline"`
}

func (StockInOutData) IndexName() string {
	return stockInOutIndexName
}

type StockInOutDoc struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	StockInOutData `bson:"inline"`
	Activity       `bson:"inline"`
}

func (StockInOutDoc) CollectionName() string {
	return stockInOutCollectionName
}
