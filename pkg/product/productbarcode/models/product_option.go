package models

import "smlcloudplatform/pkg/models"

type ProductOption struct {
	MaxSelect  uint16           `json:"maxselect" bson:"maxselect" validate:"min=0,max=60000"`
	Names      *[]models.NameX  `json:"names" bson:"names" validate:"required,min=1,unique=Code,dive"`
	IsRequired bool             `json:"isrequired" bson:"isrequired"`
	Choices    *[]ProductChoice `json:"choices" bson:"choices"`
}
