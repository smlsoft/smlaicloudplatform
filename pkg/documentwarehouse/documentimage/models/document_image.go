package models

import (
	"smlcloudplatform/pkg/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const documentImageCollectionName = "documentImages"

type DocumentImage struct {
	DocumentRef string    `json:"documentref" bson:"documentref"`
	ImageUri    string    `json:"imageuri" bson:"imageuri"`
	Module      string    `json:"module" bson:"module"`
	DocGUIDRef  string    `json:"docguidref" bson:"docguidref"`
	Status      int8      `json:"status" bson:"status"`
	UploadedBy  string    `json:"uploadedby" bson:"uploadedby"`
	UploadedAt  time.Time `json:"uploadedat" bson:"uploadedat"`
}
type DocumentImageInfo struct {
	models.DocIdentity `bson:"inline"`
	DocumentImage      `bson:"inline"`
}

func (DocumentImageInfo) CollectionName() string {
	return documentImageCollectionName
}

type DocumentImageData struct {
	models.ShopIdentity `bson:"inline"`
	DocumentImageInfo   `bson:"inline"`
}

type DocumentImageDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	DocumentImageData  `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
	models.LastUpdate  `bson:"inline"`
}

func (DocumentImageDoc) CollectionName() string {
	return documentImageCollectionName
}

type DocumentImageInfoResponse struct {
	Success bool              `json:"success"`
	Data    DocumentImageInfo `json:"data,omitempty"`
}

type DocumentImagePageResponse struct {
	Success    bool                          `json:"success"`
	Data       []DocumentImageInfo           `json:"data,omitempty"`
	Pagination models.PaginationDataResponse `json:"pagination,omitempty"`
}

type DocumentImageStatus struct {
	Status int8 `json:"status" bson:"status"`
}

type ImageStatus = int8

const (
	ImageNormal ImageStatus = iota
	ImageReject
	ImageCompleted
)
