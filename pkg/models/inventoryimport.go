package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Inventory Import
const inventoryImportCollectionName = "inventoryImports"

type InventoryImport struct {
	Inventory `bson:"inline"`
}

type InventoryImportInfo struct {
	DocIdentity     `bson:"inline" `
	InventoryImport `bson:"inline" `
}

func (InventoryImportInfo) CollectionName() string {
	return inventoryImportCollectionName
}

type InventoryImportData struct {
	ShopIdentity        `bson:"inline" `
	InventoryImportInfo `bson:"inline" `
}

type InventoryImportDoc struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	InventoryImportData `bson:"inline"`
	Activity            `bson:"inline"`
}

func (InventoryImportDoc) CollectionName() string {
	return inventoryImportCollectionName
}

// Option Import

const inventoryImportOptionCollectionName string = "inventoryOptionsImports"

type InventoryOptionMainImport struct {
	Option `bson:"inline" `
}

type InventoryOptionMainImportInfo struct {
	DocIdentity               `bson:"inline" `
	InventoryOptionMainImport `bson:"inline" `
}

func (InventoryOptionMainImportInfo) CollectionName() string {
	return inventoryImportOptionCollectionName
}

type InventoryOptionMainImportData struct {
	ShopIdentity                  `bson:"inline" `
	InventoryOptionMainImportInfo `bson:"inline" `
}

type InventoryOptionMainImportDoc struct {
	ID                            primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	InventoryOptionMainImportData `bson:"inline" `
	Activity                      `bson:"inline" `
}

func (InventoryOptionMainImportDoc) CollectionName() string {
	return inventoryImportOptionCollectionName
}

// Category Import
const categoryImportCollectionName string = "categoryImports"

type CategoryImport struct {
	Code     string `json:"code" bson:"code"`
	Category `bson:"inline"`
}

type CategoryImportInfo struct {
	DocIdentity    `bson:"inline"`
	CategoryImport `bson:"inline"`
}

func (CategoryImportInfo) CollectionName() string {
	return categoryImportCollectionName
}

type CategoryImportData struct {
	ShopIdentity       `bson:"inline"`
	CategoryImportInfo `bson:"inline"`
}

type CategoryImportDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CategoryImportData `bson:"inline"`
	Activity           `bson:"inline"`
}

func (CategoryImportDoc) CollectionName() string {
	return categoryImportCollectionName
}
