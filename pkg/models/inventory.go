package models

import (
	"time"
)

type Inventory struct {

	ProductId string `json:"product_id,omitempty"`

	ItemSku string `json:"item_sku,omitempty"`

	ProductName string `json:"product_name,omitempty"`

	Barcodes []Barcode `json:"barcodes,omitempty"`

	Pictures []Picture `json:"pictures,omitempty"`

	LatestUpdate time.Time `json:"latest_update,omitempty"`
}


type InventoryDescription struct {

	Lang string `json:"lang,omitempty"`

	ProductName string `json:"product_name,omitempty"`

	UnitName string `json:"unit_name,omitempty"`
}
