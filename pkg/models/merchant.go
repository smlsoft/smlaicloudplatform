package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Merchant struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	GuidFixed string             `json:"-" bson:"guidFixed"`
	Name1     string             `json:"name1" bson:"name1"`
	CreatedBy string             `json:"-" bson:"createdBy"`
	CreatedAt time.Time          `json:"-" bson:"createdAt"`
	UpdatedBy string             `json:"-" bson:"updatedBy,omitempty"`
	UpdatedAt time.Time          `json:"-" bson:"updatedAt,omitempty"`
	Deleted   bool               `json:"-" bson:"deleted"`
}

func (*Merchant) CollectionName() string {
	return "merchant"
}

type MerchantInfo struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	GuidFixed   string             `json:"guidFixed" bson:"guidFixed"`
	Name1       string             `json:"name1" bson:"name1"`
	TotalMember int                `json:"totalMember"`
}
