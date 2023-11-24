package services

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/productbarcode/models"
	"smlcloudplatform/pkg/product/productbarcode/repositories"
	productcategory_models "smlcloudplatform/pkg/product/productcategory/models"
	productcategory_services "smlcloudplatform/pkg/product/productcategory/services"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"strconv"
	"time"

	"github.com/samber/lo"
	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IProductBarcodeHttpService interface {
	CreateProductBarcode(shopID string, authUsername string, doc models.ProductBarcodeRequest) (string, error)
	UpdateProductBarcode(shopID string, guid string, authUsername string, doc models.ProductBarcodeRequest) error
	DeleteProductBarcode(shopID string, guid string, authUsername string) error
	DeleteProductBarcodeByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoProductBarcode(shopID string, guid string) (models.ProductBarcodeInfo, error)
	InfoProductBarcodeByBarcode(shopID string, barcode string) (models.ProductBarcodeInfo, error)
	InfoWTFArray(shopID string, codes []string) ([]interface{}, error)
	InfoWTFArrayMaster(codes []string) ([]interface{}, error)
	GetProductBarcodeByBarcodeRef(shopID string, barcodeRef string) ([]models.ProductBarcodeInfo, error)
	SearchProductBarcode(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error)
	SearchProductBarcode2(shopID string, pageable micromodels.Pageable) ([]models.ProductBarcodeSearch, common.Pagination, error)

	SearchProductBarcodeStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ProductBarcodeInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.ProductBarcode) (common.BulkImport, error)

	XSortsSave(shopID string, authUsername string, xsorts []common.XSortModifyReqesut) error
	GetProductBarcodeByBarcodes(shopID string, barcodes []string) ([]models.ProductBarcodeInfo, error)

	GetModuleName() string
	GetProductBarcodeByUnits(shopID string, unitCodes []string, pageable micromodels.Pageable) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error)
	GetProductBarcodeByGroups(shopID string, groupCodes []string, pageable micromodels.Pageable) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error)
	Export(shopID string, languageCode string) ([][]string, error)
}

type ProductBarcodeHttpService struct {
	repo          repositories.IProductBarcodeRepository
	chRepo        repositories.IProductBarcodeClickhouseRepository
	syncCacheRepo mastersync.IMasterSyncCacheRepository
	mqRepo        repositories.IProductBarcodeMessageQueueRepository
	categorySvc   productcategory_services.IProductCategoryHttpService
	services.ActivityService[models.ProductBarcodeActivity, models.ProductBarcodeDeleteActivity]
	contextTimeout time.Duration
}

func NewProductBarcodeHttpService(
	repo repositories.IProductBarcodeRepository,
	mqRepo repositories.IProductBarcodeMessageQueueRepository,
	chRepo repositories.IProductBarcodeClickhouseRepository,
	categorySvc productcategory_services.IProductCategoryHttpService,
	syncCacheRepo mastersync.IMasterSyncCacheRepository,
) *ProductBarcodeHttpService {

	contextTimeout := time.Duration(15) * time.Second

	insSvc := &ProductBarcodeHttpService{
		repo:           repo,
		chRepo:         chRepo,
		syncCacheRepo:  syncCacheRepo,
		mqRepo:         mqRepo,
		categorySvc:    categorySvc,
		contextTimeout: contextTimeout,
	}
	insSvc.ActivityService = services.NewActivityService[models.ProductBarcodeActivity, models.ProductBarcodeDeleteActivity](repo)
	return insSvc
}

