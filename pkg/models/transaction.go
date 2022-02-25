package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	Id         primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	MerchantId string              `json:"merchantId" bson:"merchantId"`
	GuidFixed  string              `json:"guidFixed,omitempty" bson:"guidFixed"`
	Items      []TransactionDetail `json:"items" bson:"items" `
	SumAmount  float64             `json:"sumAmount" bson:"sumAmount" `

	CreatedBy string    `json:"-" bson:"createdBy"`
	CreatedAt time.Time `json:"-" bson:"createdAt"`
	UpdatedBy string    `json:"-" bson:"updatedBy,omitempty"`
	UpdatedAt time.Time `json:"-" bson:"updatedAt,omitempty"`
	Deleted   bool      `json:"-" bson:"deleted"`
}

func (*Transaction) CollectionName() string {
	return "transactions"
}

type TransactionDetail struct {
	InventoryId    string  `json:"inventoryId" bson:"inventoryId"`
	ItemSku        string  `json:"itemSku,omitempty" bson:"itemSku,omitempty"`
	CategoryGuid   string  `json:"categoryGuid" bson:"categoryGuid"`
	LineNumber     int     `json:"lineNumber" bson:"lineNumber"`
	Price          float64 `json:"price" bson:"price" `
	Qty            float64 `json:"qty" bson:"qty" `
	DiscountAmount float64 `json:"discountAmount" bson:"discountAmount"`
	DiscountText   string  `json:"discountText" bson:"discountText"`
}
