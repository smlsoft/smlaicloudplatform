package models

type ApiResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message,omitempty"`
	Id         interface{} `json:"id,omitempty"`
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

type ResponseSuccessWithId struct {
	Success bool        `json:"success"`
	Id      interface{} `json:"id,omitempty"`
}
