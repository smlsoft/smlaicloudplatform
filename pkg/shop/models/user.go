package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const userCollectionName = "users"

type UserDetail struct {
	Name string `json:"name,omitempty"  validate:"required"`
}

type UsernameCode struct {
	Username string `json:"username,omitempty" bson:"username" validate:"required,gte=3"`
}

type UserPassword struct {
	Password string `json:"password,omitempty" bson:"password" validate:"required,gte=3"`
}

type UserDoc struct {
	ID           primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	UsernameCode `bson:"inline"`
	UserPassword `bson:"inline"`
	UserDetail   `bson:"inline"`
	CreatedAt    time.Time `json:"-" bson:"createdat,omitempty"`
	UpdatedAt    time.Time `json:"-" bson:"updatedat,omitempty"`
}

func (*UserDoc) CollectionName() string {
	return userCollectionName
}

type UserRequest struct {
	UsernameCode `bson:"inline"`
	UserPassword `bson:"inline"`
	UserDetail   `bson:"inline"`
}

func (*UserRequest) CollectionName() string {
	return userCollectionName
}

type UserLoginRequest struct {
	UsernameCode `bson:"inline"`
	UserPassword `bson:"inline"`
	ShopID       string `json:"shopid,omitempty"`
}

type UserProfile struct {
	UsernameCode `bson:"inline"`
	UserDetail   `bson:"inline"`
	CreatedAt    time.Time `json:"-" bson:"created_at,omitempty"`
}

type UserProfileRequest struct {
	UserDetail `bson:"inline"`
}

type UserPasswordRequest struct {
	CurrentPassword string `json:"currentpassword" bson:"currentpassword" validate:"required,gte=3"`
	NewPassword     string `json:"newpassword" bson:"newpassword" validate:"required,gte=3"`
}

type UserProfileReponse struct {
	Success bool        `json:"success"`
	Data    UserProfile `json:"data"`
}

type ShopSelectRequest struct {
	ShopID string `json:"shopid" validate:"required"`
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
