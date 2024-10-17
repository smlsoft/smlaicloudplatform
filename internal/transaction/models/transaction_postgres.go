package models

import (
	"smlcloudplatform/internal/models"
	"time"
)

type TransactionPG struct {
	GuidFixed                string `json:"guidfixed" gorm:"column:guidfixed"`
	models.ShopIdentity      `bson:"inline"`
	models.PartitionIdentity `gorm:"embedded;"`
	InquiryType              int       `json:"inquirytype" gorm:"column:inquirytype"`
	TransFlag                int8      `json:"transflag" gorm:"column:transflag" `
	DocNo                    string    `json:"docno" gorm:"column:docno;primaryKey"`
	DocDate                  time.Time `json:"docdate" gorm:"column:docdate"`
	DocRefType               int8      `json:"docreftype" gorm:"column:docreftype"`
	DocRefNo                 string    `json:"docrefno" gorm:"column:docrefno"`
	DeviceName               string    `json:"devicename" gorm:"column:devicename"`
	GuidPos                  string    `json:"guidpos" gorm:"column:guidpos"`

	DocRefDate     time.Time    `json:"docrefdate" gorm:"column:docrefdate"`
	BranchCode     string       `json:"branchcode" gorm:"column:branchcode"`
	BranchNames    models.JSONB `json:"branchnames" gorm:"column:branchnames;type:jsonb"`
	Description    string       `json:"description" gorm:"column:description"`
	TaxDocNo       string       `json:"taxdocno"  gorm:"column:taxdocno"`
	TaxDocDate     time.Time    `json:"taxdocdate" gorm:"column:taxdocdate"`
	IsCancel       bool         `json:"iscancel" gorm:"column:iscancel"`
	IsBom          bool         `json:"isbom" gorm:"column:isbom"`
	Status         int8         `json:"status" gorm:"column:status"`
	VatType        int8         `json:"vattype" gorm:"column:vattype" `
	VatRate        float64      `json:"vatrate" gorm:"column:vatrate"`
	TotalValue     float64      `json:"totalvalue" gorm:"column:totalvalue"`
	DiscountWord   string       `json:"discountword" gorm:"column:discountword"`
	DeliveryAmount float64      `json:"deliveryamount" gorm:"column:deliveryamount"`
	TotalDiscount  float64      `json:"totaldiscount" gorm:"column:totaldiscount"`
	TotalBeforeVat float64      `json:"totalbeforevat" gorm:"column:totalbeforevat"`
	TotalVatValue  float64      `json:"totalvatvalue" gorm:"column:totalvatvalue"`
	TotalExceptVat float64      `json:"totalexceptvat" gorm:"column:totalexceptvat"`
	TotalAfterVat  float64      `json:"totalaftervat" gorm:"column:totalaftervat"`
	TotalAmount    float64      `json:"totalamount" gorm:"column:totalamount"`
	GuidRef        string       `json:"guidref" gorm:"column:guidref"`
	IsManualAmount bool         `json:"ismanualamount" gorm:"column:ismanualamount"`
	AlcoholAmount  float64      `json:"alcoholamount" gorm:"column:alcoholamount"`
	OtherAmount    float64      `json:"otheramount" gorm:"column:otheramount"`
	DrinkAmount    float64      `json:"drinkamount" gorm:"column:drinkamount"`
	FoodAmount     float64      `json:"foodamount" gorm:"column:foodamount"`
	// TotalCost                float64                   `json:"totalcost" gorm:"column:totalcost"`
	// PosID                    string                    `json:"posid" gorm:"column:posid"`

	// Details                  *[]StockTransactionDetail `json:"details" gorm:"details;foreignKey:shopid,docno"`
}

