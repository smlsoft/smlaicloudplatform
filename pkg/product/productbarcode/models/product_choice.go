package models

import "smlcloudplatform/pkg/models"

type ProductChoice struct {
	RefBarcode  string          `json:"refbarcode" bson:"refbarcode"`
	RefUnitCode string          `json:"refunitcode" bson:"refunitcode"`
	Names       *[]models.NameX `json:"names" bson:"names" validate:"required,min=1,unique=Code,dive"`
	Price       float64         `json:"price" bson:"price"`
	Qty         float64         `json:"qty" bson:"qty"`
	QtyMax      float64         `json:"qtymax" bson:"qtymax"`
}
