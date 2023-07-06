package models

import (
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const documentformateCollectionName = "documentFormate"

type DocumentFormate struct {
	models.PartitionIdentity `bson:"inline"`
	DocCode                  string                   `json:"doccode" bson:"doccode" validate:"required,min=1"`
	Module                   string                   `json:"module" bson:"module" validate:"required,min=1"`
	DateFormate              string                   `json:"dateformate" bson:"dateformate" validate:"required,min=1"`
	DocNumber                int                      `json:"docnumber" bson:"docnumber" validate:"required"`
	Details                  *[]DocumentFormateDetail `json:"details" bson:"details"`
}

type DocumentFormateDetail struct {
	ActionCode         string `json:"actioncode" bson:"actioncode" validate:"required,min=1"`
	Detail             string `json:"detail" bson:"detail"`
	Debit              string `json:"debit" bson:"debit"`
	Credit             string `json:"credit" bson:"credit"`
	IsEntrySelfAccount bool   `json:"isentryselfaccount" bson:"isentryselfaccount"`
}

type DocumentFormateInfo struct {
	models.DocIdentity `bson:"inline"`
	DocumentFormate    `bson:"inline"`
}

func (DocumentFormateInfo) CollectionName() string {
	return documentformateCollectionName
}

type DocumentFormateData struct {
	models.ShopIdentity `bson:"inline"`
	DocumentFormateInfo `bson:"inline"`
}

type DocumentFormateDoc struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	DocumentFormateData `bson:"inline"`
	models.ActivityDoc  `bson:"inline"`
}

func (DocumentFormateDoc) CollectionName() string {
	return documentformateCollectionName
}

type DocumentFormateItemGuid struct {
	DocCode string `json:"doccode" bson:"doccode"`
}

func (DocumentFormateItemGuid) CollectionName() string {
	return documentformateCollectionName
}

type DocumentFormateActivity struct {
	DocumentFormateData `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (DocumentFormateActivity) CollectionName() string {
	return documentformateCollectionName
}

type DocumentFormateDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (DocumentFormateDeleteActivity) CollectionName() string {
	return documentformateCollectionName
}
