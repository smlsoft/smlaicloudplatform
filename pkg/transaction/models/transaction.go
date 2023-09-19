package models

import (
	"smlcloudplatform/pkg/models"
	"time"
)

type TransactionHeader struct {
	DocNo            string          `json:"docno" bson:"docno"`
	DocDatetime      time.Time       `json:"docdatetime" bson:"docdatetime"`
	GuidRef          string          `json:"guidref" bson:"guidref"`
	TransFlag        int             `json:"transflag" bson:"transflag"`
	DocRefType       int8            `json:"docreftype" bson:"docreftype"`
	DocRefNo         string          `json:"docrefno" bson:"docrefno"`
	DocRefDate       time.Time       `json:"docrefdate" bson:"docrefdate"`
	TaxDocDate       time.Time       `json:"taxdocdate" bson:"taxdocdate"`
	TaxDocNo         string          `json:"taxdocno" bson:"taxdocno"`
	DocType          int8            `json:"doctype" bson:"doctype"`
	InquiryType      int             `json:"inquirytype" bson:"inquirytype"`
	VatType          int8            `json:"vattype" bson:"vattype"`
	VatRate          float64         `json:"vatrate" bson:"vatrate"`
	CustCode         string          `json:"custcode" bson:"custcode"`
	CustNames        *[]models.NameX `json:"custnames" bson:"custnames"`
	Description      string          `json:"description" bson:"description"`
	DiscountWord     string          `json:"discountword" bson:"discountword"`
	TotalDiscount    float64         `json:"totaldiscount" bson:"totaldiscount"`
	TotalValue       float64         `json:"totalvalue" bson:"totalvalue"`
	TotalExceptVat   float64         `json:"totalexceptvat" bson:"totalexceptvat"`
	TotalAfterVat    float64         `json:"totalaftervat" bson:"totalaftervat"`
	TotalBeforeVat   float64         `json:"totalbeforevat" bson:"totalbeforevat"`
	TotalVatValue    float64         `json:"totalvatvalue" bson:"totalvatvalue"`
	TotalAmount      float64         `json:"totalamount" bson:"totalamount"`
	TotalCost        float64         `json:"totalcost" bson:"totalcost"`
	PosID            string          `json:"posid" bson:"posid"`
	CashierCode      string          `json:"cashiercode" bson:"cashiercode"`
	SaleCode         string          `json:"salecode" bson:"salecode"`
	SaleName         string          `json:"salename" bson:"salename"`
	MemberCode       string          `json:"membercode" bson:"membercode"`
	IsCancel         bool            `json:"iscancel" bson:"iscancel"`
	IsManualAmount   bool            `json:"ismanualamount" bson:"ismanualamount"`
	Status           int8            `json:"status" bson:"status"`
	PaymentDetail    PaymentDetail   `json:"paymentdetail" bson:"paymentdetail"`
	PaymentDetailRaw string          `json:"paymentdetailraw" bson:"paymentdetailraw"`
	PayCashAmount    float64         `json:"paycashamount" bson:"paycashamount"`

	BillTaxType                     int8      `json:"billtaxtype" bson:"billtaxtype"`
	CancelDateTime                  string    `json:"canceldatetime" bson:"canceldatetime"`
	CancelUserCode                  string    `json:"cancelusercode" bson:"cancelusercode"`
	CancelUserName                  string    `json:"cancelusername" bson:"cancelusername"`
	CancelDescription               string    `json:"canceldescription" bson:"canceldescription"`
	CancelReason                    string    `json:"cancelreason" bson:"cancelreason"`
	FullVatAddress                  string    `json:"fullvataddress" bson:"fullvataddress"`
	FullVatBranchNumber             string    `json:"fullvatbranchnumber" bson:"fullvatbranchnumber"`
	FullVatName                     string    `json:"fullvatname" bson:"fullvatname"`
	FullVatDocNumber                string    `json:"fullvatdocnumber" bson:"fullvatdocnumber"`
	FullVatTaxID                    string    `json:"fullvattaxid" bson:"fullvattaxid"`
	FullVatPrint                    bool      `json:"fullvatprint" bson:"fullvatprint"`
	IsVatRegister                   bool      `json:"isvatregister" bson:"isvatregister"`
	PrintCopyBillDateTime           []string  `json:"printcopybilldatetime" bson:"printcopybilldatetime"`
	TableNumber                     string    `json:"tablenumber" bson:"tablenumber"`
	TableOpenDateTime               time.Time `json:"tableopendatetime" bson:"tableopendatetime"`
	TableCloseDateTime              time.Time `json:"tableclosedatetime" bson:"tableclosedatetime"`
	ManCount                        int       `json:"mancount" bson:"mancount"`
	WomanCount                      int       `json:"womancount" bson:"womancount"`
	ChildCount                      int       `json:"childcount" bson:"childcount"`
	IsTableAllacrateMode            bool      `json:"istableallacratemode" bson:"istableallacratemode"`
	BuffetCode                      bool      `json:"buffetcode" bson:"buffetcode"`
	CustomerTelephone               string    `json:"customertelephone" bson:"customertelephone"`
	TotalQty                        float64   `json:"totalqty" bson:"totalqty"`
	TotalDiscountVatAmount          float64   `json:"totaldiscountvatamount" bson:"totaldiscountvatamount"`
	TotalDiscountExceptVatAmount    float64   `json:"totaldiscountexceptvatamount" bson:"totaldiscountexceptvatamount"`
	CashierName                     string    `json:"cashiername" bson:"cashiername"`
	PayCashChange                   float64   `json:"paycashchange" bson:"paycashchange"`
	SumQRCode                       float64   `json:"sumqrcode" bson:"sumqrcode"`
	SumCreditCard                   float64   `json:"sumcreditcard" bson:"sumcreditcard"`
	SumMoneyTransfer                float64   `json:"summoneytransfer" bson:"summoneytransfer"`
	SumCheque                       float64   `json:"sumcheque" bson:"sumcheque"`
	SumCoupon                       float64   `json:"sumcoupon" bson:"sumcoupon"`
	DetailDiscountFormula           string    `json:"detaildiscountformula" bson:"detaildiscountformula"`
	DetailTotalAmount               float64   `json:"detailtotalamount" bson:"detailtotalamount"`
	DetailTotalDiscount             float64   `json:"detailtotaldiscount" bson:"detailtotaldiscount"`
	RoundAmount                     float64   `json:"roundamount" bson:"roundamount"`
	TotalAmountAfterDiscount        float64   `json:"totalamountafterdiscount" bson:"totalamountafterdiscount"`
	DetailTotalAmountBeforeDiscount float64   `json:"detailtotalamountbeforediscount" bson:"detailtotalamountbeforediscount"`
}

