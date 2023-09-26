package models

import (
	"smlcloudplatform/pkg/models"
	transmodels "smlcloudplatform/pkg/transaction/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const paidCollectionName = "transactionPaid"

type Paid struct {
	models.PartitionIdentity `bson:"inline"`
	DocNo                    string                    `json:"docno" bson:"docno"`
	DocDatetime              time.Time                 `json:"docdatetime" bson:"docdatetime"`
	DocType                  int8                      `json:"doctype" bson:"doctype"`
	TransFlag                int8                      `json:"transflag" bson:"transflag"`
	CustCode                 string                    `json:"custcode" bson:"custcode"`
	CustNames                *[]models.NameX           `json:"custnames" bson:"custnames"`
	SaleCode                 string                    `json:"salecode" bson:"salecode"`
	SaleName                 string                    `json:"salename" bson:"salename"`
	TotalPaymentAmount       float64                   `json:"totalpaymentamount" bson:"totalpaymentamount"`
	TotalAmount              float64                   `json:"totalamount" bson:"totalamount"`
	TotalBalance             float64                   `json:"totalbalance" bson:"totalbalance"`
	TotalValue               float64                   `json:"totalvalue" bson:"totalvalue"`
	Details                  *[]PaidDetail             `json:"details" bson:"details"`
	PaymentDetail            transmodels.PaymentDetail `json:"paymentdetail" bson:"paymentdetail"`
	PaymentDetailRaw         string                    `json:"paymentdetailraw" bson:"paymentdetailraw"`
}

type PaidDetail struct {
	Selected      bool      `json:"selected" bson:"selected"`
	DocNo         string    `json:"docno" bson:"docno"`
	DocDatetime   time.Time `json:"docdatetime" bson:"docdatetime"`
	TransFlag     int8      `json:"transflag" bson:"transflag"`
	Value         float64   `json:"value" bson:"value"`
	Balance       float64   `json:"balance" bson:"balance"`
	PaymentAmount float64   `json:"paymentamount" bson:"paymentamount"`
}

type PaidInfo struct {
	models.DocIdentity `bson:"inline"`
	Paid               `bson:"inline"`
}

func (PaidInfo) CollectionName() string {
	return paidCollectionName
}

type PaidData struct {
	models.ShopIdentity `bson:"inline"`
	PaidInfo            `bson:"inline"`
}

type PaidDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PaidData           `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (PaidDoc) CollectionName() string {
	return paidCollectionName
}

type PaidItemGuid struct {
	DocNo string `json:"docno" bson:"docno"`
}

func (PaidItemGuid) CollectionName() string {
	return paidCollectionName
}

type PaidActivity struct {
	PaidData            `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (PaidActivity) CollectionName() string {
	return paidCollectionName
}

type PaidDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (PaidDeleteActivity) CollectionName() string {
	return paidCollectionName
}
