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
	Username string `json:"username,omitempty" validate:"required,gte=3"`
	Password string `json:"password,omitempty" validate:"required,gte=6"`
	Name     string `json:"name,omitempty" `
}

func (*UserRequest) CollectionName() string {
	return "user"
}

type UserLoginRequest struct {
	Username string `json:"username" validate:"required,gte=3"`
	Password string `json:"password" validate:"required,gte=6"`
	ShopID   string `json:"shopID,omitempty"`
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
	ShopID string `json:"shopID"`
}

const (
	ROLE_OWNER string = "OWNER"
	ROLE_ADMIN string = "ADMIN"
	ROLE_USER  string = "USER"
)

type ShopUser struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string             `json:"username" bson:"username"`
	ShopID   string             `json:"shopID" bson:"shopID"`
	Role     string             `json:"role" bson:"role"`
}

func (*ShopUser) CollectionName() string {
	return "shopUsers"
}

type ShopUserInfo struct {
	ShopID string `json:"shopID" bson:"shopID"`
	Name   string `json:"name" bson:"name"`
	Role   string `json:"role" bson:"role"`
}

func (*ShopUserInfo) CollectionName() string {
	return "shopUsers"
}
