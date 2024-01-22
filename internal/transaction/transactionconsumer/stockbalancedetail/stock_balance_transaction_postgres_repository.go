package stockbalancedetail

import (
	"smlcloudplatform/internal/transaction/models"
	"smlcloudplatform/internal/transaction/transactionconsumer/repositories"
	"smlcloudplatform/pkg/microservice"
)

type IStockReceiveTransactionDetailPGRepository interface {
	Get(shopID string, docNo string) (*models.StockBalanceTransactionDetailPG, error)
	Create(doc models.StockBalanceTransactionDetailPG) error
	Update(shopID string, docNo string, doc models.StockBalanceTransactionDetailPG) error
	Delete(shopID string, docNo string, doc models.StockBalanceTransactionDetailPG) error
	DeleteData(shopID string, docNo string, doc models.StockBalanceTransactionDetailPG) error
}

type StockReceiveTransactionDetailPGRepository struct {
	pst microservice.IPersister
	repositories.ITransactionConsumerRepository[models.StockBalanceTransactionDetailPG]
}

func NewStockReceiveTransactionDetailPGRepository(pst microservice.IPersister) IStockReceiveTransactionDetailPGRepository {

	repo := &StockReceiveTransactionDetailPGRepository{
		pst: pst,
	}

	repo.ITransactionConsumerRepository = repositories.NewTransactionConsumerRepository[models.StockBalanceTransactionDetailPG](pst)
	return repo
}

func (repo *StockReceiveTransactionDetailPGRepository) DeleteData(shopID string, docNo string, doc models.StockBalanceTransactionDetailPG) error {

	var details *[]models.StockBalanceTransactionDetailPG
	tx := repo.pst.DBClient().Begin()

	tx.Model(&models.StockBalanceTransactionDetailPG{}).Where(" shopid=? AND docno=?", shopID, docNo).Find(&details)
	for _, tmp := range *details {
		// mark delete
		tx.Delete(&models.StockBalanceTransactionDetailPG{}, tmp.ID)
	}

	err := tx.Delete(models.StockBalanceTransactionDetailPG{}, map[string]interface{}{
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
