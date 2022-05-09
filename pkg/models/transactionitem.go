package models

import "time"

type Trans struct {
	Id          uint
	ParId       string
	DBId        string
	DocDate     time.Time
	DocNumber   string
	TotalAmount float64
	TotalQty    float64
}

func (Trans) TableName() string {
	return "trans"
}

type TransItemDetail struct {
	Id              uint
	ParId           string
	DBId            string
	TransType       int8
	CalcType        int8
	DocDate         time.Time
	DocNumber       string
	LineNumber      int8
	Barcode         string
	ItemCode        string
	UnitCode        string
	Qty             float64
	Price           float64
	DiscWord        string
	DiscAmount      float64
	SumCost         float64
	VatType         int8
	TotalAmount     float64
	UnitStd         float64
	UnitDiv         float64
	QtyExtra        float64
	SumCostExtra    float64
	IsUpdated       bool
	BalanceAmount   float64
	BalanceQty      float64
	AvgCost         float64
	WhCode          string
	LocCode         string
	ParIdForCalc    string
	DBIdForCal      string
	ItemCodeForCalc string
}

func (TransItemDetail) TableName() string {
	return "item_trans"
}