type Transaction struct {
	TransactionHeader `bson:"inline"`
	Details           *[]Detail `json:"details" bson:"details"`
}

type Detail struct {
	InquiryType         int8            `json:"inquirytype" bson:"inquirytype"`
	LineNumber          int             `json:"linenumber" bson:"linenumber"`
	DocDatetime         time.Time       `json:"docdatetime" bson:"docdatetime"`
	DocRef              string          `json:"docref" bson:"docref"`
	DocRefDatetime      time.Time       `json:"docrefdatetime" bson:"docrefdatetime"`
	CalcFlag            int8            `json:"calcflag" bson:"calcflag"`
	Barcode             string          `json:"barcode" bson:"barcode"`
	ItemCode            string          `json:"itemcode" bson:"itemcode"`
	UnitCode            string          `json:"unitcode" bson:"unitcode"`
	ItemType            int8            `json:"itemtype" bson:"itemtype"`
	ItemGuid            string          `json:"itemguid" bson:"itemguid"`
	Qty                 float64         `json:"qty" bson:"qty"`
	TotalQty            float64         `json:"totalqty" bson:"totalqty"`
	Price               float64         `json:"price" bson:"price"`
	Discount            string          `json:"discount" bson:"discount"`
	DiscountAmount      float64         `json:"discountamount" bson:"discountamount"`
	TotalValueVat       float64         `json:"totalvaluevat" bson:"totalvaluevat"`
	PriceExcludeVat     float64         `json:"priceexcludevat" bson:"priceexcludevat"`
	SumAmount           float64         `json:"sumamount" bson:"sumamount"`
	SumAmountExcludeVat float64         `json:"sumamountexcludevat" bson:"sumamountexcludevat"`
	DivideValue         float64         `json:"dividevalue" bson:"dividevalue"`
	StandValue          float64         `json:"standvalue" bson:"standvalue"`
	VatType             int8            `json:"vattype" bson:"vattype"`
	Remark              string          `json:"remark" bson:"remark"`
	MultiUnit           bool            `json:"multiunit" bson:"multiunit"`
	SumOfCost           float64         `json:"sumofcost" bson:"sumofcost"`
	AverageCost         float64         `json:"averagecost" bson:"averagecost"`
	LastStatus          int8            `json:"laststatus" bson:"laststatus"`
	IsPos               int8            `json:"ispos" bson:"ispos"`
	TaxType             int8            `json:"taxtype" bson:"taxtype"`
	VatCal              int             `json:"vatcal" bson:"vatcal"`
	WhCode              string          `json:"whcode" bson:"whcode"`
	ShelfCode           string          `json:"shelfcode" bson:"shelfcode"`
	LocationCode        string          `json:"locationcode" bson:"locationcode"`
	ToWhCode            string          `json:"towhcode" bson:"towhcode"`
	ToLocationCode      string          `json:"tolocationcode" bson:"tolocationcode"`
	ItemNames           *[]models.NameX `json:"itemnames" bson:"itemnames"`
	UnitNames           *[]models.NameX `json:"unitnames" bson:"unitnames" `
	WhNames             *[]models.NameX `json:"whnames" bson:"whnames"`
	LocationNames       *[]models.NameX `json:"locationnames" bson:"locationnames"`
	ToWhNames           *[]models.NameX `json:"towhnames" bson:"towhnames"`
	ToLocationNames     *[]models.NameX `json:"tolocationnames" bson:"tolocationnames" `
	SKU                 string          `json:"sku" bson:"sku"`
	ExtraJson           string          `json:"extrajson" bson:"extrajson"`
}

