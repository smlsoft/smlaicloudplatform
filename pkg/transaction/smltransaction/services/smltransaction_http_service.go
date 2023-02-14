package services

import (
	"fmt"
	"smlcloudplatform/pkg/transaction/smltransaction/models"
	"smlcloudplatform/pkg/transaction/smltransaction/repositories"
	"smlcloudplatform/pkg/utils"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type ISMLTransactionHttpService interface {
	CreateSMLTransaction(shopID string, authUsername string, smlRequest models.SMLTransactionRequest) (string, error)
	SaveInBatch(shopID string, authUsername string, dataReq models.SMLTransactionBulkRequest) ([]string, error)
	DeleteSMLTransaction(shopID string, authUsername string, smlKeyRequest models.SMLTransactionKeyRequest) ([]string, error)
}

type SMLTransactionHttpService struct {
	repo   repositories.ISMLTransactionRepository
	mqRepo repositories.ISMLTransactionMessageQueueRepository
}

func NewSMLTransactionHttpService(repo repositories.ISMLTransactionRepository, mqRepo repositories.ISMLTransactionMessageQueueRepository) *SMLTransactionHttpService {

	insSvc := &SMLTransactionHttpService{
		repo:   repo,
		mqRepo: mqRepo,
	}

	return insSvc
}

func (svc SMLTransactionHttpService) CreateSMLTransaction(shopID string, authUsername string, smlRequest models.SMLTransactionRequest) (string, error) {
	guid, err := svc.save(shopID, authUsername, smlRequest)

	if err != nil {
		return "", err
	}

	svc.mqRepo.Save(smlRequest)

	return guid, nil
}

func (svc SMLTransactionHttpService) save(shopID string, authUsername string, smlRequest models.SMLTransactionRequest) (string, error) {
	collectionName := svc.getCollectionName(smlRequest.Collection)
	findDoc, err := svc.repo.FindByDocIndentityKey(collectionName, shopID, smlRequest.KeyID, smlRequest.Body[smlRequest.KeyID])

	if err != nil {
		return "", err
	}

	_, ok := findDoc[smlRequest.KeyID]

	if ok || findDoc[smlRequest.KeyID] != nil {
		guid, err := svc.update(shopID, authUsername, findDoc, smlRequest)

		if err != nil {
			return "", err
		}

		svc.mqRepo.Save(smlRequest)
		return guid, nil
	}

	guid, err := svc.create(shopID, authUsername, smlRequest)

	if err != nil {
		return "", err
	}

	return guid, nil
}

func (svc SMLTransactionHttpService) update(shopID string, authUsername string, findDoc map[string]interface{}, smlRequest models.SMLTransactionRequest) (string, error) {
	collectionName := svc.getCollectionName(smlRequest.Collection)

	guidFixed := fmt.Sprintf("%v", findDoc["guidfixed"])
	docData := smlRequest.Body
	docData["shopid"] = findDoc["shopid"]
	docData["guidfixed"] = guidFixed
	docData["createdby"] = findDoc["createdby"]
	docData["createdat"] = findDoc["createdat"]

	docData["updatedby"] = authUsername
	docData["updatedat"] = time.Now()

	err := svc.repo.Update(collectionName, shopID, guidFixed, docData)

	if err != nil {
		return "", err
	}

	return guidFixed, nil
}

func (svc SMLTransactionHttpService) create(shopID string, authUsername string, smlRequest models.SMLTransactionRequest) (string, error) {
	collectionName := svc.getCollectionName(smlRequest.Collection)

	docData := smlRequest.Body
	newGuidFixed := utils.NewGUID()

	docData["shopid"] = shopID
	docData["guidfixed"] = newGuidFixed
	docData["createdby"] = authUsername
	docData["createdat"] = time.Now()

	_, err := svc.repo.Create(collectionName, docData)

	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc SMLTransactionHttpService) SaveInBatch(shopID string, authUsername string, dataReq models.SMLTransactionBulkRequest) ([]string, error) {

	guids := []string{}
	tempSaveSuccess := []map[string]interface{}{}
	err := svc.repo.Transaction(func() error {
		for _, smlRequest := range dataReq.Body {
			guidFixed, err := svc.CreateSMLTransaction(shopID, authUsername, models.SMLTransactionRequest{
				Collection: dataReq.Collection,
				KeyID:      dataReq.KeyID,
				Body:       smlRequest,
			})

			if err != nil {
				return err
			}

			guids = append(guids, guidFixed)
			tempSaveSuccess = append(tempSaveSuccess, smlRequest)
		}

		return nil
	})

	if err != nil {
		return []string{}, err
	}

	err = svc.mqRepo.BulkSave(models.SMLTransactionBulkRequest{
		Collection: dataReq.Collection,
		KeyID:      dataReq.KeyID,
		Body:       tempSaveSuccess,
	})

	if err != nil {
		return guids, err
	}

	return guids, nil
}

func (svc SMLTransactionHttpService) DeleteSMLTransaction(shopID string, authUsername string, smlKeyRequest models.SMLTransactionKeyRequest) ([]string, error) {
	collectionName := svc.getCollectionName(smlKeyRequest.Collection)

	filters := map[string]interface{}{
		smlKeyRequest.KeyID: bson.M{"$in": smlKeyRequest.DeleteKeys},
	}
	err := svc.repo.Delete(collectionName, shopID, authUsername, filters)

	if err != nil {
		return []string{}, err
	}

	err = svc.mqRepo.Delete(smlKeyRequest)

	if err != nil {
		return []string{}, err
	}

	return smlKeyRequest.DeleteKeys, nil
}

func (svc SMLTransactionHttpService) getCollectionName(collectionName string) string {
	return "smlx" + strings.Trim(collectionName, " ")
}
