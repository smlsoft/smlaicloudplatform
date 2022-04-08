package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const transactionCollectionName = "transactions"
const transactionIndexName = "transactions"

type Transaction struct {
	Items     *[]TransactionDetail `json:"items" bson:"items" `
	SumAmount float64              `json:"sumamount" bson:"sumamount" `
}

type TransactionDetail struct {
	InventoryID    string  `json:"inventoryid" bson:"inventoryid"`
	ItemSku        string  `json:"itemsku,omitempty" bson:"itemsku,omitempty"`
	CategoryGuid   string  `json:"categoryguid" bson:"categoryguid"`
	LineNumber     int     `json:"linenumber" bson:"linenumber"`
	Price          float64 `json:"price" bson:"price" `
	Qty            float64 `json:"qty" bson:"qty" `
	DiscountAmount float64 `json:"discountamount" bson:"discountamount"`
	DiscountText   string  `json:"discounttext" bson:"discounttext"`
}

type TransactionInfo struct {
	DocIdentity `bson:"inline"`
	Transaction `bson:"inline"`
}

func (TransactionInfo) CollectionName() string {
	return transactionCollectionName
}

type TransactionData struct {
	ShopIdentity    `bson:"inline"`
	TransactionInfo `bson:"inline"`
}

func (TransactionData) IndexName() string {
	return transactionIndexName
}

type TransactionDoc struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	TransactionData `bson:"inline"`
	Activity        `bson:"inline"`
}

func (TransactionDoc) CollectionName() string {
	return transactionCollectionName
}
