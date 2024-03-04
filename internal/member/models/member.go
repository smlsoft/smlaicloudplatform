package models

import (
	"smlcloudplatform/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const memberCollectionName = "members"
const memberIndexName string = "members_index"

type Member struct {
	PictureUrl   string `json:"pictureurl" bson:"pictureurl"`
	Telephone    string `json:"telephone" bson:"telephone"`
	Name         string `json:"name" bson:"name"`
	Surname      string `json:"surname" bson:"surname"`
	TaxID        string `json:"taxid" bson:"taxid"`
	ContactType  int    `json:"contacttype" bson:"contacttype"`
	PersonalType int    `json:"personaltype" bson:"personaltype"`
	BranchType   int    `json:"branchtype" bson:"branchtype"`
	BranchCode   string `json:"branchcode" bson:"branchcode"`
	LineUID      string `json:"lineuid" bson:"lineuid"`
	MemberAdress `bson:"inline"`
	SubAddress   []MemberAdress `json:"subaddress" bson:"subaddress"`
}

type MemberAdress struct {
	Telephone    string `json:"telephone" bson:"telephone"`
	Address      string `json:"address" bson:"address"`
	CountryCode  string `json:"countrycode" bson:"countrycode"`
	ProvinceCode string `json:"provincecode" bson:"provincecode"`
	DistrictCode string `json:"districtcode" bson:"districtcode"`
	ZipCode      string `json:"zipcode" bson:"zipcode"`
}

type MemberInfo struct {
	models.DocIdentity `bson:"inline"`
	Member             `bson:"inline"`
}

func (MemberInfo) CollectionName() string {
	return memberCollectionName
}

type MemberData struct {
	models.ShopIdentity `bson:"inline"`
	MemberInfo          `bson:"inline"`
}
type MemberDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	MemberData         `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
	models.LastUpdate  `bson:"inline"`
}

func (MemberDoc) CollectionName() string {
	return memberCollectionName
}

type MemberIndex struct {
	models.Index `bson:"inline"`
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
	MemberData          `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (MemberActivity) CollectionName() string {
	return memberCollectionName
}

type MemberDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (MemberDeleteActivity) CollectionName() string {
	return memberCollectionName
}

type MemberLastActivityResponse struct {
	New    []MemberActivity       `json:"new" `
	Remove []MemberDeleteActivity `json:"remove"`
}

type MemberFetchUpdateResponse struct {
	Success    bool                          `json:"success"`
	Data       MemberLastActivityResponse    `json:"data,omitempty"`
	Pagination models.PaginationDataResponse `json:"pagination,omitempty"`
}

type MemberInfoResponse struct {
	Success bool       `json:"success"`
	Data    MemberInfo `json:"data,omitempty"`
}

type MemberPageResponse struct {
	Success    bool                          `json:"success"`
	Data       []MemberInfo                  `json:"data,omitempty"`
	Pagination models.PaginationDataResponse `json:"pagination,omitempty"`
}
