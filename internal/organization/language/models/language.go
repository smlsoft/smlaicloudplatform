package models

import (
	"smlaicloudplatform/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const languageCollectionName = "organizationLanguages"

type Language struct {
	models.PartitionIdentity `bson:"inline"`
	Code                     string    `json:"code" bson:"code"`
	Languages                *[]string `json:"languages" bson:"languages"`
}

type LanguageInfo struct {
	models.DocIdentity `bson:"inline"`
	Language           `bson:"inline"`
}

func (LanguageInfo) CollectionName() string {
	return languageCollectionName
}

type LanguageData struct {
	models.ShopIdentity `bson:"inline"`
	LanguageInfo        `bson:"inline"`
}

type LanguageDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	LanguageData       `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (LanguageDoc) CollectionName() string {
	return languageCollectionName
}

type LanguageItemGuid struct {
	Code string `json:"code" bson:"code"`
}

func (LanguageItemGuid) CollectionName() string {
	return languageCollectionName
}

type LanguageActivity struct {
	LanguageData        `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (LanguageActivity) CollectionName() string {
	return languageCollectionName
}

type LanguageDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (LanguageDeleteActivity) CollectionName() string {
	return languageCollectionName
}
