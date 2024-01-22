package stockbalancedetail

import (
	"encoding/json"
	"errors"
	"smlcloudplatform/internal/transaction/models"
	stockBalanceDetailModels "smlcloudplatform/internal/transaction/stockbalancedetail/models"
)

type StockBalanceTransactionPhaser struct{}

func (p StockBalanceTransactionPhaser) PhaseSingleDoc(input string) (*models.StockBalanceTransactionDetailPG, error) {

	doc := stockBalanceDetailModels.StockBalanceDetailDoc{}
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

func (p StockBalanceTransactionPhaser) PhaseMultipleDoc(input string) (*[]models.StockBalanceTransactionDetailPG, error) {
	docs := []stockBalanceDetailModels.StockBalanceDetailDoc{}
	err := json.Unmarshal([]byte(input), &docs)
	if err != nil {
		return nil, errors.New("Cannot Unmarshal StockBalance Doc Error: " + err.Error())
	}

	transactions := make([]models.StockBalanceTransactionDetailPG, len(docs))
	for i, doc := range docs {
		transaction, err := p.PhaseStockBalanceTransaction(doc)
		if err != nil {
			return nil, errors.New("Cannot Phase StockBalance Doc Error: " + err.Error())
		}
		transactions[i] = *transaction
	}
	return &transactions, nil
}

func (p StockBalanceTransactionPhaser) PhaseStockBalanceTransaction(doc stockBalanceDetailModels.StockBalanceDetailDoc) (*models.StockBalanceTransactionDetailPG, error) {

	stockDetail := models.StockBalanceTransactionDetailPG{
		TransactionDetailPG: models.TransactionDetailPG{
			DocNo:               doc.DocNo,
			ShopID:              doc.ShopID,
			LineNumber:          int8(doc.LineNumber),
			DocRef:              doc.DocRef,
			DocRefDateTime:      doc.DocRefDatetime,
			Barcode:             doc.Barcode,
			UnitCode:            doc.UnitCode,
			Qty:                 doc.Qty,
			Price:               doc.Price,
			PriceExcludeVat:     doc.PriceExcludeVat,
			Discount:            doc.Discount,
			DiscountAmount:      doc.DiscountAmount,
			SumAmount:           doc.SumAmount,
			SumAmountExcludeVat: doc.SumAmountExcludeVat,
			WhCode:              doc.WhCode,
			LocationCode:        doc.LocationCode,
			VatType:             doc.VatType,
			TaxType:             doc.TaxType,
			StandValue:          doc.StandValue,
			DivideValue:         doc.DivideValue,
			ItemType:            doc.ItemType,
			ItemGuid:            doc.ItemGuid,
			TotalValueVat:       doc.TotalValueVat,
			Remark:              doc.Remark,
			ItemNames:           *doc.ItemNames,
			WhNames:             *doc.WhNames,
			LocationNames:       *doc.LocationNames,
			UnitNames:           *doc.UnitNames,
		},
	}

	return &stockDetail, nil
}
