package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Picture struct {
	Id primitive.ObjectID `json:"id" bson:"_id,omitempty"`

	ImageUri string `json:"image_uri,omitempty" bson:"image_uri,omitempty"`
}
