package transactionservice

import (
	"errors"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
)

type ITransactionService interface {
	CreateTransaction(merchantId string, username string, trans models.Transaction) (string, error)
	UpdateTransaction(guid string, merchantId string, username string, trans models.Transaction) error
	InfoTransaction(guid string, merchantId string) (models.Transaction, error)
	SearchTransaction(merchantId string, q string, page int, limit int) ([]models.Transaction, paginate.PaginationData, error)
	SearchItemsTransaction(guid string, merchantId string, q string, page int, limit int) ([]models.Transaction, paginate.PaginationData, error)
}

type TransactionService struct {
	transactionRepository ITransactionRepository
}

func NewTransactionService(transactionRepository ITransactionRepository) ITransactionService {

	return &TransactionService{
		transactionRepository: transactionRepository,
	}
}

func (svc *TransactionService) CreateTransaction(merchantId string, username string, trans models.Transaction) (string, error) {

	sumAmount := 0.0
	for i, transDetail := range trans.Items {
		trans.Items[i].LineNumber = i + 1
		sumAmount += transDetail.Price * transDetail.Qty
	}

	newGuidFixed := utils.NewGUID()
	trans.MerchantId = merchantId
	trans.GuidFixed = newGuidFixed
	trans.SumAmount = sumAmount
	trans.Deleted = false
	trans.CreatedBy = username
	trans.CreatedAt = time.Now()

	_, err := svc.transactionRepository.Create(trans)

	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc *TransactionService) UpdateTransaction(guid string, merchantId string, username string, trans models.Transaction) error {

	findTrans, err := svc.transactionRepository.FindByGuid(guid, merchantId)

	if err != nil {
		return errors.New("guid invalid")
	}

	sumAmount := 0.0
	for i, transDetail := range findTrans.Items {
		findTrans.Items[i].LineNumber = i + 1
		sumAmount += transDetail.Price * transDetail.Qty
	}

	findTrans.Items = trans.Items
	findTrans.SumAmount = sumAmount
	findTrans.UpdatedBy = username
	findTrans.UpdatedAt = time.Now()

	err = svc.transactionRepository.Update(guid, findTrans)

	if err != nil {
		return err
	}
	return nil
}

func (svc *TransactionService) InfoTransaction(guid string, merchantId string) (models.Transaction, error) {
	trans, err := svc.transactionRepository.FindByGuid(guid, merchantId)

	if err != nil {
		return models.Transaction{}, err
	}

	return trans, nil
}

func (svc *TransactionService) SearchTransaction(merchantId string, q string, page int, limit int) ([]models.Transaction, paginate.PaginationData, error) {
	transList, pagination, err := svc.transactionRepository.FindPage(merchantId, q, page, limit)

	if err != nil {
		return transList, pagination, err
	}

	return transList, pagination, nil
}

func (svc *TransactionService) SearchItemsTransaction(guid string, merchantId string, q string, page int, limit int) ([]models.Transaction, paginate.PaginationData, error) {
	transList, pagination, err := svc.transactionRepository.FindItemsByGuidPage(guid, merchantId, q, page, limit)

	if err != nil {
		return transList, pagination, err
	}

	return transList, pagination, nil
}
