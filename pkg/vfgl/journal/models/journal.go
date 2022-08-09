package models

import (
	"smlcloudplatform/pkg/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const journalCollectionName = "journals"

type JournalBody struct {
	BatchID            string    `json:"batchId" bson:"batch" gorm:"column:batchid"`
	DocNo              string    `json:"docno" bson:"docno" gorm:"column:docno;primaryKey"`
	DocDate            time.Time `json:"docdate" bson:"docdate" format:"dateTime" gorm:"column:docdate"`
	DocumentRef        string    `json:"documentref" bson:"documentref" gorm:"column:documentref"`
	AccountPeriod      int16     `json:"accountperiod" bson:"accountperiod" gorm:"column:accountperiod"`
	AccountYear        int16     `json:"accountyear" bson:"accountyear" gorm:"column:accountyear"`
	AccountGroup       string    `json:"accountgroup" bson:"accountgroup" gorm:"column:accountgroup"`
	Amount             float64   `json:"amount" bson:"amount" gorm:"column:amount"`
	AccountDescription string    `json:"accountdescription" bson:"accountdescription" gorm:"column:accountdescription"`
	BookCode           string    `json:"bookcode" bson:"bookcode"`
	Vats               []Vat     `json:"vats" bson:"vats"`
	Taxes              []Tax     `json:"taxes" bson:"taxes"`
}

type Journal struct {
	JournalBody              `bson:"inline"`
	models.PartitionIdentity `bson:"inline"`
	AccountBook              *[]JournalDetail `json:"journaldetail" bson:"journaldetail"`
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
	DocNo string `json:"docno" bson:"docno" gorm:"docno"`
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

type TaxVatCustomer struct {
	CustTaxID    string `json:"custtaxid" bson:"custtaxid"`
	CustName     string `json:"custname" bson:"custname"`
	CustType     int8   `json:"custtype" bson:"custtype"`
	Organization int8   `json:"organization" bson:"organization"`
	BranchCode   string `json:"branchcode" bson:"branchcode"`
	Address      string `json:"address" bson:"address"`
}

type Vat struct {
	VatDocNo       string    `json:"vatdocno" bson:"vatdocno"`
	VatDate        time.Time `json:"vatdate" bson:"vatdate"`
	VatType        int8      `json:"vattype" bson:"vattype"`
	VatMode        int8      `json:"vatmode" bson:"vatmode"`
	VatPeriod      int8      `json:"vatperiod" bson:"vatperiod"`
	VatYear        int16     `json:"vatyear" bson:"vatyear" `
	VatBase        float64   `json:"vatbase" bson:"vatbase"`
	VatRate        float64   `json:"vatrate" bson:"vatrate"`
	VatAmount      float64   `json:"vatamount" bson:"vatamount"`
	ExceptVat      float64   `json:"exceptvat" bson:"exceptvat"`
	VatSubmit      bool      `json:"vatsubmit" bson:"vatsubmit"`
	Remark         string    `json:"remark" bson:"remark"`
	TaxVatCustomer `bson:"inline"`
}

type Tax struct {
	TaxDocNo       string    `json:"taxdocno" bson:"taxdocno"`
	TaxDate        time.Time `json:"taxdate" gorm:"column:taxdate"`
	TaxType        int8      `json:"taxtype" bson:"taxtype"`
	TaxAmount      float64   `json:"taxamount" bson:"taxamount"`
	TaxVatCustomer `bson:"inline"`
	Details        *[]TaxDetail `json:"details" bson:"details"`
}

type TaxDetail struct {
	TaxBase     float64 `json:"taxbase" bson:"taxbase"`
	TaxRate     float64 `json:"taxrate" bson:"taxrate"`
	TaxAmount   float64 `json:"taxamount" bson:"taxamount"`
	Description string  `json:"description" bson:"description"`
}

// Postgresql model
type JournalPg struct {
	models.ShopIdentity      `gorm:"embedded;"`
	models.PartitionIdentity `gorm:"embedded;"`
	JournalBody              `gorm:"embedded;"`
	// Docno                    string             `json:"docno" gorm:"column:docno;primaryKey"`
	// BatchID                  string             `json:"batchid" gorm:"column:batchid"`
	// DocDate                  time.Time          `json:"docdate" gorm:"column:docdate"`
	// AccountPeriod            int16              `json:"accountperiod" gorm:"column:accountperiod"`
	// AccountYear              int16              `json:"accountyear" gorm:"column:accountyear"`
	// AccountGroup             string             `json:"accountgroup" gorm:"column:accountgroup"`
	// Amount                   float64            `json:"amount" gorm:"column:amount"`
	// AccountDescription       string             `json:"accountdescription" gorm:"column:accountdescription"`
	AccountBook *[]JournalDetailPg `json:"journaldetail" gorm:"journals_detail;foreignKey:shopid,docno"`
}

func (JournalPg) TableName() string {
	return "journals"
}

type JournalDetailPg struct {
	ID                       uint   `gorm:"primarykey"`
	ShopID                   string `json:"shopid" gorm:"column:shopid"`
	models.PartitionIdentity `gorm:"embedded;"`
	Docno                    string  `json:"docno" gorm:"column:docno"`
	AccountCode              string  `json:"accountcode" gorm:"column:accountcode"`
	AccountName              string  `json:"accountname" gorm:"column:accountname"`
	DebitAmount              float64 `json:"debitamount" gorm:"column:debitamount"`
	CreditAmount             float64 `json:"creditamount" gorm:"column:creditamount"`
}

func (JournalDetailPg) TableName() string {
	return "journals_detail"
}

type JournalInfoResponse struct {
	Success bool        `json:"success"`
	Data    JournalInfo `json:"data,omitempty"`
}

type JournalPageResponse struct {
	Success    bool                          `json:"success"`
	Data       []JournalInfo                 `json:"data,omitempty"`
	Pagination models.PaginationDataResponse `json:"pagination,omitempty"`
}

func (j *JournalPg) BeforeUpdate(tx *gorm.DB) (err error) {

	// find old data
	var details *[]JournalDetailPg
	tx.Model(&JournalDetailPg{}).Where(" shopid=? AND docno=?", j.ShopID, j.DocNo).Find(&details)

	// delete unuse data
	for _, tmp := range *details {
		var foundUpdate bool = false
		for _, data := range *j.AccountBook {
			if data.ID == tmp.ID {
				foundUpdate = true
			}
		}
		if foundUpdate == false {
			// mark delete
			tx.Delete(&JournalDetailPg{}, tmp.ID)
		}
	}

	return nil
}

func (jd *JournalDetailPg) BeforeCreate(tx *gorm.DB) error {

	tx.Statement.AddClause(clause.OnConflict{
		UpdateAll: true,
	})
	return nil
}
