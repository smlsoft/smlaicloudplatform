package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const memberCollectionName = "members"
const memberIndexName string = "members_index"

type Member struct {
	Telephone    string `json:"telephone" bson:"telephone"`
	Name         string `json:"name" bson:"name"`
	Surname      string `json:"surname" bson:"surname"`
	TaxID        string `json:"taxid" bson:"taxid"`
	ContactType  int    `json:"contacttype" bson:"contacttype"`
	PersonalType int    `json:"personaltype" bson:"personaltype"`
	BranchType   int    `json:"branchtype" bson:"branchtype"`
	BranchCode   string `json:"branchcode" bson:"branchcode"`
	Address      string `json:"address" bson:"address"`
	ZipCode      string `json:"zipcode" bson:"zipcode"`
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
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	MemberData  `bson:"inline"`
	ActivityDoc `bson:"inline"`
	LastUpdate  `bson:"inline"`
}

func (MemberDoc) CollectionName() string {
	return memberCollectionName
}

type MemberIndex struct {
	Index `bson:"inline"`
}

func (MemberIndex) TableName() string {
	return memberIndexName
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

type MemberActivity struct {
	MemberData   `bson:"inline"`
	ActivityTime `bson:"inline"`
}

func (MemberActivity) CollectionName() string {
	return memberCollectionName
}

type MemberDeleteActivity struct {
	Identity     `bson:"inline"`
	ActivityTime `bson:"inline"`
}

func (MemberDeleteActivity) CollectionName() string {
	return memberCollectionName
}

type MemberLastActivityResponse struct {
	New    []MemberActivity       `json:"new" `
	Remove []MemberDeleteActivity `json:"remove"`
}

type MemberFetchUpdateResponse struct {
	Success    bool                       `json:"success"`
	Data       MemberLastActivityResponse `json:"data,omitempty"`
	Pagination PaginationDataResponse     `json:"pagination,omitempty"`
}
