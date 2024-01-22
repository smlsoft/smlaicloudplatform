package stockbalancedetail

import (
	pkgModels "smlcloudplatform/internal/models"
	"smlcloudplatform/internal/transaction/models"
)

type StockBalanceDetailStockPhaser struct{}

func (p StockBalanceDetailStockPhaser) PhaseSingleDoc(doc models.StockBalanceTransactionDetailPG) (*models.StockTransaction, error) {

	stockDetail := models.StockTransactionDetail{
		CalcFlag:            1,
		DocRef:              doc.DocRef,
		ShopID:              doc.ShopID,
		DocNo:               doc.DocNo,
		Barcode:             doc.Barcode,
		ItemType:            doc.ItemType,
		ItemGuid:            doc.ItemGuid,
		VatType:             doc.VatType,
		TaxType:             doc.TaxType,
		UnitCode:            doc.UnitCode,
		StandValue:          doc.StandValue,
		DivideValue:         doc.DivideValue,
		WhCode:              doc.WhCode,
		LocationCode:        doc.LocationCode,
		Qty:                 doc.Qty,
		Price:               doc.Price,
		PriceExcludeVat:     doc.PriceExcludeVat,
		TotalValueVat:       doc.TotalValueVat,
		SumAmount:           doc.SumAmount,
		SumAmountExcludeVat: doc.SumAmountExcludeVat,
		Discount:            doc.Discount,
		DiscountAmount:      doc.DiscountAmount,
		LineNumber:          int8(doc.LineNumber),
	}

	stockTransaction := models.StockTransaction{
		ShopIdentity: pkgModels.ShopIdentity{
			ShopID: doc.ShopID,
		},
		GuidFixed:      doc.GuidFixed,
		GuidRef:        doc.GuidRef,
		DocRefType:     doc.DocRefType,
		DocRefNo:       doc.DocRefNo,
		DocRefDate:     doc.DocRefDate,
		VatType:        doc.VatType,
		TransFlag:      60,
		InquiryType:    doc.InquiryType,
		DocNo:          doc.DocNo,
		DocDate:        doc.DocDate,
		Details:        &[]models.StockTransactionDetail{stockDetail},
		TotalValue:     doc.TotalValue,
		DiscountWord:   doc.DiscountWord,
		TotalDiscount:  doc.TotalDiscount,
		TotalBeforeVat: doc.TotalBeforeVat,
		TotalVatValue:  doc.TotalVatValue,
		TotalExceptVat: doc.TotalExceptVat,
		TotalAfterVat:  doc.TotalAfterVat,
		TotalAmount:    doc.TotalAmount,
	}
	return &stockTransaction, nil
}
