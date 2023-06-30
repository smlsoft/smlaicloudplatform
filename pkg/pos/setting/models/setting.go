package models

import (
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const settingCollectionName = "posSettings"

type Setting struct {
	models.PartitionIdentity `bson:"inline"`
	Code                     string    `json:"code" bson:"code" validate:"required"`
	DeviceNumber             string    `json:"devicenumber" bson:"devicenumber"`
	DocFormatTaxInv          string    `json:"docformattaxinv" bson:"docformattaxinv"`
	DocFormatInv             string    `json:"docformatinv" bson:"docformatinv"`
	ReceiptForm              string    `json:"receiptform" bson:"receiptform"`
	Warehouse                Warehouse `json:"warehouse" bson:"warehouse"`
	Location                 Location  `json:"location" bson:"location"`
	Branch                   Branch    `json:"branch" bson:"branch"`
}

type Branch struct {
	models.DocIdentity `bson:"inline"`
	Code               uint16          `json:"code" bson:"code"`
	Telephone          string          `json:"telephone" bson:"telephone" validate:"max=100"`
	Location           Location        `json:"location" bson:"location"`
	Names              *[]models.NameX `json:"names" bson:"names" validate:"required,min=1,unique=Code,dive"`
}

type Warehouse struct {
	models.DocIdentity `bson:"inline"`
	Code               string          `json:"code" bson:"code"`
	Names              *[]models.NameX `json:"names" bson:"names" validate:"required,min=1,unique=Code,dive"`
}

type Location struct {
	models.DocIdentity `bson:"inline"`
	Code               string          `json:"code" bson:"code"`
	Names              *[]models.NameX `json:"names" bson:"names" validate:"required,min=1,unique=Code,dive"`
}

type SettingInfo struct {
	models.DocIdentity `bson:"inline"`
	Setting            `bson:"inline"`
}

func (SettingInfo) CollectionName() string {
	return settingCollectionName
}

type SettingData struct {
	models.ShopIdentity `bson:"inline"`
	SettingInfo         `bson:"inline"`
}

type SettingDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SettingData        `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (SettingDoc) CollectionName() string {
	return settingCollectionName
}

type SettingItemGuid struct {
	Code string `json:"code" bson:"code"`
}

func (SettingItemGuid) CollectionName() string {
	return settingCollectionName
}

type SettingActivity struct {
	SettingData         `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (SettingActivity) CollectionName() string {
	return settingCollectionName
}

type SettingDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (SettingDeleteActivity) CollectionName() string {
	return settingCollectionName
}
