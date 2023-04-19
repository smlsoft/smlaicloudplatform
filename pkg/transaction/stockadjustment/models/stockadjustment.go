package models

import (
	"smlcloudplatform/pkg/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const stockadjustmentCollectionName = "transactionStockAdjustment"

type StockAdjustment struct {
	models.PartitionIdentity `bson:"inline"`
	Docno                    string         `json:"docno" bson:"docno"`
	TotalDiscount            int            `json:"totaldiscount" bson:"totaldiscount"`
	TotalBeforeVat           int            `json:"totalbeforevat" bson:"totalbeforevat"`
	GuidRef                  string         `json:"guidref" bson:"guidref"`
	DocDatetime              time.Time      `json:"docdatetime" bson:"docdatetime"`
	DocRefNo                 string         `json:"docrefno" bson:"docrefno"`
	DocRefDate               time.Time      `json:"docrefdate" bson:"docrefdate"`
	DocType                  int            `json:"doctype" bson:"doctype"`
	CustName                 []models.NameX `json:"custname" bson:"custname"`
	TotalExceptVat           int            `json:"totalexceptvat" bson:"totalexceptvat"`
	CashierCode              string         `json:"cashiercode" bson:"cashiercode"`
	Details                  []Detail       `json:"details" bson:"details"`
	InquiryType              int            `json:"inquirytype" bson:"inquirytype"`
	DiscountWord             string         `json:"discountword" bson:"discountword"`
	TotalCost                int            `json:"totalcost" bson:"totalcost"`
	TotalVatValue            float64        `json:"totalvatvalue" bson:"totalvatvalue"`
	TotalAmount              float64        `json:"totalamount" bson:"totalamount"`
	TaxDocDate               time.Time      `json:"taxdocdate" bson:"taxdocdate"`
	SaleCode                 string         `json:"salecode" bson:"salecode"`
	PosID                    string         `json:"posid" bson:"posid"`
	SaleName                 string         `json:"salename" bson:"salename"`
	MemberCode               string         `json:"membercode" bson:"membercode"`
	VatRate                  int            `json:"vatrate" bson:"vatrate"`
	TotalValue               int            `json:"totalvalue" bson:"totalvalue"`
	TaxDocNo                 string         `json:"taxdocno" bson:"taxdocno"`
	DocRefType               int            `json:"docreftype" bson:"docreftype"`
	VatType                  int            `json:"vattype" bson:"vattype"`
	CustCode                 string         `json:"custcode" bson:"custcode"`
	TotalAfterVat            float64        `json:"totalaftervat" bson:"totalaftervat"`
	TransFlag                int            `json:"transflag" bson:"transflag"`
	Status                   int            `json:"status" bson:"status"`
}

type Detail struct {
	SumAmount           int            `json:"sumamount" bson:"sumamount"`
	LocationNames       []models.NameX `json:"locationnames" bson:"locationnames"`
	SumAmountExcludeVat int            `json:"sumamountexcludevat" bson:"sumamountexcludevat"`
	DivideValue         int            `json:"dividevalue" bson:"dividevalue"`
	InquiryType         int            `json:"inquirytype" bson:"inquirytype"`
	Price               int            `json:"price" bson:"price"`
	Barcode             string         `json:"barcode" bson:"barcode"`
	UnitCode            string         `json:"unitcode" bson:"unitcode"`
	ToWhCode            string         `json:"towhcode" bson:"towhcode"`
	ToLocationCode      string         `json:"tolocationcode" bson:"tolocationcode"`
	TotalValueVat       float64        `json:"totalvaluevat" bson:"totalvaluevat"`
	ItemGuid            string         `json:"itemguid" bson:"itemguid"`
	ShelfCode           string         `json:"shelfcode" bson:"shelfcode"`
	TotalQty            int            `json:"totalqty" bson:"totalqty"`
	StandValue          int            `json:"standvalue" bson:"standvalue"`
	CalcFlag            int            `json:"calcflag" bson:"calcflag"`
	VatType             int            `json:"vattype" bson:"vattype"`
	ToWhNames           []models.NameX `json:"towhnames" bson:"towhnames"`
	ItemName            []models.NameX `json:"itemname" bson:"itemname"`
	LineNumber          int            `json:"linenumber" bson:"linenumber"`
	WhNames             []models.NameX `json:"whnames" bson:"whnames"`
	AverageCost         int            `json:"averagecost" bson:"averagecost"`
	LastStatus          int            `json:"laststatus" bson:"laststatus"`
	TaxType             int            `json:"taxtype" bson:"taxtype"`
	ItemCode            string         `json:"itemcode" bson:"itemcode"`
	IsPos               int            `json:"ispos" bson:"ispos"`
	MultiUnit           bool           `json:"multiunit" bson:"multiunit"`
	PriceExcludeVat     int            `json:"priceexcludevat" bson:"priceexcludevat"`
	LocationCode        string         `json:"locationcode" bson:"locationcode"`
	ItemType            int            `json:"itemtype" bson:"itemtype"`
	Remark              string         `json:"remark" bson:"remark"`
	Qty                 int            `json:"qty" bson:"qty"`
	Discount            string         `json:"discount" bson:"discount"`
	DocDatetime         time.Time      `json:"docdatetime" bson:"docdatetime"`
	WhCode              string         `json:"whcode" bson:"whcode"`
	ToLocationNames     []models.NameX `json:"tolocationnames" bson:"tolocationnames"`
	DiscountAmount      int            `json:"discountamount" bson:"discountamount"`
	UnitNames           []models.NameX `json:"unitnames" bson:"unitnames"`
	SumOfCost           int            `json:"sumofcost" bson:"sumofcost"`
}

type StockAdjustmentInfo struct {
	models.DocIdentity `bson:"inline"`
	StockAdjustment    `bson:"inline"`
}

func (StockAdjustmentInfo) CollectionName() string {
	return stockadjustmentCollectionName
}

type StockAdjustmentData struct {
	models.ShopIdentity `bson:"inline"`
	StockAdjustmentInfo `bson:"inline"`
}

type StockAdjustmentDoc struct {
	ID                  primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	StockAdjustmentData `bson:"inline"`
	models.ActivityDoc  `bson:"inline"`
}

func (StockAdjustmentDoc) CollectionName() string {
	return stockadjustmentCollectionName
}

type StockAdjustmentItemGuid struct {
	Docno string `json:"docno" bson:"docno"`
}

func (StockAdjustmentItemGuid) CollectionName() string {
	return stockadjustmentCollectionName
}

type StockAdjustmentActivity struct {
	StockAdjustmentData `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (StockAdjustmentActivity) CollectionName() string {
	return stockadjustmentCollectionName
}

type StockAdjustmentDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (StockAdjustmentDeleteActivity) CollectionName() string {
	return stockadjustmentCollectionName
}
