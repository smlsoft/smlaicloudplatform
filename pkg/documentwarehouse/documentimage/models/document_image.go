package models

import (
	"smlcloudplatform/pkg/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const documentImageCollectionName = "documentImages"

type DocumentImage struct {
	ImageURI        string           `json:"imageuri" bson:"imageuri"`
	Name            string           `json:"name" bson:"name"`
	IsReject        bool             `json:"isreject" bson:"isreject"`
	References      []Reference      `json:"references" bson:"references"`
	ReferenceGroups []ReferenceGroup `json:"referencegroups" bson:"referencegroups"`

	UploadedBy string    `json:"uploadedby" bson:"uploadedby"`
	UploadedAt time.Time `json:"uploadedat" bson:"uploadedat"`
	MetaFileAt time.Time `json:"metafileat" bson:"metafileat"`
}

type ReferenceGroup struct {
	GroupType  string `json:"grouptype" bson:"grouptype"`
	ParentGUID string `json:"parentguid" bson:"parentguid"`
	XOder      int    `json:"xorder" bson:"xorder"`
	XType      int    `json:"xtype" bson:"xtype"`
}

type Reference struct {
	Module string `json:"module" bson:"module"`
	DocNo  string `json:"docno" bson:"docno" `
}

type Comment struct {
	Username    string    `json:"username" bson:"username"`
	Comment     string    `json:"comment" bson:"comment"`
	CommentedAt time.Time `json:"commentedat" bson:"commentedat"`
}

type DocumentImageRequest struct {
	DocumentImage `bson:"inline"`
	Tags          *[]string `json:"tags,omitempty" bson:"tags,omitempty"`
	JobGUID       string    `json:"jobguid" bson:"jobguid"`
	PathJob       string    `json:"pathjob" bson:"pathjob"`
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

type DocumentImageItemGuid struct {
	DocumentImageGuid string `json:"categoryguid" bson:"categoryguid" gorm:"categoryguid"`
}

func (DocumentImageItemGuid) CollectionName() string {
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

type RequestDocumentImageReject struct {
	IsReject bool `json:"isreject" bson:"isreject"`
}

type DocumentImageStatus struct {
	DocGUIDRef string `json:"docguidref" bson:"docguidref"`
	Status     int8   `json:"status" bson:"status"`
}

type ImageStatus = int8

const (
	ImageNormal ImageStatus = iota
	ImageReject
	ImageCompleted
)
