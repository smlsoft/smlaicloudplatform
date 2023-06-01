package kitchen

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/restaurant/kitchen/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IKitchenService interface {
	CreateKitchen(shopID string, authUsername string, doc models.Kitchen) (string, error)
	UpdateKitchen(shopID string, guid string, authUsername string, doc models.Kitchen) error
	DeleteKitchen(shopID string, guid string, authUsername string) error
	InfoKitchen(shopID string, guid string) (models.KitchenInfo, error)
	SearchKitchen(shopID string, pageable micromodels.Pageable) ([]models.KitchenInfo, mongopagination.PaginationData, error)
	SearchKitchenStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.KitchenInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Kitchen) (common.BulkImport, error)

	// LastActivity(shopID string, action string, lastUpdatedDate time.Time, pageable micromodels.Pageable) (common.LastActivity, mongopagination.PaginationData, error)

	GetModuleName() string
}

type KitchenService struct {
	repo          KitchenRepository
	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.KitchenActivity, models.KitchenDeleteActivity]
}

func NewKitchenService(repo KitchenRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) KitchenService {

	insSvc := KitchenService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.KitchenActivity, models.KitchenDeleteActivity](repo)
	return insSvc
}

func (svc KitchenService) CreateKitchen(shopID string, authUsername string, doc models.Kitchen) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := models.KitchenDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Kitchen = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err := svc.repo.Create(docData)

	if err != nil {
		return "", err

	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc KitchenService) UpdateKitchen(shopID string, guid string, authUsername string, doc models.Kitchen) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.Kitchen = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc KitchenService) DeleteKitchen(shopID string, guid string, authUsername string) error {
	err := svc.repo.DeleteByGuidfixed(shopID, guid, authUsername)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc KitchenService) InfoKitchen(shopID string, guid string) (models.KitchenInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.KitchenInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.KitchenInfo{}, errors.New("document not found")
	}

	return findDoc.KitchenInfo, nil

}

func (svc KitchenService) SearchKitchen(shopID string, pageable micromodels.Pageable) ([]models.KitchenInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"code",
		"name1",
		"name2",
		"name3",
		"name4",
		"name5",
	}

	for i := range [5]bool{} {
		searchInFields = append(searchInFields, fmt.Sprintf("name%d", (i+1)))
	}

	docList, pagination, err := svc.repo.FindPage(shopID, searchInFields, pageable)

	if err != nil {
		return []models.KitchenInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc KitchenService) SearchKitchenStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.KitchenInfo, int, error) {
	searchInFields := []string{
		"code",
		"name1",
		"name2",
		"name3",
		"name4",
		"name5",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.KitchenInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc KitchenService) SaveInBatch(shopID string, authUsername string, dataList []models.Kitchen) (common.BulkImport, error) {

	// createDataList := []models.KitchenDoc{}
	// duplicateDataList := []models.Kitchen{}

	payloadCategoryList, payloadDuplicateCategoryList := importdata.FilterDuplicate[models.Kitchen](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadCategoryList {
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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Kitchen, models.KitchenDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadCategoryList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Kitchen) models.KitchenDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.KitchenDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Kitchen = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			dataDoc.LastUpdatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Kitchen, models.KitchenDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.KitchenDoc, error) {
			return svc.repo.FindByGuid(shopID, guid)
		},
		func(doc models.KitchenDoc) bool {
			return false
		},
		func(shopID string, authUsername string, data models.Kitchen, doc models.KitchenDoc) error {

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
	for _, doc := range payloadDuplicateCategoryList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.Code)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {
		updateDataKey = append(updateDataKey, doc.Code)
	}

	updateFailDataKey := []string{}
	for _, doc := range updateFailDataList {
		updateFailDataKey = append(updateFailDataKey, doc.Code)
	}

	svc.saveMasterSync(shopID)

	return common.BulkImport{
		Created:          createDataKey,
		Updated:          updateDataKey,
		UpdateFailed:     updateFailDataKey,
		PayloadDuplicate: payloadDuplicateDataKey,
	}, nil
}

func (svc KitchenService) getDocIDKey(doc models.Kitchen) string {
	return doc.Code
}

func (svc KitchenService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc KitchenService) GetModuleName() string {
	return "kitchen"
}
