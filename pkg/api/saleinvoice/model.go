package saleinvoice

import (
	"smlcloudplatform/pkg/models"
	"time"
)

type SaleInvoiceTable struct {
	Id                  uint `gorm:"primaryKey"`
	models.ShopIdentity `gorm:"embedded"`
	DocDate             time.Time                `json:"docdate,omitempty" gorm:"column:docdate"`
	DocNo               string                   `json:"docno,omitempty"  gorm:"column:docno"`
	TaxRate             float64                  `json:"taxrate" gorm:"column:taxrate;" `
	TaxAmount           float64                  `json:"taxamount" gorm:"column:taxamount;" `
	TaxBaseAmount       float64                  `json:"taxbaseamount" gorm:"column:taxbaseamount;" `
	DiscountAmount      float64                  `json:"discountamount" gorm:"column:discountamount;" `
	SumAmount           float64                  `json:"sumamount" gorm:"column:sumamount;" `
	Items               []SaleInvoiceDetailTable `json:"items" bson:"items" gorm:"foreignKey:SaleInvoiceId" `
}

func (*SaleInvoiceTable) TableName() string {
	return "saleinvoice"
}

type SaleInvoiceDetailTable struct {
	Id                  uint `gorm:"primaryKey"`
	models.ShopIdentity `gorm:"embedded"`
	SaleInvoiceId       uint    `json:"saleinvoiceid,omitempty"  gorm:"column:saleinvoiceid;foreignKey:Id"`
	DocNo               string  `json:"docno,omitempty"  gorm:"column:docno"`
	LineNumber          int     `json:"linenumber" gorm:"column:linenumber"`
	ItemGuid            string  `json:"itemguid,omitempty" gorm:"column:itemguid"`
	ItemCode            string  `json:"itemcode,omitempty" gorm:"column:itemcode"`
	Barcode             string  `json:"barcode" gorm:"column:barcode"`
	Name1               string  `json:"name1" gorm:"name1"` // ชื่อภาษาไทย
	ItemUnitCode        string  `json:"itemunitcode,omitempty" gorm:"column:itemunitcode"`
	ItemUnitStd         float64 `json:"itemunitstd,omitempty" gorm:"column:itemunitstd"`
	ItemUnitDiv         float64 `json:"itemunitdiv,omitempty" gorm:"column:itemunitdiv"`
	Price               float64 `json:"price" gorm:"column:price" `
	Qty                 float64 `json:"qty" gorm:"column:qty" `
	DiscountAmount      float64 `json:"discountamount" gorm:"column:discountamount"`
	DiscountText        string  `json:"discounttext" gorm:"column:discounttext"`
	Amount              float64 `json:"amount" gorm:"column:amount"`
}

func (SaleInvoiceDetailTable) TableName() string {
	return "saleinvoice_detail"
}
