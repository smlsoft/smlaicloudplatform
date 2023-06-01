package models

import (
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const branchCollectionName = "organizationBranches"

type Branch struct {
	models.PartitionIdentity `bson:"inline"`
	Code                     string          `json:"code" bson:"code"`
	Names                    *[]models.NameX `json:"names" bson:"names" validate:"required,min=1,unique=Code,dive"`
	Departments              *[]Department   `json:"departments" bson:"departments"`
	BusinessTypes            *[]string       `json:"businesstypes" bson:"businesstypes"`
	Languages                *[]string       `json:"languages" bson:"languages"`
	Contact                  Contact         `json:"contact" bson:"contact"`
}

type Contact struct {
	Address         string  `json:"address" bson:"address"`
	CountryCode     string  `json:"countrycode" bson:"countrycode"`
	ProvinceCode    string  `json:"provincecode" bson:"provincecode"`
	DistrictCode    string  `json:"districtcode" bson:"districtcode"`
	SubDistrictCode string  `json:"subdistrictcode" bson:"subdistrictcode"`
	ZipCode         string  `json:"zipcode" bson:"zipcode"`
	PhoneNumber     string  `json:"phonenumber" bson:"phonenumber"`
	Latitude        float64 `json:"latitude" bson:"latitude"`
	Longitude       float64 `json:"longitude" bson:"longitude"`
}

type BranchInfo struct {
	models.DocIdentity `bson:"inline"`
	Branch             `bson:"inline"`
}

func (BranchInfo) CollectionName() string {
	return branchCollectionName
}

type BranchData struct {
	models.ShopIdentity `bson:"inline"`
	BranchInfo          `bson:"inline"`
}

type BranchDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	BranchData         `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (BranchDoc) CollectionName() string {
	return branchCollectionName
}

type BranchItemGuid struct {
	Code string `json:"code" bson:"code"`
}

func (BranchItemGuid) CollectionName() string {
	return branchCollectionName
}

type BranchActivity struct {
	BranchData          `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (BranchActivity) CollectionName() string {
	return branchCollectionName
}

type BranchDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (BranchDeleteActivity) CollectionName() string {
	return branchCollectionName
}

type BranchInfoResponse struct {
	BranchInfo
	// Departments   []Department   `json:"departments" bson:"departments"`
	BusinessTypes []BusinessType `json:"businesstypes" bson:"businesstypes"`
}

type Department struct {
	// GuidFixed string         `json:"guidfixed"`
	Code  string         `json:"code"`
	Names []models.NameX `json:"names"`
}

type BusinessType struct {
	GuidFixed string         `json:"guidfixed"`
	Code      string         `json:"code"`
	Names     []models.NameX `json:"names"`
}
