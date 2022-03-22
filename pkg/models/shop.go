package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const shopCollectionName = "shops"

type ShopDoc struct {
	ID primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	ShopInfo
	Activity
}

func (ShopDoc) CollectionName() string {
	return shopCollectionName
}

type ShopInfo struct {
	DocIdentity
	Shop
}

func (ShopInfo) CollectionName() string {
	return shopCollectionName
}

type Shop struct {
	Name1 string `json:"name1" bson:"name1"`
}
