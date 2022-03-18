package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	Id        primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	ShopId    string              `json:"shopId" bson:"shopId"`
	GuidFixed string              `json:"guidFixed,omitempty" bson:"guidFixed"`
	Items     []TransactionDetail `json:"items" bson:"items" `
	SumAmount float64             `json:"sumAmount" bson:"sumAmount" `
	Activity
}

// CreatedBy  string              `json:"-" bson:"createdBy"`
// CreatedAt  time.Time           `json:"-" bson:"createdAt"`
// UpdatedBy  string              `json:"-" bson:"updatedBy,omitempty"`
// UpdatedAt  time.Time           `json:"-" bson:"updatedAt,omitempty"`

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

type TransactionRequest struct {
	ShopId    string              `json:"shopId" `
	GuidFixed string              `json:"guidFixed,omitempty" `
	Items     []TransactionDetail `json:"items" `
	SumAmount float64             `json:"sumAmount" `

	// CreatedBy string    `json:"createdBy" bson:"createdBy"`
	// CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	// UpdatedBy string    `json:"updatedBy" bson:"updatedBy,omitempty"`
	// UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt,omitempty"`
	// Deleted   bool      `json:"deleted" bson:"deleted"`
}

func (*TransactionRequest) IndexName() string {
	return "transaction"
}

func (transReq *TransactionRequest) MapRequest(trans Transaction) {
	transReq.ShopId = trans.ShopId
	transReq.GuidFixed = trans.GuidFixed
	transReq.Items = trans.Items
	transReq.SumAmount = trans.SumAmount
	// transReq.CreatedBy = trans.CreatedBy
	// transReq.CreatedAt = trans.CreatedAt
	// transReq.UpdatedBy = trans.UpdatedBy
	// transReq.UpdatedAt = trans.UpdatedAt
}
