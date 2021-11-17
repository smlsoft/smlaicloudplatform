package models

type Barcode struct {

	Id string `json:"id,omitempty"`

	Barcode string `json:"barcode,omitempty"`

	Unit string `json:"unit,omitempty"`

	Price float32 `json:"price,omitempty"`
}