package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const saleinvoiceCollectionName = "saleinvoices"
const saleinvoiceIndexName = "saleinvoices"

type Saleinvoice struct {
	Items     *[]SaleinvoiceDetail `json:"items" bson:"items" `
	SumAmount float64              `json:"sumamount" bson:"sumamount" `
}

type SaleinvoiceDetail struct {
	InventoryID    string  `json:"inventoryid" bson:"inventoryid"`
	ItemSku        string  `json:"itemsku,omitempty" bson:"itemsku,omitempty"`
	CategoryGuid   string  `json:"categoryguid" bson:"categoryguid"`
	LineNumber     int     `json:"linenumber" bson:"linenumber"`
	Price          float64 `json:"price" bson:"price" `
	Qty            float64 `json:"qty" bson:"qty" `
	DiscountAmount float64 `json:"discountamount" bson:"discountamount"`
	DiscountText   string  `json:"discounttext" bson:"discounttext"`
}

type SaleinvoiceInfo struct {
	DocIdentity `bson:"inline"`
	Saleinvoice `bson:"inline"`
}

func (SaleinvoiceInfo) CollectionName() string {
	return saleinvoiceCollectionName
}

type SaleinvoiceData struct {
	ShopIdentity    `bson:"inline"`
	SaleinvoiceInfo `bson:"inline"`
}

func (SaleinvoiceData) IndexName() string {
	return saleinvoiceIndexName
}

type SaleinvoiceDoc struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SaleinvoiceData `bson:"inline"`
	ActivityDoc     `bson:"inline"`
}

func (SaleinvoiceDoc) CollectionName() string {
	return saleinvoiceCollectionName
}
