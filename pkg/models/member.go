package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Member struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	GuidFixed string             `json:"guidFixed" bson:"GuidFixed"`
	ShopID    string             `json:"shopID" bson:"shop_id"`
	Name      string             `json:"name,omitempty" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Username  string             `json:"username" bson:"username"`
	Password  string             `json:"password" bson:"password"`
	Activity
}

func (*Member) CollectionName() string {
	return "member"
}
