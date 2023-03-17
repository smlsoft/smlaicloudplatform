package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/productgroup/models"
	"smlcloudplatform/pkg/product/productgroup/repositories"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IProductGroupHttpService interface {
	SaveProductGroup(shopID string, authUsername string, doc models.ProductGroup) (string, error)
	CreateProductGroup(shopID string, authUsername string, doc models.ProductGroup) (string, error)
	UpdateProductGroup(shopID string, guid string, authUsername string, doc models.ProductGroup) error
	DeleteProductGroup(shopID string, guid string, authUsername string) error
	DeleteProductGroupByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoProductGroup(shopID string, guid string) (models.ProductGroupInfo, error)
	InfoWTFArray(shopID string, unitCodes []string) ([]interface{}, error)
	SearchProductGroup(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.ProductGroupInfo, mongopagination.PaginationData, error)
	SearchProductGroupStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.ProductGroupInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.ProductGroup) (common.BulkImport, error)

	GetModuleName() string
}

type ProductGroupHttpService struct {
	repo repositories.IProductGroupRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.ProductGroupActivity, models.ProductGroupDeleteActivity]
}

func NewProductGroupHttpService(repo repositories.IProductGroupRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *ProductGroupHttpService {

	insSvc := &ProductGroupHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.ProductGroupActivity, models.ProductGroupDeleteActivity](repo)

	return insSvc
}

func (svc ProductGroupHttpService) SaveProductGroup(shopID string, authUsername string, doc models.ProductGroup) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", doc.Code)

	if err != nil {
		return "", err
	}

	if len(findDoc.Code) > 0 {

		err = svc.update(shopID, authUsername, findDoc.GuidFixed, findDoc, doc)

		if err != nil {
			return "", err
		}

		return findDoc.GuidFixed, nil

	} else {
		return svc.create(shopID, authUsername, doc)
	}
}

func (svc ProductGroupHttpService) CreateProductGroup(shopID string, authUsername string, doc models.ProductGroup) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", doc.Code)

	if err != nil {
		return "", err
	}

	if findDoc.Code != "" {
		return "", errors.New("code is exists")
	}

	svc.saveMasterSync(shopID)

	return svc.create(shopID, authUsername, doc)
}

func (svc ProductGroupHttpService) create(shopID string, authUsername string, doc models.ProductGroup) (string, error) {
	newGuidFixed := utils.NewGUID()

	docData := models.ProductGroupDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.ProductGroup = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err := svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc ProductGroupHttpService) UpdateProductGroup(shopID string, guid string, authUsername string, doc models.ProductGroup) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.ProductGroup = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.update(shopID, authUsername, guid, findDoc, doc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc ProductGroupHttpService) update(shopID string, authUsername string, guid string, findDoc models.ProductGroupDoc, docUpdate models.ProductGroup) error {
	findDoc.ProductGroup = docUpdate

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err := svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc ProductGroupHttpService) DeleteProductGroup(shopID string, guid string, authUsername string) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	err = svc.repo.DeleteByGuidfixed(shopID, guid, authUsername)
	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc ProductGroupHttpService) DeleteProductGroupByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc ProductGroupHttpService) InfoProductGroup(shopID string, guid string) (models.ProductGroupInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.ProductGroupInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.ProductGroupInfo{}, errors.New("document not found")
	}

	return findDoc.ProductGroupInfo, nil

}

func (svc ProductGroupHttpService) InfoWTFArray(shopID string, codes []string) ([]interface{}, error) {
	docList := []interface{}{}

	for _, code := range codes {
		findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", code)
		if err != nil || findDoc.ID == primitive.NilObjectID {
			// add item empty
			docList = append(docList, nil)
		} else {
			docList = append(docList, findDoc.ProductGroupInfo)
		}
	}

	return docList, nil
}

func (svc ProductGroupHttpService) SearchProductGroup(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.ProductGroupInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"names.name",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.ProductGroupInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc ProductGroupHttpService) SearchProductGroupStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.ProductGroupInfo, int, error) {
	searchInFields := []string{
		"code",
	}

	selectFields := map[string]interface{}{
		"guidfixed": 1,
		"code":      1,
	}

	if langCode != "" {
		selectFields["names"] = bson.M{"$elemMatch": bson.M{"code": langCode}}
	} else {
		selectFields["names"] = 1
	}

	docList, total, err := svc.repo.FindStep(shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.ProductGroupInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc ProductGroupHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.ProductGroup) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.ProductGroup](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.ProductGroup, models.ProductGroupDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.ProductGroup) models.ProductGroupDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.ProductGroupDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.ProductGroup = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.ProductGroup, models.ProductGroupDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.ProductGroupDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "code", guid)
		},
		func(doc models.ProductGroupDoc) bool {
			return doc.Code != ""
		},
		func(shopID string, authUsername string, data models.ProductGroup, doc models.ProductGroupDoc) error {

			doc.ProductGroup = data
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

func (svc ProductGroupHttpService) getDocIDKey(doc models.ProductGroup) string {
	return doc.Code
}

func (svc ProductGroupHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc ProductGroupHttpService) GetModuleName() string {
	return "productGroup"
}
