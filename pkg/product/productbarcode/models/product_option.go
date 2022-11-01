package models

import "smlcloudplatform/pkg/models"

type ProductOption struct {
	GUID      string           `json:"guid" bson:"guid"`
	MaxSelect uint16           `json:"maxselect" bson:"maxselect" validate:"min=0,max=60000"`
	MinSelect uint16           `json:"minselect" bson:"minselect" validate:"min=0,max=60000"`
	Names     *[]models.NameX  `json:"names" bson:"names" validate:"required,min=1,unique=Code,dive"`
	Choices   *[]ProductChoice `json:"choices" bson:"choices"`
}
