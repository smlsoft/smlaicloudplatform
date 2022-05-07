package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const saleinvoiceCollectionName = "saleinvoices"
const saleinvoiceIndexName = "saleinvoices"

type Saleinvoice struct {
	DocDate        *time.Time           `json:"docdate,omitempty" bson:"docdate,omitempty" format:"date-time" examples:"2019-10-12T07:20:50.52Z or 2019-10-12T07:20:50.52+07:00" gorm:"docdate,type:date"`
	DocNo          string               `json:"docno,omitempty"  bson:"docno,omitempty"`
	Member         *Member              `json:"member,omitempty"  bson:"member,omitempty"`
	Items          *[]SaleinvoiceDetail `json:"items" bson:"items" gorm:""`
	TotalAmount    float64              `json:"totalamount" bson:"totalamount" `
	TaxRate        float64              `json:"taxrate" bson:"taxrate" `
	TaxAmount      float64              `json:"taxamount" bson:"taxamount" `
	TaxBaseAmount  float64              `json:"taxbaseamount" bson:"taxbaseamount" `
	DiscountAmount float64              `json:"discountamount" bson:"discountamount" `
	SumAmount      float64              `json:"sumamount" bson:"sumamount" `
	Payment        Payment              `json:"payment" bson:"payment"`
}

type SaleinvoiceDetail struct {
	LineNumber     int `json:"linenumber" bson:"linenumber"`
	InventoryInfo  `bson:"inline" gorm:"embedded;"`
	Price          float64 `json:"price" bson:"price" `
	Qty            float64 `json:"qty" bson:"qty" `
	DiscountAmount float64 `json:"discountamount" bson:"discountamount"`
	DiscountText   string  `json:"discounttext" bson:"discounttext"`
	Amount         float64 `json:"amount" bson:"amount"`
}

type SaleinvoiceInfo struct {
	DocIdentity `bson:"inline"`
	Saleinvoice `bson:"inline"`
}

func (SaleinvoiceInfo) CollectionName() string {
	return saleinvoiceCollectionName
}

type SaleinvoiceData struct {
	ShopIdentity    `bson:"inline"`
	SaleinvoiceInfo `bson:"inline"`
}

func (SaleinvoiceData) IndexName() string {
	return saleinvoiceIndexName
}

type SaleinvoiceDoc struct {
	ID              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SaleinvoiceData `bson:"inline"`
	ActivityDoc     `bson:"inline"`
}

func (SaleinvoiceDoc) CollectionName() string {
	return saleinvoiceCollectionName
}

type SaleInvoiceListPageResponse struct {
	Success    bool                   `json:"success"`
	Data       []SaleinvoiceInfo      `json:"data,omitempty"`
	Pagination PaginationDataResponse `json:"pagination,omitempty"`
}
