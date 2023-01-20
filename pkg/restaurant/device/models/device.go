package models

import (
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const deviceCollectionName = "restarantDevices"

type Device struct {
	models.PartitionIdentity `bson:"inline"`
	Code                     string `json:"code" bson:"code" validate:"required" `
	Type                     int16  `json:"type" bson:"type" `
	Name1                    string `json:"name1" bson:"name1" `
	Name2                    string `json:"name2,omitempty" bson:"name2,omitempty"`
	Name3                    string `json:"name3,omitempty" bson:"name3,omitempty"`
	Name4                    string `json:"name4,omitempty" bson:"name4,omitempty"`
	Name5                    string `json:"name5,omitempty" bson:"name5,omitempty"`
}

type DeviceInfo struct {
	models.DocIdentity `bson:"inline"`
	Device             `bson:"inline"`
}

func (DeviceInfo) CollectionName() string {
	return deviceCollectionName
}

type DeviceData struct {
	models.ShopIdentity `bson:"inline"`
	DeviceInfo          `bson:"inline"`
}

type DeviceDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	DeviceData         `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (DeviceDoc) CollectionName() string {
	return deviceCollectionName
}

type DeviceItemGuid struct {
	Code string `json:"code" bson:"code"`
}

func (DeviceItemGuid) CollectionName() string {
	return deviceCollectionName
}

type DeviceActivity struct {
	DeviceData          `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (DeviceActivity) CollectionName() string {
	return deviceCollectionName
}

type DeviceDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (DeviceDeleteActivity) CollectionName() string {
	return deviceCollectionName
}
