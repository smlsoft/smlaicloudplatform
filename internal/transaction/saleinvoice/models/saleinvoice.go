package models

import (
	"smlcloudplatform/internal/models"
	transmodels "smlcloudplatform/internal/transaction/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const saleinvoiceCollectionName = "transactionSaleInvoice"

type SaleInvoice struct {
	models.PartitionIdentity `bson:"inline"`
	transmodels.Transaction  `bson:"inline"`
	IsPOS                    bool    `json:"ispos" bson:"ispos"`
	IsBom                    bool    `json:"isbom" bson:"isbom"`
	IsDelivery               bool    `json:"isdelivery" bson:"isdelivery"`
	IsTransport              bool    `json:"istransport" bson:"istransport"`
	TransportCode            string  `json:"transportcode" bson:"transportcode"`
	TransportAmount          float64 `json:"transportamount" bson:"transportamount"`

	CouponNo          string  `json:"couponno" bson:"couponno"`
	CouponAmount      float64 `json:"couponamount" bson:"couponamount"`
	CouponDescription string  `json:"coupondescription" bson:"coupondescription"`

	OrderNumber string `json:"ordernumber" bson:"ordernumber"`

	QRCode         string  `json:"qrcode" bson:"qrcode"`
	QRCodeAmount   float64 `json:"qrcodeamount" bson:"qrcodeamount"`
	DeliveryAmount float64 `json:"deliveryamount" bson:"deliveryamount"`

	ChequeNo         string  `json:"chequeno" bson:"chequeno"`
	ChequeBookNumber string  `json:"chequebooknumber" bson:"chequebooknumber"`
	ChequeBookCode   string  `json:"chequebookcode" bson:"chequebookcode"`
	ChequeDueDate    string  `json:"chequeduedate" bson:"chequeduedate"`
	ChequeAmount     float64 `json:"chequeamount" bson:"chequeamount"`

	SaleChannelCode   string  `json:"salechannelcode" bson:"salechannelcode"`
	SaleChannelGP     float64 `json:"salechannelgp" bson:"salechannelgp"`
	SaleChannelGPType int8    `json:"salechannelgptype" bson:"salechannelgptype"`
	TakeAway          int8    `json:"takeaway" bson:"takeaway"`

	SlipUrl            string   `json:"slipurl" bson:"slipurl"`
	SlipQrUrl          string   `json:"slipqrurl" bson:"slipqrurl"`
	SlipUrlHistories   []string `json:"slipurlhistories" bson:"slipurlhistories"`
	SlipQrUrlHistories []string `json:"slipqrurlhistories" bson:"slipqrurlhistories"`
	// PosID              string   `json:"posid" bson:"posid"`
	MachineCode     string `json:"machinecode" bson:"machinecode"`
	ZoneGroupNumber string `json:"zonegroupnumber" bson:"zonegroupnumber"`
}

// type SaleInvoiceDetail struct {
// 	transmodels.Detail `bson:"inline"`
// }

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
