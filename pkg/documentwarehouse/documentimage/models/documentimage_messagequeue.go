package models

type DocumentImageTaskChangeMessage struct {
	ShopID   string          `json:"shopid"`
	TaskGUID string          `json:"taskguid"`
	Event    TaskChangeEvent `json:"event"`
	Count    int             `json:"count"`
}

type DocumentImageTaskRejectMessage struct {
	ShopID   string          `json:"shopid"`
	TaskGUID string          `json:"taskguid"`
	Event    TaskRejectEvent `json:"event"`
	Count    int             `json:"count"`
}

type TaskChangeEvent int8

const (
	TaskChangePlus TaskChangeEvent = iota
	TaskChangeMinus
)

type TaskRejectEvent int8

const (
	TaskRejectPlus TaskRejectEvent = iota
	TaskRejectMinus
)
