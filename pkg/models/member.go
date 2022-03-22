package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const memberCollectionName = "members"

type Member struct {
	Telephone    string `json:"telephone" bson:"Telephone"`
	Name         string `json:"name,omitempty" bson:"name,omitempty"`
	Surname      string `json:"surname,omitempty" bson:"surname,omitempty"`
	TaxID        string `json:"TaxID,omitempty" bson:"TaxID,omitempty"`
	ContactType  int    `json:"contactType,omitempty" bson:"ContactType,omitempty"`
	PersonalType int    `json:"personalType,omitempty" bson:"PersonalType,omitempty"`
	BranchType   int    `json:"branchType,omitempty" bson:"BranchType,omitempty"`
	BranchCode   string `json:"branchCode,omitempty" bson:"branchCode,omitempty"`
	Address      string `json:"address,omitempty" bson:"address,omitempty"`
	ZipCode      string `json:"zipCode,omitempty" bson:"zipCode,omitempty"`
}

type MemberInfo struct {
	DocIdentity `bson:"inline"`
	Member      `bson:"inline"`
}

func (MemberInfo) CollectionName() string {
	return memberCollectionName
}

type MemberData struct {
	ShopIdentity `bson:"inline"`
	MemberInfo   `bson:"inline"`
}
type MemberDoc struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`
	MemberData `bson:"inline"`
	Activity   `bson:"inline"`
}

func (MemberDoc) CollectionName() string {
	return memberCollectionName
}

type MemberRequestEdit struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
}

func (MemberRequestEdit) CollectionName() string {
	return memberCollectionName
}

type MemberRequestPassword struct {
	Password string `json:"password" validate:"required,gte=6"`
}

func (MemberRequestPassword) CollectionName() string {
	return memberCollectionName
}
