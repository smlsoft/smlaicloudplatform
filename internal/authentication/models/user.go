package models

import (
	"time"

	"smlcloudplatform/internal/models"
	timezone "smlcloudplatform/internal/models/timezone"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const userCollectionName = "users"

type UserDetail struct {
	UID               string `json:"uid" bson:"uid"`
	Name              string `json:"name,omitempty"`
	Avatar            string `json:"avatar"`
	timezone.Timezone `bson:"inline"`
	YearType          string   `json:"yeartype" bson:"yeartype" validate:"max=21"`
	DedeZoom          DedeZoom `json:"dedezoom" bson:"dedezoom"`
	RegisterType      string   `json:"registertype" bson:"registertype"`
}

type DedeZoom struct {
	Email       string `json:"email" bson:"email"`
	PhoneNumber string `json:"phonenumber" bson:"phonenumber"`
	Address     string `json:"address" bson:"address"`
}

type UsernameField struct {
	Username string `json:"username,omitempty" bson:"username" validate:"required,gte=3,max=233"` // validate:"required,alphanum,gte=3,max=233"
}

type PhoneNumberField struct {
	CountryCode string `json:"countrycode" bson:"countrycode" validate:"required,max=20"`
	PhoneNumber string `json:"phonenumber" bson:"phonenumber" validate:"required,max=100"`
}

type EmailField struct {
	Email string `json:"email,omitempty" bson:"email" validate:"required,email,max=233"`
}

type UserPassword struct {
	Password string `json:"password,omitempty" bson:"password" validate:"required,gte=5,max=233"`
}

type UserDoc struct {
	ID               primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	UsernameField    `bson:"inline"`
	EmailField       `bson:"inline"`
	PhoneNumberField `bson:"inline"`
	UserPassword     `bson:"inline"`
	UserDetail       `bson:"inline"`
	CreatedAt        time.Time `json:"-" bson:"createdat,omitempty"`
	UpdatedAt        time.Time `json:"-" bson:"updatedat,omitempty"`
}

func (*UserDoc) CollectionName() string {
	return userCollectionName
}

type RegisterEmailRequest struct {
	EmailField
	UserPassword
	UserDetail
}

type RegisterPhoneNumberRequest struct {
	PhoneNumberField
	UsernameField
	UserPassword
	UserDetail
	OTPVerifyRequest
}

type ForgotPasswordPhoneNumberRequest struct {
	PhoneNumberField
	UserPassword
	OTPVerifyRequest
}

func (RegisterEmailRequest) CollectionName() string {
	return userCollectionName
}

func (RegisterPhoneNumberRequest) CollectionName() string {
	return userCollectionName
}

type UserRequest struct {
	UsernameField `bson:"inline"`
	UserPassword  `bson:"inline"`
	UserDetail    `bson:"inline"`
}

func (*UserRequest) CollectionName() string {
	return userCollectionName
}

type UserLoginRequest struct {
	UsernameField `bson:"inline"`
	UserPassword  `bson:"inline"`
	ShopID        string `json:"shopid,omitempty"`
}

type UserLoginPhoneNumberRequest struct {
	PhoneNumberField `bson:"inline"`
	UserPassword     `bson:"inline"`
	ShopID           string `json:"shopid,omitempty"`
}

type UserProfile struct {
	UsernameField `bson:"inline"`
	UserDetail    `bson:"inline"`
	CreatedAt     time.Time `json:"-" bson:"created_at,omitempty"`
}

func (UserProfile) CollectionName() string {
	return userCollectionName
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

type ShopUserBase struct {
	Username string   `json:"username" bson:"username"`
	ShopID   string   `json:"shopid" bson:"shopid"`
	Role     UserRole `json:"role" bson:"role"`
}

type ShopUser struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ShopUserBase   `bson:"inline"`
	IsFavorite     bool      `json:"isfavorite" bson:"isfavorite"`
	LastAccessedAt time.Time `json:"lastaccessedat" bson:"lastaccessedat"`
}

func (*ShopUser) CollectionName() string {
	return "shopUsers"
}

type ShopUserInfo struct {
	ShopID string `json:"shopid" bson:"shopid"`
	Name   string `json:"name" bson:"name1"`
	// Name1          string         `json:"name1" bson:"name1"`
	Names          []models.NameX `json:"names" bson:"names"`
	BranchCode     string         `json:"branchcode" bson:"branchcode"`
	Role           UserRole       `json:"role" bson:"role"`
	IsFavorite     bool           `json:"isfavorite" bson:"isfavorite"`
	LastAccessedAt time.Time      `json:"lastaccessedat" bson:"lastaccessedat"`
	CreatedBy      string         `json:"createdby" bson:"createdby"`
}

func (*ShopUserInfo) CollectionName() string {
	return "shopUsers"
}

type UserRoleRequest struct {
	ShopID       string   `json:"shopid" bson:"shopid"`
	EditUsername string   `json:"editusername" bson:"editusername"`
	Username     string   `json:"username" bson:"username"`
	Role         UserRole `json:"role" bson:"role"`
}

type ShopUserAccessLog struct {
	ID             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ShopID         string             `json:"shopid" bson:"shopid"`
	Username       string             `json:"username" bson:"username"`
	Ip             string             `json:"ip" bson:"ip"`
	LastAccessedAt time.Time          `json:"lastaccessedat" bson:"lastaccessedat"`
}

func (*ShopUserAccessLog) CollectionName() string {
	return "shopUserAccessLogs"
}

type ShopUserProfile struct {
	ShopUserBase    `bson:"inline"`
	UserProfileName string `json:"userprofilename" bson:"userprofilename"`
}
