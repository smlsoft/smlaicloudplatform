package models

import (
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const zoneDesignCollectionName = "ZoneDesigns"

type ZoneDesign struct {
	Code        string        `json:"code" bson:"code"`
	Order       int           `json:"order" bson:"order"`
	Tables      []TableDesign `json:"tables" bson:"tables"`
	Props       []PropDesign  `json:"props" bson:"props"`
	models.Name `bson:"inline"`
}

type ZoneDesignInfo struct {
	models.DocIdentity `bson:"inline"`
	ZoneDesign         `bson:"inline"`
}

func (ZoneDesignInfo) CollectionName() string {
	return zoneDesignCollectionName
}

type ZoneDesignData struct {
	models.ShopIdentity `bson:"inline"`
	ZoneDesignInfo      `bson:"inline"`
}

type ZoneDesignDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ZoneDesignData     `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
	models.LastUpdate  `bson:"inline"`
}

func (ZoneDesignDoc) CollectionName() string {
	return zoneDesignCollectionName
}
