package models

type BulkImport struct {
	Created          []string `json:"created"`
	Updated          []string `json:"updated"`
	UpdateFailed     []string `json:"updateFailed"`
	PayloadDuplicate []string `json:"payloadDuplicate"`
}

type BulkResponse struct {
	Success bool `json:"success"`
	BulkImport
}
