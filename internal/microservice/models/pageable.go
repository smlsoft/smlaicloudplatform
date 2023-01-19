package models

type Pageable struct {
	Query string   `json:"q,omitempty"`
	Page  int      `json:"page"`
	Limit int      `json:"limit,omitempty"`
	Sorts []KeyInt `json:"sorts,omitempty"`
}

type PageableStep struct {
	Query string   `json:"q,omitempty"`
	Skip  int      `json:"skip,omitempty"`
	Limit int      `json:"limit,omitempty"`
	Sorts []KeyInt `json:"sorts,omitempty"`
}

type KeyInt struct {
	Key   string
	Value int8
}
