package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/transaction/stockreturnproduct/models"
	"smlcloudplatform/pkg/transaction/stockreturnproduct/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IStockReturnProductHttpService interface {
	CreateStockReturnProduct(shopID string, authUsername string, doc models.StockReturnProduct) (string, error)
	UpdateStockReturnProduct(shopID string, guid string, authUsername string, doc models.StockReturnProduct) error
	DeleteStockReturnProduct(shopID string, guid string, authUsername string) error
	DeleteStockReturnProductByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoStockReturnProduct(shopID string, guid string) (models.StockReturnProductInfo, error)
	InfoStockReturnProductByCode(shopID string, code string) (models.StockReturnProductInfo, error)
	SearchStockReturnProduct(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.StockReturnProductInfo, mongopagination.PaginationData, error)
	SearchStockReturnProductStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.StockReturnProductInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.StockReturnProduct) (common.BulkImport, error)

	GetModuleName() string
}

type StockReturnProductHttpService struct {
	repo repositories.IStockReturnProductRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.StockReturnProductActivity, models.StockReturnProductDeleteActivity]
}

func NewStockReturnProductHttpService(repo repositories.IStockReturnProductRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *StockReturnProductHttpService {

	insSvc := &StockReturnProductHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.StockReturnProductActivity, models.StockReturnProductDeleteActivity](repo)

	return insSvc
}

func (svc StockReturnProductHttpService) CreateStockReturnProduct(shopID string, authUsername string, doc models.StockReturnProduct) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "docno", doc.Docno)

	if err != nil {
		return "", err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", errors.New("Docno is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.StockReturnProductDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.StockReturnProduct = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc StockReturnProductHttpService) UpdateStockReturnProduct(shopID string, guid string, authUsername string, doc models.StockReturnProduct) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	findDoc.StockReturnProduct = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc StockReturnProductHttpService) DeleteStockReturnProduct(shopID string, guid string, authUsername string) error {

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

func (svc StockReturnProductHttpService) DeleteStockReturnProductByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc StockReturnProductHttpService) InfoStockReturnProduct(shopID string, guid string) (models.StockReturnProductInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.StockReturnProductInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.StockReturnProductInfo{}, errors.New("document not found")
	}

	return findDoc.StockReturnProductInfo, nil
}

func (svc StockReturnProductHttpService) InfoStockReturnProductByCode(shopID string, code string) (models.StockReturnProductInfo, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "docno", code)

	if err != nil {
		return models.StockReturnProductInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.StockReturnProductInfo{}, errors.New("document not found")
	}

	return findDoc.StockReturnProductInfo, nil
}

func (svc StockReturnProductHttpService) SearchStockReturnProduct(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.StockReturnProductInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"docno",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.StockReturnProductInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc StockReturnProductHttpService) SearchStockReturnProductStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.StockReturnProductInfo, int, error) {
	searchInFields := []string{
		"docno",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, filters, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.StockReturnProductInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc StockReturnProductHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.StockReturnProduct) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.StockReturnProduct](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.StockReturnProduct, models.StockReturnProductDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.StockReturnProduct) models.StockReturnProductDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.StockReturnProductDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.StockReturnProduct = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.StockReturnProduct, models.StockReturnProductDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.StockReturnProductDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "docno", guid)
		},
		func(doc models.StockReturnProductDoc) bool {
			return doc.Docno != ""
		},
		func(shopID string, authUsername string, data models.StockReturnProduct, doc models.StockReturnProductDoc) error {

			doc.StockReturnProduct = data
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

func (svc StockReturnProductHttpService) getDocIDKey(doc models.StockReturnProduct) string {
	return doc.Docno
}

func (svc StockReturnProductHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc StockReturnProductHttpService) GetModuleName() string {
	return "stockReturnProduct"
}
