package models

import "time"

type Trans struct {
	Id          uint      `json:"id" gorm:"column:id;primaryKey"`
	ParId       string    `json:"parid" gorm:"column:parid"`
	DBId        string    `json:"dbid" gorm:"column:dbid;uniqueIndex:trans_unique_dockey"`
	DocDate     time.Time `json:"docdate" gorm:"column:docdate"`
	DocNumber   string    `json:"docnumber" gorm:"column:docnumber;uniqueIndex:trans_unique_dockey"`
	TotalAmount float64   `json:"totalamount" gorm:"column:totalamount"`
	TotalQty    float64   `json:"totalqty" gorm:"column:totalqty"`
}

func (Trans) TableName() string {
	return "trans"
}

type TransItemDetail struct {
	Id              uint      `json:"id" gorm:"column:id;primaryKey"`
	ParId           string    `json:"parid" gorm:"column:parid"`
	DBId            string    `json:"dbid" gorm:"column:dbid;uniqueIndex:transitemdetail_unique_dockey"` //
	TransType       int8      `json:"transtype" gorm:"column:transtype"`
	CalcType        int8      `json:"calctype" gorm:"column:calctype"`
	DocDate         time.Time `json:"docdate" gorm:"column:docdate"`
	DocNumber       string    `json:"docnumber" gorm:"column:docnumber;uniqueIndex:transitemdetail_unique_dockey"`   //
	LineNumber      int8      `json:"linenumber" gorm:"column:linenumber;uniqueIndex:transitemdetail_unique_dockey"` //
	Barcode         string    `json:"barcode" gorm:"column:barcode"`
	ItemCode        string    `json:"itemcode" gorm:"column:itemcode"`
	UnitCode        string    `json:"unitcode" gorm:"column:unitcode"`
	Qty             float64   `json:"qty" gorm:"column:qty"`
	Price           float64   `json:"price" gorm:"column:price"`
	DiscWord        string    `json:"discword" gorm:"column:discword"`
	DiscAmount      float64   `json:"discamount" gorm:"column:discamount"`
	SumCost         float64   `json:"sumamount" gorm:"column:sumamount"`
	VatType         int8      `json:"vattype" gorm:"column:vattype"`
	TotalAmount     float64   `json:"totalamount" gorm:"column:totalamount"`
	UnitStd         float64   `json:"unitstd" gorm:"column:unitstd"`
	UnitDiv         float64   `json:"unitdiv" gorm:"column:unitdiv"`
	QtyExtra        float64   `json:"-" gorm:"column:qtyextra"`
	SumCostExtra    float64   `json:"-" gorm:"column:sumcostextra"`
	IsUpdated       bool      `json:"isupdated" gorm:"column:isupdated"`
	BalanceAmount   float64   `json:"balanceamount" gorm:"column:balanceamount"`
	BalanceQty      float64   `json:"balanceqty" gorm:"column:balanceqty"`
	AvgCost         float64   `json:"avgcost" gorm:"column:avgcost"`
	WhCode          string    `json:"whcode" gorm:"column:whcode"`
	LocCode         string    `json:"loccode" gorm:"column:loccode"`
	ParIdForCalc    string    `json:"paridforcalc" gorm:"column:paridforcalc"`
	DBIdForCal      string    `json:"dbidforcal" gorm:"column:dbidforcal"`
	ItemCodeForCalc string    `json:"itemcodeforcalc" gorm:"column:itemcodeforcalc"`
}

func (TransItemDetail) TableName() string {
	return "item_trans"
}
