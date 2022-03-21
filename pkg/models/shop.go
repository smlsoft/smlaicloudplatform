package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ShopDoc struct {
	ID primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Shop
	Activity
}

func (*ShopDoc) CollectionName() string {
	return "shop"
}

type Shop struct {
	GuidFixed string `json:"guidFixed" bson:"GuidFixed"`
	Name1     string `json:"name1" bson:"name1"`
}

func (*Shop) CollectionName() string {
	return "shop"
}
