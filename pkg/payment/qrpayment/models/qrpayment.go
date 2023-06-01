package models

import (
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const qrpaymentCollectionName = "qrPayment"

type QrPayment struct {
	models.PartitionIdentity `bson:"inline"`
	PaymentCode              string          `json:"paymentcode" bson:"paymentcode"`
	CountryCode              string          `json:"countrycode" bson:"countrycode"`
	PaymentLogo              string          `json:"paymentlogo" bson:"paymentlogo"`
	PaymentType              int8            `json:"paymenttype" bson:"paymenttype"`
	FeeRate                  float64         `json:"feerate" bson:"feerate"`
	WalletPayType            int16           `json:"wallettype" bson:"wallettype"`
	Names                    *[]models.NameX `json:"names" bson:"names" validate:"required,min=1,unique=Code,dive"`
	BookBankCode             string          `json:"bookbankcode" bson:"bookbankcode"`
	BankCode                 string          `json:"bankcode" bson:"bankcode"`
}

type QrPaymentInfo struct {
	models.DocIdentity `bson:"inline"`
	QrPayment          `bson:"inline"`
}

func (QrPaymentInfo) CollectionName() string {
	return qrpaymentCollectionName
}

type QrPaymentData struct {
	models.ShopIdentity `bson:"inline"`
	QrPaymentInfo       `bson:"inline"`
}

type QrPaymentDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	QrPaymentData      `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (QrPaymentDoc) CollectionName() string {
	return qrpaymentCollectionName
}

type QrPaymentItemGuid struct {
	PaymentCode string `json:"paymentcode" bson:"paymentcode"`
}

func (QrPaymentItemGuid) CollectionName() string {
	return qrpaymentCollectionName
}

type QrPaymentActivity struct {
	QrPaymentData       `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (QrPaymentActivity) CollectionName() string {
	return qrpaymentCollectionName
}

type QrPaymentDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (QrPaymentDeleteActivity) CollectionName() string {
	return qrpaymentCollectionName
}
