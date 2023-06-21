package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/promotion/models"
	"smlcloudplatform/pkg/product/promotion/repositories"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IPromotionHttpService interface {
	CreatePromotion(shopID string, authUsername string, doc models.Promotion) (string, error)
	UpdatePromotion(shopID string, guid string, authUsername string, doc models.Promotion) error
	DeletePromotion(shopID string, guid string, authUsername string) error
	DeletePromotionByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoPromotion(shopID string, guid string) (models.PromotionInfo, error)
	InfoPromotionByCode(shopID string, code string) (models.PromotionInfo, error)
	SearchPromotion(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.PromotionInfo, mongopagination.PaginationData, error)
	SearchPromotionStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.PromotionInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Promotion) (common.BulkImport, error)

	GetModuleName() string
}

type PromotionHttpService struct {
	repo repositories.IPromotionRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.PromotionActivity, models.PromotionDeleteActivity]
}

func NewPromotionHttpService(repo repositories.IPromotionRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *PromotionHttpService {

	insSvc := &PromotionHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.PromotionActivity, models.PromotionDeleteActivity](repo)

	return insSvc
}

func (svc PromotionHttpService) CreatePromotion(shopID string, authUsername string, doc models.Promotion) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", doc.Code)

	if err != nil {
		return "", err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", errors.New("code is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.PromotionDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Promotion = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	if docData.Details == nil {
		docData.Details = &[]models.PromotionDetail{}
	}

	if docData.ProductBarcode.Names == nil {
		docData.ProductBarcode.Names = &[]common.NameX{}
	}

	if docData.ProductBarcode.ItemUnitNames == nil {
		docData.ProductBarcode.Names = &[]common.NameX{}
	}

	for idx, detail := range *docData.Details {
		tempDoc := (*docData.Details)[idx]

		if detail.ProductBarcode.Names == nil {
			tempDoc.ProductBarcode.Names = &[]common.NameX{}
		}

		if detail.ProductBarcode.ItemUnitNames == nil {
			tempDoc.ProductBarcode.ItemUnitNames = &[]common.NameX{}
		}
	}

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc PromotionHttpService) UpdatePromotion(shopID string, guid string, authUsername string, doc models.Promotion) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	docData := models.PromotionDoc{}

	docData.Promotion = doc

	docData.UpdatedBy = authUsername
	docData.UpdatedAt = time.Now()

	if docData.Details == nil {
		docData.Details = &[]models.PromotionDetail{}
	}

	if docData.ProductBarcode.Names == nil {
		docData.ProductBarcode.Names = &[]common.NameX{}
	}

	if docData.ProductBarcode.ItemUnitNames == nil {
		docData.ProductBarcode.Names = &[]common.NameX{}
	}

	for idx, detail := range *docData.Details {
		tempDoc := (*docData.Details)[idx]

		if detail.ProductBarcode.Names == nil {
			tempDoc.ProductBarcode.Names = &[]common.NameX{}
		}

		if detail.ProductBarcode.ItemUnitNames == nil {
			tempDoc.ProductBarcode.ItemUnitNames = &[]common.NameX{}
		}
	}
	err = svc.repo.Update(shopID, guid, docData)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc PromotionHttpService) DeletePromotion(shopID string, guid string, authUsername string) error {

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

func (svc PromotionHttpService) DeletePromotionByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc PromotionHttpService) InfoPromotion(shopID string, guid string) (models.PromotionInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.PromotionInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.PromotionInfo{}, errors.New("document not found")
	}

	return findDoc.PromotionInfo, nil
}

func (svc PromotionHttpService) InfoPromotionByCode(shopID string, code string) (models.PromotionInfo, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", code)

	if err != nil {
		return models.PromotionInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.PromotionInfo{}, errors.New("document not found")
	}

	return findDoc.PromotionInfo, nil
}

func (svc PromotionHttpService) SearchPromotion(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.PromotionInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"code",
		"name",
		"barcode",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.PromotionInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc PromotionHttpService) SearchPromotionStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.PromotionInfo, int, error) {
	searchInFields := []string{
		"code",
		"name",
		"barcode",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, filters, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.PromotionInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc PromotionHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.Promotion) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Promotion](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.Code)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "code", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.Code)
	}

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Promotion, models.PromotionDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Promotion) models.PromotionDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.PromotionDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Promotion = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Promotion, models.PromotionDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.PromotionDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "code", guid)
		},
		func(doc models.PromotionDoc) bool {
			return doc.Code != ""
		},
		func(shopID string, authUsername string, data models.Promotion, doc models.PromotionDoc) error {

			doc.Promotion = data
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
		createDataKey = append(createDataKey, doc.Code)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.Code)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {

		updateDataKey = append(updateDataKey, doc.Code)
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

func (svc PromotionHttpService) getDocIDKey(doc models.Promotion) string {
	return doc.Code
}

func (svc PromotionHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc PromotionHttpService) GetModuleName() string {
	return "promotion"
}
