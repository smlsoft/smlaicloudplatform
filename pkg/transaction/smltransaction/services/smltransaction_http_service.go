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
	repo         repositories.ISMLTransactionRepository
	mqRepo       repositories.ISMLTransactionMessageQueueRepository
	indexCreated map[string]struct{}
}

func NewSMLTransactionHttpService(repo repositories.ISMLTransactionRepository, mqRepo repositories.ISMLTransactionMessageQueueRepository) *SMLTransactionHttpService {

	insSvc := &SMLTransactionHttpService{
		repo:         repo,
		mqRepo:       mqRepo,
		indexCreated: map[string]struct{}{},
	}

	return insSvc
}

func (svc SMLTransactionHttpService) CreateSMLTransaction(shopID string, authUsername string, smlRequest models.SMLTransactionRequest) (string, error) {
	// guid, err := svc.save(shopID, authUsername, smlRequest)

	// if err != nil {
	// 	return "", err
	// }

	// svc.mqRepo.Save(smlRequest)

	// return guid, nil

	collectionName := svc.getCollectionName(smlRequest.Collection)

	_, collectionIndexExists := svc.indexCreated[collectionName]
	if !collectionIndexExists {
		_, err := svc.repo.CreateIndex(collectionName, smlRequest.KeyID)
		if err == nil {
			svc.indexCreated[collectionName] = struct{}{}
		}
	}

	_, ok := smlRequest.Body[smlRequest.KeyID]

	if !ok || smlRequest.Body[smlRequest.KeyID] == nil {
		return "", fmt.Errorf("KeyID is not found in body")
	}

	deleteFilter := bson.M{smlRequest.KeyID: smlRequest.Body[smlRequest.KeyID]}
	err := svc.repo.Delete(collectionName, shopID, authUsername, deleteFilter)

	if err != nil {
		return "", err
	}

	tempData := svc.createBody(shopID, authUsername, smlRequest.Body)
	guid, err := svc.repo.Create(collectionName, tempData)

	if err != nil {
		return "", err
	}

	err = svc.mqRepo.Save(smlRequest)

	if err != nil {
		return "", err
	}

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

func (svc SMLTransactionHttpService) createBody(shopID string, authUsername string, bodyRequest map[string]interface{}) map[string]interface{} {

	docData := bodyRequest
	newGuidFixed := utils.NewGUID()

	docData["shopid"] = shopID
	docData["guidfixed"] = newGuidFixed
	docData["createdby"] = authUsername
	docData["createdat"] = time.Now()

	return docData
}

func (svc SMLTransactionHttpService) SaveInBatch(shopID string, authUsername string, dataReq models.SMLTransactionBulkRequest) ([]string, error) {
	collectionName := svc.getCollectionName(dataReq.Collection)
	_, collectionIndexExists := svc.indexCreated[collectionName]
	if !collectionIndexExists {
		_, err := svc.repo.CreateIndex(collectionName, dataReq.KeyID)
		if err == nil {
			svc.indexCreated[collectionName] = struct{}{}
		}
	}

	identityKeys := []string{}
	tempData := []map[string]interface{}{}
	for _, smlRequest := range dataReq.Body {
		_, ok := smlRequest[dataReq.KeyID]

		if !ok || smlRequest[dataReq.KeyID] == nil {
			continue
		}

		identityKeys = append(identityKeys, smlRequest[dataReq.KeyID].(string))
		createData := svc.createBody(shopID, authUsername, smlRequest)
		tempData = append(tempData, createData)
	}

	if len(identityKeys) == 0 {
		return []string{}, fmt.Errorf("body request is empty")
	}

	err := svc.repo.Transaction(func() error {

		filters := map[string]interface{}{
			dataReq.KeyID: bson.M{"$in": identityKeys},
		}
		err := svc.repo.Delete(collectionName, shopID, authUsername, filters)

		if err != nil {
			return err
		}

		err = svc.repo.CreateInBatch(collectionName, tempData)

		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return []string{}, err
	}

	err = svc.mqRepo.Delete(models.SMLTransactionKeyRequest{
		Collection: dataReq.Collection,
		KeyID:      dataReq.KeyID,
		DeleteKeys: identityKeys,
	})

	if err != nil {
		return []string{}, err
	}

	err = svc.mqRepo.BulkSave(models.SMLTransactionBulkRequest{
		Collection: dataReq.Collection,
		KeyID:      dataReq.KeyID,
		Body:       dataReq.Body,
	})

	if err != nil {
		return identityKeys, err
	}

	return identityKeys, nil
}

func (svc SMLTransactionHttpService) SaveInBatchOld(shopID string, authUsername string, dataReq models.SMLTransactionBulkRequest) ([]string, error) {

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
