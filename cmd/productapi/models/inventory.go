package models

import "time"

type Inventory struct {
	Code     string  `json:"code" form:"code"`
	Name     string  `json:"name" form:"name"`
	ImageUri string  `json:"image_uri"`
	Barcodes Barcode `json:"barcodes"`
	LatestUpdate time.Time              `json:"latest_update"`
}

type Barcode struct {
	Barcode  string `json:"barcode"`
	UnitCode string `json:"unit_code"`
}

func (*Inventory) TableName() string {
	return "inventory"
}