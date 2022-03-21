package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ShopDoc struct {
	ID primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Identity
	Shop
	Activity
}

func (*ShopDoc) CollectionName() string {
	return "shop"
}

type Shop struct {
	Name1 string `json:"name1" bson:"name1"`
}

func (*Shop) CollectionName() string {
	return "shop"
}

type ShopInfo struct {
	GuidFixed string `json:"guidFixed" bson:"guidFixed"`
	Shop
}
