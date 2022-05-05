package restaurant

import (
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const kitchenCollectionName = "kitchens"

type Kitchen struct {
	Code     string              `json:"code" bson:"code"`
	Name1    string              `json:"name1" bson:"name1" gorm:"name1"`
	Name2    string              `json:"name2,omitempty" bson:"name2,omitempty"`
	Name3    string              `json:"name3,omitempty" bson:"name3,omitempty"`
	Name4    string              `json:"name4,omitempty" bson:"name4,omitempty"`
	Name5    string              `json:"name5,omitempty" bson:"name5,omitempty"`
	Printers *[]PrinterTerminal  `json:"printers" bson:"printers"`
	Products *[]models.Inventory `json:"products" bson:"products"`
	Zones    *[]ShopZone         `json:"zones" bson:"zones"`
	Category *models.Category    `json:"category" bson:"category"`
}

type KitchenInfo struct {
	models.DocIdentity `bson:"inline"`
	Kitchen            `bson:"inline"`
}

func (KitchenInfo) CollectionName() string {
	return kitchenCollectionName
}

type KitchenData struct {
	models.ShopIdentity `bson:"inline"`
	KitchenInfo         `bson:"inline"`
}

type KitchenDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	KitchenData        `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
	models.LastUpdate  `bson:"inline"`
}

func (KitchenDoc) CollectionName() string {
	return kitchenCollectionName
}

type KitchenItemGuid struct {
	Code string `json:"code" bson:"code" gorm:"code"`
}

func (KitchenItemGuid) CollectionName() string {
	return kitchenCollectionName
}

type KitchenActivity struct {
	KitchenData         `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (KitchenActivity) CollectionName() string {
	return kitchenCollectionName
}

type KitchenDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (KitchenDeleteActivity) CollectionName() string {
	return kitchenCollectionName
}
