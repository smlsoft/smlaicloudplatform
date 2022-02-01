package models

type ApiResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`
	Id         interface{} `json:"id,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
}
