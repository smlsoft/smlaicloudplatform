package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MemberDoc struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	Identity
	Member
	Activity
}

func (*MemberDoc) CollectionName() string {
	return "member"
}

type Member struct {
	Telephone    string `json:"telephone" bson:"Telephone"`
	Name         string `json:"name,omitempty" bson:"name,omitempty"`
	Surname      string `json:"surname,omitempty" bson:"surname,omitempty"`
	TaxId        string `json:"TaxId,omitempty" bson:"TaxId,omitempty"`
	ContactType  int    `json:"contactType,omitempty" bson:"ContactType,omitempty"`
	PersonalType int    `json:"personalType,omitempty" bson:"PersonalType,omitempty"`
	BranchType   int    `json:"branchType,omitempty" bson:"BranchType,omitempty"`
	BranchCode   string `json:"branchCode,omitempty" bson:"branchCode,omitempty"`
	Address      string `json:"address,omitempty" bson:"address,omitempty"`
	ZipCode      string `json:"zipCode,omitempty" bson:"zipCode,omitempty"`
}

func (*Member) CollectionName() string {
	return "member"
}

type MemberInfo struct {
	ID       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ShopID   string             `json:"shop_id" bson:"shop_id"`
	Name     string             `json:"name,omitempty" bson:"name"`
	Email    string             `json:"email" bson:"email"`
	Username string             `json:"username" bson:"username"`
	Activity
}

func (*MemberInfo) CollectionName() string {
	return "member"
}

type MemberRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,gte=6"`
}

func (*MemberRequest) CollectionName() string {
	return "member"
}

type MemberRequestEdit struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
}

func (*MemberRequestEdit) CollectionName() string {
	return "member"
}

type MemberRequestPassword struct {
	Password string `json:"password" validate:"required,gte=6"`
}

func (*MemberRequestPassword) CollectionName() string {
	return "member"
}
