package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BarcodeSync struct {
	ProductCode primitive.ObjectID `json:"product_code" bson:"_id`

	Barcode string `json:"barcode,omitempty" bson:"barcode,omitempty"`

	LatestUpdate time.Time `json:"latest_update,omitempty" bson:"latest_update,omitempty"`

	ImageUri string `json:"image_uri,omitempty" bson:"image_uri,omitempty"`

	Description []InventoryDescription `json:"description,omitempty" bson:"description,omitempty"`
}
