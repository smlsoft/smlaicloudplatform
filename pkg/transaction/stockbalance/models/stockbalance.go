package models

import (
	"smlcloudplatform/pkg/models"
	transmodels "smlcloudplatform/pkg/transaction/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const stockbalanceCollectionName = "transactionStockBalance"

type StockBalanceHeader struct {
	transmodels.TransactionHeader `bson:"inline"`
}

type StockBalance struct {
	models.PartitionIdentity `bson:"inline"`
	StockBalanceHeader       `bson:"inline"`
	Details                  *[]StockBalanceDetail `json:"details" bson:"details"`
}

type StockBalanceDetail struct {
	transmodels.Detail
}

type StockBalanceInfo struct {
	models.DocIdentity `bson:"inline"`
	StockBalance       `bson:"inline"`
}

func (StockBalanceInfo) CollectionName() string {
	return stockbalanceCollectionName
}

type StockBalanceData struct {
	models.ShopIdentity `bson:"inline"`
	StockBalanceInfo    `bson:"inline"`
}

type StockBalanceDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	StockBalanceData   `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (StockBalanceDoc) CollectionName() string {
	return stockbalanceCollectionName
}

type StockBalanceItemGuid struct {
	DocNo string `json:"docno" bson:"docno"`
}

func (StockBalanceItemGuid) CollectionName() string {
	return stockbalanceCollectionName
}

type StockBalanceActivity struct {
	StockBalanceData    `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (StockBalanceActivity) CollectionName() string {
	return stockbalanceCollectionName
}

type StockBalanceDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (StockBalanceDeleteActivity) CollectionName() string {
	return stockbalanceCollectionName
}
