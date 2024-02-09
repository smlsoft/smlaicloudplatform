package models

import (
	pkgModels "smlcloudplatform/internal/models"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SaleInvoiceReturnTransactionPG struct {
	TransactionPG    `gorm:"embedded;"`
	DebtorCode       string                                  `json:"creditorcode" gorm:"column:creditorcode"`
	DebtorNames      pkgModels.JSONB                         `json:"creditornames" gorm:"column:creditornames;type:jsonb"`
	TotalPayCash     float64                                 `json:"totalpaycash" gorm:"column:totalpaycash"`
	TotalPayTransfer float64                                 `json:"totalpaytransfer" gorm:"column:totalpaytransfer"`
	TotalPayCredit   float64                                 `json:"totalpaycredit" gorm:"column:totalpaycredit"`
	Items            *[]SaleInvoiceReturnTransactionDetailPG `json:"items" gorm:"items;foreignKey:shopid,docno"`

	SaleCode                     string          `json:"salecode" gorm:"column:salecode"`
	SaleName                     string          `json:"salename" gorm:"column:salename"`
	BranchCode                   string          `json:"branchcode" gorm:"column:branchcode"`
	BranchNames                  pkgModels.JSONB `json:"branchnames" gorm:"column:branchnames;type:jsonb"`
	DetailDiscountFormula        string          `json:"detaildiscountformula" gorm:"column:detaildiscountformula"`
	DetailTotalAmount            float64         `json:"detailtotalamount" gorm:"column:detailtotalamount"`
	TotalDiscountVatAmount       float64         `json:"totaldiscountvatamount" gorm:"column:totaldiscountvatamount"`
	TotalDiscountExceptVatAmount float64         `json:"totaldiscountexceptvatamount" gorm:"column:totaldiscountexceptvatamount"`
	DetailTotalDiscount          float64         `json:"detailtotaldiscount" gorm:"column:detailtotaldiscount"`
}

type SaleInvoiceReturnTransactionDetailPG struct {
	TransactionDetailPG `gorm:"embedded;"`
}

// table name

func (t *SaleInvoiceReturnTransactionPG) TableName() string {
	return "saleinvoice_return_transaction"
}

func (t *SaleInvoiceReturnTransactionDetailPG) TableName() string {
	return "saleinvoice_return_transaction_detail"
}

func (t SaleInvoiceReturnTransactionPG) HasDebtorEffectDoc() bool {

	hasCreditorEffectDoc := t.InquiryType == 0 || t.InquiryType == 1
	return hasCreditorEffectDoc
}

func (t SaleInvoiceReturnTransactionPG) HasStockEffectDoc() bool {
	hasStockEffectDoc := t.InquiryType == 0 || t.InquiryType == 2
	return hasStockEffectDoc
}

func (j *SaleInvoiceReturnTransactionPG) BeforeUpdate(tx *gorm.DB) (err error) {

	// find old data
	var details *[]SaleInvoiceReturnTransactionDetailPG
	tx.Model(&SaleInvoiceReturnTransactionDetailPG{}).Where(" shopid=? AND docno=?", j.ShopID, j.DocNo).Find(&details)

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
			tx.Delete(&SaleInvoiceReturnTransactionDetailPG{}, tmp.ID)
		}
	}

	return nil
}

func (jd *SaleInvoiceReturnTransactionDetailPG) BeforeCreate(tx *gorm.DB) error {

	tx.Statement.AddClause(clause.OnConflict{
		UpdateAll: true,
	})
	return nil
}

func (s *SaleInvoiceReturnTransactionPG) CompareTo(other *SaleInvoiceReturnTransactionPG) bool {

	diff := cmp.Diff(s, other,
		cmpopts.IgnoreFields(SaleInvoiceReturnTransactionDetailPG{}, "ID"),
	)

	if diff == "" {
		return true
	}

	return false
}
