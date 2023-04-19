package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/transaction/purchase/models"
	"smlcloudplatform/pkg/transaction/purchase/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IPurchaseHttpService interface {
	CreatePurchase(shopID string, authUsername string, doc models.Purchase) (string, error)
	UpdatePurchase(shopID string, guid string, authUsername string, doc models.Purchase) error
	DeletePurchase(shopID string, guid string, authUsername string) error
	DeletePurchaseByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoPurchase(shopID string, guid string) (models.PurchaseInfo, error)
	InfoPurchaseByCode(shopID string, code string) (models.PurchaseInfo, error)
	SearchPurchase(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.PurchaseInfo, mongopagination.PaginationData, error)
	SearchPurchaseStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.PurchaseInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Purchase) (common.BulkImport, error)

	GetModuleName() string
}

type PurchaseHttpService struct {
	repo repositories.IPurchaseRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.PurchaseActivity, models.PurchaseDeleteActivity]
}

func NewPurchaseHttpService(repo repositories.IPurchaseRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *PurchaseHttpService {

	insSvc := &PurchaseHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.PurchaseActivity, models.PurchaseDeleteActivity](repo)

	return insSvc
}

func (svc PurchaseHttpService) CreatePurchase(shopID string, authUsername string, doc models.Purchase) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "docno", doc.Docno)

	if err != nil {
		return "", err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", errors.New("Docno is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.PurchaseDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Purchase = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc PurchaseHttpService) UpdatePurchase(shopID string, guid string, authUsername string, doc models.Purchase) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	findDoc.Purchase = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc PurchaseHttpService) DeletePurchase(shopID string, guid string, authUsername string) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	err = svc.repo.DeleteByGuidfixed(shopID, guid, authUsername)
	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc PurchaseHttpService) DeletePurchaseByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc PurchaseHttpService) InfoPurchase(shopID string, guid string) (models.PurchaseInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.PurchaseInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.PurchaseInfo{}, errors.New("document not found")
	}

	return findDoc.PurchaseInfo, nil
}

func (svc PurchaseHttpService) InfoPurchaseByCode(shopID string, code string) (models.PurchaseInfo, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "docno", code)

	if err != nil {
		return models.PurchaseInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.PurchaseInfo{}, errors.New("document not found")
	}

	return findDoc.PurchaseInfo, nil
}

func (svc PurchaseHttpService) SearchPurchase(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.PurchaseInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"docno",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.PurchaseInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc PurchaseHttpService) SearchPurchaseStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.PurchaseInfo, int, error) {
	searchInFields := []string{
		"docno",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.PurchaseInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc PurchaseHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.Purchase) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Purchase](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.Docno)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "docno", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.Docno)
	}

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Purchase, models.PurchaseDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Purchase) models.PurchaseDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.PurchaseDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Purchase = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Purchase, models.PurchaseDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.PurchaseDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "docno", guid)
		},
		func(doc models.PurchaseDoc) bool {
			return doc.Docno != ""
		},
		func(shopID string, authUsername string, data models.Purchase, doc models.PurchaseDoc) error {

			doc.Purchase = data
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
		createDataKey = append(createDataKey, doc.Docno)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.Docno)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {

		updateDataKey = append(updateDataKey, doc.Docno)
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

func (svc PurchaseHttpService) getDocIDKey(doc models.Purchase) string {
	return doc.Docno
}

func (svc PurchaseHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc PurchaseHttpService) GetModuleName() string {
	return "purchase"
}
