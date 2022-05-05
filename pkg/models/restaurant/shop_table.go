package restaurant

import (
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const shopTableCollectionName = "shopTables"

type ShopTable struct {
	Number string    `json:"number" bson:"number"`
	Name1  string    `json:"name1" bson:"name1" gorm:"name1"`
	Name2  string    `json:"name2,omitempty" bson:"name2,omitempty"`
	Name3  string    `json:"name3,omitempty" bson:"name3,omitempty"`
	Name4  string    `json:"name4,omitempty" bson:"name4,omitempty"`
	Name5  string    `json:"name5,omitempty" bson:"name5,omitempty"`
	Seat   int8      `json:"seat" bson:"seat"`
	Zone   *ShopZone `json:"zone" bson:"zone"`
}

type ShopTableInfo struct {
	models.DocIdentity `bson:"inline"`
	ShopTable          `bson:"inline"`
}

func (ShopTableInfo) CollectionName() string {
	return shopTableCollectionName
}

type ShopTableData struct {
	models.ShopIdentity `bson:"inline"`
	ShopTableInfo       `bson:"inline"`
}

type ShopTableDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ShopTableData      `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
	models.LastUpdate  `bson:"inline"`
}

func (ShopTableDoc) CollectionName() string {
	return shopTableCollectionName
}

// Extra
type ShopTableItemGuid struct {
	Code string `json:"code" bson:"code" gorm:"code"`
}

func (ShopTableItemGuid) CollectionName() string {
	return shopTableCollectionName
}

type ShopTableActivity struct {
	ShopTableData       `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (ShopTableActivity) CollectionName() string {
	return shopTableCollectionName
}

type ShopTableDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (ShopTableDeleteActivity) CollectionName() string {
	return shopTableCollectionName
}
