package models

import (
	"smlcloudplatform/pkg/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const bookbankCollectionName = "bookBank"

type BookBank struct {
	models.PartitionIdentity `bson:"inline"`
	PassBook                 string          `json:"passbook" bson:"passbook"`
	BankCode                 string          `json:"bankcode" bson:"bankcode"`
	Names                    *[]models.NameX `json:"names" bson:"names" validate:"required,min=1,unique=Code,dive"`
}

type BookBankInfo struct {
	models.DocIdentity `bson:"inline"`
	BookBank           `bson:"inline"`
}

func (BookBankInfo) CollectionName() string {
	return bookbankCollectionName
}

type BookBankData struct {
	models.ShopIdentity `bson:"inline"`
	BookBankInfo        `bson:"inline"`
}

type BookBankDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	BookBankData       `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (BookBankDoc) CollectionName() string {
	return bookbankCollectionName
}

type BookBankItemGuid struct {
	PassBook string `json:"passbook" bson:"passbook"`
}

func (BookBankItemGuid) CollectionName() string {
	return bookbankCollectionName
}

type BookBankActivity struct {
	BookBankData        `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (BookBankActivity) CollectionName() string {
	return bookbankCollectionName
}

type BookBankDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (BookBankDeleteActivity) CollectionName() string {
	return bookbankCollectionName
}
