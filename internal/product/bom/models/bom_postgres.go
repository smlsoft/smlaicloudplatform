package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"smlcloudplatform/internal/models"
)

type BOMProductBarcodePg struct {
	BarcodeGuidFixed string       `json:"guidfixed" gorm:"column:guidfixed"`
	Level            int          `json:"level" gorm:"column:level"`
	Names            models.JSONB `json:"names" gorm:"column:names;type:jsonb"`
	ItemUnitCode     string       `json:"itemunitcode" gorm:"column:itemunitcode"`
	ItemUnitNames    models.JSONB `json:"itemunitnames" gorm:"column:itemunitnames;type:jsonb"`
	Barcode          string       `json:"barcode" gorm:"column:barcode" validate:"required,min=1"`
	Condition        bool         `json:"condition" gorm:"column:condition"`
	DivideValue      float64      `json:"dividevalue" gorm:"column:dividevalue"`
	StandValue       float64      `json:"standvalue" gorm:"column:standvalue"`
	Qty              float64      `json:"qty" gorm:"column:qty"`
}

type ProductBarcodeBOMViewPG struct {
	BOMProductBarcode `gorm:"embedded;"`
	ImageURI          string    `json:"imageuri" gor:"column:imageuri"`
	BOM               BOMViewPg `json:"bom" gor:"column:bom"`
}

func (b *ProductBarcodeBOMViewPG) TableName() string {
	return "productBarcodeBOMs"
}

type BOMViewPg []ProductBarcodeBOMViewPG

// Value Marshal
func (a BOMViewPg) Value() (driver.Value, error) {

	j, err := json.Marshal(a)
	return j, err
	//return json.Marshal(a)
}

// Scan Unmarshal
func (a *BOMViewPg) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
