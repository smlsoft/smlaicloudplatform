package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Inventory struct {
	ProductId primitive.ObjectID `json:"-" bson:"_id,omitempty"`

	ItemSku string `json:"item_sku,omitempty" bson:"item_sku,omitempty"`

	ProductName string `json:"product_name,omitempty" bson:"product_name,omitempty"`

	Barcodes []Barcode `json:"barcodes,omitempty" bson:"barcodes,omitempty"`

	Pictures []Picture `json:"pictures,omitempty" bson:"pictures,omitempty"`

	LatestUpdate time.Time `json:"-" bson:"latest_update,omitempty"`
}

func (*Inventory) CollectionName() string {
	return "inventory"
}

type InventoryDescription struct {
	Lang string `json:"lang,omitempty" bson:"lang,omitempty"`

	ProductName string `json:"product_name,omitempty" bson:"product_name,omitempty"`

	UnitName string `json:"unit_name,omitempty" bson:"unit_name,omitempty"`
}
