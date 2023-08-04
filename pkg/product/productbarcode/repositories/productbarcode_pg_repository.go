package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/productbarcode/models"

	"gorm.io/gorm"
)

type IProductBarcodePGRepository interface {
	Get(shopID string, barcode string) (*models.ProductBarcodePg, error)
	Create(doc *models.ProductBarcodePg) error
	Update(shopID string, barcode string, doc *models.ProductBarcodePg) error
	Delete(shopID string, barcode string) error
}

type ProductBarcodePGRepository struct {
	pst microservice.IPersister
}

func NewProductBarcodePGRepository(pst microservice.IPersister) *ProductBarcodePGRepository {
	return &ProductBarcodePGRepository{
		pst: pst,
	}
}

func (repo *ProductBarcodePGRepository) Get(shopID string, barcode string) (*models.ProductBarcodePg, error) {
	var result models.ProductBarcodePg
	_, err := repo.pst.First(&result, "shopid=? AND barcode=?", shopID, barcode)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &result, nil
}

func (repo *ProductBarcodePGRepository) Create(doc *models.ProductBarcodePg) error {
	err := repo.pst.Create(doc)
	if err != nil {
		return err
	}
	return nil
}

func (repo *ProductBarcodePGRepository) Update(shopID string, barcode string, doc *models.ProductBarcodePg) error {
	err := repo.pst.Update(&doc, map[string]interface{}{
		"shopid":  shopID,
		"barcode": barcode,
	})

	if err != nil {
		return err
	}
	return nil
}

func (repo *ProductBarcodePGRepository) Delete(shopID string, barcode string) error {

	err := repo.pst.Delete(&models.ProductBarcodePg{}, map[string]interface{}{
		"shopid":  shopID,
		"barcode": barcode,
	})

	if err != nil {
		return err
	}
	return nil
}
