package models

import (
	modelsCustomergroup "smlcloudplatform/pkg/customershop/customergroup/models"
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const customerCollectionName = "customershopCustomers"

type Customer struct {
	models.PartitionIdentity `bson:"inline"`
	Code                     string          `json:"code" bson:"code" validate:"required,min=1"`
	PersonalType             int8            `json:"personaltype" bson:"personaltype"`
	Images                   *[]Image        `json:"images" bson:"images"`
	Names                    *[]models.NameX `json:"names" bson:"names" validate:"required,min=1,unique=Code,dive"`

	AddressForBilling  CustomerAddress    `json:"addressforbilling" bson:"addressforbilling"`
	AddressForShipping *[]CustomerAddress `json:"addressforshipping" bson:"addressforshipping"`
	TaxId              string             `json:"taxid" bson:"taxid"`
	Email              string             `json:"email" bson:"email"`
	CustomerType       int                `json:"customertype" bson:"customertype"`
	BranchNumber       string             `json:"branchnumber" bson:"branchnumber"`
}

type CustomerAddress struct {
	GUID            string          `json:"guid" bson:"guid"`
	Address         *[]string       `json:"address" bson:"address"`
	CountryCode     string          `json:"countrycode" bson:"countrycode"`
	ProvinceCode    string          `json:"provincecode" bson:"provincecode"`
	DistrictCode    string          `json:"districtcode" bson:"districtcode"`
	SubDistrictCode string          `json:"subdistrictcode" bson:"subdistrictcode"`
	ZipCode         string          `json:"zipcode" bson:"zipcode"`
	ContactNames    *[]models.NameX `json:"contactnames" bson:"contactnames"`
	PhonePrimary    string          `json:"phoneprimary" bson:"phoneprimary"`
	PhoneSecondary  string          `json:"phonesecondary" bson:"phonesecondary"`
	Latitude        float64         `json:"latitude" bson:"latitude"`
	Longitude       float64         `json:"longitude" bson:"longitude"`
}

type Image struct {
	XOrder int    `json:"xorder" bson:"xorder"`
	URI    string `json:"uri" bson:"uri"`
}

type CustomerRequest struct {
	Customer
	Groups []CustomerGroupRequest `json:"groups"`
}

type CustomerGroupRequest struct {
	GuidFixed string `json:"guidfixed"`
}

type CustomerInfo struct {
	models.DocIdentity `bson:"inline"`
	Customer           `bson:"inline"`
	Groups             *[]modelsCustomergroup.CustomerGroupInfo `json:"groups" bson:"groups"`
}

func (CustomerInfo) CollectionName() string {
	return customerCollectionName
}

type CustomerData struct {
	models.ShopIdentity `bson:"inline"`
	CustomerInfo        `bson:"inline"`
}

type CustomerDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CustomerData       `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (CustomerDoc) CollectionName() string {
	return customerCollectionName
}

type CustomerItemGuid struct {
	Code string `json:"code" bson:"code"`
}

func (CustomerItemGuid) CollectionName() string {
	return customerCollectionName
}

type CustomerActivity struct {
	CustomerData        `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (CustomerActivity) CollectionName() string {
	return customerCollectionName
}

type CustomerDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (CustomerDeleteActivity) CollectionName() string {
	return customerCollectionName
}
