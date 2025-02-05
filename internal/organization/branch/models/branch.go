package models

import (
	"smlaicloudplatform/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const branchCollectionName = "organizationBranches"

type Branch struct {
	models.PartitionIdentity `bson:"inline"`
	Code                     string             `json:"code" bson:"code"`
	Names                    *[]models.NameX    `json:"names" bson:"names" validate:"required,min=1,unique=Code,dive"`
	Departments              *[]Department      `json:"departments" bson:"departments"`
	DataLanguage             string             `json:"datalanguage" bson:"datalanguage"`
	Address                  []models.NameX     `json:"address" bson:"addressx"`
	PhoneNumber              string             `json:"phonenumber" bson:"phonenumber"`
	Latitude                 float64            `json:"latitude" bson:"latitude"`
	Longitude                float64            `json:"longitude" bson:"longitude"`
	TaxID                    string             `json:"taxid" bson:"taxid"`
	BusinessType             BranchBusinessType `json:"businesstype" bson:"businesstype"`
	BusinessTypes            *[]string          `json:"businesstypes" bson:"businesstypes"`
	VatRate                  float64            `json:"vatrate" bson:"vatrate"`
	IsHeadOffice             bool               `json:"isheadoffice" bson:"isheadoffice"`
	IsTaxByAddress           bool               `json:"istaxbyaddress" bson:"istaxbyaddress"`
	TaxAddress               []models.NameX     `json:"taxaddress" bson:"taxaddress"`
	Email                    string             `json:"email" bson:"email"`
	VatTypeSale              int8               `json:"vattypesale" bson:"vattypesale"`
}

type BranchBusinessType struct {
	models.DocIdentity `bson:"inline"`
	Code               string          `json:"code" bson:"code"`
	Names              *[]models.NameX `json:"names" bson:"names"`
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
