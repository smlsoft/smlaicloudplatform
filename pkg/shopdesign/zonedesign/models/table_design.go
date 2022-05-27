package models

import "smlcloudplatform/pkg/models"

type TableDesign struct {
	Index        int `json:"index" bson:"index" validate:"required"`
	ObjectDesign `bson:"inline"`
	Number       string `json:"number" bson:"number"`
	Chair        int    `json:"chair" bson:"chair" validate:"required"`
	models.Name  `bson:"inline"`
}
