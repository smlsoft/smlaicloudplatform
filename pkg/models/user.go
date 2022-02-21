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

	Merchants []UserMerchant `json:"merchants " bson:"merchants"`
}

func (*User) CollectionName() string {
	return "user"
}

type UserMerchant struct {
	MerchantId string   `json:"merchantId" bson:"merchantId"`
	Role       UserRole `json:"role" bson:"role"`
}

type UserRole string

const (
	ROLE_OWNER  UserRole = "OWNER"
	ROLE_ADMIN  UserRole = "ADMIN"
	ROLE_MEMBER UserRole = "MEMBER"
)
