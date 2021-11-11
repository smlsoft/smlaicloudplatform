package product

import "time"

type Inventory struct {
	Code     string  `json:"code" form:"code"`
	Name     string  `json:"name" form:"name"`
	ImageUri string  `json:"image_uri"`
	Barcodes Barcode `json:"barcodes"`
}

type Barcode struct {
	Barcode  string `json:"barcode"`
	UnitCode string `json:"unit_code"`
}

type BarcodeSync struct {
	Barcode      string                 `json:"barcode"`
	LatestUpdate time.Time              `json:"latest_update"`
	ImageUri     string                 `json:"image_uri"`
	Description  BarcodeSyncDescription `json:"description"`
}

type BarcodeSyncDescription struct {
	Language    string `json:"lang"`
	ProductName string `json:"product_name"`
	UnitName    string `json:"unit_name"`
}
