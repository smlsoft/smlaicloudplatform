package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/transaction/stockreceiveproduct/models"
	"smlcloudplatform/pkg/transaction/stockreceiveproduct/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IStockReceiveProductHttpService interface {
	CreateStockReceiveProduct(shopID string, authUsername string, doc models.StockReceiveProduct) (string, error)
	UpdateStockReceiveProduct(shopID string, guid string, authUsername string, doc models.StockReceiveProduct) error
	DeleteStockReceiveProduct(shopID string, guid string, authUsername string) error
	DeleteStockReceiveProductByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoStockReceiveProduct(shopID string, guid string) (models.StockReceiveProductInfo, error)
	InfoStockReceiveProductByCode(shopID string, code string) (models.StockReceiveProductInfo, error)
	SearchStockReceiveProduct(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.StockReceiveProductInfo, mongopagination.PaginationData, error)
	SearchStockReceiveProductStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.StockReceiveProductInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.StockReceiveProduct) (common.BulkImport, error)

	GetModuleName() string
}

type StockReceiveProductHttpService struct {
	repo repositories.IStockReceiveProductRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.StockReceiveProductActivity, models.StockReceiveProductDeleteActivity]
}

func NewStockReceiveProductHttpService(repo repositories.IStockReceiveProductRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *StockReceiveProductHttpService {

	insSvc := &StockReceiveProductHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.StockReceiveProductActivity, models.StockReceiveProductDeleteActivity](repo)

	return insSvc
}

func (svc StockReceiveProductHttpService) CreateStockReceiveProduct(shopID string, authUsername string, doc models.StockReceiveProduct) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "docno", doc.Docno)

	if err != nil {
		return "", err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", errors.New("Docno is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.StockReceiveProductDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.StockReceiveProduct = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc StockReceiveProductHttpService) UpdateStockReceiveProduct(shopID string, guid string, authUsername string, doc models.StockReceiveProduct) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	findDoc.StockReceiveProduct = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc StockReceiveProductHttpService) DeleteStockReceiveProduct(shopID string, guid string, authUsername string) error {

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

func (svc StockReceiveProductHttpService) DeleteStockReceiveProductByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc StockReceiveProductHttpService) InfoStockReceiveProduct(shopID string, guid string) (models.StockReceiveProductInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.StockReceiveProductInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.StockReceiveProductInfo{}, errors.New("document not found")
	}

	return findDoc.StockReceiveProductInfo, nil
}

func (svc StockReceiveProductHttpService) InfoStockReceiveProductByCode(shopID string, code string) (models.StockReceiveProductInfo, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "docno", code)

	if err != nil {
		return models.StockReceiveProductInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.StockReceiveProductInfo{}, errors.New("document not found")
	}

	return findDoc.StockReceiveProductInfo, nil
}

func (svc StockReceiveProductHttpService) SearchStockReceiveProduct(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.StockReceiveProductInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"docno",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.StockReceiveProductInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc StockReceiveProductHttpService) SearchStockReceiveProductStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.StockReceiveProductInfo, int, error) {
	searchInFields := []string{
		"docno",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.StockReceiveProductInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc StockReceiveProductHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.StockReceiveProduct) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.StockReceiveProduct](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.StockReceiveProduct, models.StockReceiveProductDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.StockReceiveProduct) models.StockReceiveProductDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.StockReceiveProductDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.StockReceiveProduct = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.StockReceiveProduct, models.StockReceiveProductDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.StockReceiveProductDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "docno", guid)
		},
		func(doc models.StockReceiveProductDoc) bool {
			return doc.Docno != ""
		},
		func(shopID string, authUsername string, data models.StockReceiveProduct, doc models.StockReceiveProductDoc) error {

			doc.StockReceiveProduct = data
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

func (svc StockReceiveProductHttpService) getDocIDKey(doc models.StockReceiveProduct) string {
	return doc.Docno
}

func (svc StockReceiveProductHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc StockReceiveProductHttpService) GetModuleName() string {
	return "stockReceiveProduct"
}
