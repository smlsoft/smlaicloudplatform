package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/saleinvoice/models"

	"gorm.io/gorm/clause"
)

type ISaleinvoicePgRepository interface {
	CreateInBatch(docList []models.SaleinvoicePg) error
	Create(doc models.SaleinvoicePg) error
	Update(shopID string, docNo string, doc models.SaleinvoicePg) error
	Delete(shopID string, docNo string) error
	Get(shopID string, docNo string) (*models.SaleinvoicePg, error)
}

type SaleinvoicePgRepository struct {
	pst microservice.IPersister
}

func NewSaleinvoicePgRepository(pst microservice.IPersister) SaleinvoicePgRepository {
	return SaleinvoicePgRepository{
		pst: pst,
	}
}

func (repo SaleinvoicePgRepository) CreateInBatch(docList []models.SaleinvoicePg) error {
	err := repo.pst.CreateInBatch(docList, len(docList))
	if err != nil {
		return err
	}
	return nil
}

func (repo SaleinvoicePgRepository) Create(doc models.SaleinvoicePg) error {
	err := repo.pst.Create(doc)
	if err != nil {
		return err
	}
	return nil
}

func (repo SaleinvoicePgRepository) Update(shopID string, docNo string, doc models.SaleinvoicePg) error {

	err := repo.pst.Update(&doc, map[string]interface{}{
		"shopid": shopID,
		"docno":  docNo,
	})

	if err != nil {
		return err
	}
	return nil
}

func (repo SaleinvoicePgRepository) Delete(shopID string, docNo string) error {
	err := repo.pst.Delete(models.SaleinvoicePg{}, map[string]interface{}{
		"shopid": shopID,
		"docno":  docNo,
	})

	if err != nil {
		return err
	}
	return nil
}

func (repo SaleinvoicePgRepository) Get(shopID string, docNo string) (*models.SaleinvoicePg, error) {

	var data models.SaleinvoicePg

	err := repo.pst.DBClient().Preload(clause.Associations).
		Where("shopid=? AND docno=?", shopID, docNo).
		First(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}
