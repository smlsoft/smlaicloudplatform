package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/productbarcode/models"
	"smlcloudplatform/pkg/product/productbarcode/repositories"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/samber/lo"
	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IProductBarcodeHttpService interface {
	CreateProductBarcode(shopID string, authUsername string, doc models.ProductBarcode) (string, error)
	UpdateProductBarcode(shopID string, guid string, authUsername string, doc models.ProductBarcode) error
	DeleteProductBarcode(shopID string, guid string, authUsername string) error
	DeleteProductBarcodeByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoProductBarcode(shopID string, guid string) (models.ProductBarcodeInfo, error)
	InfoWTFArray(shopID string, codes []string) ([]interface{}, error)
	InfoWTFArrayMaster(codes []string) ([]interface{}, error)
	SearchProductBarcode(shopID string, pageable micromodels.Pageable) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error)
	SearchProductBarcodeStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.ProductBarcodeInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.ProductBarcode) (common.BulkImport, error)

	XSortsSave(shopID string, authUsername string, xsorts []common.XSortModifyReqesut) error

	GetModuleName() string
}

type ProductBarcodeHttpService struct {
	repo          repositories.IProductBarcodeRepository
	syncCacheRepo mastersync.IMasterSyncCacheRepository
	mqRepo        repositories.IProductBarcodeMessageQueueRepository
	services.ActivityService[models.ProductBarcodeActivity, models.ProductBarcodeDeleteActivity]
}

func NewProductBarcodeHttpService(repo repositories.IProductBarcodeRepository, mqRepo repositories.IProductBarcodeMessageQueueRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *ProductBarcodeHttpService {

	insSvc := &ProductBarcodeHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
		mqRepo:        mqRepo,
	}
	insSvc.ActivityService = services.NewActivityService[models.ProductBarcodeActivity, models.ProductBarcodeDeleteActivity](repo)
	return insSvc
}

func (svc ProductBarcodeHttpService) CreateProductBarcode(shopID string, authUsername string, doc models.ProductBarcode) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "barcode", doc.Barcode)

	if err != nil {
		return "", err
	}

	if findDoc.Barcode != "" {
		return "", errors.New("barcode is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.ProductBarcodeDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.ProductBarcode = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	if doc.Options != nil {
		options := *doc.Options
		for idxOpt := range options {
			option := &options[idxOpt]
			if len(option.GUID) < 1 {
				option.GUID = utils.NewGUID()
			}

			choices := *option.Choices
			for idxChoice := range choices {
				choice := &choices[idxChoice]
				if len(choice.GUID) < 1 {
					choice.GUID = utils.NewGUID()
				}
			}
		}
	}

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	err = svc.mqRepo.Create(docData)
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

	err = svc.mqRepo.Update(findDoc)
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

	err = svc.mqRepo.Delete(findDoc)
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

func (svc ProductBarcodeHttpService) InfoWTFArray(shopID string, codes []string) ([]interface{}, error) {
	docList := []interface{}{}

	for _, code := range codes {
		findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "barcode", code)
		if err != nil || findDoc.ID == primitive.NilObjectID {
			// add item empty
			docList = append(docList, nil)
		} else {
			docList = append(docList, findDoc.ProductBarcodeInfo)
		}
	}

	return docList, nil
}

func (svc ProductBarcodeHttpService) InfoWTFArrayMaster(codes []string) ([]interface{}, error) {
	docList := []interface{}{}

	findDocList, err := svc.repo.FindMasterInCodes(codes)

	if err != nil {
		return []interface{}{}, err
	}

	for _, code := range codes {
		findDoc, ok := lo.Find(findDocList, func(item models.ProductBarcodeInfo) bool {
			return item.Barcode == code
		})
		if !ok {
			// add item empty
			docList = append(docList, nil)
		} else {
			docList = append(docList, findDoc)
		}
	}

	return docList, nil
}

func (svc ProductBarcodeHttpService) SearchProductBarcode(shopID string, pageable micromodels.Pageable) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"barcode",
		"names.name",
	}

	docList, pagination, err := svc.repo.FindPage(shopID, searchInFields, pageable)

	if err != nil {
		return []models.ProductBarcodeInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc ProductBarcodeHttpService) SearchProductBarcodeStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.ProductBarcodeInfo, int, error) {
	searchInFields := []string{
		"barcode",
		"names.name",
	}

	selectFields := map[string]interface{}{
		"guidfixed":       1,
		"barcode":         1,
		"groupcode":       1,
		"groupname":       1,
		"itemcode":        1,
		"itemunitcode":    1,
		"itemunitnames":   1,
		"prices":          1,
		"imageuri":        1,
		"options":         1,
		"xsorts":          1,
		"images":          1,
		"useimageorcolor": 1,
		"colorselect":     1,
		"colorselecthex":  1,
		"branches":        1,
	}

	if langCode != "" {
		selectFields["names"] = bson.M{"$elemMatch": bson.M{"code": langCode}}
		selectFields["itemunitnames"] = bson.M{"$elemMatch": bson.M{"code": langCode}}
	} else {
		selectFields["names"] = 1
		selectFields["itemunitnames"] = 1
	}

	docList, total, err := svc.repo.FindStep(shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

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

	if len(createDataList) > 0 {
		err = svc.mqRepo.CreateInBatch(createDataList)
		if err != nil {
			return common.BulkImport{}, err
		}
	}

	if len(updateSuccessDataList) > 0 {
		err = svc.mqRepo.UpdateInBatch(updateSuccessDataList)
		if err != nil {
			return common.BulkImport{}, err
		}
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

func (svc ProductBarcodeHttpService) XSortsSave(shopID string, authUsername string, xsorts []common.XSortModifyReqesut) error {
	for _, xsort := range xsorts {
		if len(xsort.GUIDFixed) < 1 {
			continue
		}
		findDoc, err := svc.repo.FindByGuid(shopID, xsort.GUIDFixed)

		if err != nil {
			return err
		}

		if len(findDoc.GuidFixed) < 1 {
			continue
		}

		if findDoc.XSorts == nil {
			findDoc.XSorts = &[]common.XSort{}
		}

		dictXSorts := map[string]common.XSort{}

		for _, tempXSort := range *findDoc.XSorts {
			dictXSorts[tempXSort.Code] = tempXSort
		}

		dictXSorts[xsort.Code] = common.XSort{
			Code:   xsort.Code,
			XOrder: xsort.XOrder,
		}

		tempXSorts := []common.XSort{}

		for _, tempXSort := range dictXSorts {
			tempXSorts = append(tempXSorts, tempXSort)
		}

		findDoc.XSorts = &tempXSorts
		findDoc.UpdatedBy = authUsername
		findDoc.UpdatedAt = time.Now()

		err = svc.repo.Update(shopID, findDoc.GuidFixed, findDoc)

		if err != nil {
			return err
		}

	}

	svc.saveMasterSync(shopID)

	return nil

}

func (svc ProductBarcodeHttpService) DeleteProductBarcodeByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc ProductBarcodeHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc ProductBarcodeHttpService) GetModuleName() string {
	return "productbarcode"
}
