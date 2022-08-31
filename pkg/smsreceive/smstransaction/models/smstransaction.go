package models

import (
	"smlcloudplatform/pkg/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const smstransactionCollectionName = "smsTransactions"

type SmsTransaction struct {
	models.PartitionIdentity `bson:"inline"`
	TransId                  string    `json:"transid" bson:"transid"`
	Address                  string    `json:"address" bson:"address"`
	Body                     string    `json:"body" bson:"body"`
	SendedAt                 time.Time `json:"sendedat" bson:"sendedat"`
	Status                   int8      `json:"status" bson:"status"`
}

type SmsTransactionInfo struct {
	models.DocIdentity `bson:"inline"`
	SmsTransaction     `bson:"inline"`
}

func (SmsTransactionInfo) CollectionName() string {
	return smstransactionCollectionName
}

type SmsTransactionData struct {
	models.ShopIdentity `bson:"inline"`
	SmsTransactionInfo  `bson:"inline"`
}

type SmsTransactionDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SmsTransactionData `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (SmsTransactionDoc) CollectionName() string {
	return smstransactionCollectionName
}

type SmsTransactionItemGuid struct {
	DocNo string `json:"docno" bson:"docno" gorm:"docno"`
}

func (SmsTransactionItemGuid) CollectionName() string {
	return smstransactionCollectionName
}

type SmsTransactionActivity struct {
	SmsTransactionData  `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (SmsTransactionActivity) CollectionName() string {
	return smstransactionCollectionName
}

type SmsTransactionDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (SmsTransactionDeleteActivity) CollectionName() string {
	return smstransactionCollectionName
}

type SmsTransactionInfoResponse struct {
	Success bool               `json:"success"`
	Data    SmsTransactionInfo `json:"data,omitempty"`
}

type SmsTransactionPageResponse struct {
	Success    bool                          `json:"success"`
	Data       []SmsTransactionInfo          `json:"data,omitempty"`
	Pagination models.PaginationDataResponse `json:"pagination,omitempty"`
}
