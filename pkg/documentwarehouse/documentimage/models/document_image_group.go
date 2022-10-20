package models

import (
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const documentImageGroupCollectionName = "documentImageGroup"

type DocumentImageGroup struct {
	// DocumentRef     string            `json:"documentref" bson:"documentref"`
	Title           string            `json:"title" bson:"title"`
	References      *[]Reference      `json:"references,omitempty" bson:"references,omitempty"`
	Tags            *[]string         `json:"tags,omitempty" bson:"tags,omitempty"`
	ImageReferences *[]ImageReference `json:"imagereferences,omitempty" bson:"imagereferences,omitempty"`
}

type ImageReference struct {
	XOrder            int    `json:"xorder" bson:"xorder"`
	DocumentImageGUID string `json:"documentimageguid" bson:"documentimageguid"`
}

func (DocumentImageGroup) CollectionName() string {
	return documentImageGroupCollectionName
}

type DocumentImageGroupInfo struct {
	models.DocIdentity `bson:"inline"`
	DocumentImageGroup `bson:"inline"`
}

func (DocumentImageGroupInfo) CollectionName() string {
	return documentImageGroupCollectionName
}

type DocumentImageGroupData struct {
	models.ShopIdentity    `bson:"inline"`
	DocumentImageGroupInfo `bson:"inline"`
}

type DocumentImageGroupDoc struct {
	ID                     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	DocumentImageGroupData `bson:"inline"`
	models.ActivityDoc     `bson:"inline"`
	models.LastUpdate      `bson:"inline"`
}

func (DocumentImageGroupDoc) CollectionName() string {
	return documentImageGroupCollectionName
}
