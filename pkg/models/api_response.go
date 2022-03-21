package models

type ApiResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`
	ID         interface{} `json:"id,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

type AuthResponseFailed struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

type ResponseSuccessWithID struct {
	Success bool        `json:"success"`
	ID      interface{} `json:"id,omitempty"`
}
