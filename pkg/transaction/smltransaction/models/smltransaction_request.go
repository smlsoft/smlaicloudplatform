package models

type SMLTransactionRequest struct {
	Collection string                 `json:"collection" validate:"required,min=1,max=100"`
	KeyName    string                 `json:"keyname" validate:"required,min=1,max=100"`
	Body       map[string]interface{} `json:"body" validate:"required"`
}

type SMLTransactionBulkRequest struct {
	Collection string                   `json:"collection" validate:"required,min=1,max=100"`
	KeyName    string                   `json:"keyname" validate:"required,min=1,max=100"`
	Body       []map[string]interface{} `json:"body" validate:"required"`
}

type SMLTransactionKeyRequest struct {
	Collection string   `json:"collection" validate:"required,min=1,max=100"`
	KeyName    string   `json:"keyname" validate:"required,min=1,max=100"`
	Keys       []string `json:"keys" validate:"required"`
}
