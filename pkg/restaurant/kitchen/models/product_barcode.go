package models

import "smlcloudplatform/pkg/models"

type ProductBarcode struct {
	Barcode  string           `json:"barcode" bson:"barcode"`
	Kitchens []KitchenBarcode `json:"kitchens" bson:"kitchens"`
}

type KitchenBarcode struct {
	Code  string          `json:"code" bson:"code"`
	Names *[]models.NameX `json:"names" bson:"names" validate:"required,min=1,unique=Code,dive"`
}
