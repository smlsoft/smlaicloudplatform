package models

import (
	common "smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const stockAdjustmentCollectionName = "stockAdjustments"
const stockAdjustmentIndexName = "stockAdjustments"

type StockAdjustment struct {
	Items     *[]StockAdjustmentDetail `json:"items" bson:"items" `
	SumAmount float64                  `json:"sumamount" bson:"sumamount" `
}

type StockAdjustmentDetail struct {
	InventoryID    string  `json:"inventoryid" bson:"inventoryid"`
	ItemSku        string  `json:"itemsku,omitempty" bson:"itemsku,omitempty"`
	CategoryGuid   string  `json:"categoryguid" bson:"categoryguid"`
	LineNumber     int     `json:"linenumber" bson:"linenumber"`
	Price          float64 `json:"price" bson:"price" `
	Qty            float64 `json:"qty" bson:"qty" `
	DiscountAmount float64 `json:"discountamount" bson:"discountamount"`
	DiscountText   string  `json:"discounttext" bson:"discounttext"`
}

type StockAdjustmentInfo struct {
	common.DocIdentity `bson:"inline"`
	StockAdjustment    `bson:"inline"`
}

func (StockAdjustmentInfo) CollectionName() string {
	return stockAdjustmentCollectionName
}

type StockAdjustmentData struct {
	common.ShopIdentity `bson:"inline"`
	StockAdjustmentInfo `bson:"inline"`
}

func (StockAdjustmentData) IndexName() string {
	return stockAdjustmentIndexName
}

type StockAdjustmentDoc struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	StockAdjustmentData `bson:"inline"`
	common.ActivityDoc  `bson:"inline"`
}

func (StockAdjustmentDoc) CollectionName() string {
	return stockAdjustmentCollectionName
}
