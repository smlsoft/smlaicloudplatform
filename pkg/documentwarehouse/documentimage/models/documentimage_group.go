package models

import (
	"smlcloudplatform/pkg/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const documentImageGroupCollectionName = "documentImageGroups"

const (
	IMAGE_PENDING = iota
	IMAGE_CHECKED
	IMAGE_REJECT
	IMAGE_BANNED
	IMAGE_REJECT_KEYING
)

type DocumentImageGroup struct {
	// DocumentRef     string            `json:"documentref" bson:"documentref"`
	Title               string            `json:"title" bson:"title"`
	References          []Reference       `json:"references" bson:"references"`
	Tags                *[]string         `json:"tags,omitempty" bson:"tags,omitempty"`
	ImageReferences     *[]ImageReference `json:"imagereferences" bson:"imagereferences"`
	UploadedBy          string            `json:"uploadedby" bson:"uploadedby"`
	UploadedAt          time.Time         `json:"uploadedat" bson:"uploadedat"`
	Status              int8              `json:"status" bson:"status"`
	Description         string            `json:"description" bson:"description"`
	TaskGUID            string            `json:"taskguid" bson:"taskguid" validate:"required,min=1"`
	PathTask            string            `json:"pathtask" bson:"pathtask"`
	IsTaskCompleted     bool              `json:"iscompleted" bson:"iscompleted"`
	RejectFromGroupGUID string            `json:"rejectfromgroupguid" bson:"rejectfromgroupguid"`
	XOrder              int               `json:"xorder" bson:"xorder"`
	RejectRemark        string            `json:"rejectremark" bson:"rejectremark"`
}

type DocumentImageGroupBody struct {
	DocumentImageGroup `bson:"inline"`
	ImageReferences    *[]ImageReferenceBody `json:"imagereferences,omitempty" bson:"imagereferences,omitempty"`
}

type ImageReferenceBody struct {
	XOrder            int    `json:"xorder" bson:"xorder"`
	DocumentImageGUID string `json:"documentimageguid" bson:"documentimageguid"`
}
type ImageReference struct {
	ImageReferenceBody `bson:",inline"`
	ImageURI           string `json:"imageuri" bson:"imageuri"`
	ImageEditURI       string `json:"imageedituri" bson:"imageedituri"`
	Name               string `json:"name" bson:"name"`
	// IsReject           bool      `json:"isreject" bson:"isreject"`
	UploadedBy string    `json:"uploadedby" bson:"uploadedby"`
	UploadedAt time.Time `json:"uploadedat" bson:"uploadedat"`
	MetaFileAt time.Time `json:"metafileat" bson:"metafileat"`
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

type Status struct {
	Status int8 `json:"status"`
}