func (svc ProductBarcodeHttpService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc ProductBarcodeHttpService) CreateProductBarcode(shopID string, authUsername string, docReq models.ProductBarcodeRequest) (string, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "barcode", docReq.Barcode)

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
	docData.ProductBarcode = docReq.ToProductBarcode()

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	if docReq.Options != nil {
		options := *docReq.Options
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

	docData.RefBarcodes, err = svc.prepareRefBarcode(ctx, shopID, docReq.RefBarcodes)

	if err != nil {
		return "", err
	}

	docData.BOM, err = svc.prepareBOM(ctx, shopID, docReq.BOM)

	if err != nil {
		return "", err
	}

	_, err = svc.repo.Create(ctx, docData)

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

func (svc ProductBarcodeHttpService) UpdateProductBarcode(shopID string, guid string, authUsername string, docReq models.ProductBarcodeRequest) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	docData := findDoc

	docData.ProductBarcode = docReq.ToProductBarcode()

	docData.Barcode = findDoc.Barcode
	docData.ItemCode = findDoc.ItemCode

	docData.UpdatedBy = authUsername
	docData.UpdatedAt = time.Now()

	docData.RefBarcodes, err = svc.prepareRefBarcode(ctx, shopID, docReq.RefBarcodes)

	if err != nil {
		return err
	}

	docData.BOM, err = svc.prepareBOM(ctx, shopID, docReq.BOM)

	if err != nil {
		return err
	}

	err = svc.updateMetaInRefBarcode(ctx, shopID, docData)

	if err != nil {
		return err
	}

	err = svc.updateMetaInBOMBarcode(ctx, shopID, docData)

	if err != nil {
		return err
	}

	err = svc.repo.Update(ctx, shopID, guid, docData)

	if err != nil {
		return err
	}

	err = svc.mqRepo.Update(docData)
	if err != nil {
		return err
	}

	categoryBarcode := productcategory_models.CodeXSort{
		Barcode:          docData.Barcode,
		Code:             docData.ItemCode,
		Names:            docData.Names,
		UnitCode:         docData.ItemUnitCode,
		UnitNames:        docData.ItemUnitNames,
		ManufacturerGUID: docData.ManufacturerGUID,
	}

	svc.categorySvc.UpdateBarcode(shopID, categoryBarcode)

	svc.saveMasterSync(shopID)

	return nil
}

func (svc ProductBarcodeHttpService) prepareRefBarcode(ctx context.Context, shopID string, barcodes []models.BarcodeRequest) (*[]models.RefProductBarcode, error) {
	tempChildrenBarcodes := []string{}
	tempRefBarcodes := map[string]models.BarcodeRequest{}

	for _, item := range barcodes {
		tempChildrenBarcodes = append(tempChildrenBarcodes, item.Barcode)
		tempRefBarcodes[item.Barcode] = item
	}

	findChildrenDocs, err := svc.repo.FindByDocIndentityGuids(ctx, shopID, "barcode", tempChildrenBarcodes)

	if err != nil {
		return &[]models.RefProductBarcode{}, err
	}

	tempBarcodes := []models.RefProductBarcode{}
	for _, childDoc := range findChildrenDocs {
		tempRef := childDoc.ToRefBarcode()

		tempRef.Condition = tempRefBarcodes[tempRef.Barcode].Condition
		tempRef.StandValue = tempRefBarcodes[tempRef.Barcode].StandValue
		tempRef.DivideValue = tempRefBarcodes[tempRef.Barcode].DivideValue
		tempRef.Qty = tempRefBarcodes[tempRef.Barcode].Qty

		tempBarcodes = append(tempBarcodes, tempRef)
	}

	return &tempBarcodes, nil
}

func (svc ProductBarcodeHttpService) prepareBOM(ctx context.Context, shopID string, barcodes []models.BOMRequest) (*[]models.BOMProductBarcode, error) {
	tempChildrenBarcodes := []string{}
	tempBOM := map[string]models.BOMRequest{}

	for _, item := range barcodes {
		tempChildrenBarcodes = append(tempChildrenBarcodes, item.Barcode)
		tempBOM[item.Barcode] = item
	}

	findChildrenDocs, err := svc.repo.FindByDocIndentityGuids(ctx, shopID, "barcode", tempChildrenBarcodes)

	if err != nil {
		return &[]models.BOMProductBarcode{}, err
	}

	tempBarcodes := []models.BOMProductBarcode{}
	for _, childDoc := range findChildrenDocs {
		temp := childDoc.ToBOM()

		temp.Condition = tempBOM[temp.Barcode].Condition
		temp.StandValue = tempBOM[temp.Barcode].StandValue
		temp.DivideValue = tempBOM[temp.Barcode].DivideValue
		temp.Qty = tempBOM[temp.Barcode].Qty

		tempBarcodes = append(tempBarcodes, temp)
	}

	return &tempBarcodes, nil
}

