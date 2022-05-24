package category

import (
	"errors"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"sync"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ICategoryService interface {
	CreateCategory(shopID string, authUsername string, category models.Category) (string, error)
	UpdateCategory(shopID string, guid string, authUsername string, category models.Category) error
	DeleteCategory(shopID string, guid string, authUsername string) error
	InfoCategory(shopID string, guid string) (models.CategoryInfo, error)
	SearchCategory(shopID string, q string, page int, limit int) ([]models.CategoryInfo, mongopagination.PaginationData, error)
	LastActivity(shopID string, lastUpdatedDate time.Time, page int, limit int) (models.LastActivity, mongopagination.PaginationData, error)
	SaveInBatch(shopID string, authUsername string, categories []models.Category) (models.BulkImport, error)
}

type CategoryService struct {
	repo         ICategoryRepository
	guidRepo     repositories.GuidRepository[models.CategoryItemGuid]
	activityRepo repositories.ActivityRepository[models.CategoryActivity, models.CategoryDeleteActivity]
}

func NewCategoryService(
	categoryRepository ICategoryRepository,
	guidRepo repositories.GuidRepository[models.CategoryItemGuid],
	activityRepo repositories.ActivityRepository[models.CategoryActivity, models.CategoryDeleteActivity],
) CategoryService {
	return CategoryService{
		repo:         categoryRepository,
		guidRepo:     guidRepo,
		activityRepo: activityRepo,
	}
}

func (svc CategoryService) CreateCategory(shopID string, authUsername string, category models.Category) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := models.CategoryDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Category = category

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err := svc.repo.Create(docData)

	if err != nil {
		return "", err
	}
	return newGuidFixed, nil
}

func (svc CategoryService) UpdateCategory(shopID string, guid string, authUsername string, category models.Category) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.Category = category

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(guid, findDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc CategoryService) DeleteCategory(shopID string, guid string, authUsername string) error {
	err := svc.repo.Delete(shopID, guid, authUsername)

	if err != nil {
		return err
	}
	return nil
}

func (svc CategoryService) InfoCategory(shopID string, guid string) (models.CategoryInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.CategoryInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.CategoryInfo{}, errors.New("document not found")
	}

	return findDoc.CategoryInfo, nil

}

func (svc CategoryService) SearchCategory(shopID string, q string, page int, limit int) ([]models.CategoryInfo, mongopagination.PaginationData, error) {
	docList, pagination, err := svc.repo.FindPage(shopID, q, page, limit)

	if err != nil {
		return []models.CategoryInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc CategoryService) LastActivity(shopID string, lastUpdatedDate time.Time, page int, limit int) (models.LastActivity, mongopagination.PaginationData, error) {
	var wg sync.WaitGroup

	wg.Add(1)
	var deleteDocList []models.CategoryDeleteActivity
	var pagination1 mongopagination.PaginationData
	var err1 error

	go func() {
		deleteDocList, pagination1, err1 = svc.activityRepo.FindDeletedPage(shopID, lastUpdatedDate, page, limit)
		wg.Done()
	}()

	wg.Add(1)
	var createAndUpdateDocList []models.CategoryActivity
	var pagination2 mongopagination.PaginationData
	var err2 error

	go func() {
		createAndUpdateDocList, pagination2, err2 = svc.activityRepo.FindCreatedOrUpdatedPage(shopID, lastUpdatedDate, page, limit)
		wg.Done()
	}()

	wg.Wait()

	if err1 != nil {
		return models.LastActivity{}, pagination1, err1
	}

	if err2 != nil {
		return models.LastActivity{}, pagination2, err2
	}

	lastActivity := models.LastActivity{}

	lastActivity.Remove = &deleteDocList
	lastActivity.New = &createAndUpdateDocList

	pagination := pagination1

	if pagination.Total < pagination2.Total {
		pagination = pagination2
	}

	return lastActivity, pagination, nil
}

func (svc CategoryService) SaveInBatch(shopID string, authUsername string, dataList []models.Category) (models.BulkImport, error) {

	createDataList := []models.CategoryDoc{}
	duplicateDataList := []models.Category{}

	payloadCategoryList, payloadDuplicateCategoryList := importdata.FilterDuplicate[models.Category](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadCategoryList {
		itemCodeGuidList = append(itemCodeGuidList, doc.CategoryGuid)
	}

	findItemGuid, err := svc.guidRepo.FindInItemGuid(shopID, "code", itemCodeGuidList)

	if err != nil {
		return models.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.CategoryGuid)
	}

	duplicateDataList, createDataList = importdata.PreparePayloadData[models.Category, models.CategoryDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadCategoryList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Category) models.CategoryDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.CategoryDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Category = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			dataDoc.LastUpdatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Category, models.CategoryDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.CategoryDoc, error) {
			return svc.repo.FindByGuid(shopID, guid)
		},
		func(doc models.CategoryDoc) bool {
			return false
		},
		func(shopID string, authUsername string, data models.Category, doc models.CategoryDoc) error {

			return nil
		},
	)

	if len(createDataList) > 0 {
		err = svc.repo.CreateInBatch(createDataList)

		if err != nil {
			return models.BulkImport{}, err
		}
	}
	createDataKey := []string{}

	for _, doc := range createDataList {
		createDataKey = append(createDataKey, doc.CategoryGuid)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateCategoryList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.CategoryGuid)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {
		updateDataKey = append(updateDataKey, doc.CategoryGuid)
	}

	updateFailDataKey := []string{}
	for _, doc := range updateFailDataList {
		updateFailDataKey = append(updateFailDataKey, doc.CategoryGuid)
	}

	return models.BulkImport{
		Created:          createDataKey,
		Updated:          updateDataKey,
		UpdateFailed:     updateFailDataKey,
		PayloadDuplicate: payloadDuplicateDataKey,
	}, nil
}

func (svc CategoryService) getDocIDKey(doc models.Category) string {
	return doc.CategoryGuid
}
