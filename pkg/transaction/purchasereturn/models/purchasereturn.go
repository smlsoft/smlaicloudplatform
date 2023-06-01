package models

import (
	"smlcloudplatform/pkg/models"
	transmodels "smlcloudplatform/pkg/transaction/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const purchasereturnCollectionName = "transactionPurchaseReturn"

type PurchaseReturn struct {
	models.PartitionIdentity `bson:"inline"`
	transmodels.Transaction  `bson:"inline"`
}

type PurchaseReturnInfo struct {
	models.DocIdentity `bson:"inline"`
	PurchaseReturn     `bson:"inline"`
}

func (PurchaseReturnInfo) CollectionName() string {
	return purchasereturnCollectionName
}

type PurchaseReturnData struct {
	models.ShopIdentity `bson:"inline"`
	PurchaseReturnInfo  `bson:"inline"`
}

type PurchaseReturnDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PurchaseReturnData `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (PurchaseReturnDoc) CollectionName() string {
	return purchasereturnCollectionName
}

type PurchaseReturnItemGuid struct {
	Docno string `json:"docno" bson:"docno"`
}

func (PurchaseReturnItemGuid) CollectionName() string {
	return purchasereturnCollectionName
}

type PurchaseReturnActivity struct {
	PurchaseReturnData  `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (PurchaseReturnActivity) CollectionName() string {
	return purchasereturnCollectionName
}

type PurchaseReturnDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (PurchaseReturnDeleteActivity) CollectionName() string {
	return purchasereturnCollectionName
}
