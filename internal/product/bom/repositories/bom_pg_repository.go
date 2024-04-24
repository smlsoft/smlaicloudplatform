package repositories

import (
	"smlcloudplatform/internal/product/bom/models"
	"smlcloudplatform/internal/transaction/transactionconsumer/repositories"
	"smlcloudplatform/pkg/microservice"

	"gorm.io/gorm/clause"
)

type IBomPgRepository interface {
	Get(shopID string, docNo string) (*models.ProductBarcodeBOMViewPG, error)
	Create(doc models.ProductBarcodeBOMViewPG) error
	Update(shopID string, docNo string, doc models.ProductBarcodeBOMViewPG) error
	Delete(shopID string, docNo string, doc models.ProductBarcodeBOMViewPG) error
}

type BomPgRepository struct {
	pst microservice.IPersister
	repositories.ITransactionConsumerRepository[models.ProductBarcodeBOMViewPG]
}

func NewBomPgRepository(pst microservice.IPersister) *BomPgRepository {

	repo := &BomPgRepository{
		pst: pst,
	}

	return repo
}

func (repo BomPgRepository) Get(shopID string, docNo string) (*models.ProductBarcodeBOMViewPG, error) {

	var data models.ProductBarcodeBOMViewPG
	err := repo.pst.DBClient().Preload(clause.Associations).
		Where("shopid=? AND docno=?", shopID, docNo).
		First(&data).Error

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (repo BomPgRepository) Create(doc models.ProductBarcodeBOMViewPG) error {
	err := repo.pst.Create(&doc)
	if err != nil {
		return err
	}
	return nil
}

func (repo BomPgRepository) Update(shopID string, docNo string, doc models.ProductBarcodeBOMViewPG) error {

	err := repo.pst.Update(&doc, map[string]interface{}{
		"shopid": shopID,
		"docno":  docNo,
	})

	if err != nil {
		return err
	}
	return nil
}

func (repo BomPgRepository) Delete(shopID string, docNo string, doc models.ProductBarcodeBOMViewPG) error {

	tx := repo.pst.DBClient().Begin()

	err := tx.Delete(&doc, map[string]interface{}{
		"shopid": shopID,
		"docno":  docNo,
	}).Error

	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
