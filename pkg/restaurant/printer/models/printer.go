package models

import (
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const printerCollectionName = "printers"

type Printer struct {
	Code    string `json:"code" bson:"code"`
	Name1   string `json:"name1" bson:"name1" gorm:"name1"`
	Name2   string `json:"name2,omitempty" bson:"name2,omitempty"`
	Name3   string `json:"name3,omitempty" bson:"name3,omitempty"`
	Name4   string `json:"name4,omitempty" bson:"name4,omitempty"`
	Name5   string `json:"name5,omitempty" bson:"name5,omitempty"`
	Address string `json:"address" bson:"address" `
	Type    int8   `json:"type" bson:"type"`
}

type PrinterInfo struct {
	models.DocIdentity `bson:"inline"`
	Printer            `bson:"inline"`
}

func (PrinterInfo) CollectionName() string {
	return printerCollectionName
}

type PrinterData struct {
	models.ShopIdentity `bson:"inline"`
	PrinterInfo         `bson:"inline"`
}

type PrinterDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	PrinterData        `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
	models.LastUpdate  `bson:"inline"`
}

func (PrinterDoc) CollectionName() string {
	return printerCollectionName
}

//Extra

type PrinterItemGuid struct {
	Code string `json:"code" bson:"code" gorm:"code"`
}

func (PrinterItemGuid) CollectionName() string {
	return printerCollectionName
}

type PrinterActivity struct {
	PrinterData         `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (PrinterActivity) CollectionName() string {
	return printerCollectionName
}

type PrinterDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (PrinterDeleteActivity) CollectionName() string {
	return printerCollectionName
}

type PrinterInfoResponse struct {
	Success bool        `json:"success"`
	Data    PrinterInfo `json:"data,omitempty"`
}

type PrinterPageResponse struct {
	Success    bool                          `json:"success"`
	Data       []PrinterInfo                 `json:"data,omitempty"`
	Pagination models.PaginationDataResponse `json:"pagination,omitempty"`
}

type PrinterLastActivityResponse struct {
	New    []PrinterActivity       `json:"new" `
	Remove []PrinterDeleteActivity `json:"remove"`
}

type PrinterFetchUpdateResponse struct {
	Success    bool                          `json:"success"`
	Data       PrinterLastActivityResponse   `json:"data,omitempty"`
	Pagination models.PaginationDataResponse `json:"pagination,omitempty"`
}
