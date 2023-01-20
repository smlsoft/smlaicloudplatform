package models

import (
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const jobCollectionName = "fileFolder"

type Job struct {
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

type JobInfo struct {
	models.DocIdentity `bson:"inline"`
	Job                `bson:"inline"`
}

func (JobInfo) CollectionName() string {
	return jobCollectionName
}

type JobData struct {
	models.ShopIdentity `bson:"inline"`
	JobInfo             `bson:"inline"`
}

type JobDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	JobData            `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (JobDoc) CollectionName() string {
	return jobCollectionName
}

type JobItemGuid struct {
	Name string `json:"name" bson:"name"`
}

func (JobItemGuid) CollectionName() string {
	return jobCollectionName
}

type JobActivity struct {
	JobData             `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (JobActivity) CollectionName() string {
	return jobCollectionName
}

type JobDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (JobDeleteActivity) CollectionName() string {
	return jobCollectionName
}
