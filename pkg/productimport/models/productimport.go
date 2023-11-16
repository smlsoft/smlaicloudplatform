package models

import (
	"math"
	"time"
)

type ProductImportRaw struct {
	Barcode     string  `json:"barcode" ch:"barcode"`
	Name        string  `json:"name" ch:"name"`
	UnitCode    string  `json:"unitcode" ch:"unitcode"`
	Price       float64 `json:"price" ch:"price"`
	PriceMember float64 `json:"pricemember" ch:"pricemember"`
	IsDuplicate bool    `json:"isduplicate" ch:"isduplicate"`
	IsExist     bool    `json:"isexist" ch:"isexist"`
}

type ProductImport struct {
	TaskID    string  `json:"taskid" ch:"taskid"`
	RowNumber float64 `json:"rownumber" ch:"rownumber"`
	ProductImportRaw
}

type ProductImportInfo struct {
	GUIDFixed string `json:"guidfixed" ch:"guidfixed"`
	ShopID    string `json:"shopid" ch:"shopid"`
	ProductImport
}

type ProductImportDoc struct {
	ProductImportInfo
	CreatedAt time.Time `json:"createdat" ch:"createdat"`
	CreatedBy string    `json:"createdby" ch:"createdby"`
}

func (ProductImportDoc) TableName() string {
	return "productbarcodeimport"
}

type TaskStatus int8

const (
	TaskStatusPending TaskStatus = iota
	TaskStatusProcessing
	TaskStatusDone
	TaskStatusError
	TaskStatusSaveSucceded
	TaskStatusSaveFailed
	TaskStatusNotFound
)

type PaginationData struct {
	Total     int64 `json:"total"`
	Page      int64 `json:"page"`
	PerPage   int64 `json:"perPage"`
	Prev      int64 `json:"prev"`
	Next      int64 `json:"next"`
	TotalPage int64 `json:"totalPage"`
}

func (p *PaginationData) Build() {
	totalPage := math.Ceil(float64(p.Total) / float64(p.PerPage))
	p.TotalPage = int64(totalPage)

	if p.Page == 0 {
		p.Page = 1
	}

	if p.Page > 1 {
		p.Prev = p.Page - 1
	}

	if p.Page < p.TotalPage {
		p.Next = p.Page + 1
	}
}
