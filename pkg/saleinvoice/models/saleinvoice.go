package models

import (
	"smlcloudplatform/pkg/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const saleinvoiceCollectionName = "saleinvoices"
const saleinvoiceIndexName = "saleinvoices"

type Saleinvoice struct {
	DocDate        *time.Time           `json:"docdate,omitempty" bson:"docdate,omitempty" format:"date-time" examples:"2019-10-12T07:20:50.52Z or 2019-10-12T07:20:50.52+07:00" gorm:"docdate,type:date"`
	DocNo          string               `json:"docno,omitempty"  bson:"docno,omitempty"`
	Member         *models.Member       `json:"member,omitempty"  bson:"member,omitempty"`
	Items          *[]SaleinvoiceDetail `json:"items" bson:"items" gorm:""`
	TotalAmount    float64              `json:"totalamount" bson:"totalamount" `
	TaxRate        float64              `json:"taxrate" bson:"taxrate" `
	TaxAmount      float64              `json:"taxamount" bson:"taxamount" `
	TaxBaseAmount  float64              `json:"taxbaseamount" bson:"taxbaseamount" `
	DiscountAmount float64              `json:"discountamount" bson:"discountamount" `
	SumAmount      float64              `json:"sumamount" bson:"sumamount" `
	Payment        models.Payment       `json:"payment" bson:"payment"`
	LocalDate      string               `json:"localdate" bson:"localdate"`
	LocalTime      string               `json:"localtime" bson:"localtime"`
}

type SaleinvoiceDetail struct {
	LineNumber           int `json:"linenumber" bson:"linenumber"`
	models.InventoryInfo `bson:"inline" gorm:"embedded;"`
	Price                float64 `json:"price" bson:"price"`
	Qty                  float64 `json:"qty" bson:"qty" `
	DiscountAmount       float64 `json:"discountamount" bson:"discountamount"`
	DiscountText         string  `json:"discounttext" bson:"discounttext"`
	Amount               float64 `json:"amount" bson:"amount"`
}

type SaleinvoiceInfo struct {
	models.DocIdentity `bson:"inline"`
	Saleinvoice        `bson:"inline"`
}

func (SaleinvoiceInfo) CollectionName() string {
	return saleinvoiceCollectionName
}

type SaleinvoiceData struct {
	models.ShopIdentity `bson:"inline"`
	SaleinvoiceInfo     `bson:"inline"`
}

func (SaleinvoiceData) IndexName() string {
	return saleinvoiceIndexName
}

type SaleinvoiceDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SaleinvoiceData    `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (SaleinvoiceDoc) CollectionName() string {
	return saleinvoiceCollectionName
}

type SaleinvoiceListPageResponse struct {
	Success    bool                          `json:"success"`
	Data       []SaleinvoiceInfo             `json:"data,omitempty"`
	Pagination models.PaginationDataResponse `json:"pagination,omitempty"`
}

type SaleinvoicePg struct {
	Id                  uint `gorm:"primaryKey"`
	models.ShopIdentity `gorm:"embedded"`
	DocDate             time.Time             `json:"docdate,omitempty" gorm:"column:docdate"`
	DocNo               string                `json:"docno,omitempty"  gorm:"column:docno"`
	TaxRate             float64               `json:"taxrate" gorm:"column:taxrate;" `
	TaxAmount           float64               `json:"taxamount" gorm:"column:taxamount;" `
	TaxBaseAmount       float64               `json:"taxbaseamount" gorm:"column:taxbaseamount;" `
	DiscountAmount      float64               `json:"discountamount" gorm:"column:discountamount;" `
	SumAmount           float64               `json:"sumamount" gorm:"column:sumamount;" `
	Items               []SaleinvoiceDetailPg `json:"items" bson:"items" gorm:"foreignKey:SaleinvoiceId" `
}

func (*SaleinvoicePg) TableName() string {
	return "saleinvoice"
}

type SaleinvoiceDetailPg struct {
	Id                  uint `gorm:"primaryKey"`
	models.ShopIdentity `gorm:"embedded"`
	SaleinvoiceId       uint    `json:"saleinvoiceid,omitempty"  gorm:"column:saleinvoiceid;foreignKey:Id"`
	DocNo               string  `json:"docno,omitempty"  gorm:"column:docno"`
	LineNumber          int     `json:"linenumber" gorm:"column:linenumber"`
	ItemGuid            string  `json:"itemguid,omitempty" gorm:"column:itemguid"`
	ItemCode            string  `json:"itemcode,omitempty" gorm:"column:itemcode"`
	Barcode             string  `json:"barcode" gorm:"column:barcode"`
	Name1               string  `json:"name1" gorm:"name1"` // ชื่อภาษาไทย
	ItemUnitCode        string  `json:"itemunitcode,omitempty" gorm:"column:itemunitcode"`
	ItemUnitStd         float64 `json:"itemunitstd,omitempty" gorm:"column:itemunitstd"`
	ItemUnitDiv         float64 `json:"itemunitdiv,omitempty" gorm:"column:itemunitdiv"`
	Price               float64 `json:"price" gorm:"column:price" `
	Qty                 float64 `json:"qty" gorm:"column:qty" `
	DiscountAmount      float64 `json:"discountamount" gorm:"column:discountamount"`
	DiscountText        string  `json:"discounttext" gorm:"column:discounttext"`
	Amount              float64 `json:"amount" gorm:"column:amount"`
}

func (SaleinvoiceDetailPg) TableName() string {
	return "saleinvoice_detail"
}
