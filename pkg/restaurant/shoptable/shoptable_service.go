package shoptable

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/restaurant/shoptable/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IShopTableService interface {
	CreateShopTable(shopID string, authUsername string, doc models.ShopTable) (string, error)
	UpdateShopTable(shopID string, guid string, authUsername string, doc models.ShopTable) error
	DeleteShopTable(shopID string, guid string, authUsername string) error
	InfoShopTable(shopID string, guid string) (models.ShopTableInfo, error)
	SearchShopTable(shopID string, pageable micromodels.Pageable) ([]models.ShopTableInfo, mongopagination.PaginationData, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.ShopTable) (common.BulkImport, error)

	LastActivity(shopID string, action string, lastUpdatedDate time.Time, pageable micromodels.Pageable) (common.LastActivity, mongopagination.PaginationData, error)
	GetModuleName() string
}

type ShopTableService struct {
	repo          ShopTableRepository
	syncCacheRepo mastersync.IMasterSyncCacheRepository

	services.ActivityService[models.ShopTableActivity, models.ShopTableDeleteActivity]
}

func NewShopTableService(repo ShopTableRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) ShopTableService {
	insSvc := ShopTableService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.ShopTableActivity, models.ShopTableDeleteActivity](repo)
	return insSvc
}

func (svc ShopTableService) CreateShopTable(shopID string, authUsername string, doc models.ShopTable) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := models.ShopTableDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.ShopTable = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err := svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc ShopTableService) UpdateShopTable(shopID string, guid string, authUsername string, doc models.ShopTable) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.ShopTable = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc ShopTableService) DeleteShopTable(shopID string, guid string, authUsername string) error {
	err := svc.repo.DeleteByGuidfixed(shopID, guid, authUsername)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc ShopTableService) InfoShopTable(shopID string, guid string) (models.ShopTableInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.ShopTableInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.ShopTableInfo{}, errors.New("document not found")
	}

	return findDoc.ShopTableInfo, nil

}

func (svc ShopTableService) SearchShopTable(shopID string, pageable micromodels.Pageable) ([]models.ShopTableInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"guidfixed",
		"code",
	}

	for i := range [5]bool{} {
		searchInFields = append(searchInFields, fmt.Sprintf("name%d", (i+1)))
	}

	docList, pagination, err := svc.repo.FindPage(shopID, searchInFields, pageable)

	if err != nil {
		return []models.ShopTableInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc ShopTableService) SaveInBatch(shopID string, authUsername string, dataList []models.ShopTable) (common.BulkImport, error) {

	// createDataList := []models.ShopTableDoc{}
	// duplicateDataList := []models.ShopTable{}

	payloadCategoryList, payloadDuplicateCategoryList := importdata.FilterDuplicate[models.ShopTable](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadCategoryList {
		itemCodeGuidList = append(itemCodeGuidList, doc.Number)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "code", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.Code)
	}

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.ShopTable, models.ShopTableDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadCategoryList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.ShopTable) models.ShopTableDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.ShopTableDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.ShopTable = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			dataDoc.LastUpdatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.ShopTable, models.ShopTableDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.ShopTableDoc, error) {
			return svc.repo.FindByGuid(shopID, guid)
		},
		func(doc models.ShopTableDoc) bool {
			return false
		},
		func(shopID string, authUsername string, data models.ShopTable, doc models.ShopTableDoc) error {

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
		createDataKey = append(createDataKey, doc.Number)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateCategoryList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.Number)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {
		updateDataKey = append(updateDataKey, doc.Number)
	}

	updateFailDataKey := []string{}
	for _, doc := range updateFailDataList {
		updateFailDataKey = append(updateFailDataKey, doc.Number)
	}

	svc.saveMasterSync(shopID)

	return common.BulkImport{
		Created:          createDataKey,
		Updated:          updateDataKey,
		UpdateFailed:     updateFailDataKey,
		PayloadDuplicate: payloadDuplicateDataKey,
	}, nil
}

func (svc ShopTableService) getDocIDKey(doc models.ShopTable) string {
	return doc.Number
}

func (svc ShopTableService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc ShopTableService) GetModuleName() string {
	return "shoptable"
}
