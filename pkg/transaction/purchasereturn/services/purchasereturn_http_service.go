package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/transaction/purchasereturn/models"
	"smlcloudplatform/pkg/transaction/purchasereturn/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IPurchaseReturnHttpService interface {
	CreatePurchaseReturn(shopID string, authUsername string, doc models.PurchaseReturn) (string, error)
	UpdatePurchaseReturn(shopID string, guid string, authUsername string, doc models.PurchaseReturn) error
	DeletePurchaseReturn(shopID string, guid string, authUsername string) error
	DeletePurchaseReturnByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoPurchaseReturn(shopID string, guid string) (models.PurchaseReturnInfo, error)
	InfoPurchaseReturnByCode(shopID string, code string) (models.PurchaseReturnInfo, error)
	SearchPurchaseReturn(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.PurchaseReturnInfo, mongopagination.PaginationData, error)
	SearchPurchaseReturnStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.PurchaseReturnInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.PurchaseReturn) (common.BulkImport, error)

	GetModuleName() string
}

type PurchaseReturnHttpService struct {
	repo repositories.IPurchaseReturnRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.PurchaseReturnActivity, models.PurchaseReturnDeleteActivity]
}

func NewPurchaseReturnHttpService(repo repositories.IPurchaseReturnRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *PurchaseReturnHttpService {

	insSvc := &PurchaseReturnHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.PurchaseReturnActivity, models.PurchaseReturnDeleteActivity](repo)

	return insSvc
}

func (svc PurchaseReturnHttpService) CreatePurchaseReturn(shopID string, authUsername string, doc models.PurchaseReturn) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "docno", doc.Docno)

	if err != nil {
		return "", err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", errors.New("Docno is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.PurchaseReturnDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.PurchaseReturn = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc PurchaseReturnHttpService) UpdatePurchaseReturn(shopID string, guid string, authUsername string, doc models.PurchaseReturn) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	findDoc.PurchaseReturn = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc PurchaseReturnHttpService) DeletePurchaseReturn(shopID string, guid string, authUsername string) error {

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

func (svc PurchaseReturnHttpService) DeletePurchaseReturnByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc PurchaseReturnHttpService) InfoPurchaseReturn(shopID string, guid string) (models.PurchaseReturnInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.PurchaseReturnInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.PurchaseReturnInfo{}, errors.New("document not found")
	}

	return findDoc.PurchaseReturnInfo, nil
}

func (svc PurchaseReturnHttpService) InfoPurchaseReturnByCode(shopID string, code string) (models.PurchaseReturnInfo, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "docno", code)

	if err != nil {
		return models.PurchaseReturnInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.PurchaseReturnInfo{}, errors.New("document not found")
	}

	return findDoc.PurchaseReturnInfo, nil
}

func (svc PurchaseReturnHttpService) SearchPurchaseReturn(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.PurchaseReturnInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"docno",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.PurchaseReturnInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc PurchaseReturnHttpService) SearchPurchaseReturnStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.PurchaseReturnInfo, int, error) {
	searchInFields := []string{
		"docno",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, filters, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.PurchaseReturnInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc PurchaseReturnHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.PurchaseReturn) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.PurchaseReturn](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.PurchaseReturn, models.PurchaseReturnDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.PurchaseReturn) models.PurchaseReturnDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.PurchaseReturnDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.PurchaseReturn = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.PurchaseReturn, models.PurchaseReturnDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.PurchaseReturnDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "docno", guid)
		},
		func(doc models.PurchaseReturnDoc) bool {
			return doc.Docno != ""
		},
		func(shopID string, authUsername string, data models.PurchaseReturn, doc models.PurchaseReturnDoc) error {

			doc.PurchaseReturn = data
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

func (svc PurchaseReturnHttpService) getDocIDKey(doc models.PurchaseReturn) string {
	return doc.Docno
}

func (svc PurchaseReturnHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc PurchaseReturnHttpService) GetModuleName() string {
	return "purchaseReturn"
}
