package models

import (
	common "smlcloudplatform/pkg/models"
	invModel "smlcloudplatform/pkg/product/inventory/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const barcodemasterCollectionName string = "barcodeMaster"
const barcodemasterTableName string = "barcodemaster"
const barcodemasterIndexName string = "barcodemaster_index"

type BarcodeMaster struct {
	invModel.Inventory `bson:"inline" gorm:"embedded;"`
}

type BarcodeMasterItemGuid struct {
	ItemGuid string `json:"itemguid,omitempty" bson:"itemguid,omitempty"`
}

func (BarcodeMasterItemGuid) CollectionName() string {
	return barcodemasterCollectionName
}

type BarcodeMasterImage struct {
	DocID string `json:"-" bson:"-" gorm:"docid;primaryKey"`
	Uri   string `json:"uri" bson:"uri" gorm:"uri;primaryKey"`
}

type BarcodeMasterTag struct {
	DocID string `json:"-" bson:"-" gorm:"docid;primaryKey"`
	Name  string `json:"name" bson:"name" gorm:"name;primaryKey"`
}

type BarcodeMasterInfo struct {
	common.DocIdentity `bson:"inline" gorm:"embedded;"`
	BarcodeMaster      `bson:"inline" gorm:"embedded;"`
}

func (BarcodeMasterInfo) CollectionName() string {
	return barcodemasterCollectionName
}

type BarcodeMasterData struct {
	common.ShopIdentity `bson:"inline" gorm:"embedded;"`
	BarcodeMasterInfo   `bson:"inline" gorm:"embedded;"`
}

func (BarcodeMasterData) TableName() string {
	return barcodemasterTableName
}

type BarcodeMasterDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	BarcodeMasterData  `bson:"inline"`
	common.ActivityDoc `bson:"inline"`
	common.LastUpdate  `bson:"inline"`
}

func (BarcodeMasterDoc) CollectionName() string {
	return barcodemasterCollectionName
}

type BarcodeMasterActivity struct {
	BarcodeMasterData `bson:"inline"`
	CreatedAt         *time.Time `json:"createdat,omitempty" bson:"createdat,omitempty"`
	UpdatedAt         *time.Time `json:"updatedat,omitempty" bson:"updatedat,omitempty"`
	DeletedAt         *time.Time `json:"deletedat,omitempty" bson:"deletedat,omitempty"`
}

func (BarcodeMasterActivity) CollectionName() string {
	return barcodemasterCollectionName
}

type BarcodeMasterDeleteActivity struct {
	common.Identity `bson:"inline"`
	CreatedAt       *time.Time `json:"createdat,omitempty" bson:"createdat,omitempty"`
	UpdatedAt       *time.Time `json:"updatedat,omitempty" bson:"updatedat,omitempty"`
	DeletedAt       *time.Time `json:"deletedat,omitempty" bson:"deletedat,omitempty"`
}

func (BarcodeMasterDeleteActivity) CollectionName() string {
	return barcodemasterCollectionName
}

type BarcodeMasterIndex struct {
	common.Index `bson:"inline"`
}

func (BarcodeMasterIndex) TableName() string {
	return barcodemasterIndexName
}

// for swagger gen

type BarcodeMasterBulkImport struct {
	Created          []string `json:"created"`
	Updated          []string `json:"updated"`
	UpdateFailed     []string `json:"updateFailed"`
	PayloadDuplicate []string `json:"payloadDuplicate"`
}

type BarcodeMasterBulkReponse struct {
	Success bool `json:"success"`
	BarcodeMasterBulkImport
}

type BarcodeMasterPageResponse struct {
	Success    bool                          `json:"success"`
	Data       []BarcodeMasterInfo           `json:"data,omitempty"`
	Pagination common.PaginationDataResponse `json:"pagination,omitempty"`
}

type BarcodeMasterInfoResponse struct {
	Success bool              `json:"success"`
	Data    BarcodeMasterInfo `json:"data,omitempty"`
}

type BarcodeMasterBulkInsertResponse struct {
	Success    bool     `json:"success"`
	Created    []string `json:"created"`
	Updated    []string `json:"updated"`
	Failed     []string `json:"updateFailed"`
	Duplicated []string `json:"payloadDuplicate"`
}

type BarcodeMasterLastActivityResponse struct {
	New    []BarcodeMasterActivity       `json:"new" `
	Remove []BarcodeMasterDeleteActivity `json:"remove"`
}

type BarcodeMasterFetchUpdateResponse struct {
	Success    bool                              `json:"success"`
	Data       BarcodeMasterLastActivityResponse `json:"data,omitempty"`
	Pagination common.PaginationDataResponse     `json:"pagination,omitempty"`
}
