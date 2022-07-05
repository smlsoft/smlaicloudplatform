package models

import (
	common "smlcloudplatform/pkg/models"
	categoryModel "smlcloudplatform/pkg/product/category/models"
	inventoryModel "smlcloudplatform/pkg/product/inventory/models"
	optionModel "smlcloudplatform/pkg/product/option/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Inventory Import
const inventoryImportCollectionName = "inventoryImports"

type InventoryImport struct {
	inventoryModel.Inventory `bson:"inline"`
}

type InventoryImportInfo struct {
	common.DocIdentity `bson:"inline" `
	InventoryImport    `bson:"inline" `
}

func (InventoryImportInfo) CollectionName() string {
	return inventoryImportCollectionName
}

type InventoryImportData struct {
	common.ShopIdentity `bson:"inline" `
	InventoryImportInfo `bson:"inline" `
}

type InventoryImportDoc struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	InventoryImportData `bson:"inline"`
	common.ActivityDoc  `bson:"inline"`
}

func (InventoryImportDoc) CollectionName() string {
	return inventoryImportCollectionName
}

// Option Import

const inventoryImportOptionCollectionName string = "inventoryOptionsImports"

type InventoryOptionMainImport struct {
	optionModel.Option `bson:"inline" `
}

type InventoryOptionMainImportInfo struct {
	common.DocIdentity        `bson:"inline" `
	InventoryOptionMainImport `bson:"inline" `
}

func (InventoryOptionMainImportInfo) CollectionName() string {
	return inventoryImportOptionCollectionName
}

type InventoryOptionMainImportData struct {
	common.ShopIdentity           `bson:"inline" `
	InventoryOptionMainImportInfo `bson:"inline" `
}

type InventoryOptionMainImportDoc struct {
	ID                            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	InventoryOptionMainImportData `bson:"inline" `
	common.ActivityDoc            `bson:"inline" `
}

func (InventoryOptionMainImportDoc) CollectionName() string {
	return inventoryImportOptionCollectionName
}

// Category Import
const categoryImportCollectionName string = "categoryImports"

type CategoryImport struct {
	categoryModel.Category `bson:"inline"`
}

type CategoryImportInfo struct {
	common.DocIdentity `bson:"inline"`
	CategoryImport     `bson:"inline"`
}

func (CategoryImportInfo) CollectionName() string {
	return categoryImportCollectionName
}

type CategoryImportData struct {
	common.ShopIdentity `bson:"inline"`
	CategoryImportInfo  `bson:"inline"`
}

type CategoryImportDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CategoryImportData `bson:"inline"`
	common.ActivityDoc `bson:"inline"`
}

func (CategoryImportDoc) CollectionName() string {
	return categoryImportCollectionName
}

type CategoryImportPageResponse struct {
	Success    bool                          `json:"success"`
	Data       []CategoryImportInfo          `json:"data,omitempty"`
	Pagination common.PaginationDataResponse `json:"pagination,omitempty"`
}
