package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/config"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/unit/models"
	"smlcloudplatform/pkg/product/unit/repositories"
	"smlcloudplatform/pkg/requestapi"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"strings"
	"time"

	"github.com/samber/lo"
	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUnitHttpService interface {
	CreateUnit(shopID string, authUsername string, doc models.Unit) (string, error)
	UpdateUnit(shopID string, guid string, authUsername string, doc models.Unit) error
	UpdateFieldUnit(shopID string, guid string, authUsername string, doc models.Unit) error
	DeleteUnit(shopID string, guid string, authHeader string, authUsername string) error
	DeleteUnitByGUIDs(shopID string, authHeader string, authUsername string, GUIDs []string) error
	InfoUnit(shopID string, guid string) (models.UnitInfo, error)
	InfoUnitWTFArray(shopID string, unitCodes []string) ([]interface{}, error)
	InfoWTFArrayMaster(codes []string) ([]interface{}, error)
	SearchUnit(shopID string, codeFilters []string, pageable micromodels.Pageable) ([]models.UnitInfo, mongopagination.PaginationData, error)
	SearchUnitLimit(shopID string, langCode string, codeFilters []string, pageableStep micromodels.PageableStep) ([]models.UnitInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Unit) (common.BulkImport, error)

	GetModuleName() string
}

type UnitHttpService struct {
	repo              repositories.IUnitRepository
	syncCacheRepo     mastersync.IMasterSyncCacheRepository
	unitServiceConfig config.IUnitServiceConfig

	services.ActivityService[models.UnitActivity, models.UnitDeleteActivity]
}

func NewUnitHttpService(repo repositories.IUnitRepository, unitServiceConfig config.IUnitServiceConfig, syncCacheRepo mastersync.IMasterSyncCacheRepository) *UnitHttpService {

	insSvc := &UnitHttpService{
		repo:              repo,
		unitServiceConfig: unitServiceConfig,
		syncCacheRepo:     syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.UnitActivity, models.UnitDeleteActivity](repo)
	return insSvc
}

func (svc UnitHttpService) CreateUnit(shopID string, authUsername string, doc models.Unit) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "unitcode", doc.UnitCode)

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

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc UnitHttpService) UpdateUnit(shopID string, guid string, authUsername string, doc models.Unit) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

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

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc UnitHttpService) UpdateFieldUnit(shopID string, guid string, authUsername string, doc models.Unit) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

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

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc UnitHttpService) existsUnitRefInProduct(authHeader, unitCode string) (bool, error) {
	products, err := svc.getProductByUnit(authHeader, []string{unitCode})

	if err != nil {
		return true, fmt.Errorf("error check unit ref product: %s", err.Error())
	}

	if len(products) > 0 {
		return true, fmt.Errorf("unit code %s is ref by product", unitCode)
	}

	return false, nil
}

func (svc UnitHttpService) deleteByUnitCode(shopID, guid, authHeader, authUsername string) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	existsInProduct, _ := svc.existsUnitRefInProduct(authHeader, findDoc.UnitCode)

	if existsInProduct {
		return fmt.Errorf("unit code \"%s\" is referenced", findDoc.UnitCode)
	}

	err = svc.repo.DeleteByGuidfixed(shopID, guid, authUsername)
	if err != nil {
		return err
	}

	return nil
}

func (svc UnitHttpService) DeleteUnit(shopID, guid, authHeader, authUsername string) error {

	err := svc.deleteByUnitCode(shopID, guid, authHeader, authUsername)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc UnitHttpService) DeleteUnitByGUIDs(shopID, authHeader, authUsername string, GUIDs []string) error {

	// deleteFilterQuery := map[string]interface{}{
	// 	"guidfixed": bson.M{"$in": GUIDs},
	// }

	// err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	// if err != nil {
	// 	return err
	// }

	for idx, guid := range GUIDs {
		if idx == 1 {
			return errors.New("test rollback")
		}

		err := svc.DeleteUnit(shopID, guid, authHeader, authUsername)

		if err != nil {
			return err
		}
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc UnitHttpService) InfoUnit(shopID string, guid string) (models.UnitInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.UnitInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.UnitInfo{}, errors.New("document not found")
	}

	return findDoc.UnitInfo, nil

}

func (svc UnitHttpService) InfoUnitWTFArray(shopID string, unitCodes []string) ([]interface{}, error) {

	docList := []interface{}{}

	for _, unitCode := range unitCodes {
		findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "unitcode", unitCode)
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
	docList := []interface{}{}

	findDocList, err := svc.repo.FindMasterInCodes(codes)

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
	searchInFields := []string{
		"unitcode",
		"names.name",
	}

	filters := map[string]interface{}{}
	if len(codeFilters) > 0 {
		filters["unitcode"] = bson.M{"$in": codeFilters}
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.UnitInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc UnitHttpService) SearchUnitLimit(shopID string, langCode string, codeFilters []string, pageableStep micromodels.PageableStep) ([]models.UnitInfo, int, error) {
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

	docList, total, err := svc.repo.FindStep(shopID, filters, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.UnitInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc UnitHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.Unit) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Unit](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.UnitCode)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "unitcode", itemCodeGuidList)

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
			return svc.repo.FindByDocIndentityGuid(shopID, "unitcode", guid)
		},
		func(doc models.UnitDoc) bool {
			return doc.UnitCode != ""
		},
		func(shopID string, authUsername string, data models.Unit, doc models.UnitDoc) error {

			doc.Unit = data
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

func (svc UnitHttpService) getProductByUnit(authHeader string, unitCodes []string) ([]interface{}, error) {

	reqCodes := strings.Join(unitCodes, ",")

	url := fmt.Sprintf("%s/product/barcode/units?codes=%s", svc.unitServiceConfig.ProductHost(), reqCodes)

	return requestapi.Get(url, authHeader)
}
