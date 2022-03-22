package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const transactionCollectionName = "transactions"
const transactionIndexName = "transactions"

type Transaction struct {
	Items     []TransactionDetail `json:"items" bson:"items" `
	SumAmount float64             `json:"sumAmount" bson:"sumAmount" `
}

type TransactionDetail struct {
	InventoryID    string  `json:"inventoryID" bson:"inventoryID"`
	ItemSku        string  `json:"itemSku,omitempty" bson:"itemSku,omitempty"`
	CategoryGuid   string  `json:"categoryGuid" bson:"categoryGuid"`
	LineNumber     int     `json:"lineNumber" bson:"lineNumber"`
	Price          float64 `json:"price" bson:"price" `
	Qty            float64 `json:"qty" bson:"qty" `
	DiscountAmount float64 `json:"discountAmount" bson:"discountAmount"`
	DiscountText   string  `json:"discountText" bson:"discountText"`
}

type TransactionInfo struct {
	DocIdentity
	Transaction
}

func (TransactionInfo) CollectionName() string {
	return transactionCollectionName
}

type TransactionData struct {
	ShopIdentity
	TransactionInfo
}

func (TransactionData) IndexName() string {
	return transactionIndexName
}

type TransactionDoc struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	TransactionData
	Activity
}

func (TransactionDoc) CollectionName() string {
	return transactionCollectionName
}
