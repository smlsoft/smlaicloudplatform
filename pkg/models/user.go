package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `json:"-" bson:"_id,omitempty"`

	Username string `json:"username" bson:"username"`

	Password string `json:"password" bson:"password"`

	Name string `json:"name,omitempty" bson:"name"`

	CreatedAt time.Time `json:"-" bson:"created_at,omitempty"`
}

func (*User) CollectionName() string {
	return "user"
}

type UserProfile struct {
	Username string `json:"username" bson:"username"`

	Name string `json:"name,omitempty" bson:"name"`

	CreatedAt time.Time `json:"-" bson:"created_at,omitempty"`
}

type MerchantSelectRequest struct {
	MerchantId string `json:"merchantId"`
}

type UserRole string

const (
	ROLE_OWNER UserRole = "OWNER"
	ROLE_ADMIN UserRole = "ADMIN"
	ROLE_User  UserRole = "User"
)

type MerchantUser struct {
	Id         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username   string             `json:"username" bson:"username"`
	MerchantId string             `json:"merchantId" bson:"merchantId"`
	Role       UserRole           `json:"role" bson:"role"`
}

func (*MerchantUser) CollectionName() string {
	return "merchantUsers"
}
