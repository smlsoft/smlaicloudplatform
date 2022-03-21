package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ShopDoc struct {
	ID primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	ShopInfo
	Activity
}

func (*ShopDoc) CollectionName() string {
	return "shop"
}

type ShopInfo struct {
	DocIdentity
	Shop
}

type Shop struct {
	Name1 string `json:"name1" bson:"name1"`
}

func (*Shop) CollectionName() string {
	return "shop"
}
