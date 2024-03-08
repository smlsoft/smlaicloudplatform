package models

import "smlcloudplatform/internal/models"

type ProductBarcodeBOM struct {
	Level         int             `json:"level" gorm:"column:level"`
	MainBarcode   string          `json:"mainbarcode" gorm:"column:mainbarcode"`
	GuidFixed     string          `json:"guidfixed" gorm:"column:guidfixed"`
	Names         *[]models.NameX `json:"names" gorm:"column:names"`
	ItemUnitCode  string          `json:"itemunitcode" gorm:"column:itemunitcode"`
	ItemUnitNames *[]models.NameX `json:"itemunitnames" gorm:"column:itemunitnames"`
	Barcode       string          `json:"barcode" gorm:"column:barcode"`
	Condition     bool            `json:"condition" gorm:"column:condition"`
	DivideValue   float64         `json:"dividevalue" gorm:"column:dividevalue"`
	StandValue    float64         `json:"standvalue" gorm:"column:standvalue"`
	Qty           float64         `json:"qty" gorm:"column:qty"`
}
