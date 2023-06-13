package models

import (
	"smlcloudplatform/pkg/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StockTransaction struct {
	models.ShopIdentity `bson:"inline"`
	DocNo               string                    `json:"docno" gorm:"column:docno;primaryKey"`
	DocDate             time.Time                 `json:"docdate" bson:"docdate" format:"dateTime" gorm:"column:docdate"`
	TransFlag           int8                      `json:"transflag" gorm:"column:transflag" `
	Details             *[]StockTransactionDetail `json:"details" gorm:"details;foreignKey:shopid,docno"`
}

func (StockTransaction) TableName() string {
	return "stock_transaction"
}

type StockTransactionDetail struct {
	ID                       uint   `gorm:"primarykey"`
	ShopID                   string `json:"shopid" gorm:"column:shopid"`
	models.PartitionIdentity `gorm:"embedded;"`
	DocNo                    string  `json:"docno" gorm:"column:docno"`
	Barcode                  string  `json:"barcode" gorm:"column:barcode"`
	RefBarcode               string  `json:"refbarcode" gorm:"column:refbarcode"`
	Qty                      float64 `json:"qty" gorm:"column:qty"`
	StandValue               float64 `json:"standvalue" gorm:"column:standvalue"`
	DivideValue              float64 `json:"dividevalue" gorm:"column:dividevalue"`
	CalcFlag                 int8    `json:"calcflag" gorm:"column:calcflag"`
	SumOfCost                float64 `json:"sumofcost" gorm:"column:sumofcost"`
	AverageCost              float64 `json:"averagecost" gorm:"column:averagecost"`
	LineNumber               int8    `json:"linenumber" gorm:"column:linenumber"`
	CostPerUnit              float64 `json:"costperunit" gorm:"column:costperunit"`
	TotalCost                float64 `json:"totalcost" gorm:"column:totalcost"`
}

func (StockTransactionDetail) TableName() string {
	return "stock_transaction_detail"
}

func (j *StockTransaction) BeforeUpdate(tx *gorm.DB) (err error) {

	// find old data
	var details *[]StockTransactionDetail
	tx.Model(&StockTransactionDetail{}).Where(" shopid=? AND docno=?", j.ShopID, j.DocNo).Find(&details)

	// delete unuse data
	for _, tmp := range *details {
		var foundUpdate bool = false
		for _, data := range *j.Details {
			if data.ID == tmp.ID {
				foundUpdate = true
			}
		}
		if foundUpdate == false {
			// mark delete
			tx.Delete(&StockTransactionDetail{}, tmp.ID)
		}
	}

	return nil
}

func (jd *StockTransactionDetail) BeforeCreate(tx *gorm.DB) error {

	tx.Statement.AddClause(clause.OnConflict{
		UpdateAll: true,
	})
	return nil
}
