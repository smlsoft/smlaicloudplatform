package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Username  string             `json:"username" bson:"username"`
	Password  string             `json:"password" bson:"password"`
	Name      string             `json:"name,omitempty" bson:"name"`
	CreatedAt time.Time          `json:"-" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"-" bson:"updatedAt,omitempty"`
}

func (*User) CollectionName() string {
	return "user"
}

type UserRequest struct {
	Username string `json:"username,omitempty" `

	Password string `json:"password,omitempty" `

	Name string `json:"name,omitempty" `
}

func (*UserRequest) CollectionName() string {
	return "user"
}

type UserLoginRequest struct {
	Username string `json:"username,omitempty" `
	Password string `json:"password,omitempty" `
	ShopId   string `json:"shopId,omitempty"`
}

type UserProfile struct {
	Username  string    `json:"username" bson:"username"`
	Name      string    `json:"name,omitempty" bson:"name"`
	CreatedAt time.Time `json:"-" bson:"created_at,omitempty"`
}

type UserPasswordRequest struct {
	CurrentPassword string `json:"currentPassword" bson:"currentPassword"`
	NewPassword     string `json:"newPassword" bson:"newPassword"`
}

type ShopSelectRequest struct {
	ShopId string `json:"shopId"`
}

const (
	ROLE_OWNER string = "OWNER"
	ROLE_ADMIN string = "ADMIN"
	ROLE_USER  string = "USER"
)

type ShopUser struct {
	Id       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	ShopId   string             `json:"shopId" bson:"shopId"`
	Role     string             `json:"role" bson:"role"`
}

func (*ShopUser) CollectionName() string {
	return "shopUsers"
}

type ShopUserInfo struct {
	ShopId string `json:"shopId" bson:"shopId"`
	Name   string `json:"name" bson:"name"`
	Role   string `json:"role" bson:"role"`
}

func (*ShopUserInfo) CollectionName() string {
	return "shopUsers"
}
