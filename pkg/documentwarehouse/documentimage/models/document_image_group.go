package models

import "time"

const documentImageGroupCollectionName = "documentImages"

type DocumentImageGroup struct {
	DocumentRef    string                      `json:"documentref" bson:"documentref"`
	DocumentImages *[]DocumentImageGroupDetail `json:"documentimages" bson:"documentimages"`
}

func (DocumentImageGroup) CollectionName() string {
	return documentImageGroupCollectionName
}

type DocumentImageGroupDetail struct {
	GuidFixed  string    `json:"guidfixed" bson:"guidfixed"`
	Name       string    `json:"name" bson:"name"`
	ImageUri   string    `json:"imageuri" bson:"imageuri"`
	Module     string    `json:"module" bson:"module"`
	DocGUIDRef string    `json:"docguidref" bson:"docguidref"`
	Status     int8      `json:"status" bson:"status"`
	UploadedBy string    `json:"uploadedby" bson:"uploadedby"`
	UploadedAt time.Time `json:"uploadedat" bson:"uploadedat"`
}

type DocumentImageGroupRequest struct {
	DocumentRef    string   `json:"documentref" bson:"documentref"`
	DocumentImages []string `json:"documentimages" bson:"documentimages"`
}

func (DocumentImageGroupRequest) CollectionName() string {
	return documentImageGroupCollectionName
}
