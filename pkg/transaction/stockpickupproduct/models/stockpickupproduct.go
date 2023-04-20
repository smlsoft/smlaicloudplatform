package models

import (
	"smlcloudplatform/pkg/models"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const stockpickupproductCollectionName = "transactionStockPickupProduct"

type StockPickupProduct struct {
	models.PartitionIdentity `bson:"inline"`
	Docno                    string         `json:"docno" bson:"docno"`
	TotalDiscount            float64        `json:"totaldiscount" bson:"totaldiscount"`
	TotalBeforeVat           float64        `json:"totalbeforevat" bson:"totalbeforevat"`
	GuidRef                  string         `json:"guidref" bson:"guidref"`
	DocDatetime              time.Time      `json:"docdatetime" bson:"docdatetime"`
	DocRefNo                 string         `json:"docrefno" bson:"docrefno"`
	DocRefDate               time.Time      `json:"docrefdate" bson:"docrefdate"`
	DocType                  int8           `json:"doctype" bson:"doctype"`
	CustName                 []models.NameX `json:"custname" bson:"custname"`
	TotalExceptVat           float64        `json:"totalexceptvat" bson:"totalexceptvat"`
	CashierCode              string         `json:"cashiercode" bson:"cashiercode"`
	Details                  []Detail       `json:"details" bson:"details"`
	InquiryType              int            `json:"inquirytype" bson:"inquirytype"`
	DiscountWord             string         `json:"discountword" bson:"discountword"`
	TotalCost                float64        `json:"totalcost" bson:"totalcost"`
	TotalVatValue            float64        `json:"totalvatvalue" bson:"totalvatvalue"`
	TotalAmount              float64        `json:"totalamount" bson:"totalamount"`
	TaxDocDate               time.Time      `json:"taxdocdate" bson:"taxdocdate"`
	SaleCode                 string         `json:"salecode" bson:"salecode"`
	PosID                    string         `json:"posid" bson:"posid"`
	SaleName                 string         `json:"salename" bson:"salename"`
	MemberCode               string         `json:"membercode" bson:"membercode"`
	VatRate                  float64        `json:"vatrate" bson:"vatrate"`
	TotalValue               float64        `json:"totalvalue" bson:"totalvalue"`
	TaxDocNo                 string         `json:"taxdocno" bson:"taxdocno"`
	DocRefType               int8           `json:"docreftype" bson:"docreftype"`
	VatType                  int8           `json:"vattype" bson:"vattype"`
	CustCode                 string         `json:"custcode" bson:"custcode"`
	TotalAfterVat            float64        `json:"totalaftervat" bson:"totalaftervat"`
	TransFlag                int            `json:"transflag" bson:"transflag"`
	Status                   int8           `json:"status" bson:"status"`
}

type Detail struct {
	SumAmount           float64        `json:"sumamount" bson:"sumamount"`
	LocationNames       []models.NameX `json:"locationnames" bson:"locationnames"`
	SumAmountExcludeVat float64        `json:"sumamountexcludevat" bson:"sumamountexcludevat"`
	DivideValue         int            `json:"dividevalue" bson:"dividevalue"`
	StandValue          int            `json:"standvalue" bson:"standvalue"`
	InquiryType         int8           `json:"inquirytype" bson:"inquirytype"`
	Price               float64        `json:"price" bson:"price"`
	Barcode             string         `json:"barcode" bson:"barcode"`
	UnitCode            string         `json:"unitcode" bson:"unitcode"`
	ToWhCode            string         `json:"towhcode" bson:"towhcode"`
	ToLocationCode      string         `json:"tolocationcode" bson:"tolocationcode"`
	TotalValueVat       float64        `json:"totalvaluevat" bson:"totalvaluevat"`
	ItemGuid            string         `json:"itemguid" bson:"itemguid"`
	ShelfCode           string         `json:"shelfcode" bson:"shelfcode"`
	TotalQty            float64        `json:"totalqty" bson:"totalqty"`
	CalcFlag            int8           `json:"calcflag" bson:"calcflag"`
	VatType             int8           `json:"vattype" bson:"vattype"`
	ToWhNames           []models.NameX `json:"towhnames" bson:"towhnames"`
	ItemName            []models.NameX `json:"itemname" bson:"itemname"`
	LineNumber          int            `json:"linenumber" bson:"linenumber"`
	WhNames             []models.NameX `json:"whnames" bson:"whnames"`
	AverageCost         float64        `json:"averagecost" bson:"averagecost"`
	LastStatus          int8           `json:"laststatus" bson:"laststatus"`
	TaxType             int8           `json:"taxtype" bson:"taxtype"`
	ItemCode            string         `json:"itemcode" bson:"itemcode"`
	IsPos               int8           `json:"ispos" bson:"ispos"`
	MultiUnit           bool           `json:"multiunit" bson:"multiunit"`
	PriceExcludeVat     float64        `json:"priceexcludevat" bson:"priceexcludevat"`
	LocationCode        string         `json:"locationcode" bson:"locationcode"`
	ItemType            int8           `json:"itemtype" bson:"itemtype"`
	Remark              string         `json:"remark" bson:"remark"`
	Qty                 float64        `json:"qty" bson:"qty"`
	Discount            string         `json:"discount" bson:"discount"`
	DocDatetime         time.Time      `json:"docdatetime" bson:"docdatetime"`
	WhCode              string         `json:"whcode" bson:"whcode"`
	ToLocationNames     []models.NameX `json:"tolocationnames" bson:"tolocationnames"`
	DiscountAmount      float64        `json:"discountamount" bson:"discountamount"`
	UnitNames           []models.NameX `json:"unitnames" bson:"unitnames"`
	SumOfCost           float64        `json:"sumofcost" bson:"sumofcost"`
}

type StockPickupProductInfo struct {
	models.DocIdentity `bson:"inline"`
	StockPickupProduct `bson:"inline"`
}

func (StockPickupProductInfo) CollectionName() string {
	return stockpickupproductCollectionName
}

type StockPickupProductData struct {
	models.ShopIdentity    `bson:"inline"`
	StockPickupProductInfo `bson:"inline"`
}

type StockPickupProductDoc struct {
	ID                     primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	StockPickupProductData `bson:"inline"`
	models.ActivityDoc     `bson:"inline"`
}

func (StockPickupProductDoc) CollectionName() string {
	return stockpickupproductCollectionName
}

type StockPickupProductItemGuid struct {
	Docno string `json:"docno" bson:"docno"`
}

func (StockPickupProductItemGuid) CollectionName() string {
	return stockpickupproductCollectionName
}

type StockPickupProductActivity struct {
	StockPickupProductData `bson:"inline"`
	models.ActivityTime    `bson:"inline"`
}

func (StockPickupProductActivity) CollectionName() string {
	return stockpickupproductCollectionName
}

type StockPickupProductDeleteActivity struct {
	models.Identity     `bson:"inline"`
	models.ActivityTime `bson:"inline"`
}

func (StockPickupProductDeleteActivity) CollectionName() string {
	return stockpickupproductCollectionName
}
