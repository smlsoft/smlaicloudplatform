package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Member struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	GuidFixed string `json:"guidFixed" bson:"GuidFixed"`

	MerchantID string `json:"merchantID" bson:"merchant_id"`

	Name string `json:"name,omitempty" bson:"name"`

	Email string `json:"email" bson:"email"`

	Username string `json:"username" bson:"username"`

	Password string `json:"password" bson:"password"`

	CreatedBy string `json:"-" bson:"created_by"`

	CreatedAt time.Time `json:"-" bson:"created_at,omitempty"`

	UpdatedBy string `json:"-" bson:"updatedBy"`

	UpdatedAt time.Time `json:"-" bson:"updatedAt"`
}

func (*Member) CollectionName() string {
	return "member"
}
