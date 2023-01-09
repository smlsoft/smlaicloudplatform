package models

import (
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const staffCollectionName = "restaurantStaffs"

type Staff struct {
	models.PartitionIdentity `bson:"inline"`
	Code                     string `json:"code" bson:"code"`
	Name1                    string `json:"name1" bson:"name1" `
	Name2                    string `json:"name2,omitempty" bson:"name2,omitempty"`
	Name3                    string `json:"name3,omitempty" bson:"name3,omitempty"`
	Name4                    string `json:"name4,omitempty" bson:"name4,omitempty"`
	Name5                    string `json:"name5,omitempty" bson:"name5,omitempty"`
}

type StaffInfo struct {
	models.DocIdentity `bson:"inline"`
	Staff              `bson:"inline"`
}

func (StaffInfo) CollectionName() string {
	return staffCollectionName
}

type StaffData struct {
	models.ShopIdentity `bson:"inline"`
	StaffInfo           `bson:"inline"`
}

type StaffDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	StaffData          `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (StaffDoc) CollectionName() string {
	return staffCollectionName
}

type StaffItemGuid struct {
	Code string `json:"code" bson:"code"`
}

func (StaffItemGuid) CollectionName() string {
	return staffCollectionName
}

type StaffActivity struct {
	StaffData           `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (StaffActivity) CollectionName() string {
	return staffCollectionName
}

type StaffDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (StaffDeleteActivity) CollectionName() string {
	return staffCollectionName
}