type TransactionDetailPG struct {
	ID                       uint   `gorm:"primarykey"`
	ShopID                   string `json:"shopid" gorm:"column:shopid"`
	GuidFixed                string `json:"guidfixed" gorm:"column:guidfixed"`
	models.PartitionIdentity `gorm:"embedded;"`
	DocNo                    string       `json:"docno" gorm:"column:docno"`
	LineNumber               int8         `json:"linenumber" gorm:"column:linenumber"`
	Barcode                  string       `json:"barcode" gorm:"column:barcode"`
	ItemNames                models.JSONB `json:"itemnames" gorm:"column:itemnames;type:jsonb"`
	UnitCode                 string       `json:"unitcode" gorm:"column:unitcode"`
	Qty                      float64      `json:"qty" gorm:"column:qty"`
	Price                    float64      `json:"price" gorm:"column:price"`
	PriceExcludeVat          float64      `json:"priceexcludevat" gorm:"column:priceexcludevat"`
	Discount                 string       `json:"discount" gorm:"column:discount"`
	DiscountAmount           float64      `json:"discountamount" gorm:"column:discountamount"`
	SumAmount                float64      `json:"sumamount" gorm:"column:sumamount"`
	SumAmountExcludeVat      float64      `json:"sumamountexcludevat" gorm:"column:sumamountexcludevat"`
	SumAmountChoice          float64      `json:"sumamountchoice" gorm:"column:sumamountchoice"`
	RefGuid                  string       `json:"refguid" gorm:"column:refguid"`
	WhCode                   string       `json:"whcode" gorm:"column:whcode"`
	WhNames                  models.JSONB `json:"whnames" gorm:"column:whnames;type:jsonb"`
	LocationCode             string       `json:"locationcode" gorm:"column:locationcode"`
	LocationNames            models.JSONB `json:"locationnames" gorm:"column:locationnames;type:jsonb"`
	VatCal                   int8         `json:"vatcal" gorm:"column:vatcal"`
	FoodType                 int8         `json:"foodtype" gorm:"column:foodtype"`
	VatType                  int8         `json:"vattype" gorm:"column:vattype"`
	TaxType                  int8         `json:"taxtype" gorm:"column:taxtype"`
	IsChoice                 int8         `json:"ischoice" gorm:"column:ischoice"`
	StandValue               float64      `json:"standvalue" gorm:"column:standvalue"`
	DivideValue              float64      `json:"dividevalue" gorm:"column:dividevalue"`
	ItemType                 int8         `json:"itemtype" gorm:"column:itemtype"`
	ItemGuid                 string       `json:"itemguid" gorm:"column:itemguid"`
	TotalValueVat            float64      `json:"totalvaluevat" gorm:"column:totalvaluevat"`
	DocRef                   string       `json:"docref" gorm:"column:docref"`
	DocRefDateTime           time.Time    `json:"docrefdatetime" gorm:"column:docrefdatetime"`
	Remark                   string       `json:"remark" gorm:"column:remark"`
	WhCodeDestination        string       `json:"whcodedestination" gorm:"column:whcodedestination"`
	WhDestinationNames       models.JSONB `json:"whcodedestinationnames" gorm:"column:whcodedestinationnames;type:jsonb"`
	LocationCodeDestination  string       `json:"locationcodedestination" gorm:"column:locationcodedestination"`
	LocationDestination      models.JSONB `json:"locationdestination" gorm:"column:locationdestination;type:jsonb"`
	UnitNames                models.JSONB `json:"unitnames" gorm:"column:unitnames;type:jsonb"`
	GroupCode                string       `json:"groupcode" gorm:"column:groupcode"`
	GroupNames               models.JSONB `json:"groupnames" gorm:"column:groupnames;type:jsonb"`
	DocDate                  time.Time    `json:"docdate" gorm:"column:docdate"`
	// Barcode                  string  `json:"barcode" gorm:"column:barcode"`
	// RefBarcode               string  `json:"refbarcode" gorm:"column:refbarcode"`
	// UnitCode                 string  `json:"unitcode" gorm:"column:unitcode"`
	// Qty                      float64 `json:"qty" gorm:"column:qty"`
	// StandValue               float64 `json:"standvalue" gorm:"column:standvalue"`
	// DivideValue              float64 `json:"dividevalue" gorm:"column:dividevalue"`
	// WhCode                   string  `json:"whcode" gorm:"whcode"`
	// LocationCode             string  `json:"locationcode" gorm:"locationcode"`
	// ItemType                 int8    `json:"itemtype" gorm:"column:itemtype"`
	// Price                    float64 `json:"price" gorm:"column:price"`
	// Discount                 string  `json:"discount" gorm:"column:discount"`
	// DiscountAmount           float64 `json:"discountamount" gorm:"column:discountamount"`
	// SumAmount                float64 `json:"sumamount" gorm:"column:sumamount"`
	// CalcFlag                 int8    `json:"calcflag" gorm:"column:calcflag"`
	// SumOfCost                float64 `json:"sumofcost" gorm:"column:sumofcost"`
	// AverageCost              float64 `json:"averagecost" gorm:"column:averagecost"`
	// CostPerUnit              float64 `json:"costperunit" gorm:"column:costperunit"`
	// TotalCost                float64 `json:"totalcost" gorm:"column:totalcost"`
	// SumAmountExcludeVat      float64 `json:"sumamountexcludevat" gorm:"column:sumamountexcludevat"`
	// TotalValueVat            float64 `json:"totalvaluevat" gorm:"column:totalvaluevat"`
	// ItemGuid                 string  `json:"itemguid" gorm:"column:itemguid"`
	// VatType                  int8    `json:"vattype" gorm:"column:vattype"`
	// TaxType                  int8    `json:"taxtype" gorm:"column:taxtype"`
	// DocRef                   string  `json:"docref" gorm:"column:docref"`
}

// func (j *StockTransaction) BeforeUpdate(tx *gorm.DB) (err error) {

// 	// find old data
// 	var details *[]StockTransactionDetail
// 	tx.Model(&StockTransactionDetail{}).Where(" shopid=? AND docno=?", j.ShopID, j.DocNo).Find(&details)

// 	// delete un use data
// 	for _, tmp := range *details {
// 		var foundUpdate bool = false
// 		for _, data := range *j.Details {
// 			if data.ID == tmp.ID {
// 				foundUpdate = true
// 			}
// 		}
// 		if foundUpdate == false {
// 			// mark delete
// 			tx.Delete(&StockTransactionDetail{}, tmp.ID)
// 		}
// 	}

// 	return nil
// }

// func (jd *StockTransactionDetail) BeforeCreate(tx *gorm.DB) error {

// 	tx.Statement.AddClause(clause.OnConflict{
// 		UpdateAll: true,
// 	})
// 	return nil
// }

// func (s *StockTransaction) CompareTo(other *StockTransaction) bool {

// 	diff := cmp.Diff(s, other,
// 		cmpopts.IgnoreFields(StockTransaction{}, "TotalCost"),
// 		cmpopts.IgnoreFields(StockTransactionDetail{}, "ID", "SumOfCost", "AverageCost"),
// 	)

// 	if diff == "" {
// 		return true
// 	}

// 	return false
// }
