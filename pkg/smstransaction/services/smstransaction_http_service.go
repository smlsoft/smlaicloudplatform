package services

import (
	"errors"
	"smlcloudplatform/pkg/smstransaction/models"
	"smlcloudplatform/pkg/smstransaction/repositories"
	"smlcloudplatform/pkg/utils"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ISmsTransactionHttpService interface {
	CreateSmsTransaction(shopID string, authUsername string, doc models.SmsTransaction) (string, error)
	UpdateSmsTransaction(guid string, shopID string, authUsername string, doc models.SmsTransaction) error
	DeleteSmsTransaction(guid string, shopID string, authUsername string) error
	InfoSmsTransaction(guid string, shopID string) (models.SmsTransactionInfo, error)
	SearchSmsTransaction(shopID string, q string, page int, limit int, sort map[string]int) ([]models.SmsTransactionInfo, mongopagination.PaginationData, error)
}

type SmsTransactionHttpService struct {
	repo repositories.SmsTransactionRepository
}

func NewSmsTransactionHttpService(repo repositories.SmsTransactionRepository) SmsTransactionHttpService {

	return SmsTransactionHttpService{
		repo: repo,
	}
}

func (svc SmsTransactionHttpService) CreateSmsTransaction(shopID string, authUsername string, doc models.SmsTransaction) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentiryGuid(shopID, "docno", doc.TransId)

	if err != nil {
		return "", err
	}

	if findDoc.TransId != "" {
		return "", errors.New("TransId is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.SmsTransactionDoc{}
	docData.ShopID = shopID
	docData.SmsTransaction = doc

	docData.GuidFixed = newGuidFixed
	docData.TransId = utils.NewGUID()

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc SmsTransactionHttpService) UpdateSmsTransaction(guid string, shopID string, authUsername string, doc models.SmsTransaction) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	tempTransId := findDoc.TransId

	findDoc.SmsTransaction = doc
	findDoc.TransId = tempTransId

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc SmsTransactionHttpService) DeleteSmsTransaction(guid string, shopID string, authUsername string) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	err = svc.repo.DeleteByGuidfixed(shopID, guid, authUsername)
	if err != nil {
		return err
	}
	return nil
}

func (svc SmsTransactionHttpService) InfoSmsTransaction(guid string, shopID string) (models.SmsTransactionInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.SmsTransactionInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.SmsTransactionInfo{}, errors.New("document not found")
	}

	return findDoc.SmsTransactionInfo, nil

}

func (svc SmsTransactionHttpService) SearchSmsTransaction(shopID string, q string, page int, limit int, sort map[string]int) ([]models.SmsTransactionInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"guidfixed",
		"docno",
	}

	docList, pagination, err := svc.repo.FindPageSort(shopID, searchCols, q, page, limit, sort)

	if err != nil {
		return []models.SmsTransactionInfo{}, pagination, err
	}

	return docList, pagination, nil
}
