package models

import (
	"smlcloudplatform/pkg/models"
	transmodels "smlcloudplatform/pkg/transaction/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const payCollectionName = "pay"

type Pay struct {
	models.PartitionIdentity `bson:"inline"`
	DocNo                    string                    `json:"docno" bson:"docno"`
	DocDate                  time.Time                 `json:"docdate" bson:"docdate"`
	DocType                  int8                      `json:"doctype" bson:"doctype"`
	TransFlag                int8                      `json:"transflag" bson:"transflag"`
	CustCode                 string                    `json:"custcode" bson:"custcode"`
	CustNames                models.NameX              `json:"custnames" bson:"custnames"`
	SaleCode                 string                    `json:"salecode" bson:"salecode"`
	SaleName                 string                    `json:"salename" bson:"salename"`
	TotalPaymentAmount       float64                   `json:"totalpaymentamount" bson:"totalpaymentamount"`
	TotalAmount              float64                   `json:"totalamount" bson:"totalamount"`
	TotalBalance             float64                   `json:"totalbalance" bson:"totalbalance"`
	TotalValue               float64                   `json:"totalvalue" bson:"totalvalue"`
	Details                  *[]PayDetail              `json:"details" bson:"details"`
	PaymentDetail            transmodels.PaymentDetail `json:"paymentdetail" bson:"paymentdetail"`
}

type PayDetail struct {
	Selected      bool      `json:"selected" bson:"selected"`
	DocNo         string    `json:"docno" bson:"docno"`
	DocDate       time.Time `json:"docdate" bson:"docdate"`
	TransFlag     int8      `json:"transflag" bson:"transflag"`
	Value         float64   `json:"value" bson:"value"`
	Balance       float64   `json:"balance" bson:"balance"`
	PaymentAmount float64   `json:"paymentamount" bson:"paymentamount"`
}

type PayInfo struct {
	models.DocIdentity `bson:"inline"`
	Pay                `bson:"inline"`
}

func (PayInfo) CollectionName() string {
	return payCollectionName
}

type PayData struct {
	models.ShopIdentity `bson:"inline"`
	PayInfo             `bson:"inline"`
}

type PayDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PayData            `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (PayDoc) CollectionName() string {
	return payCollectionName
}

type PayItemGuid struct {
	DocNo string `json:"docno" bson:"docno"`
}

func (PayItemGuid) CollectionName() string {
	return payCollectionName
}

type PayActivity struct {
	PayData             `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (PayActivity) CollectionName() string {
	return payCollectionName
}

type PayDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (PayDeleteActivity) CollectionName() string {
	return payCollectionName
}
