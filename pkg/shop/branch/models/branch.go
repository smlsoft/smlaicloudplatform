package models

import (
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const branchCollectionName = "branch"

type Branch struct {
	models.PartitionIdentity `bson:"inline"`
	Code                     uint16 `json:"code" bson:"code" validate:"required"`
	Name                     string `json:"name" bson:"name" validate:"required,min=1,max=100"`
	Telephone                string `json:"telephone" bson:"telephone"`
	Lat                      string `json:"lat" bson:"lat"`
	Lng                      string `json:"lng" bson:"lng"`
}

type BranchInfo struct {
	models.DocIdentity `bson:"inline"`
	Branch             `bson:"inline"`
}

func (BranchInfo) CollectionName() string {
	return branchCollectionName
}

type BranchData struct {
	models.ShopIdentity `bson:"inline"`
	BranchInfo          `bson:"inline"`
}

type BranchDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	BranchData         `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (BranchDoc) CollectionName() string {
	return branchCollectionName
}

type BranchItemGuid struct {
	Code string `json:"code" bson:"code"`
}

func (BranchItemGuid) CollectionName() string {
	return branchCollectionName
}

type BranchActivity struct {
	BranchData          `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (BranchActivity) CollectionName() string {
	return branchCollectionName
}

type BranchDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (BranchDeleteActivity) CollectionName() string {
	return branchCollectionName
}
