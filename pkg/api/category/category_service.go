package category

import (
	"errors"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"sync"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ICategoryService interface {
	CreateCategory(shopID string, authUsername string, category models.Category) (string, error)
	UpdateCategory(guid string, shopID string, authUsername string, category models.Category) error
	DeleteCategory(guid string, shopID string, authUsername string) error
	InfoCategory(guid string, shopID string) (models.CategoryInfo, error)
	SearchCategory(shopID string, q string, page int, limit int) ([]models.CategoryInfo, paginate.PaginationData, error)
	LastActivityCategory(shopID string, lastUpdatedDate time.Time, page int, limit int) (models.LastActivity, paginate.PaginationData, error)
	CreateInBatch(shopID string, authUsername string, categories []models.Category) (models.CategoryBulkImport, error)
}

type CategoryService struct {
	repo ICategoryRepository
}

func NewCategoryService(categoryRepository ICategoryRepository) CategoryService {
	return CategoryService{
		repo: categoryRepository,
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

func (svc CategoryService) UpdateCategory(guid string, shopID string, authUsername string, category models.Category) error {

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

func (svc CategoryService) DeleteCategory(guid string, shopID string, authUsername string) error {
	err := svc.repo.Delete(shopID, guid, authUsername)

	if err != nil {
		return err
	}
	return nil
}

func (svc CategoryService) InfoCategory(guid string, shopID string) (models.CategoryInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.CategoryInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.CategoryInfo{}, errors.New("document not found")
	}

	return findDoc.CategoryInfo, nil

}

func (svc CategoryService) SearchCategory(shopID string, q string, page int, limit int) ([]models.CategoryInfo, paginate.PaginationData, error) {
	docList, pagination, err := svc.repo.FindPage(shopID, q, page, limit)

	if err != nil {
		return []models.CategoryInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc CategoryService) LastActivityCategory(shopID string, lastUpdatedDate time.Time, page int, limit int) (models.LastActivity, paginate.PaginationData, error) {
	var wg sync.WaitGroup

	wg.Add(1)
	var deleteDocList []models.CategoryDeleteActivity
	var pagination1 paginate.PaginationData
	var err1 error

	go func() {
		deleteDocList, pagination1, err1 = svc.repo.FindDeletedPage(shopID, lastUpdatedDate, page, limit)
		wg.Done()
	}()

	wg.Add(1)
	var createAndUpdateDocList []models.CategoryActivity
	var pagination2 paginate.PaginationData
	var err2 error

	go func() {
		createAndUpdateDocList, pagination2, err2 = svc.repo.FindCreatedOrUpdatedPage(shopID, lastUpdatedDate, page, limit)
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

	if pagination.TotalPage < pagination2.TotalPage {
		pagination = pagination2
	}

	return lastActivity, pagination, nil
}

func (svc CategoryService) CreateInBatch(shopID string, authUsername string, categories []models.Category) (models.CategoryBulkImport, error) {

	createDataList := []models.CategoryDoc{}
	duplicateDataList := []models.Category{}

	payloadCategoryList, payloadDuplicateCategoryList := filterDuplicateCategory(categories)

	itemCodeGuidList := []string{}
	for _, category := range payloadCategoryList {
		itemCodeGuidList = append(itemCodeGuidList, category.CategoryGuid)
	}

	findItemGuid, err := svc.repo.FindByCategoryGuidList(shopID, itemCodeGuidList)

	if err != nil {
		return models.CategoryBulkImport{}, err
	}

	duplicateDataList, createDataList = preparePayloadDataCategory(shopID, authUsername, findItemGuid, payloadCategoryList)

	updateSuccessDataList, updateFailDataList := updateOnDuplicateCategory(shopID, authUsername, duplicateDataList, svc.repo)

	if len(createDataList) > 0 {
		err = svc.repo.CreateInBatch(createDataList)

		if err != nil {
			return models.CategoryBulkImport{}, err
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

	return models.CategoryBulkImport{
		Created:          createDataKey,
		Updated:          updateDataKey,
		UpdateFailed:     updateFailDataKey,
		PayloadDuplicate: payloadDuplicateDataKey,
	}, nil
}

func filterDuplicateCategory(categories []models.Category) (itemTemp []models.Category, itemDuplicate []models.Category) {
	tempFilterDict := map[string]models.Category{}
	for _, category := range categories {
		if _, ok := tempFilterDict[category.CategoryGuid]; ok {
			itemDuplicate = append(itemDuplicate, category)

		}
		tempFilterDict[category.CategoryGuid] = category
	}

	for _, inventory := range tempFilterDict {
		itemTemp = append(itemTemp, inventory)
	}

	return itemTemp, itemDuplicate
}

func updateOnDuplicateCategory(shopID string, authUsername string, duplicateDataList []models.Category, repo ICategoryRepository) ([]models.CategoryDoc, []models.Category) {
	updateSuccessDataList := []models.CategoryDoc{}
	updateFailDataList := []models.Category{}

	for _, doc := range duplicateDataList {
		findDoc, err := repo.FindByCategoryGuid(shopID, doc.CategoryGuid)

		if err != nil || findDoc.ID == primitive.NilObjectID {
			updateFailDataList = append(updateFailDataList, doc)
			continue
		}

		findDoc.Category = doc

		findDoc.UpdatedBy = authUsername
		findDoc.UpdatedAt = time.Now()
		findDoc.LastUpdatedAt = time.Now()

		err = repo.Update(findDoc.GuidFixed, findDoc)

		if err != nil {
			updateFailDataList = append(updateFailDataList, doc)
			continue
		}

		updateSuccessDataList = append(updateSuccessDataList, findDoc)
	}
	return updateSuccessDataList, updateFailDataList
}

func preparePayloadDataCategory(shopID string, authUsername string, itemGuidList []models.CategoryItemCategoryGuid, payloadCategoryList []models.Category) ([]models.Category, []models.CategoryDoc) {
	tempItemGuidDict := make(map[string]bool)
	duplicateDataList := []models.Category{}
	createDataList := []models.CategoryDoc{}

	for _, itemGuid := range itemGuidList {
		tempItemGuidDict[itemGuid.CategoryGuid] = true
	}

	for _, categories := range payloadCategoryList {

		if _, ok := tempItemGuidDict[categories.CategoryGuid]; ok {
			duplicateDataList = append(duplicateDataList, categories)
		} else {
			newGuid := utils.NewGUID()

			dataDoc := models.CategoryDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Category = categories

			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = time.Now()
			dataDoc.LastUpdatedAt = time.Now()

			createDataList = append(createDataList, dataDoc)
		}
	}
	return duplicateDataList, createDataList
}
