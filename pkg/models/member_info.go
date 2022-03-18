package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MemberInfo struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ShopID   string             `json:"shop_id" bson:"shop_id"`
	Name     string             `json:"name,omitempty" bson:"name"`
	Email    string             `json:"email" bson:"email"`
	Username string             `json:"username" bson:"username"`
	Activity
}

func (*MemberInfo) CollectionName() string {
	return "member"
}
