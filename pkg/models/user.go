package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const userCollectionName = "users"

type User struct {
	ID        primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username"`
	Password  string             `json:"password" bson:"password"`
	Name      string             `json:"name,omitempty" bson:"name"`
	CreatedAt time.Time          `json:"-" bson:"createdat,omitempty"`
	UpdatedAt time.Time          `json:"-" bson:"updatedat,omitempty"`
}

func (*User) CollectionName() string {
	return userCollectionName
}

type UserRequest struct {
	Username string `json:"username,omitempty" validate:"required,gte=3"`
	Password string `json:"password,omitempty" validate:"required,gte=6"`
	Name     string `json:"name,omitempty" `
}

func (*UserRequest) CollectionName() string {
	return userCollectionName
}

type UserLoginRequest struct {
	Username string `json:"username" validate:"required,gte=3"`
	Password string `json:"password" validate:"required,gte=6"`
	ShopID   string `json:"shopid,omitempty"`
}

type UserProfile struct {
	Username  string    `json:"username" bson:"username"`
	Name      string    `json:"name,omitempty" bson:"name"`
	CreatedAt time.Time `json:"-" bson:"created_at,omitempty"`
}

type UserPasswordRequest struct {
	CurrentPassword string `json:"currentpassword" bson:"currentpassword"`
	NewPassword     string `json:"newpassword" bson:"newpassword"`
}

type ShopSelectRequest struct {
	ShopID string `json:"shopid"`
}

const (
	ROLE_OWNER string = "OWNER"
	ROLE_ADMIN string = "ADMIN"
	ROLE_USER  string = "USER"
)

type ShopUser struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	ShopID   string             `json:"shopid" bson:"shopid"`
	Role     string             `json:"role" bson:"role"`
}

func (*ShopUser) CollectionName() string {
	return "shopUsers"
}

type ShopUserInfo struct {
	ShopID string `json:"shopid" bson:"shopid"`
	Name   string `json:"name" bson:"name"`
	Role   string `json:"role" bson:"role"`
}

func (*ShopUserInfo) CollectionName() string {
	return "shopUsers"
}
