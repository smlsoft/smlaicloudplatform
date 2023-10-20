package models

import (
	"smlcloudplatform/pkg/models"
)

const deviceCollectionName = "orderDevices"

type OrderDevice struct {
	models.PartitionIdentity `bson:"inline"`
	ID                       string          `json:"id" bson:"id"`
	Names                    *[]models.NameX `json:"names" bson:"names"`
	ActivePin                string          `json:"activepin" bson:"activepin"`
	SettingCode              string          `json:"settingcode" bson:"settingcode"`
}

type OrderDeviceInfo struct {
	models.DocIdentity `bson:"inline"`
	OrderDevice        `bson:"inline"`
}

func (OrderDeviceInfo) CollectionName() string {
	return deviceCollectionName
}

type OrderDeviceData struct {
	models.ShopIdentity `bson:"inline"`
	OrderDeviceInfo     `bson:"inline"`
}

type OrderDeviceDoc struct {
	OrderDeviceData    `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (OrderDeviceDoc) CollectionName() string {
	return deviceCollectionName
}

type OrderDeviceItemGuid struct {
	ID string `json:"id" bson:"id"`
}

func (OrderDeviceItemGuid) CollectionName() string {
	return deviceCollectionName
}

type OrderDeviceActivity struct {
	OrderDeviceData     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (OrderDeviceActivity) CollectionName() string {
	return deviceCollectionName
}

type OrderDeviceDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (OrderDeviceDeleteActivity) CollectionName() string {
	return deviceCollectionName
}
