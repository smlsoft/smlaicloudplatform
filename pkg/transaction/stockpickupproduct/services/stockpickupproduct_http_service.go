package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/transaction/stockpickupproduct/models"
	"smlcloudplatform/pkg/transaction/stockpickupproduct/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IStockPickupProductHttpService interface {
	CreateStockPickupProduct(shopID string, authUsername string, doc models.StockPickupProduct) (string, error)
	UpdateStockPickupProduct(shopID string, guid string, authUsername string, doc models.StockPickupProduct) error
	DeleteStockPickupProduct(shopID string, guid string, authUsername string) error
	DeleteStockPickupProductByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoStockPickupProduct(shopID string, guid string) (models.StockPickupProductInfo, error)
	InfoStockPickupProductByCode(shopID string, code string) (models.StockPickupProductInfo, error)
	SearchStockPickupProduct(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.StockPickupProductInfo, mongopagination.PaginationData, error)
	SearchStockPickupProductStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.StockPickupProductInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.StockPickupProduct) (common.BulkImport, error)

	GetModuleName() string
}

type StockPickupProductHttpService struct {
	repo repositories.IStockPickupProductRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.StockPickupProductActivity, models.StockPickupProductDeleteActivity]
}

func NewStockPickupProductHttpService(repo repositories.IStockPickupProductRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *StockPickupProductHttpService {

	insSvc := &StockPickupProductHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.StockPickupProductActivity, models.StockPickupProductDeleteActivity](repo)

	return insSvc
}

func (svc StockPickupProductHttpService) CreateStockPickupProduct(shopID string, authUsername string, doc models.StockPickupProduct) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "docno", doc.Docno)

	if err != nil {
		return "", err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", errors.New("Docno is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.StockPickupProductDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.StockPickupProduct = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc StockPickupProductHttpService) UpdateStockPickupProduct(shopID string, guid string, authUsername string, doc models.StockPickupProduct) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	findDoc.StockPickupProduct = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc StockPickupProductHttpService) DeleteStockPickupProduct(shopID string, guid string, authUsername string) error {

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

func (svc StockPickupProductHttpService) DeleteStockPickupProductByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc StockPickupProductHttpService) InfoStockPickupProduct(shopID string, guid string) (models.StockPickupProductInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.StockPickupProductInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.StockPickupProductInfo{}, errors.New("document not found")
	}

	return findDoc.StockPickupProductInfo, nil
}

func (svc StockPickupProductHttpService) InfoStockPickupProductByCode(shopID string, code string) (models.StockPickupProductInfo, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "docno", code)

	if err != nil {
		return models.StockPickupProductInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.StockPickupProductInfo{}, errors.New("document not found")
	}

	return findDoc.StockPickupProductInfo, nil
}

func (svc StockPickupProductHttpService) SearchStockPickupProduct(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.StockPickupProductInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"docno",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.StockPickupProductInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc StockPickupProductHttpService) SearchStockPickupProductStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.StockPickupProductInfo, int, error) {
	searchInFields := []string{
		"docno",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.StockPickupProductInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc StockPickupProductHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.StockPickupProduct) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.StockPickupProduct](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.StockPickupProduct, models.StockPickupProductDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.StockPickupProduct) models.StockPickupProductDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.StockPickupProductDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.StockPickupProduct = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.StockPickupProduct, models.StockPickupProductDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.StockPickupProductDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "docno", guid)
		},
		func(doc models.StockPickupProductDoc) bool {
			return doc.Docno != ""
		},
		func(shopID string, authUsername string, data models.StockPickupProduct, doc models.StockPickupProductDoc) error {

			doc.StockPickupProduct = data
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

func (svc StockPickupProductHttpService) getDocIDKey(doc models.StockPickupProduct) string {
	return doc.Docno
}

func (svc StockPickupProductHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc StockPickupProductHttpService) GetModuleName() string {
	return "stockPickupProduct"
}
