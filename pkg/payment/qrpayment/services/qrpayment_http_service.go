package services

import (
	"context"
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
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
	SearchQrPayment(shopID string, pageable micromodels.Pageable) ([]models.QrPaymentInfo, mongopagination.PaginationData, error)
	SearchQrPaymentStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.QrPaymentInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.QrPayment) (common.BulkImport, error)

	GetModuleName() string
}

type QrPaymentHttpService struct {
	repo repositories.IQrPaymentRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.QrPaymentActivity, models.QrPaymentDeleteActivity]
	contextTimeout time.Duration
}

func NewQrPaymentHttpService(repo repositories.IQrPaymentRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *QrPaymentHttpService {

	contextTimeout := time.Duration(15) * time.Second

	insSvc := &QrPaymentHttpService{
		repo:           repo,
		syncCacheRepo:  syncCacheRepo,
		contextTimeout: contextTimeout,
	}

	insSvc.ActivityService = services.NewActivityService[models.QrPaymentActivity, models.QrPaymentDeleteActivity](repo)

	return insSvc
}

func (svc QrPaymentHttpService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc QrPaymentHttpService) CreateQrPayment(shopID string, authUsername string, doc models.QrPayment) (string, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "paymentcode", doc.PaymentCode)

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

	_, err = svc.repo.Create(ctx, docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc QrPaymentHttpService) UpdateQrPayment(shopID string, guid string, authUsername string, doc models.QrPayment) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.QrPayment = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(ctx, shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc QrPaymentHttpService) DeleteQrPayment(shopID string, guid string, authUsername string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	err = svc.repo.DeleteByGuidfixed(ctx, shopID, guid, authUsername)
	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc QrPaymentHttpService) DeleteQrPaymentByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(ctx, shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc QrPaymentHttpService) InfoQrPayment(shopID string, guid string) (models.QrPaymentInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return models.QrPaymentInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.QrPaymentInfo{}, errors.New("document not found")
	}

	return findDoc.QrPaymentInfo, nil

}

func (svc QrPaymentHttpService) SearchQrPayment(shopID string, pageable micromodels.Pageable) ([]models.QrPaymentInfo, mongopagination.PaginationData, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"paymentcode",
		"names.name",
	}

	docList, pagination, err := svc.repo.FindPage(ctx, shopID, searchInFields, pageable)

	if err != nil {
		return []models.QrPaymentInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc QrPaymentHttpService) SearchQrPaymentStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.QrPaymentInfo, int, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"paymentcode",
		"names.name",
	}

	selectFields := map[string]interface{}{
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
		selectFields["names"] = bson.M{"$elemMatch": bson.M{"code": langCode}}
	} else {
		selectFields["names"] = 1
	}

	docList, total, err := svc.repo.FindStep(ctx, shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.QrPaymentInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc QrPaymentHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.QrPayment) (common.BulkImport, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.QrPayment](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.PaymentCode)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(ctx, shopID, "paymentcode", itemCodeGuidList)

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
			return svc.repo.FindByDocIndentityGuid(ctx, shopID, "paymentcode", guid)
		},
		func(doc models.QrPaymentDoc) bool {
			return doc.PaymentCode != ""
		},
		func(shopID string, authUsername string, data models.QrPayment, doc models.QrPaymentDoc) error {

			doc.QrPayment = data
			doc.UpdatedBy = authUsername
			doc.UpdatedAt = time.Now()

			err = svc.repo.Update(ctx, shopID, doc.GuidFixed, doc)
			if err != nil {
				return nil
			}
			return nil
		},
	)

	if len(createDataList) > 0 {
		err = svc.repo.CreateInBatch(ctx, createDataList)

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
