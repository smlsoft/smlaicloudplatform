package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/ordertype/models"
	"smlcloudplatform/pkg/product/ordertype/repositories"
	productbarcode_repositories "smlcloudplatform/pkg/product/productbarcode/repositories"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IOrderTypeHttpService interface {
	CreateOrderType(shopID string, authUsername string, doc models.OrderType) (string, error)
	UpdateOrderType(shopID string, guid string, authUsername string, doc models.OrderType) error
	DeleteOrderType(shopID string, guid string, authUsername string) error
	DeleteOrderTypeByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoOrderType(shopID string, guid string) (models.OrderTypeInfo, error)
	InfoOrderTypeByCode(shopID string, code string) (models.OrderTypeInfo, error)
	SearchOrderType(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.OrderTypeInfo, mongopagination.PaginationData, error)
	SearchOrderTypeStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.OrderTypeInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.OrderType) (common.BulkImport, error)

	GetModuleName() string
}

type OrderTypeHttpService struct {
	repo               repositories.IOrderTypeRepository
	repoProductBarcode productbarcode_repositories.IProductBarcodeRepository
	syncCacheRepo      mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.OrderTypeActivity, models.OrderTypeDeleteActivity]
}

func NewOrderTypeHttpService(
	repo repositories.IOrderTypeRepository,
	repoProductBarcode productbarcode_repositories.IProductBarcodeRepository,
	syncCacheRepo mastersync.IMasterSyncCacheRepository,
) *OrderTypeHttpService {

	insSvc := &OrderTypeHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.OrderTypeActivity, models.OrderTypeDeleteActivity](repo)

	return insSvc
}

func (svc OrderTypeHttpService) CreateOrderType(shopID string, authUsername string, doc models.OrderType) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", doc.Code)

	if err != nil {
		return "", err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", errors.New("code is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.OrderTypeDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.OrderType = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc OrderTypeHttpService) UpdateOrderType(shopID string, guid string, authUsername string, doc models.OrderType) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	findDoc.OrderType = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc OrderTypeHttpService) DeleteOrderType(shopID string, guid string, authUsername string) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return nil
	}

	existsInProduct, _ := svc.existsOrderTypeRefInProduct(shopID, []string{guid})

	if existsInProduct {
		return fmt.Errorf("\"%s\" is referenced in product barcode", findDoc.Code)
	}

	err = svc.repo.DeleteByGuidfixed(shopID, guid, authUsername)
	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc OrderTypeHttpService) DeleteOrderTypeByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	existsInProduct, _ := svc.existsOrderTypeRefInProduct(shopID, GUIDs)

	if existsInProduct {
		return fmt.Errorf("referenced in product")
	}

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc OrderTypeHttpService) InfoOrderType(shopID string, guid string) (models.OrderTypeInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.OrderTypeInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.OrderTypeInfo{}, errors.New("document not found")
	}

	return findDoc.OrderTypeInfo, nil
}

func (svc OrderTypeHttpService) InfoOrderTypeByCode(shopID string, code string) (models.OrderTypeInfo, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", code)

	if err != nil {
		return models.OrderTypeInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.OrderTypeInfo{}, errors.New("document not found")
	}

	return findDoc.OrderTypeInfo, nil
}

func (svc OrderTypeHttpService) SearchOrderType(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.OrderTypeInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"code",
		"names.name",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.OrderTypeInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc OrderTypeHttpService) SearchOrderTypeStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.OrderTypeInfo, int, error) {
	searchInFields := []string{
		"code",
		"names.name",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.OrderTypeInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc OrderTypeHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.OrderType) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.OrderType](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.OrderType, models.OrderTypeDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.OrderType) models.OrderTypeDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.OrderTypeDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.OrderType = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.OrderType, models.OrderTypeDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.OrderTypeDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "code", guid)
		},
		func(doc models.OrderTypeDoc) bool {
			return doc.Code != ""
		},
		func(shopID string, authUsername string, data models.OrderType, doc models.OrderTypeDoc) error {

			doc.OrderType = data
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

func (svc OrderTypeHttpService) getDocIDKey(doc models.OrderType) string {
	return doc.Code
}

func (svc OrderTypeHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc OrderTypeHttpService) GetModuleName() string {
	return "ordertype"
}

func (svc OrderTypeHttpService) existsOrderTypeRefInProduct(shopID string, GUIDs []string) (bool, error) {
	docCount, err := svc.repoProductBarcode.CountByOrderTypes(shopID, GUIDs)
	if err != nil {
		return true, err
	}

	if docCount > 0 {
		return true, fmt.Errorf("referenced in product barcode")
	}

	return false, nil
}
