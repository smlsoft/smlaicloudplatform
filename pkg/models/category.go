package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const categoryCollectionName = "categories"

type Category struct {
	CategoryGuid string `json:"categoryguid" bson:"categoryguid" gorm:"categoryguid" validate:"required"`
	ParentGuid   string `json:"parentguid"  bson:"parentguid"`
	Name1        string `json:"name1" bson:"name1" validate:"required"`
	Name2        string `json:"name2,omitempty" bson:"name2,omitempty"`
	Name3        string `json:"name3,omitempty" bson:"name3,omitempty"`
	Name4        string `json:"name4,omitempty" bson:"name4,omitempty"`
	Name5        string `json:"name5,omitempty" bson:"name5,omitempty"`
	Image        string `json:"image" bson:"image,omitempty"`
	XOrder       int8   `json:"xorder" bson:"xorder,omitempty"`
}

type CategoryInfo struct {
	DocIdentity `bson:"inline"`
	Category    `bson:"inline"`
}

func (CategoryInfo) CollectionName() string {
	return categoryCollectionName
}

type CategoryData struct {
	ShopIdentity `bson:"inline"`
	CategoryInfo `bson:"inline"`
}

type CategoryDoc struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CategoryData `bson:"inline"`
	ActivityDoc  `bson:"inline"`
	LastUpdate   `bson:"inline"`
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
	Identity  `bson:"inline"`
	CreatedAt *time.Time `json:"createdat,omitempty" bson:"createdat,omitempty"`
	UpdatedAt *time.Time `json:"updatedat,omitempty" bson:"updatedat,omitempty"`
	DeletedAt *time.Time `json:"deletedat,omitempty" bson:"deletedat,omitempty"`
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
	Success    bool                         `json:"success"`
	Data       CategoryLastActivityResponse `json:"data,omitempty"`
	Pagination PaginationDataResponse       `json:"pagination,omitempty"`
}

type CategoryPageResponse struct {
	Success    bool                   `json:"success"`
	Data       []CategoryInfo         `json:"data,omitempty"`
	Pagination PaginationDataResponse `json:"pagination,omitempty"`
}

type CategoryInfoResponse struct {
	Success bool         `json:"success"`
	Data    CategoryInfo `json:"data,omitempty"`
}