func (svc ProductBarcodeHttpService) updateMetaInRefBarcode(ctx context.Context, shopID string, docData models.ProductBarcodeDoc) error {

	findDocs, err := svc.repo.FindByRefBarcode(ctx, shopID, docData.Barcode)
	if err != nil {
		return err
	}

	for _, findDoc := range findDocs {
		tempRefBarcodes := []models.RefProductBarcode{}
		for _, refBarcode := range *findDoc.RefBarcodes {
			if refBarcode.Barcode == docData.Barcode {
				refBarcode.Names = docData.Names
				refBarcode.ItemUnitCode = docData.ItemUnitCode
				refBarcode.ItemUnitNames = docData.ItemUnitNames
			}

			tempRefBarcodes = append(tempRefBarcodes, refBarcode)
		}

		findDoc.RefBarcodes = &tempRefBarcodes

		err = svc.repo.Update(ctx, shopID, findDoc.GuidFixed, findDoc)
		if err != nil {
			return err
		}
	}

	return nil
}

func (svc ProductBarcodeHttpService) updateMetaInBOMBarcode(ctx context.Context, shopID string, docData models.ProductBarcodeDoc) error {

	findDocs, err := svc.repo.FindByBOMBarcode(ctx, shopID, docData.Barcode)
	if err != nil {
		return err
	}

	for _, findDoc := range findDocs {
		tempBOMBarcodes := []models.BOMProductBarcode{}
		for _, refBarcode := range *findDoc.BOM {
			if refBarcode.Barcode == docData.Barcode {
				refBarcode.Names = docData.Names
				refBarcode.ItemUnitCode = docData.ItemUnitCode
				refBarcode.ItemUnitNames = docData.ItemUnitNames
			}

			tempBOMBarcodes = append(tempBOMBarcodes, refBarcode)
		}

		findDoc.BOM = &tempBOMBarcodes

		err = svc.repo.Update(ctx, shopID, findDoc.GuidFixed, findDoc)
		if err != nil {
			return err
		}
	}

	return nil
}

func (svc ProductBarcodeHttpService) DeleteProductBarcode(shopID string, guid string, authUsername string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	countRef, err := svc.repo.CountByRefBarcode(ctx, shopID, findDoc.Barcode)

	if err != nil {
		return err
	}

	if countRef > 0 {
		return errors.New("document has refenced")
	}

	countBOM, err := svc.repo.CountByBOM(ctx, shopID, findDoc.Barcode)

	if err != nil {
		return err
	}

	if countBOM > 0 {
		return errors.New("document has refenced")
	}

	err = svc.repo.DeleteByGuidfixed(ctx, shopID, guid, authUsername)
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

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return models.ProductBarcodeInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.ProductBarcodeInfo{}, errors.New("document not found")
	}

	return findDoc.ProductBarcodeInfo, nil
}

func (svc ProductBarcodeHttpService) InfoProductBarcodeByBarcode(shopID string, barcode string) (models.ProductBarcodeInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "barcode", barcode)

	if err != nil {
		return models.ProductBarcodeInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.ProductBarcodeInfo{}, errors.New("document not found")
	}

	return findDoc.ProductBarcodeInfo, nil
}

