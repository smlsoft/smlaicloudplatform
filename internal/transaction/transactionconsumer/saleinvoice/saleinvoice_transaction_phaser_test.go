package saleinvoice_test

import (
	pkgModels "smlcloudplatform/internal/models"
	"smlcloudplatform/internal/transaction/models"
	"smlcloudplatform/internal/transaction/transactionconsumer/saleinvoice"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func SaleInvoiceTransactionStruct() models.SaleInvoiceTransactionPG {

	codeTh := "th"
	nameTh := "ลูกค้าทั่วไป"

	give := models.SaleInvoiceTransactionPG{
		TransactionPG: models.TransactionPG{
			ShopIdentity: pkgModels.ShopIdentity{
				ShopID: "2Eh6e3pfWvXTp0yV3CyFEhKPjdI",
			},
			GuidFixed:      "2TKOzSqEElEKNuIacaMHxbc4GgU",
			TransFlag:      44,
			DocNo:          "a91d29f5-67af-4334-8999-8bc49ed73b4a",
			DocDate:        time.Date(2023, 7, 31, 7, 29, 28, 0, time.UTC),
			GuidRef:        "zzzzz",
			DocRefType:     4,
			DocRefNo:       "REFNO",
			DocRefDate:     time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
			TaxDocNo:       "5cca85fe-6804-4813-bf0d-fdccb3551175",
			TaxDocDate:     time.Date(2023, 7, 31, 7, 29, 28, 0, time.UTC),
			Description:    "POS",
			InquiryType:    1,
			VatRate:        7,
			VatType:        1,
			DiscountWord:   "100",
			TotalDiscount:  100,
			TotalValue:     2000,
			TotalBeforeVat: 2,
			TotalExceptVat: 1000,
			TotalVatValue:  51.02678028444716,
			TotalAfterVat:  2,
			TotalAmount:    2000,
		},
		Items: &[]models.SaleInvoiceTransactionDetailPG{
			{
				TransactionDetailPG: models.TransactionDetailPG{
					DocRef:              "--",
					DocRefDateTime:      time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
					DocNo:               "a91d29f5-67af-4334-8999-8bc49ed73b4a",
					ShopID:              "2Eh6e3pfWvXTp0yV3CyFEhKPjdI",
					LineNumber:          1,
					ItemGuid:            "-",
					Barcode:             "8850086130359",
					UnitCode:            "ซอง",
					Qty:                 5,
					Price:               6,
					PriceExcludeVat:     99,
					Discount:            "2",
					DiscountAmount:      2,
					SumAmount:           1250,
					SumAmountExcludeVat: 1245,
					TotalValueVat:       75,
					WhCode:              "POSWH000",
					LocationCode:        "POSLC000",
					VatType:             1,
					TaxType:             0,
					StandValue:          1,
					DivideValue:         1,
					ItemType:            0,
					Remark:              "-",
					VatCal:              0,
				},
			},
		},
		TotalPayCash:     500,
		TotalPayCredit:   0,
		TotalPayTransfer: 0,
		DebtorCode:       "POS001",
		DebtorNames: []pkgModels.NameX{
			{
				Code:     &codeTh,
				Name:     &nameTh,
				IsAuto:   false,
				IsDelete: false,
			},
		},
	}

	return give
}

func TestSaleInvoiceTransactionPhaser(t *testing.T) {

	giveInput := `{
		"id": "000000000000000000000000",
		"shopid": "2Eh6e3pfWvXTp0yV3CyFEhKPjdI",
		"guidfixed": "2TKOzSqEElEKNuIacaMHxbc4GgU",
		"docno": "a91d29f5-67af-4334-8999-8bc49ed73b4a",
		"docdatetime": "2023-07-31T07:29:28.000Z",
		"transflag": 12,
		"guidref": "zzzzz",
		"docreftype": 4,
		"docrefno": "REFNO",
		"docrefdate": "0001-01-01T00:00:00Z",
		"doctype": 1,
		"taxdocno": "5cca85fe-6804-4813-bf0d-fdccb3551175",
		"taxdocdate": "2023-07-31T07:29:28.000Z",
		"inquirytype": 1,
		"vatrate": 7,
		"vattype": 1,
		"discountword": "100",
		"totaldiscount": 100,
		"totalvalue": 2000,
		"totalbeforevat": 2,
		"totalaftervat": 2,
		"totalvatvalue": 51.02678028444716,
		"totalexceptvat": 1000,
		"totalamount": 2000,
		"totalcost": 0,
		"salecode": "",
		"posid": "",
		"salename": "",
		"membercode": "",
		"description": "POS",
		"cashiercode": "ADMIN09",
		"details": [
			{
				"linenumber": 1,
				"docdatetime": "2023-07-31T07:29:28.566Z",
				"itemguid": "-",
				"docref": "--",
				"docrefdatetime": "0001-01-01T00:00:00Z",
				"barcode": "8850086130359",
				"unitcode": "ซอง",
				"whcode": "POSWH000",
				"locationcode": "POSLC000",
				"price": 6,
				"itemtype": 0,
				"remark": "-",
				"itemcode": "",
				"priceexcludevat": 99,
				"qty": 5,
				"discount": "2",
				"discountamount": 2,
				"totalvaluevat": 75,
				"sumamount": 1250,
				"sumamountexcludevat": 1245,
				"dividevalue": 1,
				"standvalue": 1,
				"vattype": 0,
				"inquirytype": 0,
				"towhcode": "",
				"tolocationcode": "",
				"shelfcode": "",
				"totalqty": 5,
				"calcflag": 0,
				"towhnames": [],
				"locationnames": [],
				"itemnames": [
					{
						"code": "th",
						"name": "[{\"code\":\"th\",\"name\":\"โอวัลติน ซอง\"},{\"code\":\"en\",\"name\":\"\"}]",
						"isauto": false,
						"isdelete": false
					}
				],
				"whnames": [],
				"averagecost": 0,
				"taxtype": 0,
				"laststatus": 0,
				"ispos": 1,
				"multiunit": false,
				

				
				"tolocationnames": [],
				"unitnames": [
					{
						"code": "th",
						"name": "[{\"code\":\"th\",\"name\":\"ซอง\"}]",
						"isauto": false,
						"isdelete": false
					}
				],
				"sumofcost": 0,
				"vatcal": 0
			}
		],
		"custcode": "POS001",
		"custnames": [
			{
				"code": "th",
				"name": "ลูกค้าทั่วไป",
				"isauto": false,
				"isdelete": false
			}
		],
		"status": 0,
		"iscancel": false,
		"ismanualamount": false,
		"ispos": true,
		"paymentdetail": {
			"cashamounttext": "",
			"cashamount": 500,
			"paymentcreditcards": [],
			"paymenttransfers": []
		},
		"couponno": "",
		"couponamount": 0,
		"coupondescription": "",
		"qrcode": "",
		"qrcodeamount": 0,
		"chequeno": "",
		"chequebooknumber": "",
		"chequebookcode": "",
		"chequeduedate": "",
		"chequeamount": 0,
		"paymentdetailraw": "{\"cashamount\":500.0,\"cashamounttext\":\"\",\"paymentcreditcards\":[],\"paymenttransfers\":[]}"
	}
	`

	phaser := saleinvoice.SalesInvoiceTransactionPhaser{}
	want := SaleInvoiceTransactionStruct()

	get, err := phaser.PhaseSingleDoc(giveInput)

	assert.Nil(t, err)

	assert.Equal(t, get.ShopID, want.ShopID, "shopid")
	assert.Equal(t, get.GuidFixed, want.GuidFixed, "guidfixed")
	assert.Equal(t, get.TransFlag, want.TransFlag, "transflag")
	assert.Equal(t, get.DocNo, want.DocNo, "docno")
	assert.Equal(t, get.DocDate, want.DocDate, "docdate")
	assert.Equal(t, get.GuidRef, want.GuidRef, "guidref")
	assert.Equal(t, get.DocRefType, want.DocRefType, "docreftype")
	assert.Equal(t, get.DocRefNo, want.DocRefNo, "docrefno")
	assert.Equal(t, get.DocRefDate, want.DocRefDate, "docrefdate")
	assert.Equal(t, get.TaxDocNo, want.TaxDocNo, "taxdocno")
	assert.Equal(t, get.TaxDocDate, want.TaxDocDate, "taxdocdate")
	assert.Equal(t, get.Description, want.Description, "description")
	assert.Equal(t, get.InquiryType, want.InquiryType, "inquirytype")
	assert.Equal(t, get.VatRate, want.VatRate, "vatrate")
	assert.Equal(t, get.VatType, want.VatType, "vattype")
	assert.Equal(t, get.DiscountWord, want.DiscountWord, "discountword")
	assert.Equal(t, get.TotalDiscount, want.TotalDiscount, "totaldiscount")
	assert.Equal(t, get.TotalValue, want.TotalValue, "totalvalue")
	assert.Equal(t, get.TotalBeforeVat, want.TotalBeforeVat, "totalbeforevat")
	assert.Equal(t, get.TotalExceptVat, want.TotalExceptVat, "totalexceptvat")
	assert.Equal(t, get.TotalVatValue, want.TotalVatValue, "totalvatvalue")
	assert.Equal(t, get.TotalAfterVat, want.TotalAfterVat, "totalaftervat")
	assert.Equal(t, get.TotalAmount, want.TotalAmount, "totalamount")

	assert.Equal(t, get.DebtorCode, want.DebtorCode, "creditorcode")
	assert.Equal(t, *get.DebtorNames[0].Name, "ลูกค้าทั่วไป", "creditorname")

	// detail
	assert.Equal(t, (*get.Items)[0].DocNo, (*want.Items)[0].DocNo, "item.docno")
	assert.Equal(t, (*get.Items)[0].ShopID, (*want.Items)[0].ShopID, "item.shopid")
	assert.Equal(t, (*get.Items)[0].LineNumber, (*want.Items)[0].LineNumber, "item.linenumber")
	assert.Equal(t, (*get.Items)[0].ItemGuid, (*want.Items)[0].ItemGuid, "item.itemguid")
	assert.Equal(t, (*get.Items)[0].Barcode, (*want.Items)[0].Barcode, "item.barcode")
	assert.Equal(t, (*get.Items)[0].UnitCode, (*want.Items)[0].UnitCode, "item.unitcode")
	assert.Equal(t, (*get.Items)[0].WhCode, (*want.Items)[0].WhCode, "item.unitname")
	assert.Equal(t, (*get.Items)[0].LocationCode, (*want.Items)[0].LocationCode, "item.unitname")
	assert.Equal(t, (*get.Items)[0].Qty, (*want.Items)[0].Qty, "item.unitname")
	assert.Equal(t, (*get.Items)[0].Price, (*want.Items)[0].Price, "item.price")
	assert.Equal(t, (*get.Items)[0].PriceExcludeVat, (*want.Items)[0].PriceExcludeVat, "item.PriceExcludeVat")
	assert.Equal(t, (*get.Items)[0].Discount, (*want.Items)[0].Discount, "item.discount")
	assert.Equal(t, (*get.Items)[0].DiscountAmount, (*want.Items)[0].DiscountAmount, "item.discountamount")
	assert.Equal(t, (*get.Items)[0].SumAmount, (*want.Items)[0].SumAmount, "item.sumamount")
	assert.Equal(t, (*get.Items)[0].SumAmountExcludeVat, (*want.Items)[0].SumAmountExcludeVat, "item.sumamountexcludevat")
	assert.Equal(t, (*get.Items)[0].StandValue, (*want.Items)[0].StandValue, "item.standvalue")
	assert.Equal(t, (*get.Items)[0].DivideValue, (*want.Items)[0].DivideValue, "item.dividevalue")
	// assert.Equal(t, *((*get.Items)[0].ItemNames[0]).Name, "โอวัลติน ซอง", "item.ItemNames")
	assert.Equal(t, (*get.Items)[0].ItemType, (*want.Items)[0].ItemType, "item.itemtype")
	assert.Equal(t, (*get.Items)[0].Remark, (*want.Items)[0].Remark, "item.remark")
	assert.Equal(t, (*get.Items)[0].DocRef, (*want.Items)[0].DocRef, "item.docref")
	assert.Equal(t, (*get.Items)[0].DocRefDateTime, (*want.Items)[0].DocRefDateTime, "item.docrefdatetime")
	assert.Equal(t, (*get.Items)[0].VatCal, (*want.Items)[0].VatCal, "item.vatcal")

	// payment
	assert.Equal(t, get.TotalPayCash, want.TotalPayCash, "totalpaycash")
	assert.Equal(t, get.TotalPayCredit, want.TotalPayCredit, "totalpaycredit")
	assert.Equal(t, get.TotalPayTransfer, want.TotalPayTransfer, "totalpaytransfer")
}
