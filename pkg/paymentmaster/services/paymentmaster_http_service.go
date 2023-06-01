package services

import (
	"errors"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/paymentmaster/models"
	"smlcloudplatform/pkg/paymentmaster/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IPaymentMasterHttpService interface {
	CreatePaymentMaster(shopID string, authUsername string, doc models.PaymentMaster) (string, error)
	UpdatePaymentMaster(guid string, shopID string, authUsername string, doc models.PaymentMaster) error
	DeletePaymentMaster(guid string, shopID string, authUsername string) error
	InfoPaymentMaster(guid string, shopID string) (models.PaymentMasterInfo, error)
	SearchPaymentMaster(shopID string, q string) ([]models.PaymentMasterInfo, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.PaymentMaster) (common.BulkImport, error)
}

type PaymentMasterHttpService struct {
	repo repositories.PaymentMasterRepository
}

func NewPaymentMasterHttpService(repo repositories.PaymentMasterRepository) PaymentMasterHttpService {

	return PaymentMasterHttpService{
		repo: repo,
	}
}

func (svc PaymentMasterHttpService) CreatePaymentMaster(shopID string, authUsername string, doc models.PaymentMaster) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "paymentcode", doc.PaymentCode)

	if err != nil {
		return "", err
	}

	if findDoc.PaymentCode != "" {
		return "", errors.New("PaymentCode is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.PaymentMasterDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.PaymentMaster = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc PaymentMasterHttpService) UpdatePaymentMaster(guid string, shopID string, authUsername string, doc models.PaymentMaster) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.PaymentMaster = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc PaymentMasterHttpService) DeletePaymentMaster(guid string, shopID string, authUsername string) error {

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

func (svc PaymentMasterHttpService) InfoPaymentMaster(guid string, shopID string) (models.PaymentMasterInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.PaymentMasterInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.PaymentMasterInfo{}, errors.New("document not found")
	}

	return findDoc.PaymentMasterInfo, nil

}

func (svc PaymentMasterHttpService) SearchPaymentMaster(shopID string, q string) ([]models.PaymentMasterInfo, error) {
	searchInFields := []string{
		"guidfixed",
		"paymentcode",
	}

	docList, err := svc.repo.Find(shopID, searchInFields, q)

	if err != nil {
		return []models.PaymentMasterInfo{}, err
	}

	return docList, nil
}

func (svc PaymentMasterHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.PaymentMaster) (common.BulkImport, error) {

	createDataList := []models.PaymentMasterDoc{}
	duplicateDataList := []models.PaymentMaster{}

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.PaymentMaster](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.PaymentCode)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "paymentcode", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.PaymentCode)
	}

	duplicateDataList, createDataList = importdata.PreparePayloadData[models.PaymentMaster, models.PaymentMasterDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.PaymentMaster) models.PaymentMasterDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.PaymentMasterDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.PaymentMaster = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.PaymentMaster, models.PaymentMasterDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.PaymentMasterDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "paymentcode", guid)
		},
		func(doc models.PaymentMasterDoc) bool {
			if doc.PaymentCode != "" {
				return true
			}
			return false
		},
		func(shopID string, authUsername string, data models.PaymentMaster, doc models.PaymentMasterDoc) error {

			doc.PaymentMaster = data
			doc.UpdatedBy = authUsername
			doc.UpdatedAt = time.Now()

			err = svc.repo.Update(shopID, doc.GuidFixed, doc)
			if err != nil {
				return nil
			}
			return nil
		},
	)

	if len(createDataList) > 0 {
		err = svc.repo.CreateInBatch(createDataList)

		if err != nil {
			return common.BulkImport{}, err
		}

	}

	createDataKey := []string{}

	for _, doc := range createDataList {
		createDataKey = append(createDataKey, doc.PaymentCode)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.PaymentCode)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {
		updateDataKey = append(updateDataKey, doc.PaymentCode)
	}

	updateFailDataKey := []string{}
	for _, doc := range updateFailDataList {
		updateFailDataKey = append(updateFailDataKey, svc.getDocIDKey(doc))
	}

	return common.BulkImport{
		Created:          createDataKey,
		Updated:          updateDataKey,
		UpdateFailed:     updateFailDataKey,
		PayloadDuplicate: payloadDuplicateDataKey,
	}, nil
}

func (svc PaymentMasterHttpService) getDocIDKey(doc models.PaymentMaster) string {
	return doc.PaymentCode
}
