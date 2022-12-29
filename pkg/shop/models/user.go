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
	Username string `json:"username,omitempty" bson:"username" validate:"required,gte=5,max=233"`
}

type UserPassword struct {
	Password string `json:"password,omitempty" bson:"password" validate:"required,gte=5,max=233"`
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
	CurrentPassword string `json:"currentpassword" bson:"currentpassword" validate:"required,gte=5"`
	NewPassword     string `json:"newpassword" bson:"newpassword" validate:"required,gte=5"`
}

type UserProfileReponse struct {
	Success bool        `json:"success"`
	Data    UserProfile `json:"data"`
}

type ShopSelectRequest struct {
	ShopID string `json:"shopid" validate:"required"`
}

type UserRole = uint8

const (
	ROLE_USER  UserRole = iota // "USER"
	ROLE_ADMIN                 // "ADMIN"
	ROLE_OWNER                 // "OWNER"

	ROLE_SYSTEM = 255 // APP MANAGER
)

type ShopUser struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username       string             `json:"username" bson:"username"`
	ShopID         string             `json:"shopid" bson:"shopid"`
	Role           UserRole           `json:"role" bson:"role"`
	IsFavorite     bool               `json:"isfavorite" bson:"isfavorite"`
	LastAccessedAt time.Time          `json:"lastaccessedat" bson:"lastaccessedat"`
}

func (*ShopUser) CollectionName() string {
	return "shopUsers"
}

type ShopUserInfo struct {
	ShopID         string    `json:"shopid" bson:"shopid"`
	Name           string    `json:"name" bson:"name"`
	BranchCode     string    `json:"branchcode" bson:"branchcode"`
	Role           UserRole  `json:"role" bson:"role"`
	IsFavorite     bool      `json:"isfavorite" bson:"isfavorite"`
	LastAccessedAt time.Time `json:"lastaccessedat" bson:"lastaccessedat"`
}

func (*ShopUserInfo) CollectionName() string {
	return "shopUsers"
}

type UserRoleRequest struct {
	ShopID   string   `json:"shopid" bson:"shopid"`
	Username string   `json:"username" bson:"username"`
	Role     UserRole `json:"role" bson:"role"`
}

// func (u UserRole) EqualString(userRoleStr string)  bool {
// 	switch u {
// 		case
// 	}
// }
