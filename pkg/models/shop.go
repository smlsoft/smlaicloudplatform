package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Shop struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty" default:"-"`
	GuidFixed string             `json:"-" bson:"guidFixed"`
	Name1     string             `json:"name1" bson:"name1"`
	Activity
}

func (*Shop) CollectionName() string {
	return "shop"
}

type ShopInfo struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	GuidFixed string             `json:"guidFixed" bson:"guidFixed"`
	Name1     string             `json:"name1" bson:"name1"`
}
