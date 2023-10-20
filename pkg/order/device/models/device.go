package models

import (
	"smlcloudplatform/pkg/models"
)

const deviceCollectionName = "device"

type Device struct {
	models.PartitionIdentity `bson:"inline"`
	ID                       string          `json:"id" bson:"id"`
	Names                    *[]models.NameX `json:"names" bson:"names"`
	ActivePin                string          `json:"activepin" bson:"activepin"`
	SettingCode              string          `json:"settingcode" bson:"settingcode"`
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
	DeviceData         `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (DeviceDoc) CollectionName() string {
	return deviceCollectionName
}

type DeviceItemGuid struct {
	ID string `json:"id" bson:"id"`
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
