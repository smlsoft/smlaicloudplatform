package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const categoryCollectionName = "categories"

type Category struct {
	Name1 string `json:"name1" bson:"name1"`
	Name2 string `json:"name2" bson:"name2"`
	Name3 string `json:"name3" bson:"name3"`
	Name4 string `json:"name4" bson:"name4"`
	Name5 string `json:"name5" bson:"name5"`
	Image string `json:"image" bson:"image"`
}

type CategoryInfo struct {
	DocIdentity
	Category
}

func (CategoryInfo) CollectionName() string {
	return categoryCollectionName
}

type CategoryData struct {
	ShopIdentity
	CategoryInfo
}

type CategoryDoc struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	CategoryData
	Activity
}

func (CategoryData) CollectionName() string {
	return categoryCollectionName
}
