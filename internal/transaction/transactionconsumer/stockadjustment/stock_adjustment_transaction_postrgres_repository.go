package stockadjustment

import (
	"smlcloudplatform/internal/transaction/models"
	"smlcloudplatform/internal/transaction/transactionconsumer/repositories"
	"smlcloudplatform/pkg/microservice"
)

type IStockAdjustmentTransactionPGRepository interface {
	Get(shopID string, docNo string) (*models.StockAdjustmentTransactionPG, error)
	Create(doc models.StockAdjustmentTransactionPG) error
	Update(shopID string, docNo string, doc models.StockAdjustmentTransactionPG) error
	Delete(shopID string, docNo string, doc models.StockAdjustmentTransactionPG) error
	DeleteData(shopID string, docNo string, doc models.StockAdjustmentTransactionPG) error
}

type StockAdjustmentTransactionPGRepository struct {
	pst microservice.IPersister
	repositories.ITransactionConsumerRepository[models.StockAdjustmentTransactionPG]
}

func NewStockAdjustmentTransactionPGRepository(pst microservice.IPersister) IStockAdjustmentTransactionPGRepository {

	repo := &StockAdjustmentTransactionPGRepository{
		pst: pst,
	}

	repo.ITransactionConsumerRepository = repositories.NewTransactionConsumerRepository[models.StockAdjustmentTransactionPG](pst)
	return repo
}

func (repo *StockAdjustmentTransactionPGRepository) DeleteData(shopID string, docNo string, doc models.StockAdjustmentTransactionPG) error {

	var details *[]models.StockAdjustmentTransactionDetailPG
	tx := repo.pst.DBClient().Begin()

	tx.Model(&models.StockAdjustmentTransactionDetailPG{}).Where(" shopid=? AND docno=?", shopID, docNo).Find(&details)
	for _, tmp := range *details {
		// mark delete
		tx.Delete(&models.StockAdjustmentTransactionDetailPG{}, tmp.ID)
	}

	err := tx.Delete(models.StockAdjustmentTransactionPG{}, map[string]interface{}{
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
