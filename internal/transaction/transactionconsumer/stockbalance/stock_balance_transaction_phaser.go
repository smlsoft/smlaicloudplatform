package stockbalance

import (
	"encoding/json"
	"errors"
	pkgModels "smlcloudplatform/internal/models"
	"smlcloudplatform/internal/transaction/models"
	stockBalanceModels "smlcloudplatform/internal/transaction/stockbalance/models"
)

type StockBalanceTransactionPhaser struct{}

func (p StockBalanceTransactionPhaser) PhaseSingleDoc(input string) (*models.StockBalanceTransactionPG, error) {

	doc := stockBalanceModels.StockBalanceDoc{}
	err := json.Unmarshal([]byte(input), &doc)
	if err != nil {
		return nil, errors.New("Cannot Unmarshal StockBalance Doc Error: " + err.Error())
	}

	transaction, err := p.PhaseStockBalanceTransaction(doc)
	if err != nil {
		return nil, errors.New("Cannot Phase StockBalance Doc Error: " + err.Error())
	}
	return transaction, nil

}

func (p StockBalanceTransactionPhaser) PhaseMultipleDoc(input string) (*[]models.StockBalanceTransactionPG, error) {
	docs := []stockBalanceModels.StockBalanceDoc{}
	err := json.Unmarshal([]byte(input), &docs)
	if err != nil {
		return nil, errors.New("Cannot Unmarshal StockBalance Doc Error: " + err.Error())
	}

	transactions := make([]models.StockBalanceTransactionPG, len(docs))
	for i, doc := range docs {
		transaction, err := p.PhaseStockBalanceTransaction(doc)
		if err != nil {
			return nil, errors.New("Cannot Phase StockBalance Doc Error: " + err.Error())
		}
		transactions[i] = *transaction
	}
	return &transactions, nil
}

func (p StockBalanceTransactionPhaser) PhaseStockBalanceTransaction(doc stockBalanceModels.StockBalanceDoc) (*models.StockBalanceTransactionPG, error) {

	details := make([]models.StockBalanceTransactionDetailPG, len(*doc.Transaction.Details))

	for i, detail := range *doc.Transaction.Details {
		stockDetail := models.StockBalanceTransactionDetailPG{
			TransactionDetailPG: models.TransactionDetailPG{
				DocNo:               doc.DocNo,
				ShopID:              doc.ShopID,
				LineNumber:          int8(detail.LineNumber),
				DocRef:              detail.DocRef,
				DocRefDateTime:      detail.DocRefDatetime,
				Barcode:             detail.Barcode,
				UnitCode:            detail.UnitCode,
				Qty:                 detail.Qty,
				Price:               detail.Price,
				PriceExcludeVat:     detail.PriceExcludeVat,
				Discount:            detail.Discount,
				DiscountAmount:      detail.DiscountAmount,
				SumAmount:           detail.SumAmount,
				SumAmountExcludeVat: detail.SumAmountExcludeVat,
				WhCode:              detail.WhCode,
				LocationCode:        detail.LocationCode,
				VatType:             detail.VatType,
				TaxType:             detail.TaxType,
				StandValue:          detail.StandValue,
				DivideValue:         detail.DivideValue,
				ItemType:            detail.ItemType,
				ItemGuid:            detail.ItemGuid,
				TotalValueVat:       detail.TotalValueVat,
				Remark:              detail.Remark,
				ItemNames:           *detail.ItemNames,
				WhNames:             *detail.WhNames,
				LocationNames:       *detail.LocationNames,
				UnitNames:           *detail.UnitNames,
			},
		}
		details[i] = stockDetail
	}

	stockTransaction := models.StockBalanceTransactionPG{
		TransactionPG: models.TransactionPG{
			ShopIdentity: pkgModels.ShopIdentity{
				ShopID: doc.ShopID,
			},
			GuidFixed:      doc.GuidFixed,
			GuidRef:        doc.GuidRef,
			InquiryType:    doc.InquiryType,
			TransFlag:      54,
			DocNo:          doc.DocNo,
			DocDate:        doc.DocDatetime,
			TaxDocDate:     doc.Transaction.TaxDocDate,
			TaxDocNo:       doc.Transaction.TaxDocNo,
			VatType:        doc.Transaction.VatType,
			VatRate:        doc.Transaction.VatRate,
			DocRefType:     doc.Transaction.DocRefType,
			DocRefNo:       doc.Transaction.DocRefNo,
			DocRefDate:     doc.Transaction.DocRefDate,
			Description:    doc.Transaction.Description,
			TotalValue:     doc.Transaction.TotalValue,
			DiscountWord:   doc.Transaction.DiscountWord,
			TotalDiscount:  doc.Transaction.TotalDiscount,
			TotalBeforeVat: doc.Transaction.TotalBeforeVat,
			TotalVatValue:  doc.Transaction.TotalVatValue,
			TotalAfterVat:  doc.Transaction.TotalAfterVat,
			TotalExceptVat: doc.Transaction.TotalExceptVat,
			TotalAmount:    doc.Transaction.TotalAmount,
		},
		Items: &details,
	}
	return &stockTransaction, nil

}
