package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/transaction/stocktransfer/models"
	"smlcloudplatform/pkg/transaction/stocktransfer/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IStockTransferHttpService interface {
	CreateStockTransfer(shopID string, authUsername string, doc models.StockTransfer) (string, error)
	UpdateStockTransfer(shopID string, guid string, authUsername string, doc models.StockTransfer) error
	DeleteStockTransfer(shopID string, guid string, authUsername string) error
	DeleteStockTransferByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoStockTransfer(shopID string, guid string) (models.StockTransferInfo, error)
	InfoStockTransferByCode(shopID string, code string) (models.StockTransferInfo, error)
	SearchStockTransfer(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.StockTransferInfo, mongopagination.PaginationData, error)
	SearchStockTransferStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.StockTransferInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.StockTransfer) (common.BulkImport, error)

	GetModuleName() string
}

type StockTransferHttpService struct {
	repo repositories.IStockTransferRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.StockTransferActivity, models.StockTransferDeleteActivity]
}

func NewStockTransferHttpService(repo repositories.IStockTransferRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *StockTransferHttpService {

	insSvc := &StockTransferHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.StockTransferActivity, models.StockTransferDeleteActivity](repo)

	return insSvc
}

func (svc StockTransferHttpService) CreateStockTransfer(shopID string, authUsername string, doc models.StockTransfer) (string, error) {

	findDoc, err := svc.repo.FindDocOne(shopID, doc.Docno, doc.TransFlag)

	if err != nil {
		return "", err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", errors.New("docno and trans flag is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.StockTransferDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.StockTransfer = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc StockTransferHttpService) UpdateStockTransfer(shopID string, guid string, authUsername string, doc models.StockTransfer) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	findExists, err := svc.repo.FindDocOne(shopID, doc.Docno, doc.TransFlag)

	if err != nil {
		return err
	}

	if findExists.Docno != findDoc.Docno && findExists.TransFlag != findDoc.TransFlag && len(findExists.GuidFixed) > 0 {
		return errors.New("docno and trans flag is exists")
	}

	findDoc.StockTransfer = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc StockTransferHttpService) DeleteStockTransfer(shopID string, guid string, authUsername string) error {

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

func (svc StockTransferHttpService) DeleteStockTransferByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc StockTransferHttpService) InfoStockTransfer(shopID string, guid string) (models.StockTransferInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.StockTransferInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.StockTransferInfo{}, errors.New("document not found")
	}

	return findDoc.StockTransferInfo, nil
}

func (svc StockTransferHttpService) InfoStockTransferByCode(shopID string, code string) (models.StockTransferInfo, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "docno", code)

	if err != nil {
		return models.StockTransferInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.StockTransferInfo{}, errors.New("document not found")
	}

	return findDoc.StockTransferInfo, nil
}

func (svc StockTransferHttpService) SearchStockTransfer(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.StockTransferInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"docno",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.StockTransferInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc StockTransferHttpService) SearchStockTransferStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.StockTransferInfo, int, error) {
	searchInFields := []string{
		"docno",
	}

	selectFields := map[string]interface{}{}

	if langCode != "" {
		selectFields["names"] = bson.M{"$elemMatch": bson.M{"code": langCode}}
	} else {
		selectFields["names"] = 1
	}

	docList, total, err := svc.repo.FindStep(shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.StockTransferInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc StockTransferHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.StockTransfer) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.StockTransfer](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.StockTransfer, models.StockTransferDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.StockTransfer) models.StockTransferDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.StockTransferDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.StockTransfer = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.StockTransfer, models.StockTransferDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.StockTransferDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "docno", guid)
		},
		func(doc models.StockTransferDoc) bool {
			return doc.Docno != ""
		},
		func(shopID string, authUsername string, data models.StockTransfer, doc models.StockTransferDoc) error {

			doc.StockTransfer = data
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

func (svc StockTransferHttpService) getDocIDKey(doc models.StockTransfer) string {
	return doc.Docno
}

func (svc StockTransferHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc StockTransferHttpService) GetModuleName() string {
	return "stockTransfer"
}
