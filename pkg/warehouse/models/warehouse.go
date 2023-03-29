package models

import (
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const warehouseCollectionName = "warehouse"

type Warehouse struct {
	models.PartitionIdentity `bson:"inline"`
	Code                     string          `json:"code" bson:"code"`
	Names                    *[]models.NameX `json:"names" bson:"names" validate:"required,min=1,unique=Code,dive"`
	Location                 *[]Location     `json:"location" bson:"location" validate:"omitempty,unique=Code,dive"`
}

type Location struct {
	Code  string          `json:"code" bson:"code"`
	Names *[]models.NameX `json:"names" bson:"names" validate:"required,min=1,unique=Code,dive"`
	Shelf *[]Shelf        `json:"shelf" bson:"shelf" validate:"omitempty,unique=Code,dive"`
}

type Shelf struct {
	Code string `json:"code" bson:"code"`
	Name string `json:"name" bson:"name" validate:"required,min=1"`
}

type WarehouseInfo struct {
	models.DocIdentity `bson:"inline"`
	Warehouse          `bson:"inline"`
}

func (WarehouseInfo) CollectionName() string {
	return warehouseCollectionName
}

type WarehouseData struct {
	models.ShopIdentity `bson:"inline"`
	WarehouseInfo       `bson:"inline"`
}

type WarehouseDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	WarehouseData      `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (WarehouseDoc) CollectionName() string {
	return warehouseCollectionName
}

type WarehouseItemGuid struct {
	Code string `json:"code" bson:"code"`
}

func (WarehouseItemGuid) CollectionName() string {
	return warehouseCollectionName
}

type WarehouseActivity struct {
	WarehouseData       `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (WarehouseActivity) CollectionName() string {
	return warehouseCollectionName
}

type WarehouseDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (WarehouseDeleteActivity) CollectionName() string {
	return warehouseCollectionName
}
