package restaurant

import (
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const shopZoneCollectionName = "shopzones"

type ShopZone struct {
	Code  string `json:"code" bson:"code"`
	Name1 string `json:"name1" bson:"name1" gorm:"name1"`
	Name2 string `json:"name2,omitempty" bson:"name2,omitempty"`
	Name3 string `json:"name3,omitempty" bson:"name3,omitempty"`
	Name4 string `json:"name4,omitempty" bson:"name4,omitempty"`
	Name5 string `json:"name5,omitempty" bson:"name5,omitempty"`
}

type ShopZoneInfo struct {
	models.DocIdentity `bson:"inline"`
	ShopZone           `bson:"inline"`
}

func (ShopZoneInfo) CollectionName() string {
	return shopZoneCollectionName
}

type ShopZoneData struct {
	models.ShopIdentity `bson:"inline"`
	ShopZoneInfo        `bson:"inline"`
}

type ShopZoneDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ShopZoneData       `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (ShopZoneDoc) CollectionName() string {
	return shopZoneCollectionName
}
