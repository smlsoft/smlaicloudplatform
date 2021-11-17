package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Barcode struct {
	Id primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	Barcode string `json:"barcode,omitempty" bson:"barcode,omitempty"`

	Unit string `json:"unit,omitempty" bson:"unit,omitempty"`

	Price float32 `json:"price" bson:"price"`
}
