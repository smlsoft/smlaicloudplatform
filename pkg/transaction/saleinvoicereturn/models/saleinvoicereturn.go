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
	Docno string `json:"docno" bson:"docno"`
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