func (svc ProductBarcodeHttpService) InfoWTFArray(shopID string, codes []string) ([]interface{}, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	docList := []interface{}{}

	for _, code := range codes {
		findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "barcode", code)
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

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	docList := []interface{}{}

	findDocList, err := svc.repo.FindMasterInCodes(ctx, codes)

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

func (svc ProductBarcodeHttpService) GetProductBarcodeByBarcodeRef(shopID string, barcodeRef string) ([]models.ProductBarcodeInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDocs, err := svc.repo.FindByRefBarcode(ctx, shopID, barcodeRef)
	if err != nil {
		return nil, err
	}

	resultDocs := []models.ProductBarcodeInfo{}
	for _, findDoc := range findDocs {
		resultDocs = append(resultDocs, findDoc.ProductBarcodeInfo)
	}

	return resultDocs, nil
}

func (svc ProductBarcodeHttpService) GetProductBarcodeByBarcodes(shopID string, barcodes []string) ([]models.ProductBarcodeInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	results, err := svc.repo.FindByBarcodes(ctx, shopID, barcodes)
	if err != nil {
		return nil, err
	}

	tempResults := map[string]models.ProductBarcodeInfo{}
	for _, result := range results {
		tempResults[result.Barcode] = result
	}

	resultDocs := []models.ProductBarcodeInfo{}
	for _, barcode := range barcodes {

		temp, ok := tempResults[barcode]
		if ok {
			resultDocs = append(resultDocs, temp)
		}
	}

	return resultDocs, nil

}

func (svc ProductBarcodeHttpService) SearchProductBarcode(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"barcode",
		"names.name",
		"itemcode",
		"groupcode",
		"groupnames.name",
		"itemunitnames.name",
	}

	isalacarte, ok := filters["isalacarte"]
	if ok {
		if !isalacarte.(bool) {
			delete(filters, "isalacarte")
			filters["$or"] = []bson.M{
				{"isalacarte": false},
				{"isalacarte": bson.M{"$exists": false}},
			}
		}

	}

	docList, pagination, err := svc.repo.FindPageFilter(ctx, shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.ProductBarcodeInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc ProductBarcodeHttpService) SearchProductBarcode2(shopID string, pageable micromodels.Pageable) ([]models.ProductBarcodeSearch, common.Pagination, error) {

	//fixed shopID
	shopID = "2Eh6e3pfWvXTp0yV3CyFEhKPjdI"
	docList, pagination, err := svc.chRepo.Search(shopID, pageable)

	if err != nil {
		return []models.ProductBarcodeSearch{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc ProductBarcodeHttpService) SearchProductBarcodeStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ProductBarcodeInfo, int, error) {
	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"barcode",
		"names.name",
		"itemcode",
		"groupcode",
		"groupnames.name",
		"itemunitnames.name",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(ctx, shopID, filters, searchInFields, selectFields, pageableStep)
	if err != nil {
		return []models.ProductBarcodeInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc ProductBarcodeHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.ProductBarcode) (common.BulkImport, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.ProductBarcode](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.Barcode)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(ctx, shopID, "barcode", itemCodeGuidList)

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
			return svc.repo.FindByDocIndentityGuid(ctx, shopID, "barcode", guid)
		},
		func(doc models.ProductBarcodeDoc) bool {
			return doc.Barcode != ""
		},
		func(shopID string, authUsername string, dataReq models.ProductBarcode, doc models.ProductBarcodeDoc) error {

			docReq := models.ProductBarcodeRequest{}
			docReq.ProductBarcodeBase = dataReq.ProductBarcodeBase

			tempBarcodes := []models.BarcodeRequest{}

			for _, docBarcode := range *dataReq.RefBarcodes {
				tempBarcodes = append(tempBarcodes, models.BarcodeRequest{
					Barcode:     docBarcode.Barcode,
					Condition:   docBarcode.Condition,
					StandValue:  docBarcode.StandValue,
					DivideValue: docBarcode.DivideValue,
				})
			}
			docReq.RefBarcodes = tempBarcodes

			svc.UpdateProductBarcode(shopID, doc.GuidFixed, authUsername, docReq)

			// err = svc.repo.Update(ctx, shopID, doc.GuidFixed, doc)
			// if err != nil {
			// 	return nil
			// }
			return nil
		},
	)

	if len(createDataList) > 0 {
		svc.repo.Transaction(ctx, func(ctx context.Context) error {
			for _, doc := range createDataList {
				docReq := models.ProductBarcodeRequest{}
				docReq.ProductBarcodeBase = doc.ProductBarcodeBase

				tempBarcodes := []models.BarcodeRequest{}

				for _, docBarcode := range *doc.RefBarcodes {
					tempBarcodes = append(tempBarcodes, models.BarcodeRequest{
						Barcode:     docBarcode.Barcode,
						Condition:   docBarcode.Condition,
						StandValue:  docBarcode.StandValue,
						DivideValue: docBarcode.DivideValue,
					})
				}

				docReq.RefBarcodes = tempBarcodes
				_, err = svc.CreateProductBarcode(shopID, authUsername, docReq)

				if err != nil {
					return err
				}
			}

			return nil
		})

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
	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	for _, xsort := range xsorts {
		if len(xsort.GUIDFixed) < 1 {
			continue
		}
		findDoc, err := svc.repo.FindByGuid(ctx, shopID, xsort.GUIDFixed)

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

		err = svc.repo.Update(ctx, shopID, findDoc.GuidFixed, findDoc)

		if err != nil {
			return err
		}

	}

	svc.saveMasterSync(shopID)

	return nil

}

func (svc ProductBarcodeHttpService) DeleteProductBarcodeByGUIDs(shopID string, authUsername string, GUIDs []string) error {
	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	countRefBarcode, err := svc.repo.CountByRefGuids(ctx, shopID, GUIDs)

	if err != nil {
		return err
	}

	if countRefBarcode > 0 {
		return errors.New("document has refenced")
	}

	countBOM, err := svc.repo.CountByBOMGuids(ctx, shopID, GUIDs)

	if err != nil {
		return err
	}

	if countBOM > 0 {
		return errors.New("document has refenced")
	}

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err = svc.repo.Delete(ctx, shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc ProductBarcodeHttpService) GetProductBarcodeByUnits(shopID string, unitCodes []string, pageable micromodels.Pageable) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	if len(unitCodes) < 1 {
		return []models.ProductBarcodeInfo{}, mongopagination.PaginationData{}, nil
	}

	results, pagination, err := svc.repo.FindPageByUnits(ctx, shopID, unitCodes, pageable)

	if err != nil {
		return []models.ProductBarcodeInfo{}, mongopagination.PaginationData{}, err
	}

	return results, pagination, nil
}

func (svc ProductBarcodeHttpService) GetProductBarcodeByGroups(shopID string, groupCodes []string, pageable micromodels.Pageable) ([]models.ProductBarcodeInfo, mongopagination.PaginationData, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	if len(groupCodes) < 1 {
		return []models.ProductBarcodeInfo{}, mongopagination.PaginationData{}, nil
	}

	results, pagination, err := svc.repo.FindPageByGroups(ctx, shopID, groupCodes, pageable)

	if err != nil {
		return []models.ProductBarcodeInfo{}, mongopagination.PaginationData{}, err
	}

	return results, pagination, nil
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

func (svc ProductBarcodeHttpService) Export(shopID string, languageCode string) ([][]string, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	filters := bson.M{}

	docs, err := svc.repo.Find(ctx, shopID, filters)

	if err != nil {
		return [][]string{}, err
	}

	results := [][]string{
		{
			"รหัสบาร์โค้ด",
			"ชื่อสินค้า",
			"รหัสหน่วยนับ",
			"ชื่อหน่วยนับ",
			"ราคาขาย",
			"ประเภทสินค้า",
			"กลุ่มสินค้า",
			"ชื่อกลุ่มสินค้า",
		},
	}

	temp := prepareDataToCSV(languageCode, docs)
	results = append(results, temp...)

	return results, nil
}

func dataToCSV(data []models.ProductBarcodeDoc, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write your data to the CSV file
	for _, value := range data {
		langCode := "th"

		productName := getName(value.Names, langCode)
		unitName := getName(value.ItemUnitNames, langCode)
		groupName := getName(value.GroupNames, langCode)

		itemType := strconv.Itoa(int(value.ItemType))
		writer.Write([]string{value.Barcode, productName, value.ItemUnitCode, unitName, itemType, value.GroupCode, groupName}) // Adjust fields as per your struct
	}
	return nil
}
func prepareDataToCSV(languageCode string, data []models.ProductBarcodeDoc) [][]string {

	results := [][]string{}

	for _, value := range data {
		langCode := languageCode

		productName := getName(value.Names, langCode)
		unitName := getName(value.ItemUnitNames, langCode)
		groupName := getName(value.GroupNames, langCode)

		price := "0"

		if value.Prices != nil && len(*value.Prices) > 0 {
			for _, priceItem := range *value.Prices {
				if priceItem.KeyNumber == 0 {
					price = fmt.Sprintf("%.2f", priceItem.Price)
					break
				}
			}
		}

		itemType := strconv.Itoa(int(value.ItemType))
		results = append(results, []string{value.Barcode, productName, value.ItemUnitCode, unitName, price, itemType, value.GroupCode, groupName})
	}

	return results
}

func getName(names *[]common.NameX, langCode string) string {
	if names == nil {
		return ""
	}

	for _, name := range *names {
		if *name.Code == langCode {
			return *name.Name
		}
	}

	return ""
}
