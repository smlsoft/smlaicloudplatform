package models

import (
	"smlcloudplatform/internal/models"
	pkgModels "smlcloudplatform/internal/models"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type DebtorTransactionPG struct {
	models.ShopIdentity      `bson:"inline"`
	models.PartitionIdentity `gorm:"embedded;"`
	GuidFixed                string          `json:"guidfixed" gorm:"column:guidfixed"`
	DocNo                    string          `json:"docno" gorm:"column:docno;primaryKey"`
	DocDate                  time.Time       `json:"docdate" gorm:"column:docdate"`
	DebtorCode               string          `json:"creditorcode" gorm:"column:creditorcode"`
	DebtorNames              pkgModels.JSONB `json:"debtornames" gorm:"column:debtornames"`
	InquiryType              int             `json:"inquirytype" gorm:"column:inquirytype"`
	TransFlag                int8            `json:"transflag" gorm:"column:transflag" `
	TotalValue               float64         `json:"totalvalue" gorm:"column:totalvalue"`
	TotalBeforeVat           float64         `json:"totalbeforevat" gorm:"column:totalbeforevat"`
	TotalVatValue            float64         `json:"totalvatvalue" gorm:"column:totalvatvalue"`
	TotalExceptVat           float64         `json:"totalexceptvat" gorm:"column:totalexceptvat"`
	TotalAfterVat            float64         `json:"totalaftervat" gorm:"column:totalaftervat"`
	TotalAmount              float64         `json:"totalamount" gorm:"column:totalamount"`
	PaidAmount               float64         `json:"paidamount" gorm:"column:paidamount"`
	BalanceAmount            float64         `json:"balanceamount" gorm:"column:balanceamount"`
	Status                   int8            `json:"status" gorm:"column:status"`
	IsCancel                 bool            `json:"iscancel" gorm:"column:iscancel"`
}

func (d *DebtorTransactionPG) TableName() string {
	return "debtor_transaction"
}

func (m *DebtorTransactionPG) CompareTo(other *DebtorTransactionPG) bool {

	diff := cmp.Diff(m, other,
		// cmpopts.IgnoreFields(PurchaseTransactionPG{}, "TotalCost"),
		cmpopts.IgnoreFields(DebtorTransactionPG{}, "PaidAmount", "BalanceAmount", "Status"),
	)

	if diff == "" {
		return true
	}

	return false
}
