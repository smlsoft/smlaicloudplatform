package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"smlaicloudplatform/internal/logger"
	mastersync "smlaicloudplatform/internal/mastersync/repositories"
	common "smlaicloudplatform/internal/models"
	productbarcode_repositories "smlaicloudplatform/internal/product/productbarcode/repositories"
	"smlaicloudplatform/internal/product/unit/models"

	"smlaicloudplatform/internal/product/unit/repositories"
	"smlaicloudplatform/internal/services"
	"smlaicloudplatform/internal/utils"
	"smlaicloudplatform/internal/utils/importdata"
	micromodels "smlaicloudplatform/pkg/microservice/models"
	"time"

	"github.com/samber/lo"
	"github.com/smlsoft/mongopagination"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUnitHttpService interface {
	CreateUnit(shopID string, authUsername string, doc models.Unit) (string, error)
	UpdateUnit(shopID string, guid string, authUsername string, doc models.Unit) error
	UpdateFieldUnit(shopID string, guid string, authUsername string, doc models.Unit) error
	DeleteUnit(shopID string, guid string, authUsername string) error
	DeleteUnitByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoUnit(shopID string, guid string) (models.UnitInfo, error)
	InfoUnitWTFArray(shopID string, unitCodes []string) ([]interface{}, error)
	InfoWTFArrayMaster(codes []string) ([]interface{}, error)
	SearchUnit(shopID string, codeFilters []string, pageable micromodels.Pageable) ([]models.UnitInfo, mongopagination.PaginationData, error)
	SearchUnitLimit(shopID string, langCode string, codeFilters []string, pageableStep micromodels.PageableStep) ([]models.UnitInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Unit) (common.BulkImport, error)
	GetModuleName() string
	ImportUnitsFromFile(file []byte, shopID string, authUsername string) (string, error)
}

type UnitHttpService struct {
	repo               repositories.IUnitRepository
	repoProductBarcode productbarcode_repositories.IProductBarcodeRepository
	syncCacheRepo      mastersync.IMasterSyncCacheRepository
	repoMessageQueue   repositories.IUnitMessageQueueRepository

	services.ActivityService[models.UnitActivity, models.UnitDeleteActivity]
	contextTimeout time.Duration
}

func NewUnitHttpService(
	repo repositories.IUnitRepository,
	repoProductBarcode productbarcode_repositories.IProductBarcodeRepository,
	repoMessageQueue repositories.IUnitMessageQueueRepository,
	syncCacheRepo mastersync.IMasterSyncCacheRepository,
) *UnitHttpService {

	contextTimeout := time.Duration(15) * time.Second

	insSvc := &UnitHttpService{
		repo:               repo,
		repoProductBarcode: repoProductBarcode,
		repoMessageQueue:   repoMessageQueue,
		syncCacheRepo:      syncCacheRepo,
		contextTimeout:     contextTimeout,
	}

	insSvc.ActivityService = services.NewActivityService[models.UnitActivity, models.UnitDeleteActivity](repo)
	return insSvc
}

func (svc UnitHttpService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc *UnitHttpService) ImportUnitsFromFile(file []byte, shopID string, authUsername string) (string, error) {
	f, err := excelize.OpenReader(bytes.NewReader(file))
	if err != nil {
		return "", fmt.Errorf("failed to open Excel file: %v", err)
	}

	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return "", fmt.Errorf("failed to read rows: %v", err)
	}

	if len(rows) < 2 {
		return "", fmt.Errorf("the Excel file is empty or missing headers")
	}

	headers := rows[0]
	var existingUnitCodes []string
	var newUnits []models.UnitDoc

	ctx, cancel := svc.getContextTimeout()
	defer cancel()

	for _, row := range rows[1:] {
		entry := map[string]string{}
		for i, value := range row {
			if i < len(headers) {
				entry[headers[i]] = value
			}
		}

		unitCode, ok := entry["code"]
		if !ok || unitCode == "" {
			continue
		}

		findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "unitcode", unitCode)
		if err != nil {
			return "", fmt.Errorf("error checking existing unit code %s: %v", unitCode, err)
		}

		if findDoc.UnitCode != "" {
			existingUnitCodes = append(existingUnitCodes, unitCode)
			continue
		}

		var names []common.NameX
		for i, lang := range headers {
			if lang != "code" && i < len(row) {
				name := row[i]
				if name != "" {
					copiedLang := lang
					copiedName := name
					names = append(names, common.NameX{
						Code:     &copiedLang,
						Name:     &copiedName,
						IsAuto:   false,
						IsDelete: false,
					})
				}
			}
		}

		newUnit := models.UnitDoc{
			UnitData: models.UnitData{

				ShopIdentity: common.ShopIdentity{ShopID: shopID},
				UnitInfo: models.UnitInfo{
					DocIdentity: common.DocIdentity{GuidFixed: utils.NewGUID()},
					Unit: models.Unit{
						UnitCode: unitCode,
						Names:    &names,
						UnitName: common.UnitName{UnitName1: "Import by " + authUsername},
					},
				},
			},
			ActivityDoc: common.ActivityDoc{
				CreatedBy: "Import by " + authUsername,
				CreatedAt: time.Now(),
			},
		}
		newUnits = append(newUnits, newUnit)
	}

	if len(newUnits) > 0 {
		err = svc.repo.CreateInBatch(ctx, newUnits)
		if err != nil {
			return "", fmt.Errorf("failed to insert new units: %v", err)
		}
	}

	return fmt.Sprintf("Existing unit codes: %v", existingUnitCodes), nil
}

