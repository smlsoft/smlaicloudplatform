package services

import (
	"errors"
	micromodels "smlcloudplatform/internal/microservice/models"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/product/models"
	"smlcloudplatform/pkg/product/product/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IProductHttpService interface {
	CreateProduct(shopID string, authUsername string, doc models.Product) (string, error)
	UpdateProduct(shopID string, guid string, authUsername string, doc models.Product) error
	DeleteProduct(shopID string, guid string, authUsername string) error
	DeleteProductByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoProduct(shopID string, guid string) (models.ProductInfo, error)
	SearchProduct(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.ProductInfo, mongopagination.PaginationData, error)
	SearchProductStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ProductInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Product) (common.BulkImport, error)
}

type ProductHttpService struct {
	repo   repositories.IProductRepository
	mqRepo repositories.IProductMessageQueueRepository
}

func NewProductHttpService(repo repositories.IProductRepository, mqRepo repositories.IProductMessageQueueRepository) *ProductHttpService {

	return &ProductHttpService{
		repo:   repo,
		mqRepo: mqRepo,
	}
}

func (svc ProductHttpService) CreateProduct(shopID string, authUsername string, doc models.Product) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "itemcode", doc.ItemCode)

	if err != nil {
		return "", err
	}

	if findDoc.ItemCode != "" {
		return "", errors.New("ItemCode is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.ProductDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Product = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	err = svc.mqRepo.Create(docData)
	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc ProductHttpService) UpdateProduct(shopID string, guid string, authUsername string, doc models.Product) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.Product = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	err = svc.mqRepo.Update(findDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc ProductHttpService) DeleteProduct(shopID string, guid string, authUsername string) error {

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

	err = svc.mqRepo.Delete(findDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc ProductHttpService) DeleteProductByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	findDocAllTemp, err := svc.repo.FindByGuids(shopID, GUIDs)
	if err != nil {
		return err
	}

	err = svc.mqRepo.DeleteInBatch(findDocAllTemp)

	if err != nil {
		return err
	}

	return nil
}

func (svc ProductHttpService) InfoProduct(shopID string, guid string) (models.ProductInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.ProductInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.ProductInfo{}, errors.New("document not found")
	}

	return findDoc.ProductInfo, nil

}

func (svc ProductHttpService) SearchProduct(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.ProductInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"itemcode",
		"names.name",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.ProductInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc ProductHttpService) SearchProductStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ProductInfo, int, error) {
	searchInFields := []string{
		"itemcode",
		"names.name",
	}

	selectFields := map[string]interface{}{
		"guidfixed":       1,
		"itemcode":        1,
		"categoryguid":    1,
		"barcodes":        1,
		"names":           1,
		"multiunit":       1,
		"useserialnumber": 1,
		"units":           1,
		"unitcost":        1,
		"itemstocktype":   1,
		"itemtype":        1,
		"vattype":         1,
		"issumpoint":      1,
		"images":          1,
		"prices":          1,
	}

	if langCode != "" {
		selectFields["names"] = bson.M{"$elemMatch": bson.M{"code": langCode}}
	} else {
		selectFields["names"] = 1
	}

	docList, total, err := svc.repo.FindStep(shopID, filters, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.ProductInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc ProductHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.Product) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Product](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.ItemCode)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "itemcode", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.ItemCode)
	}

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Product, models.ProductDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Product) models.ProductDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.ProductDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Product = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Product, models.ProductDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.ProductDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "itemcode", guid)
		},
		func(doc models.ProductDoc) bool {
			return doc.ItemCode != ""
		},
		func(shopID string, authUsername string, data models.Product, doc models.ProductDoc) error {

			doc.Product = data
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
		createDataKey = append(createDataKey, doc.ItemCode)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.ItemCode)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {
		updateDataKey = append(updateDataKey, doc.ItemCode)
	}

	updateFailDataKey := []string{}
	for _, doc := range updateFailDataList {
		updateFailDataKey = append(updateFailDataKey, svc.getDocIDKey(doc))
	}

	err = svc.mqRepo.CreateInBatch(createDataList)

	if err != nil {
		return common.BulkImport{}, err
	}

	err = svc.mqRepo.UpdateInBatch(updateSuccessDataList)

	if err != nil {
		return common.BulkImport{}, err
	}

	return common.BulkImport{
		Created:          createDataKey,
		Updated:          updateDataKey,
		UpdateFailed:     updateFailDataKey,
		PayloadDuplicate: payloadDuplicateDataKey,
	}, nil
}

func (svc ProductHttpService) getDocIDKey(doc models.Product) string {
	return doc.ItemCode
}
