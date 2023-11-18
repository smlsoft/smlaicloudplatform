package models

import (
	"smlcloudplatform/pkg/models"
	transmodels "smlcloudplatform/pkg/transaction/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const saleinvoiceCollectionName = "transactionSaleInvoice"

type SaleInvoice struct {
	models.PartitionIdentity      `bson:"inline"`
	transmodels.TransactionHeader `bson:"inline"`
	Details                       *[]SaleInvoiceDetail `json:"details" bson:"details"`
	IsPOS                         bool                 `json:"ispos" bson:"ispos"`

	CouponNo          string  `json:"couponno" bson:"couponno"`
	CouponAmount      float64 `json:"couponamount" bson:"couponamount"`
	CouponDescription string  `json:"coupondescription" bson:"coupondescription"`

	QRCode       string  `json:"qrcode" bson:"qrcode"`
	QRCodeAmount float64 `json:"qrcodeamount" bson:"qrcodeamount"`

	ChequeNo         string  `json:"chequeno" bson:"chequeno"`
	ChequeBookNumber string  `json:"chequebooknumber" bson:"chequebooknumber"`
	ChequeBookCode   string  `json:"chequebookcode" bson:"chequebookcode"`
	ChequeDueDate    string  `json:"chequeduedate" bson:"chequeduedate"`
	ChequeAmount     float64 `json:"chequeamount" bson:"chequeamount"`

	SaleChannelCode   string  `json:"salechannelcode" bson:"csalechannelode"`
	SaleChannelGP     float64 `json:"salechannelgp" bson:"salechannelgp"`
	SaleChannelGPType int8    `json:"salechannelgptype" bson:"salechannelgptype"`
	TakeAway          int8    `json:"takeaway" bson:"takeaway"`
}

type SaleInvoiceDetail struct {
	ManufacturerGUID   string `json:"manufacturerguid" bson:"manufacturerguid"`
	transmodels.Detail `bson:"inline"`
}

type SaleInvoiceInfo struct {
	models.DocIdentity `bson:"inline"`
	SaleInvoice        `bson:"inline"`
}

func (SaleInvoiceInfo) CollectionName() string {
	return saleinvoiceCollectionName
}

type SaleInvoiceData struct {
	models.ShopIdentity `bson:"inline"`
	SaleInvoiceInfo     `bson:"inline"`
}

type SaleInvoiceDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SaleInvoiceData    `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (SaleInvoiceDoc) CollectionName() string {
	return saleinvoiceCollectionName
}

type SaleInvoiceItemGuid struct {
	DocNo string `json:"docno" bson:"docno"`
}

func (SaleInvoiceItemGuid) CollectionName() string {
	return saleinvoiceCollectionName
}

type SaleInvoiceActivity struct {
	SaleInvoiceData     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (SaleInvoiceActivity) CollectionName() string {
	return saleinvoiceCollectionName
}

type SaleInvoiceDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (SaleInvoiceDeleteActivity) CollectionName() string {
	return saleinvoiceCollectionName
}
