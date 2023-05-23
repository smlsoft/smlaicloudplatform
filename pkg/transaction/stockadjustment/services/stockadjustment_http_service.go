package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/transaction/stockadjustment/models"
	"smlcloudplatform/pkg/transaction/stockadjustment/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IStockAdjustmentHttpService interface {
	CreateStockAdjustment(shopID string, authUsername string, doc models.StockAdjustment) (string, error)
	UpdateStockAdjustment(shopID string, guid string, authUsername string, doc models.StockAdjustment) error
	DeleteStockAdjustment(shopID string, guid string, authUsername string) error
	DeleteStockAdjustmentByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoStockAdjustment(shopID string, guid string) (models.StockAdjustmentInfo, error)
	InfoStockAdjustmentByCode(shopID string, code string) (models.StockAdjustmentInfo, error)
	SearchStockAdjustment(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.StockAdjustmentInfo, mongopagination.PaginationData, error)
	SearchStockAdjustmentStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.StockAdjustmentInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.StockAdjustment) (common.BulkImport, error)

	GetModuleName() string
}

type StockAdjustmentHttpService struct {
	repo repositories.IStockAdjustmentRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.StockAdjustmentActivity, models.StockAdjustmentDeleteActivity]
}

func NewStockAdjustmentHttpService(repo repositories.IStockAdjustmentRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *StockAdjustmentHttpService {

	insSvc := &StockAdjustmentHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.StockAdjustmentActivity, models.StockAdjustmentDeleteActivity](repo)

	return insSvc
}

func (svc StockAdjustmentHttpService) CreateStockAdjustment(shopID string, authUsername string, doc models.StockAdjustment) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "docno", doc.Docno)

	if err != nil {
		return "", err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", errors.New("Docno is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.StockAdjustmentDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.StockAdjustment = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc StockAdjustmentHttpService) UpdateStockAdjustment(shopID string, guid string, authUsername string, doc models.StockAdjustment) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	findDoc.StockAdjustment = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc StockAdjustmentHttpService) DeleteStockAdjustment(shopID string, guid string, authUsername string) error {

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

func (svc StockAdjustmentHttpService) DeleteStockAdjustmentByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc StockAdjustmentHttpService) InfoStockAdjustment(shopID string, guid string) (models.StockAdjustmentInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.StockAdjustmentInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.StockAdjustmentInfo{}, errors.New("document not found")
	}

	return findDoc.StockAdjustmentInfo, nil
}

func (svc StockAdjustmentHttpService) InfoStockAdjustmentByCode(shopID string, code string) (models.StockAdjustmentInfo, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "docno", code)

	if err != nil {
		return models.StockAdjustmentInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.StockAdjustmentInfo{}, errors.New("document not found")
	}

	return findDoc.StockAdjustmentInfo, nil
}

func (svc StockAdjustmentHttpService) SearchStockAdjustment(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.StockAdjustmentInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"docno",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.StockAdjustmentInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc StockAdjustmentHttpService) SearchStockAdjustmentStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.StockAdjustmentInfo, int, error) {
	searchInFields := []string{
		"docno",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, filters, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.StockAdjustmentInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc StockAdjustmentHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.StockAdjustment) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.StockAdjustment](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.StockAdjustment, models.StockAdjustmentDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.StockAdjustment) models.StockAdjustmentDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.StockAdjustmentDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.StockAdjustment = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.StockAdjustment, models.StockAdjustmentDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.StockAdjustmentDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "docno", guid)
		},
		func(doc models.StockAdjustmentDoc) bool {
			return doc.Docno != ""
		},
		func(shopID string, authUsername string, data models.StockAdjustment, doc models.StockAdjustmentDoc) error {

			doc.StockAdjustment = data
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

func (svc StockAdjustmentHttpService) getDocIDKey(doc models.StockAdjustment) string {
	return doc.Docno
}

func (svc StockAdjustmentHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc StockAdjustmentHttpService) GetModuleName() string {
	return "stockAdjustment"
}