type PaymentDetail struct {
	CashAmountText     string               `json:"cashamounttext" bson:"cashamounttext"`
	CashAmount         float64              `json:"cashamount" bson:"cashamount"`
	PaymentCreditCards *[]PaymentCreditCard `json:"paymentcreditcards" bson:"paymentcreditcards"`
	PaymentTransfers   *[]PaymentTransfer   `json:"paymenttransfers" bson:"paymenttransfers"`
}

type PaymentCreditCard struct {
	DocDatetime   time.Time `json:"docdatetime" bson:"docdatetime"`
	CardNumber    string    `json:"cardnumber" bson:"cardnumber"`
	Amount        float64   `json:"amount" bson:"amount"`
	ChargeWord    string    `json:"chargeword" bson:"chargeword"`
	ChargeValue   float64   `json:"chargevalue" bson:"chargevalue"`
	TotalNetWorth float64   `json:"totalnetworth" bson:"totalnetworth"`
}

type PaymentTransfer struct {
	DocDatetime   time.Time       `json:"docdatetime" bson:"docdatetime"`
	BankCode      string          `json:"bankcode" bson:"bankcode"`
	BankNames     *[]models.NameX `json:"banknames" bson:"banknames"`
	AccountNumber string          `json:"accountnumber" bson:"accountnumber"`
	Amount        float64         `json:"amount" bson:"amount"`
}
