package models

type ProductBarcodeBranchRequest struct {
	Branch   ProductBarcodeBranch `json:"branch"`
	Products []string             `json:"products"`
}

type ProductBarcodeBusinessTypeRequest struct {
	BusinessType ProductBarcodeBusinessType `json:"businesstype"`
	Products     []string                   `json:"products"`
}
