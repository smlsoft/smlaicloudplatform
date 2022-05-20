package vfgl

import (
	"smlcloudplatform/pkg/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const journalCollectionName = "journals"

type Journal struct {
	Docno              string          `json:"docno" bson:"docno"`
	DocDate            time.Time       `json:"docdate" bson:"docdate"`
	AccountPeriod      int16           `json:"accountperiod" bson:"accountperiod"`
	AccountYear        int16           `json:"accountyear" bson:"accountyear"`
	AccountGroup       string          `json:"accountgroup" bson:"accountgroup"`
	AccountBook        []JournalDetail `json:"journaldetail" bson:"journaldetail"`
	Amount             float64         `json:"amount" bson:"amount"`
	AccountDescription string          `json:"accountdescription" bson:"accountdescription"`
}

type JournalDetail struct {
	AccountCode  string  `json:"accountcode" bson:"accountcode"`
	AccountName  string  `json:"accountname" bson:"accountname"`
	DebitAmount  float64 `json:"debitamount" bson:"debitamount"`
	CreditAmount float64 `json:"creditamount" bson:"creditamount"`
}

type JournalInfo struct {
	models.DocIdentity `bson:"inline"`
	Journal            `bson:"inline"`
}

func (JournalInfo) CollectionName() string {
	return journalCollectionName
}

type JournalData struct {
	models.ShopIdentity `bson:"inline"`
	JournalInfo         `bson:"inline"`
}

type JournalDoc struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	JournalData        `bson:"inline"`
	models.ActivityDoc `bson:"inline"`
}

func (JournalDoc) CollectionName() string {
	return journalCollectionName
}

type JournalItemGuid struct {
	Docno string `json:"docno" bson:"docno" gorm:"docno"`
}

func (JournalItemGuid) CollectionName() string {
	return journalCollectionName
}

type JournalActivity struct {
	JournalData         `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (JournalActivity) CollectionName() string {
	return journalCollectionName
}

type JournalDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (JournalDeleteActivity) CollectionName() string {
	return journalCollectionName
}

// Postgresql model
type JournalPg struct {
	Docno              string    `json:"docno" gorm:"docno"`
	DocDate            time.Time `json:"docdate" gorm:"docdate"`
	AccountPeriod      int16     `json:"accountperiod" gorm:"accountperiod"`
	AccountYear        int16     `json:"accountyear" gorm:"accountyear"`
	AccountGroup       string    `json:"accountgroup" gorm:"accountgroup"`
	Amount             float64   `json:"amount" gorm:"amount"`
	AccountDescription string    `json:"accountdescription" gorm:"accountdescription"`
}

type JournalDetailPg struct {
	AccountCode  string  `json:"accountcode" gorm:"accountcode"`
	AccountName  string  `json:"accountname" gorm:"accountname"`
	DebitAmount  float64 `json:"debitamount" gorm:"debitamount"`
	CreditAmount float64 `json:"creditamount" gorm:"creditamount"`
}
