package usecases

import (
	"encoding/json"
	"smlcloudplatform/internal/product/productbarcode/models"
)

type IProductBarcodePhaser interface {
	PhaseProductBarcodeDoc(doc *models.ProductBarcodeDoc) (*models.ProductBarcodePg, error)
}

type ProductBarcodePhaser struct{}

func (ProductBarcodePhaser) PhaseProductBarcodeDoc(doc *models.ProductBarcodeDoc) (*models.ProductBarcodePg, error) {

	j, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	var pgDoc models.ProductBarcodePg
	err = json.Unmarshal(j, &pgDoc)
	if err != nil {
		return nil, err
	}

	pgDoc.UnitCode = doc.ItemUnitCode
	return &pgDoc, nil
}
