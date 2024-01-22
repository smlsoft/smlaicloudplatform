package services

import (
	"context"
	"fmt"
	"smlcloudplatform/internal/product/productbarcode/models"
)

func (svc ProductBarcodeHttpService) InfoBomView(shopID string, barcode string) (models.ProductBarcodeBOMView, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	doc, err := svc.repo.FindByBarcode(ctx, shopID, barcode)

	if err != nil {
		return models.ProductBarcodeBOMView{}, err
	}

	if len(doc.ProductBarcode.Barcode) == 0 {
		return models.ProductBarcodeBOMView{}, fmt.Errorf("barcode not found")
	}

	bomView := models.ProductBarcodeBOMView{}
	bomView.FromProductBarcode(doc.ProductBarcodeData)

	if doc.BOM != nil && len(*doc.BOM) > 0 {
		err = svc.buildBOMView(ctx, 1, shopID, doc.BOM, &bomView.BOM)
		if err != nil {
			return models.ProductBarcodeBOMView{}, err
		}
	}

	return bomView, nil
}

func (svc ProductBarcodeHttpService) buildBOMView(ctx context.Context, currentLevel int, shopID string, BOMs *[]models.BOMProductBarcode, bomView *[]models.ProductBarcodeBOMView) error {

	if currentLevel > 21 {
		return fmt.Errorf("BOM level is too deep")
	}

	currentLevel += 1

	for _, bom := range *BOMs {

		tempDoc, err := svc.repo.FindByBarcode(ctx, shopID, bom.Barcode)

		if err != nil {
			return err
		}

		tempBOMView := models.ProductBarcodeBOMView{}
		tempBOMView.FromProductBOM(tempDoc.ProductBarcodeData, bom)

		if tempDoc.BOM != nil && len(*tempDoc.BOM) > 0 {
			err = svc.buildBOMView(ctx, currentLevel, shopID, tempDoc.BOM, &tempBOMView.BOM)

			if err != nil {
				return err
			}
		}

		*bomView = append(*bomView, tempBOMView)

	}

	return nil
}
