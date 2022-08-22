package models

import (
	"smlcloudplatform/pkg/models"

	common "smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const bankmasterCollectionName = "bankMaster"

type BankMaster struct {
	models.PartitionIdentity `bson:"inline"`
	BankCode                 string `json:"bankcode" bson:"bankcode"`
	CountryCode              string `json:"countrycode" bson:"countrycode"`
	BankLogo                 string `json:"banklogo" bson:"banklogo"`
	common.Name              `bson:"inline"`
}

type BankMasterInfo struct {
	models.DocIdentity `bson:"inline"`
	BankMaster         `bson:"inline"`
}

func (BankMasterInfo) CollectionName() string {
	return bankmasterCollectionName
}

type BankMasterData struct {
	models.ShopIdentity `bson:"inline"`
	BankMasterInfo      `bson:"inline"`
}

type BankMasterDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	BankMasterData     `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (BankMasterDoc) CollectionName() string {
	return bankmasterCollectionName
}

type BankMasterItemGuid struct {
	BankCode string `json:"bankcode" bson:"bankcode" gorm:"bankcode"`
}

func (BankMasterItemGuid) CollectionName() string {
	return bankmasterCollectionName
}

type BankMasterActivity struct {
	BankMasterData      `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (BankMasterActivity) CollectionName() string {
	return bankmasterCollectionName
}

type BankMasterDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (BankMasterDeleteActivity) CollectionName() string {
	return bankmasterCollectionName
}

type BankMasterInfoResponse struct {
	Success bool           `json:"success"`
	Data    BankMasterInfo `json:"data,omitempty"`
}

type BankMasterPageResponse struct {
	Success    bool                          `json:"success"`
	Data       []BankMasterInfo              `json:"data,omitempty"`
	Pagination models.PaginationDataResponse `json:"pagination,omitempty"`
}
