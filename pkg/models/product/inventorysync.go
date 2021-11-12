package product

import "time"

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
