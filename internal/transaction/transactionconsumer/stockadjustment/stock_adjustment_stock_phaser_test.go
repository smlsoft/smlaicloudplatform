package stockadjustment_test

import (
	pkgModels "smlaicloudplatform/internal/models"
	"smlaicloudplatform/internal/transaction/models"
	"smlaicloudplatform/internal/transaction/transactionconsumer/stockadjustment"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStockAdjustmentStockTransactionPhaser(t *testing.T) {
	give := wantStockAdjustmentTransactionStruct()
	want := models.StockTransaction{
		GuidFixed: "2PxeTSlssQvMZS8MViihtgYOC0w",
		ShopIdentity: pkgModels.ShopIdentity{
			ShopID: "2PrIIqTWxoBXv16K310sNwfHmfY",
		},
		TransFlag:      66,
		InquiryType:    0,
		DocNo:          "PO2305111645E3FE",
		DocDate:        time.Date(2023, 5, 11, 9, 45, 46, 0, time.UTC),
		GuidRef:        "1a14456e-c859-422e-8f3c-ca9787202f0c",
		DocRefType:     0,
		DocRefNo:       "",
		DocRefDate:     time.Date(2023, 5, 18, 9, 45, 39, 0, time.UTC),
		VatType:        0,
		VatRate:        0,
		DiscountWord:   "",
		TotalDiscount:  0,
		TotalValue:     0,
		TotalBeforeVat: 0,
		TotalExceptVat: 0,
		TotalVatValue:  0,
		TotalAfterVat:  0,
		TotalAmount:    0,
		TotalCost:      0,
		Status:         0,
		IsCancel:       false,
		Details: &[]models.StockTransactionDetail{
			{
				ShopID:              "2PrIIqTWxoBXv16K310sNwfHmfY",
				DocNo:               "PO2305111645E3FE",
				Barcode:             "BARCODE015",
				UnitCode:            "PCE",
				Qty:                 5,
				Price:               0,
				PriceExcludeVat:     0,
				Discount:            "",
				DiscountAmount:      0,
				SumAmount:           0,
				SumAmountExcludeVat: 0,
				StandValue:          1,
				DivideValue:         1,
				WhCode:              "00000",
				LocationCode:        "",
				CalcFlag:            1,
				VatType:             0,
				DocRef:              "",
			},
		},
	}

	stockPhaser := stockadjustment.StockAdjustmentStockPhaser{}
	get, err := stockPhaser.PhaseSingleDoc(give)

	assert.Nil(t, err)

	assert.Equal(t, want.GuidFixed, get.GuidFixed, "GuidFixed")
	assert.Equal(t, want.ShopID, get.ShopID, "ShopID")
	assert.Equal(t, want.TransFlag, get.TransFlag, "TransFlag")
	assert.Equal(t, want.InquiryType, get.InquiryType, "InquiryType")
	assert.Equal(t, want.DocNo, get.DocNo, "DocNo")
	assert.Equal(t, want.DocDate, get.DocDate, "DocDate")
	assert.Equal(t, want.GuidRef, get.GuidRef, "GuidRef")
	assert.Equal(t, want.DocRefType, get.DocRefType, "DocRefType")
	assert.Equal(t, want.DocRefNo, get.DocRefNo, "DocRefNo")
	assert.Equal(t, want.DocRefDate, get.DocRefDate, "DocRefDate")
	assert.Equal(t, want.VatType, get.VatType, "VatType")
	assert.Equal(t, want.DiscountWord, get.DiscountWord, "DiscountWord")
	assert.Equal(t, want.TotalDiscount, get.TotalDiscount, "TotalDiscount")
	assert.Equal(t, want.TotalValue, get.TotalValue, "TotalValue")
	assert.Equal(t, want.TotalBeforeVat, get.TotalBeforeVat, "TotalBeforeVat")
	assert.Equal(t, want.TotalExceptVat, get.TotalExceptVat, "TotalExceptVat")
	assert.Equal(t, want.TotalVatValue, get.TotalVatValue, "TotalVatValue")
	assert.Equal(t, want.TotalAfterVat, get.TotalAfterVat, "TotalAfterVat")
	assert.Equal(t, want.TotalAmount, get.TotalAmount, "TotalAmount")
	assert.Equal(t, want.TotalCost, get.TotalCost, "TotalCost")
	assert.Equal(t, want.Status, get.Status, "Status")
	assert.Equal(t, want.IsCancel, get.IsCancel, "IsCancel")

	assert.Equal(t, (*want.Details)[0].ShopID, (*get.Details)[0].ShopID, "Details[0].ShopID")
	assert.Equal(t, (*want.Details)[0].DocNo, (*get.Details)[0].DocNo, "Details[0].DocNo")
	assert.Equal(t, (*want.Details)[0].Barcode, (*get.Details)[0].Barcode, "Details[0].Barcode")
	assert.Equal(t, (*want.Details)[0].UnitCode, (*get.Details)[0].UnitCode, "Details[0].UnitCode")
	assert.Equal(t, (*want.Details)[0].Qty, (*get.Details)[0].Qty, "Details[0].Qty")
	assert.Equal(t, (*want.Details)[0].Price, (*get.Details)[0].Price, "Details[0].Price")
	assert.Equal(t, (*want.Details)[0].PriceExcludeVat, (*get.Details)[0].PriceExcludeVat, "Details[0].PriceExcludeVat")
	assert.Equal(t, (*want.Details)[0].Discount, (*get.Details)[0].Discount, "Details[0].Discount")
	assert.Equal(t, (*want.Details)[0].DiscountAmount, (*get.Details)[0].DiscountAmount, "Details[0].DiscountAmount")
	assert.Equal(t, (*want.Details)[0].SumAmount, (*get.Details)[0].SumAmount, "Details[0].SumAmount")
	assert.Equal(t, (*want.Details)[0].SumAmountExcludeVat, (*get.Details)[0].SumAmountExcludeVat, "Details[0].SumAmountExcludeVat")
	assert.Equal(t, (*want.Details)[0].StandValue, (*get.Details)[0].StandValue, "Details[0].StandValue")
	assert.Equal(t, (*want.Details)[0].DivideValue, (*get.Details)[0].DivideValue, "Details[0].DivideValue")
	assert.Equal(t, (*want.Details)[0].WhCode, (*get.Details)[0].WhCode, "Details[0].WhCode")
	assert.Equal(t, (*want.Details)[0].LocationCode, (*get.Details)[0].LocationCode, "Details[0].LocationCode")
	assert.Equal(t, (*want.Details)[0].CalcFlag, (*get.Details)[0].CalcFlag, "Details[0].CalcFlag")
	assert.Equal(t, (*want.Details)[0].VatType, (*get.Details)[0].VatType, "Details[0].VatType")
	assert.Equal(t, (*want.Details)[0].DocRef, (*get.Details)[0].DocRef, "Details[0].DocRef")
}
