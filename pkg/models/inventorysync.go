package models

import (
	"time"
)

type BarcodeSync struct {

	ProductCode string `json:"product_code,omitempty"`

	Barcode string `json:"barcode,omitempty"`

	LatestUpdate time.Time `json:"latest_update,omitempty"`

	ImageUri string `json:"image_uri,omitempty"`

	Description []InventoryDescription `json:"description,omitempty"`
}
