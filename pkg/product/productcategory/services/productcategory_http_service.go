package services

import (
	"errors"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/productcategory/models"
	"smlcloudplatform/pkg/product/productcategory/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IProductCategoryHttpService interface {
	CreateProductCategory(shopID string, authUsername string, doc models.ProductCategory) (string, error)
	UpdateProductCategory(shopID string, guid string, authUsername string, doc models.ProductCategory) error
	DeleteProductCategory(shopID string, guid string, authUsername string) error
	InfoProductCategory(shopID string, guid string) (models.ProductCategoryInfo, error)
	SearchProductCategory(shopID string, q string, page int, limit int, sort map[string]int) ([]models.ProductCategoryInfo, mongopagination.PaginationData, error)
	SearchProductCategoryStep(shopID string, langCode string, q string, skip int, limit int, sort map[string]int) ([]models.ProductCategoryInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.ProductCategory) (common.BulkImport, error)
}

type ProductCategoryHttpService struct {
	repo repositories.IProductCategoryRepository
}

func NewProductCategoryHttpService(repo repositories.IProductCategoryRepository) *ProductCategoryHttpService {

	return &ProductCategoryHttpService{
		repo: repo,
	}
}

func (svc ProductCategoryHttpService) CreateProductCategory(shopID string, authUsername string, doc models.ProductCategory) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", doc.Code)

	if err != nil {
		return "", err
	}

	if findDoc.Code != "" {
		return "", errors.New("Code is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.ProductCategoryDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.ProductCategory = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc ProductCategoryHttpService) UpdateProductCategory(shopID string, guid string, authUsername string, doc models.ProductCategory) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.ProductCategory = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc ProductCategoryHttpService) DeleteProductCategory(shopID string, guid string, authUsername string) error {

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

	return nil
}

func (svc ProductCategoryHttpService) InfoProductCategory(shopID string, guid string) (models.ProductCategoryInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.ProductCategoryInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.ProductCategoryInfo{}, errors.New("document not found")
	}

	return findDoc.ProductCategoryInfo, nil

}

func (svc ProductCategoryHttpService) SearchProductCategory(shopID string, q string, page int, limit int, sort map[string]int) ([]models.ProductCategoryInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"guidfixed",
		"code",
	}

	docList, pagination, err := svc.repo.FindPageSort(shopID, searchCols, q, page, limit, sort)

	if err != nil {
		return []models.ProductCategoryInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc ProductCategoryHttpService) SearchProductCategoryStep(shopID string, langCode string, q string, skip int, limit int, sort map[string]int) ([]models.ProductCategoryInfo, int, error) {
	searchCols := []string{
		"guidfixed",
		"code",
	}

	projectQuery := map[string]interface{}{
		"guidfixed": 1,
		"code":      1,
	}

	if langCode != "" {
		projectQuery["names"] = bson.M{"$elemMatch": bson.M{"code": langCode}}
	} else {
		projectQuery["names"] = 1
	}

	docList, total, err := svc.repo.FindLimit(shopID, searchCols, q, skip, limit, sort, projectQuery)

	if err != nil {
		return []models.ProductCategoryInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc ProductCategoryHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.ProductCategory) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.ProductCategory](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.ProductCategory, models.ProductCategoryDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.ProductCategory) models.ProductCategoryDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.ProductCategoryDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.ProductCategory = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.ProductCategory, models.ProductCategoryDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.ProductCategoryDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "code", guid)
		},
		func(doc models.ProductCategoryDoc) bool {
			return doc.Code != ""
		},
		func(shopID string, authUsername string, data models.ProductCategory, doc models.ProductCategoryDoc) error {

			doc.ProductCategory = data
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

	return common.BulkImport{
		Created:          createDataKey,
		Updated:          updateDataKey,
		UpdateFailed:     updateFailDataKey,
		PayloadDuplicate: payloadDuplicateDataKey,
	}, nil
}

func (svc ProductCategoryHttpService) getDocIDKey(doc models.ProductCategory) string {
	return doc.Code
}
