package models

import "smlcloudplatform/pkg/models"

type ProductChoice struct {
	GUID           string          `json:"guid" bson:"guid"`
	RefBarcode     string          `json:"refbarcode" bson:"refbarcode"`
	RefProductCode string          `json:"refproductcode" bson:"refproductcode"`
	Names          *[]models.NameX `json:"names" bson:"names" validate:"required,min=1,unique=Code,dive"`
	RefUnitCode    string          `json:"refunitcode" bson:"refunitcode"`
	RefUnitNames   *[]models.NameX `json:"refunitnames" bson:"refunitnames" validate:"required,min=1,unique=Code,dive"`
	Price          *string         `json:"price" bson:"price"`
	Qty            float64         `json:"qty" bson:"qty"`
	// QtyMax         float64         `json:"qtymax" bson:"qtymax"`
	IsStock   bool `json:"isstock" bson:"isstock"`
	IsDefault bool `json:"isdefault" bson:"isdefault"`
}
