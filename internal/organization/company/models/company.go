package models

import (
	"smlaicloudplatform/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const companyCollectionName = "organizationCompany"

type Company struct {
	models.PartitionIdentity `bson:"inline"`
	Code                     string          `json:"code" bson:"code"`
	CompanyNames             *[]models.NameX `json:"companynames" bson:"companynames"`
	ApiKey                   string          `json:"apikey" bson:"apikey"`
	ImageURI                 string          `json:"imageuri" bson:"imageuri"`
	LogoURI                  string          `json:"logouri" bson:"logouri"`
}

type CompanyBusinessType struct {
	models.DocIdentity `bson:"inline"`
	Code               string          `json:"code" bson:"code"`
	Names              *[]models.NameX `json:"names" bson:"names"`
}

type CompanyInfo struct {
	models.DocIdentity `bson:"inline"`
	Company            `bson:"inline"`
}

func (CompanyInfo) CollectionName() string {
	return companyCollectionName
}

type CompanyData struct {
	models.ShopIdentity `bson:"inline"`
	CompanyInfo         `bson:"inline"`
}

type CompanyDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CompanyData        `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (CompanyDoc) CollectionName() string {
	return companyCollectionName
}

type CompanyItemGuid struct {
	Code string `json:"code" bson:"code"`
}

func (CompanyItemGuid) CollectionName() string {
	return companyCollectionName
}

type CompanyActivity struct {
	CompanyData         `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (CompanyActivity) CollectionName() string {
	return companyCollectionName
}

type CompanyDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (CompanyDeleteActivity) CollectionName() string {
	return companyCollectionName
}

type CompanyInfoResponse struct {
	CompanyInfo
}
