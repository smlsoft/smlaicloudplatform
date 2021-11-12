package product

import "time"

type Inventory struct {
	Code     string  `json:"code" form:"code"`
	Name     string  `json:"name" form:"name"`
	ImageUri string  `json:"image_uri"`
	Barcodes Barcode `json:"barcodes"`
	LatestUpdate time.Time              `json:"latest_update"`
}


func (*Inventory) TableName() string {
	return "inventory"
}