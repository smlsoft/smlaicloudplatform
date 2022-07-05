package models

import (
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const shopCollectionName = "shops"

type ShopDoc struct {
	ID                 primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	ShopInfo           `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (ShopDoc) CollectionName() string {
	return shopCollectionName
}

type ShopInfo struct {
	models.DocIdentity `bson:"inline"`
	Shop               `bson:"inline"`
}

func (ShopInfo) CollectionName() string {
	return shopCollectionName
}

type Shop struct {
	Name1     string `json:"name1" bson:"name1"`
	Telephone string `json:"telephone" bson:"telephone"`
}
