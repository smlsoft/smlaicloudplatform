package services

import (
	"errors"
	"fmt"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/payment/qrpayment/models"
	"smlcloudplatform/pkg/payment/qrpayment/repositories"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IQrPaymentHttpService interface {
	CreateQrPayment(shopID string, authUsername string, doc models.QrPayment) (string, error)
	UpdateQrPayment(shopID string, guid string, authUsername string, doc models.QrPayment) error
	DeleteQrPayment(shopID string, guid string, authUsername string) error
	DeleteQrPaymentByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoQrPayment(shopID string, guid string) (models.QrPaymentInfo, error)
	SearchQrPayment(shopID string, q string, page int, limit int, sort map[string]int) ([]models.QrPaymentInfo, mongopagination.PaginationData, error)
	SearchQrPaymentStep(shopID string, langCode string, q string, skip int, limit int, sort map[string]int) ([]models.QrPaymentInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.QrPayment) (common.BulkImport, error)

	GetModuleName() string
}

type QrPaymentHttpService struct {
	repo repositories.IQrPaymentRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.QrPaymentActivity, models.QrPaymentDeleteActivity]
}

func NewQrPaymentHttpService(repo repositories.IQrPaymentRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *QrPaymentHttpService {

	insSvc := &QrPaymentHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.QrPaymentActivity, models.QrPaymentDeleteActivity](repo)

	return insSvc
}

func (svc QrPaymentHttpService) CreateQrPayment(shopID string, authUsername string, doc models.QrPayment) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "paymentcode", doc.PaymentCode)

	if err != nil {
		return "", err
	}

	if findDoc.PaymentCode != "" {
		return "", errors.New("PaymentCode is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.QrPaymentDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.QrPayment = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc QrPaymentHttpService) UpdateQrPayment(shopID string, guid string, authUsername string, doc models.QrPayment) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.QrPayment = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc QrPaymentHttpService) DeleteQrPayment(shopID string, guid string, authUsername string) error {

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

	svc.saveMasterSync(shopID)

	return nil
}

func (svc QrPaymentHttpService) DeleteQrPaymentByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc QrPaymentHttpService) InfoQrPayment(shopID string, guid string) (models.QrPaymentInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.QrPaymentInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.QrPaymentInfo{}, errors.New("document not found")
	}

	return findDoc.QrPaymentInfo, nil

}

func (svc QrPaymentHttpService) SearchQrPayment(shopID string, q string, page int, limit int, sort map[string]int) ([]models.QrPaymentInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"guidfixed",
		"paymentcode",
	}

	docList, pagination, err := svc.repo.FindPageSort(shopID, searchCols, q, page, limit, sort)

	if err != nil {
		return []models.QrPaymentInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc QrPaymentHttpService) SearchQrPaymentStep(shopID string, langCode string, q string, skip int, limit int, sort map[string]int) ([]models.QrPaymentInfo, int, error) {
	searchCols := []string{
		"guidfixed",
		"paymentcode",
	}

	projectQuery := map[string]interface{}{
		"guidfixed":    1,
		"paymentcode":  1,
		"countrycode":  1,
		"paymentlogo":  1,
		"paymenttype":  1,
		"feerate":      1,
		"wallettype":   1,
		"bookbankcode": 1,
		"bankcode":     1,
	}

	if langCode != "" {
		projectQuery["names"] = bson.M{"$elemMatch": bson.M{"code": langCode}}
	} else {
		projectQuery["names"] = 1
	}

	docList, total, err := svc.repo.FindLimit(shopID, map[string]interface{}{}, searchCols, q, skip, limit, sort, projectQuery)

	if err != nil {
		return []models.QrPaymentInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc QrPaymentHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.QrPayment) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.QrPayment](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.QrPayment, models.QrPaymentDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.QrPayment) models.QrPaymentDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.QrPaymentDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.QrPayment = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.QrPayment, models.QrPaymentDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.QrPaymentDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "paymentcode", guid)
		},
		func(doc models.QrPaymentDoc) bool {
			return doc.PaymentCode != ""
		},
		func(shopID string, authUsername string, data models.QrPayment, doc models.QrPaymentDoc) error {

			doc.QrPayment = data
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

	svc.saveMasterSync(shopID)

	return common.BulkImport{
		Created:          createDataKey,
		Updated:          updateDataKey,
		UpdateFailed:     updateFailDataKey,
		PayloadDuplicate: payloadDuplicateDataKey,
	}, nil
}

func (svc QrPaymentHttpService) getDocIDKey(doc models.QrPayment) string {
	return doc.PaymentCode
}

func (svc QrPaymentHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc QrPaymentHttpService) GetModuleName() string {
	return "qrpayment"
}