func (svc UnitHttpService) CreateUnit(shopID string, authUsername string, doc models.Unit) (string, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "unitcode", doc.UnitCode)

	if err != nil {
		return "", err
	}

	if findDoc.UnitCode != "" {
		return "", errors.New("unit code is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.UnitDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Unit = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(ctx, docData)

	if err != nil {
		return "", err
	}

	go func() {
		svc.repoMessageQueue.Create(docData)
		svc.saveMasterSync(shopID)
	}()
	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc UnitHttpService) UpdateUnit(shopID string, guid string, authUsername string, doc models.Unit) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	tempCode := findDoc.UnitCode

	findDoc.Unit = doc

	//
	findDoc.UnitCode = tempCode

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(ctx, shopID, guid, findDoc)

	if err != nil {
		return err
	}

	go func() {
		err := svc.repoMessageQueue.Update(findDoc)

		if err != nil {
			logger.GetLogger().Error(err)
		}

		svc.saveMasterSync(shopID)
	}()

	return nil
}

func (svc UnitHttpService) UpdateFieldUnit(shopID string, guid string, authUsername string, doc models.Unit) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	temp := map[string]common.NameX{}

	for _, v := range *findDoc.Names {
		temp[*v.Code] = v
	}

	for _, v := range *doc.Names {
		temp[*v.Code] = v
	}

	tempNames := []common.NameX{}

	for _, v := range temp {
		tempNames = append(tempNames, v)
	}

	lo.Filter[common.NameX](tempNames, func(n common.NameX, i int) bool {
		notDelete := !n.IsDelete
		return notDelete
	})

	findDoc.Unit.Names = &tempNames

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(ctx, shopID, guid, findDoc)

	if err != nil {
		return err
	}

	go func() {

		err := svc.repoMessageQueue.Update(findDoc)

		if err != nil {
			logger.GetLogger().Error(err)
		}

		svc.saveMasterSync(shopID)
	}()

	return nil
}

func (svc UnitHttpService) existsUnitRefInProduct(shopID, unitCode string) (bool, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	docCount, err := svc.repoProductBarcode.CountByUnitCodes(ctx, shopID, []string{unitCode})

	if err != nil {
		return true, err
	}

	if docCount > 0 {
		return true, fmt.Errorf("unit code %s is referenced by product", unitCode)
	}

	return false, nil
}

func (svc UnitHttpService) deleteByUnitCode(shopID, guid, authUsername string) (models.UnitDoc, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return findDoc, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return findDoc, nil
	}

	existsInProduct, err := svc.existsUnitRefInProduct(shopID, findDoc.UnitCode)

	if existsInProduct {
		return findDoc, err
	}

	err = svc.repo.DeleteByGuidfixed(ctx, shopID, guid, authUsername)
	if err != nil {
		return findDoc, err
	}

	return findDoc, nil
}

func (svc UnitHttpService) DeleteUnit(shopID, guid, authUsername string) error {

	doc, err := svc.deleteByUnitCode(shopID, guid, authUsername)

	if err != nil {
		return err
	}

	go func() {
		svc.repoMessageQueue.Delete(doc)
		svc.saveMasterSync(shopID)
	}()

	return nil
}

