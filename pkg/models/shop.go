package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ShopDoc struct {
	ID primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	ShopInfo
	Activity
}

func (ShopDoc) CollectionName() string {
	return "shops"
}

type ShopInfo struct {
	DocIdentity
	Shop
}

func (ShopInfo) CollectionName() string {
	return "shops"
}

type Shop struct {
	Name1 string `json:"name1" bson:"name1"`
}

func (Shop) CollectionName() string {
	return "shops"
}
