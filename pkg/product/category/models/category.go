package models

import (
	"smlcloudplatform/pkg/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const categoryCollectionName = "categories"

type Category struct {
	CategoryGuid string `json:"categoryguid" bson:"categoryguid" gorm:"categoryguid" validate:"required,max=100"`
	ParentGuid   string `json:"parentguid"  bson:"parentguid" validate:"max=100"`
	models.Name  `bson:"inline"`
	Image        string `json:"image" bson:"image,omitempty"`
	XOrder       int8   `json:"xorder" bson:"xorder,omitempty" validate:"min=-125,max=125"`
	Code         string `json:"code" bson:"code" validate:"max=100"`
}

type CategoryInfo struct {
	models.DocIdentity `bson:"inline"`
	Category           `bson:"inline"`
}

func (CategoryInfo) CollectionName() string {
	return categoryCollectionName
}

type CategoryData struct {
	models.ShopIdentity `bson:"inline"`
	CategoryInfo        `bson:"inline"`
}

type CategoryDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CategoryData       `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
	models.LastUpdate  `bson:"inline"`
}

func (CategoryDoc) CollectionName() string {
	return categoryCollectionName
}

type CategoryActivity struct {
	CategoryData `bson:"inline"`
	CreatedAt    *time.Time `json:"createdat,omitempty" bson:"createdat,omitempty"`
	UpdatedAt    *time.Time `json:"updatedat,omitempty" bson:"updatedat,omitempty"`
	DeletedAt    *time.Time `json:"deletedat,omitempty" bson:"deletedat,omitempty"`
}

func (CategoryActivity) CollectionName() string {
	return categoryCollectionName
}

type CategoryDeleteActivity struct {
	models.Identity `bson:"inline"`
	CreatedAt       *time.Time `json:"createdat,omitempty" bson:"createdat,omitempty"`
	UpdatedAt       *time.Time `json:"updatedat,omitempty" bson:"updatedat,omitempty"`
	DeletedAt       *time.Time `json:"deletedat,omitempty" bson:"deletedat,omitempty"`
}

func (CategoryDeleteActivity) CollectionName() string {
	return categoryCollectionName
}

type CategoryItemGuid struct {
	CategoryGuid string `json:"categoryguid" bson:"categoryguid" gorm:"categoryguid"`
}

func (CategoryItemGuid) CollectionName() string {
	return categoryCollectionName
}

//for swagger gen

type CategoryBulkImport struct {
	Created          []string `json:"created"`
	Updated          []string `json:"updated"`
	UpdateFailed     []string `json:"updateFailed"`
	PayloadDuplicate []string `json:"payloadDuplicate"`
}

type CategoryBulkReponse struct {
	Success bool `json:"success"`
	CategoryBulkImport
}

type CategoryLastActivityResponse struct {
	New    []CategoryActivity       `json:"new" `
	Remove []CategoryDeleteActivity `json:"remove"`
}

type CategoryFetchUpdateResponse struct {
	Success    bool                          `json:"success"`
	Data       CategoryLastActivityResponse  `json:"data,omitempty"`
	Pagination models.PaginationDataResponse `json:"pagination,omitempty"`
}

type CategoryPageResponse struct {
	Success    bool                          `json:"success"`
	Data       []CategoryInfo                `json:"data,omitempty"`
	Pagination models.PaginationDataResponse `json:"pagination,omitempty"`
}

type CategoryInfoResponse struct {
	Success bool         `json:"success"`
	Data    CategoryInfo `json:"data,omitempty"`
}
