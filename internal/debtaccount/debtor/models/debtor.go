package models

import (
	groupModels "smlcloudplatform/internal/debtaccount/debtorgroup/models"
	"smlcloudplatform/internal/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const debtorCollectionName = "debtors"

type Debtor struct {
	models.PartitionIdentity `bson:"inline"`
	Code                     string `json:"code" bson:"code"`

	PersonalType int8            `json:"personaltype" bson:"personaltype"`
	Images       *[]Image        `json:"images" bson:"images"`
	Names        *[]models.NameX `json:"names" bson:"names" validate:"required,min=1,unique=Code,dive"`

	AddressForBilling  Address    `json:"addressforbilling" bson:"addressforbilling"`
	AddressForShipping *[]Address `json:"addressforshipping" bson:"addressforshipping"`
	TaxId              string     `json:"taxid" bson:"taxid"`
	Email              string     `json:"email" bson:"email"`
	CustomerType       int        `json:"customertype" bson:"customertype"`
	BranchNumber       string     `json:"branchnumber" bson:"branchnumber"`
	FundCode           string     `json:"fundcode" bson:"fundcode"`
	CreditDay          int        `json:"creditday" bson:"creditday"`
	IsMember           bool       `json:"ismember" bson:"ismember"`
	GroupGUIDs         *[]string  `json:"-" bson:"groups"`
	Auth               DebtorAuth `json:"auth" bson:"auth"`
}

type DebtorAuth struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

type Address struct {
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

type DebtorRequest struct {
	Debtor
	Groups []string `json:"groups"`
}

type DebtorInfo struct {
	models.DocIdentity `bson:"inline"`
	Debtor             `bson:"inline"`
	Groups             *[]groupModels.DebtorGroupInfo `json:"groups" bson:"-"`
}

func (DebtorInfo) CollectionName() string {
	return debtorCollectionName
}

type DebtorData struct {
	models.ShopIdentity `bson:"inline"`
	DebtorInfo          `bson:"inline"`
}

type DebtorDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	DebtorData         `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (DebtorDoc) CollectionName() string {
	return debtorCollectionName
}

type DebtorItemGuid struct {
	Code string `json:"code" bson:"code"`
}

func (DebtorItemGuid) CollectionName() string {
	return debtorCollectionName
}

type DebtorActivity struct {
	DebtorData          `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (DebtorActivity) CollectionName() string {
	return debtorCollectionName
}

type DebtorDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (DebtorDeleteActivity) CollectionName() string {
	return debtorCollectionName
}
