package transaction

import (
	"errors"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ITransactionService interface {
	CreateTransaction(shopID string, username string, trans *models.Transaction) (string, error)
	UpdateTransaction(guid string, shopID string, username string, trans models.Transaction) error
	DeleteTransaction(guid string, shopID string, username string) error
	InfoTransaction(guid string, shopID string) (models.Transaction, error)
	SearchTransaction(shopID string, q string, page int, limit int) ([]models.Transaction, paginate.PaginationData, error)
	SearchItemsTransaction(guid string, shopID string, q string, page int, limit int) ([]models.Transaction, paginate.PaginationData, error)
}

type TransactionService struct {
	transactionRepository ITransactionRepository
	mqRepo                ITransactionMQRepository
}

func NewTransactionService(transactionRepository ITransactionRepository, mqRepo ITransactionMQRepository) TransactionService {

	return TransactionService{
		transactionRepository: transactionRepository,
		mqRepo:                mqRepo,
	}
}

func (svc TransactionService) CreateTransaction(shopID string, username string, trans *models.Transaction) (string, error) {

	sumAmount := 0.0
	for i, transDetail := range trans.Items {
		trans.Items[i].LineNumber = i + 1
		sumAmount += transDetail.Price * transDetail.Qty
	}

	newGuidFixed := utils.NewGUID()
	trans.ShopID = shopID
	trans.GuidFixed = newGuidFixed
	trans.SumAmount = sumAmount
	trans.Deleted = false
	trans.CreatedBy = username
	trans.CreatedAt = time.Now()

	idx, err := svc.transactionRepository.Create(*trans)

	if err != nil {
		return "", err
	}

	trans.ID = idx

	transReq := &models.TransactionRequest{}
	transReq.MapRequest(*trans)

	err = svc.mqRepo.Create(*transReq)
	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc TransactionService) UpdateTransaction(guid string, shopID string, username string, trans models.Transaction) error {

	findDoc, err := svc.transactionRepository.FindByGuid(guid, shopID)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("guid invalid")
	}

	sumAmount := 0.0
	for i, transDetail := range trans.Items {
		findDoc.Items[i].LineNumber = i + 1
		sumAmount += transDetail.Price * transDetail.Qty
	}

	findDoc.Items = trans.Items
	findDoc.SumAmount = sumAmount
	findDoc.UpdatedBy = username
	findDoc.UpdatedAt = time.Now()

	err = svc.transactionRepository.Update(guid, findDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc TransactionService) DeleteTransaction(guid string, shopID string, username string) error {

	err := svc.transactionRepository.Delete(guid, shopID)
	if err != nil {
		return err
	}

	return nil
}

func (svc TransactionService) InfoTransaction(guid string, shopID string) (models.Transaction, error) {
	trans, err := svc.transactionRepository.FindByGuid(guid, shopID)

	if err != nil {
		return models.Transaction{}, err
	}

	return trans, nil
}

func (svc TransactionService) SearchTransaction(shopID string, q string, page int, limit int) ([]models.Transaction, paginate.PaginationData, error) {
	transList, pagination, err := svc.transactionRepository.FindPage(shopID, q, page, limit)

	if err != nil {
		return transList, pagination, err
	}

	return transList, pagination, nil
}

func (svc TransactionService) SearchItemsTransaction(guid string, shopID string, q string, page int, limit int) ([]models.Transaction, paginate.PaginationData, error) {
	transList, pagination, err := svc.transactionRepository.FindItemsByGuidPage(guid, shopID, q, page, limit)

	if err != nil {
		return transList, pagination, err
	}

	return transList, pagination, nil
}
