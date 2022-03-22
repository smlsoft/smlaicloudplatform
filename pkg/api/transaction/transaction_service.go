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
	CreateTransaction(shopID string, username string, trans models.Transaction) (string, error)
	UpdateTransaction(guid string, shopID string, username string, trans models.Transaction) error
	DeleteTransaction(guid string, shopID string, username string) error
	InfoTransaction(guid string, shopID string) (models.TransactionInfo, error)
	SearchTransaction(shopID string, q string, page int, limit int) ([]models.TransactionInfo, paginate.PaginationData, error)
	SearchItemsTransaction(guid string, shopID string, q string, page int, limit int) ([]models.TransactionInfo, paginate.PaginationData, error)
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

func (svc TransactionService) CreateTransaction(shopID string, username string, trans models.Transaction) (string, error) {

	sumAmount := 0.0
	for i, transDetail := range trans.Items {
		trans.Items[i].LineNumber = i + 1
		sumAmount += transDetail.Price * transDetail.Qty
	}

	newGuidFixed := utils.NewGUID()

	transDoc := models.TransactionDoc{}
	transDoc.ShopID = shopID
	transDoc.GuidFixed = newGuidFixed
	transDoc.SumAmount = sumAmount
	transDoc.Transaction = trans

	transDoc.CreatedBy = username
	transDoc.CreatedAt = time.Now()

	_, err := svc.transactionRepository.Create(transDoc)

	if err != nil {
		return "", err
	}

	transData := models.TransactionData{
		ShopIdentity:    transDoc.ShopIdentity,
		TransactionInfo: transDoc.TransactionInfo,
	}

	err = svc.mqRepo.Create(transData)
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

	findDoc.Transaction = trans
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

func (svc TransactionService) InfoTransaction(guid string, shopID string) (models.TransactionInfo, error) {
	trans, err := svc.transactionRepository.FindByGuid(guid, shopID)

	if err != nil {
		return models.TransactionInfo{}, err
	}

	return trans.TransactionInfo, nil
}

func (svc TransactionService) SearchTransaction(shopID string, q string, page int, limit int) ([]models.TransactionInfo, paginate.PaginationData, error) {
	transList, pagination, err := svc.transactionRepository.FindPage(shopID, q, page, limit)

	if err != nil {
		return transList, pagination, err
	}

	return transList, pagination, nil
}

func (svc TransactionService) SearchItemsTransaction(guid string, shopID string, q string, page int, limit int) ([]models.TransactionInfo, paginate.PaginationData, error) {
	transList, pagination, err := svc.transactionRepository.FindItemsByGuidPage(guid, shopID, q, page, limit)

	if err != nil {
		return transList, pagination, err
	}

	return transList, pagination, nil
}
