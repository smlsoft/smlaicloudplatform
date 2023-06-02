package models

import (
	"smlcloudplatform/pkg/models"
	transmodels "smlcloudplatform/pkg/transaction/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const saleinvoiceCollectionName = "transactionSaleInvoice"

type SaleInvoice struct {
	models.PartitionIdentity `bson:"inline"`
	transmodels.Transaction  `bson:"inline"`
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
