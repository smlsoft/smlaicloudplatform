package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MemberInfo struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	MerchantID string `json:"merchant_id" bson:"merchant_id"`

	Name string `json:"name,omitempty" bson:"name"`

	Email string `json:"email" bson:"email"`

	Username string `json:"username" bson:"username"`

	CreatedBy string `json:"-" bson:"created_by"`

	CreatedAt time.Time `json:"-" bson:"created_at,omitempty"`
}

func (*MemberInfo) CollectionName() string {
	return "member"
}
