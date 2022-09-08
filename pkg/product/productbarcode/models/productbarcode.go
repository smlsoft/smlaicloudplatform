package models

import (
	common "smlcloudplatform/pkg/models"
	invModel "smlcloudplatform/pkg/product/inventory/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const productbarcodeCollectionName string = "inventories"
const productbarcodeTableName string = "inventories"
const productbarcodeIndexName string = "inventories_index"

type ProductBarcode struct {
	invModel.Inventory `bson:"inline" gorm:"embedded;"`
}

type ProductBarcodeItemGuid struct {
	ItemGuid string `json:"itemguid,omitempty" bson:"itemguid,omitempty"`
}

func (ProductBarcodeItemGuid) CollectionName() string {
	return productbarcodeCollectionName
}

type ProductBarcodeImage struct {
	DocID string `json:"-" bson:"-" gorm:"docid;primaryKey"`
	Uri   string `json:"uri" bson:"uri" gorm:"uri;primaryKey"`
}

type ProductBarcodeTag struct {
	DocID string `json:"-" bson:"-" gorm:"docid;primaryKey"`
	Name  string `json:"name" bson:"name" gorm:"name;primaryKey"`
}

type ProductBarcodeInfo struct {
	common.DocIdentity `bson:"inline" gorm:"embedded;"`
	ProductBarcode     `bson:"inline" gorm:"embedded;"`
	Unit               *invModel.Unit    `json:"unit,omitempty" bson:"unit,omitempty"`
	BarcodeDetail      *invModel.Barcode `json:"barcodedetail,omitempty" bson:"barcodedetail,omitempty"`
}

func (ProductBarcodeInfo) CollectionName() string {
	return productbarcodeCollectionName
}

type ProductBarcodeData struct {
	common.ShopIdentity `bson:"inline" gorm:"embedded;"`
	ProductBarcodeInfo  `bson:"inline" gorm:"embedded;"`
}

func (ProductBarcodeData) TableName() string {
	return productbarcodeTableName
}

type ProductBarcodeDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProductBarcodeData `bson:"inline"`
	common.ActivityDoc `bson:"inline"`
	common.LastUpdate  `bson:"inline"`
}

func (ProductBarcodeDoc) CollectionName() string {
	return productbarcodeCollectionName
}

type ProductBarcodeActivity struct {
	ProductBarcodeData `bson:"inline"`
	CreatedAt          *time.Time `json:"createdat,omitempty" bson:"createdat,omitempty"`
	UpdatedAt          *time.Time `json:"updatedat,omitempty" bson:"updatedat,omitempty"`
	DeletedAt          *time.Time `json:"deletedat,omitempty" bson:"deletedat,omitempty"`
}

func (ProductBarcodeActivity) CollectionName() string {
	return productbarcodeCollectionName
}

type ProductBarcodeDeleteActivity struct {
	common.Identity `bson:"inline"`
	CreatedAt       *time.Time `json:"createdat,omitempty" bson:"createdat,omitempty"`
	UpdatedAt       *time.Time `json:"updatedat,omitempty" bson:"updatedat,omitempty"`
	DeletedAt       *time.Time `json:"deletedat,omitempty" bson:"deletedat,omitempty"`
}

func (ProductBarcodeDeleteActivity) CollectionName() string {
	return productbarcodeCollectionName
}

type ProductBarcodeIndex struct {
	common.Index `bson:"inline"`
}

func (ProductBarcodeIndex) TableName() string {
	return productbarcodeIndexName
}

// for swagger gen

type ProductBarcodeBulkImport struct {
	Created          []string `json:"created"`
	Updated          []string `json:"updated"`
	UpdateFailed     []string `json:"updateFailed"`
	PayloadDuplicate []string `json:"payloadDuplicate"`
}

type ProductBarcodeBulkReponse struct {
	Success bool `json:"success"`
	ProductBarcodeBulkImport
}

type ProductBarcodePageResponse struct {
	Success    bool                          `json:"success"`
	Data       []ProductBarcodeInfo          `json:"data,omitempty"`
	Pagination common.PaginationDataResponse `json:"pagination,omitempty"`
}

type ProductBarcodeInfoResponse struct {
	Success bool               `json:"success"`
	Data    ProductBarcodeInfo `json:"data,omitempty"`
}

type ProductBarcodeBulkInsertResponse struct {
	Success    bool     `json:"success"`
	Created    []string `json:"created"`
	Updated    []string `json:"updated"`
	Failed     []string `json:"updateFailed"`
	Duplicated []string `json:"payloadDuplicate"`
}

type ProductBarcodeLastActivityResponse struct {
	New    []ProductBarcodeActivity       `json:"new" `
	Remove []ProductBarcodeDeleteActivity `json:"remove"`
}

type ProductBarcodeFetchUpdateResponse struct {
	Success    bool                               `json:"success"`
	Data       ProductBarcodeLastActivityResponse `json:"data,omitempty"`
	Pagination common.PaginationDataResponse      `json:"pagination,omitempty"`
}
