package models

import (
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const filefolderCollectionName = "fileFolder"

type FileFolder struct {
	models.PartitionIdentity `bson:"inline"`
	Name                     string    `json:"name" bson:"name"`
	Module                   string    `json:"module" bson:"module"`
	Status                   int8      `json:"status" bson:"status"`
	ParentGUIDFixed          string    `json:"parentguidfixed" bson:"parentguidfixed"`
	Path                     string    `json:"path" bson:"path"`
	IsFavorit                bool      `json:"isfavorit" bson:"isfavorit"`
	Tags                     *[]string `json:"tags" bson:"tags"`
	ToTal                    int       `json:"total" bson:"total"`
}

type FileFolderInfo struct {
	models.DocIdentity `bson:"inline"`
	FileFolder         `bson:"inline"`
}

func (FileFolderInfo) CollectionName() string {
	return filefolderCollectionName
}

type FileFolderData struct {
	models.ShopIdentity `bson:"inline"`
	FileFolderInfo      `bson:"inline"`
}

type FileFolderDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FileFolderData     `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (FileFolderDoc) CollectionName() string {
	return filefolderCollectionName
}

type FileFolderItemGuid struct {
	Name string `json:"name" bson:"name"`
}

func (FileFolderItemGuid) CollectionName() string {
	return filefolderCollectionName
}

type FileFolderActivity struct {
	FileFolderData      `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (FileFolderActivity) CollectionName() string {
	return filefolderCollectionName
}

type FileFolderDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (FileFolderDeleteActivity) CollectionName() string {
	return filefolderCollectionName
}
