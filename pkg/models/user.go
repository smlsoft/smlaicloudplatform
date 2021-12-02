package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `json:"-" bson:"_id,omitempty"`

	Username string `json:"username,omitempty" bson:"username"`

	Password string `json:"password,omitempty" bson:"password"`

	Name string `json:"name,omitempty" bson:"name"`

	CreatedAt time.Time `json:"-" bson:"created_at,omitempty"`
}

func (*User) CollectionName() string {
	return "user"
}
