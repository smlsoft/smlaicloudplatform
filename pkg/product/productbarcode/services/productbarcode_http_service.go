package services

import (
	"errors"
	"fmt"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/productbarcode/models"
	"smlcloudplatform/pkg/product/productbarcode/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"sync"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IProductBarcodeHttpService interface {
	CreateProductBarcode(shopID string, authUsername string, doc models.ProductBarcode) (string, error)
	UpdateProductBarcode(shopID string, guid string, authUsername string, doc models.ProductBarcode) error
	DeleteProductBarcode(shopID string, guid string, authUsername string) error
	InfoProductBarcode(shopID string, guid string) (models.ProductBarcodeInfo, error)
	SearchProductBarcode(shopID string, q string, page int, limit int, sort map[string]int) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error)
	SearchProductBarcodeStep(shopID string, langCode string, q string, skip int, limit int, sort map[string]int) ([]models.ProductBarcodeInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.ProductBarcode) (common.BulkImport, error)
}

type ProductBarcodeHttpService struct {
	repo      repositories.IProductBarcodeRepository
	cacheRepo mastersync.IMasterSyncCacheRepository
}

func NewProductBarcodeHttpService(repo repositories.IProductBarcodeRepository, cacheRepo mastersync.IMasterSyncCacheRepository) *ProductBarcodeHttpService {

	return &ProductBarcodeHttpService{
		repo:      repo,
		cacheRepo: cacheRepo,
	}
}

func (svc ProductBarcodeHttpService) CreateProductBarcode(shopID string, authUsername string, doc models.ProductBarcode) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "barcode", doc.Barcode)

	if err != nil {
		return "", err
	}

	if findDoc.Barcode != "" {
		return "", errors.New("Barcode is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.ProductBarcodeDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.ProductBarcode = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc ProductBarcodeHttpService) UpdateProductBarcode(shopID string, guid string, authUsername string, doc models.ProductBarcode) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.ProductBarcode = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc ProductBarcodeHttpService) DeleteProductBarcode(shopID string, guid string, authUsername string) error {

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

func (svc ProductBarcodeHttpService) InfoProductBarcode(shopID string, guid string) (models.ProductBarcodeInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.ProductBarcodeInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.ProductBarcodeInfo{}, errors.New("document not found")
	}

	return findDoc.ProductBarcodeInfo, nil

}

func (svc ProductBarcodeHttpService) SearchProductBarcode(shopID string, q string, page int, limit int, sort map[string]int) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"guidfixed",
		"barcode",
	}

	docList, pagination, err := svc.repo.FindPageSort(shopID, searchCols, q, page, limit, sort)

	if err != nil {
		return []models.ProductBarcodeInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc ProductBarcodeHttpService) SearchProductBarcodeStep(shopID string, langCode string, q string, skip int, limit int, sort map[string]int) ([]models.ProductBarcodeInfo, int, error) {
	searchCols := []string{
		"guidfixed",
		"barcode",
	}

	projectQuery := map[string]interface{}{
		"guidfixed":    1,
		"barcode":      1,
		"itemcode":     1,
		"categoryguid": 1,
		"names":        1,
		"itemunitcode": 1,
		"prices":       1,
		"imageuri":     1,
		"options":      1,
	}

	if langCode != "" {
		projectQuery["names"] = bson.M{"$elemMatch": bson.M{"code": langCode}}
		projectQuery["itemunitnames"] = bson.M{"$elemMatch": bson.M{"code": langCode}}
	} else {
		projectQuery["names"] = 1
		projectQuery["itemunitnames"] = 1
	}

	docList, total, err := svc.repo.FindLimit(shopID, searchCols, q, skip, limit, sort, projectQuery)

	if err != nil {
		return []models.ProductBarcodeInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc ProductBarcodeHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.ProductBarcode) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.ProductBarcode](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.Barcode)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "barcode", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.Barcode)
	}

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.ProductBarcode, models.ProductBarcodeDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.ProductBarcode) models.ProductBarcodeDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.ProductBarcodeDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.ProductBarcode = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.ProductBarcode, models.ProductBarcodeDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.ProductBarcodeDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "barcode", guid)
		},
		func(doc models.ProductBarcodeDoc) bool {
			return doc.Barcode != ""
		},
		func(shopID string, authUsername string, data models.ProductBarcode, doc models.ProductBarcodeDoc) error {

			doc.ProductBarcode = data
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
		createDataKey = append(createDataKey, doc.Barcode)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.Barcode)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {

		updateDataKey = append(updateDataKey, doc.Barcode)
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

func (svc ProductBarcodeHttpService) getDocIDKey(doc models.ProductBarcode) string {
	return doc.Barcode
}

func (svc ProductBarcodeHttpService) LastActivity(shopID string, lastUpdatedDate time.Time, page int, limit int) (common.LastActivity, mongopagination.PaginationData, error) {
	var wg sync.WaitGroup

	wg.Add(1)
	var deleteDocList []models.ProductBarcodeDeleteActivity
	var pagination1 mongopagination.PaginationData
	var err1 error

	go func() {
		deleteDocList, pagination1, err1 = svc.repo.FindDeletedPage(shopID, lastUpdatedDate, page, limit)
		wg.Done()
	}()

	wg.Add(1)
	var createAndUpdateDocList []models.ProductBarcodeActivity
	var pagination2 mongopagination.PaginationData
	var err2 error

	go func() {
		createAndUpdateDocList, pagination2, err2 = svc.repo.FindCreatedOrUpdatedPage(shopID, lastUpdatedDate, page, limit)
		wg.Done()
	}()

	wg.Wait()

	if err1 != nil {
		return common.LastActivity{}, pagination1, err1
	}

	if err2 != nil {
		return common.LastActivity{}, pagination2, err2
	}

	lastActivity := common.LastActivity{}

	lastActivity.Remove = &deleteDocList
	lastActivity.New = &createAndUpdateDocList

	pagination := pagination1

	if pagination.Total < pagination2.Total {
		pagination = pagination2
	}

	return lastActivity, pagination, nil
}

func (svc ProductBarcodeHttpService) saveMasterSync(shopID string) {
	if svc.cacheRepo != nil {
		err := svc.cacheRepo.Save(shopID)

		if err != nil {
			fmt.Println("save member master cache error :: " + err.Error())
		}
	}
}
