package models

import (
	"smlcloudplatform/pkg/models"
	transmodels "smlcloudplatform/pkg/transaction/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const saleinvoicereturnCollectionName = "transactionSaleInvoiceReturn"

type SaleInvoiceReturn struct {
	models.PartitionIdentity `bson:"inline"`
	transmodels.Transaction  `bson:"inline"`
	IsPOS                    bool `json:"ispos" bson:"ispos"`

	CouponNo          string  `json:"couponno" bson:"couponno"`
	CouponAmount      float64 `json:"couponamount" bson:"couponamount"`
	CouponDescription string  `json:"coupondescription" bson:"coupondescription"`

	QRCode       string  `json:"qrcode" bson:"qrcode"`
	QRCodeAmount float64 `json:"qrcodeamount" bson:"qrcodeamount"`

	ChequeNo         string                       `json:"chequeno" bson:"chequeno"`
	ChequeBookNumber string                       `json:"chequebooknumber" bson:"chequebooknumber"`
	ChequeBookCode   string                       `json:"chequebookcode" bson:"chequebookcode"`
	ChequeDueDate    string                       `json:"chequeduedate" bson:"chequeduedate"`
	ChequeAmount     float64                      `json:"chequeamount" bson:"chequeamount"`
	SaleInvoice      SaleInvoiceReturnSaleChannel `json:"saleinvoice,omitempty" bson:"saleinvoice,omitempty"`
}

type SaleInvoiceReturnSaleChannel struct {
	Code   string  `json:"code" bson:"code"`
	Name   string  `json:"name" bson:"name"`
	GP     float64 `json:"gp" bson:"gp"`
	GPType int8    `json:"gptype" bson:"gptype"`
}

type SaleInvoiceReturnInfo struct {
	models.DocIdentity `bson:"inline"`
	SaleInvoiceReturn  `bson:"inline"`
}

func (SaleInvoiceReturnInfo) CollectionName() string {
	return saleinvoicereturnCollectionName
}

type SaleInvoiceReturnData struct {
	models.ShopIdentity   `bson:"inline"`
	SaleInvoiceReturnInfo `bson:"inline"`
}

type SaleInvoiceReturnDoc struct {
	ID                    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SaleInvoiceReturnData `bson:"inline"`
	models.ActivityDoc    `bson:"inline"`
}

func (SaleInvoiceReturnDoc) CollectionName() string {
	return saleinvoicereturnCollectionName
}

type SaleInvoiceReturnItemGuid struct {
	DocNo string `json:"docno" bson:"docno"`
}

func (SaleInvoiceReturnItemGuid) CollectionName() string {
	return saleinvoicereturnCollectionName
}

type SaleInvoiceReturnActivity struct {
	SaleInvoiceReturnData `bson:"inline"`
	models.ActivityTime   `bson:"inline"`
}

func (SaleInvoiceReturnActivity) CollectionName() string {
	return saleinvoicereturnCollectionName
}

type SaleInvoiceReturnDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (SaleInvoiceReturnDeleteActivity) CollectionName() string {
	return saleinvoicereturnCollectionName
}