func (svc UnitHttpService) DeleteUnitByGUIDs(shopID, authUsername string, GUIDs []string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDocs, err := svc.repo.FindByGuids(ctx, shopID, GUIDs)

	if err != nil {
		return err
	}

	if len(findDocs) == 0 {
		return nil
	}

	unitCodes := []string{}
	for _, v := range findDocs {
		unitCodes = append(unitCodes, v.UnitCode)
	}

	docCount, err := svc.repoProductBarcode.CountByUnitCodes(ctx, shopID, unitCodes)

	if err != nil {
		return err
	}

	if docCount > 0 {
		return fmt.Errorf("unit code is referenced by product")
	}

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err = svc.repo.Delete(ctx, shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	go func() {
		err := svc.repoMessageQueue.DeleteInBatch(findDocs)

		if err != nil {
			logger.GetLogger().Error(err)
		}

		svc.saveMasterSync(shopID)
	}()

	return nil
}

func (svc UnitHttpService) InfoUnit(shopID string, guid string) (models.UnitInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return models.UnitInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.UnitInfo{}, errors.New("document not found")
	}

	return findDoc.UnitInfo, nil

}

func (svc UnitHttpService) InfoUnitWTFArray(shopID string, unitCodes []string) ([]interface{}, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	docList := []interface{}{}

	for _, unitCode := range unitCodes {
		findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "unitcode", unitCode)
		if err != nil || findDoc.ID == primitive.NilObjectID {
			// add item empty
			emptyDoc := models.UnitInfo{}
			emptyDoc.UnitCode = unitCode
			docList = append(docList, nil)
		} else {
			docList = append(docList, findDoc.UnitInfo)
		}
	}

	return docList, nil
}

func (svc UnitHttpService) InfoWTFArrayMaster(codes []string) ([]interface{}, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	docList := []interface{}{}

	findDocList, err := svc.repo.FindMasterInCodes(ctx, codes)

	if err != nil {
		return []interface{}{}, err
	}

	for _, code := range codes {
		findDoc, ok := lo.Find(findDocList, func(item models.UnitInfo) bool {
			return item.UnitCode == code
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

func (svc UnitHttpService) SearchUnit(shopID string, codeFilters []string, pageable micromodels.Pageable) ([]models.UnitInfo, mongopagination.PaginationData, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"unitcode",
		"names.name",
	}

	filters := map[string]interface{}{}
	if len(codeFilters) > 0 {
		filters["unitcode"] = bson.M{"$in": codeFilters}
	}

	docList, pagination, err := svc.repo.FindPageFilter(ctx, shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.UnitInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc UnitHttpService) SearchUnitLimit(shopID string, langCode string, codeFilters []string, pageableStep micromodels.PageableStep) ([]models.UnitInfo, int, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"unitcode",
		"names.name",
	}

	selectFields := map[string]interface{}{
		"guidfixed": 1,
		"unitcode":  1,
	}

	if langCode != "" {
		selectFields["names"] = bson.M{"$elemMatch": bson.M{"code": langCode}}
	} else {
		selectFields["names"] = 1
	}

	filters := map[string]interface{}{}
	if len(codeFilters) > 0 {

		filters["$or"] = []interface{}{
			bson.M{"unitcode": bson.M{"$in": codeFilters}},
		}
	}

	docList, total, err := svc.repo.FindStep(ctx, shopID, filters, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.UnitInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc UnitHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.Unit) (common.BulkImport, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Unit](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.UnitCode)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(ctx, shopID, "unitcode", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.UnitCode)
	}

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Unit, models.UnitDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Unit) models.UnitDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.UnitDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Unit = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Unit, models.UnitDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.UnitDoc, error) {
			return svc.repo.FindByDocIndentityGuid(ctx, shopID, "unitcode", guid)
		},
		func(doc models.UnitDoc) bool {
			return doc.UnitCode != ""
		},
		func(shopID string, authUsername string, data models.Unit, doc models.UnitDoc) error {

			doc.Unit = data
			doc.UpdatedBy = authUsername
			doc.UpdatedAt = time.Now()

			err = svc.repo.Update(ctx, shopID, doc.GuidFixed, doc)
			if err != nil {
				return nil
			}
			return nil
		},
	)

	if len(createDataList) > 0 {
		err = svc.repo.CreateInBatch(ctx, createDataList)

		if err != nil {
			return common.BulkImport{}, err
		}

	}

	createDataKey := []string{}

	for _, doc := range createDataList {
		createDataKey = append(createDataKey, doc.UnitCode)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.UnitCode)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {

		updateDataKey = append(updateDataKey, doc.UnitCode)
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

func (svc UnitHttpService) getDocIDKey(doc models.Unit) string {
	return doc.UnitCode
}

func (svc UnitHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc UnitHttpService) GetModuleName() string {
	return "productunit"
}
