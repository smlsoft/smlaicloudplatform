package models

type Pageable struct {
	Q     string         `json:"q,omitempty"`
	Page  int            `json:"page"`
	Limit int            `json:"limit,omitempty"`
	Sorts map[string]int `json:"sorts,omitempty"`
}
