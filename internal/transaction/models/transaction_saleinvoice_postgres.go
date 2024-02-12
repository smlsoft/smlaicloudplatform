package models

import (
	pkgModels "smlcloudplatform/internal/models"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SaleInvoiceTransactionPG struct {
	TransactionPG    `gorm:"embedded;"`
	DebtorCode       string          `json:"creditorcode" gorm:"column:creditorcode"`
	DebtorNames      pkgModels.JSONB `json:"creditornames" gorm:"column:creditornames;type:jsonb"`
	TotalPayCash     float64         `json:"totalpaycash" gorm:"column:totalpaycash"`
	TotalPayTransfer float64         `json:"totalpaytransfer" gorm:"column:totalpaytransfer"`
	TotalPayCredit   float64         `json:"totalpaycredit" gorm:"column:totalpaycredit"`
	IsPOS            bool            `json:"ispos" gorm:"column:ispos"`

	SaleCode                     string          `json:"salecode" gorm:"column:salecode"`
	SaleName                     string          `json:"salename" gorm:"column:salename" `
	BranchCode                   string          `json:"branchcode" gorm:"column:branchcode" `
	BranchNames                  pkgModels.JSONB `json:"branchnames" gorm:"column:branchnames;type:jsonb"`
	DetailDiscountFormula        string          `json:"detaildiscountformula" gorm:"column:detaildiscountformula"`
	DetailTotalAmount            float64         `json:"detailtotalamount" gorm:"column:detailtotalamount"`
	TotalDiscountVatAmount       float64         `json:"totaldiscountvatamount" gorm:"column:totaldiscountvatamount"`
	TotalDiscountExceptVatAmount float64         `json:"totaldiscountexceptvatamount" gorm:"column:totaldiscountexceptvatamount"`
	DetailTotalDiscount          float64         `json:"detailtotaldiscount" gorm:"column:detailtotaldiscount"`

	Items *[]SaleInvoiceTransactionDetailPG `json:"items" gorm:"items;foreignKey:shopid,docno"`
}

type SaleInvoiceTransactionDetailPG struct {
	TransactionDetailPG `gorm:"embedded;"`
	ManufacturerGUID    string          `json:"manufacturerguid" gorm:"column:manufacturerguid"`
	ManufacturerCode    string          `json:"manufacturercode" gorm:"column:manufacturercode"`
	ManufacturerNames   pkgModels.JSONB `json:"manufacturernames" gorm:"column:manufacturernames"`
}

// tableName
func (SaleInvoiceTransactionPG) TableName() string {
	return "saleinvoice_transaction"
}

func (SaleInvoiceTransactionDetailPG) TableName() string {
	return "saleinvoice_transaction_detail"
}

func (t SaleInvoiceTransactionPG) HasCreditorEffectDoc() bool {

	hasCreditorEffectDoc := t.InquiryType == 0
	return hasCreditorEffectDoc
}

func (j *SaleInvoiceTransactionPG) BeforeUpdate(tx *gorm.DB) (err error) {

	// find old data
	var details *[]SaleInvoiceTransactionDetailPG
	tx.Model(&SaleInvoiceTransactionDetailPG{}).Where(" shopid=? AND docno=?", j.ShopID, j.DocNo).Find(&details)

	// delete un use data
	for _, tmp := range *details {
		var foundUpdate bool = false
		for _, data := range *j.Items {
			if data.ID == tmp.ID {
				foundUpdate = true
			}
		}
		if foundUpdate == false {
			// mark delete
			tx.Delete(&SaleInvoiceTransactionDetailPG{}, tmp.ID)
		}
	}

	return nil
}

func (jd *SaleInvoiceTransactionDetailPG) BeforeCreate(tx *gorm.DB) error {

	tx.Statement.AddClause(clause.OnConflict{
		UpdateAll: true,
	})
	return nil
}

func (s *SaleInvoiceTransactionPG) CompareTo(other *SaleInvoiceTransactionPG) bool {

	diff := cmp.Diff(s, other,
		cmpopts.IgnoreFields(SaleInvoiceTransactionDetailPG{}, "ID"),
	)

	if diff == "" {
		return true
	}

	return false
}
