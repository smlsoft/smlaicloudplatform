package models

import (
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const unitCollectionName = "unit"

type Unit struct {
	models.PartitionIdentity `bson:"inline"`
	UnitCode                 string `json:"unitcode" bson:"unitcode" `
	models.Name              `bson:"inline"`
	ItemUnitSTD              float64 `json:"itemunitstd" bson:"itemunitstd" `
	ItemUnitDIV              float64 `json:"itemunitdiv" bson:"itemunitdiv" `
	IsUnitCost               bool    `json:"isunitcost" bson:"isunitcost"`
	IsUnitStandard           bool    `json:"isunitstandard" bson:"isunitstandard"`
}

type UnitInfo struct {
	models.DocIdentity `bson:"inline"`
	Unit               `bson:"inline"`
}

func (UnitInfo) CollectionName() string {
	return unitCollectionName
}

type UnitData struct {
	models.ShopIdentity `bson:"inline"`
	UnitInfo            `bson:"inline"`
}

type UnitDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UnitData           `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (UnitDoc) CollectionName() string {
	return unitCollectionName
}

type UnitItemGuid struct {
	DocNo string `json:"docno" bson:"docno"`
}

func (UnitItemGuid) CollectionName() string {
	return unitCollectionName
}

type UnitActivity struct {
	UnitData            `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (UnitActivity) CollectionName() string {
	return unitCollectionName
}

type UnitDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (UnitDeleteActivity) CollectionName() string {
	return unitCollectionName
}
