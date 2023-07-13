package services

import (
	"context"
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/producttype/models"
	"smlcloudplatform/pkg/product/producttype/repositories"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IProductTypeHttpService interface {
	CreateProductType(shopID string, authUsername string, doc models.ProductType) (string, error)
	UpdateProductType(shopID string, guid string, authUsername string, doc models.ProductType) error
	DeleteProductType(shopID string, guid string, authUsername string) error
	DeleteProductTypeByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoProductType(shopID string, guid string) (models.ProductTypeInfo, error)
	InfoProductTypeByCode(shopID string, code string) (models.ProductTypeInfo, error)
	SearchProductType(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.ProductTypeInfo, mongopagination.PaginationData, error)
	SearchProductTypeStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.ProductTypeInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.ProductType) (common.BulkImport, error)

	GetModuleName() string
}

type ProductTypeHttpService struct {
	repo repositories.IProductTypeRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.ProductTypeActivity, models.ProductTypeDeleteActivity]
	contextTimeout time.Duration
}

func NewProductTypeHttpService(
	repo repositories.IProductTypeRepository,
	syncCacheRepo mastersync.IMasterSyncCacheRepository,
	contextTimeout time.Duration,
) *ProductTypeHttpService {

	insSvc := &ProductTypeHttpService{
		repo:           repo,
		syncCacheRepo:  syncCacheRepo,
		contextTimeout: contextTimeout,
	}

	insSvc.ActivityService = services.NewActivityService[models.ProductTypeActivity, models.ProductTypeDeleteActivity](repo)

	return insSvc
}

func (svc ProductTypeHttpService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc ProductTypeHttpService) CreateProductType(shopID string, authUsername string, doc models.ProductType) (string, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "code", doc.Code)

	if err != nil {
		return "", err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", errors.New("Code is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.ProductTypeDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.ProductType = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(ctx, docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc ProductTypeHttpService) UpdateProductType(shopID string, guid string, authUsername string, doc models.ProductType) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	findDoc.ProductType = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(ctx, shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc ProductTypeHttpService) DeleteProductType(shopID string, guid string, authUsername string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	err = svc.repo.DeleteByGuidfixed(ctx, shopID, guid, authUsername)
	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc ProductTypeHttpService) DeleteProductTypeByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(ctx, shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc ProductTypeHttpService) InfoProductType(shopID string, guid string) (models.ProductTypeInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return models.ProductTypeInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.ProductTypeInfo{}, errors.New("document not found")
	}

	return findDoc.ProductTypeInfo, nil
}

func (svc ProductTypeHttpService) InfoProductTypeByCode(shopID string, code string) (models.ProductTypeInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "code", code)

	if err != nil {
		return models.ProductTypeInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.ProductTypeInfo{}, errors.New("document not found")
	}

	return findDoc.ProductTypeInfo, nil
}

func (svc ProductTypeHttpService) SearchProductType(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.ProductTypeInfo, mongopagination.PaginationData, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"code",
		"names.name",
	}

	docList, pagination, err := svc.repo.FindPageFilter(ctx, shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.ProductTypeInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc ProductTypeHttpService) SearchProductTypeStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.ProductTypeInfo, int, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"code",
		"names.name",
	}

	selectFields := map[string]interface{}{}

	/*
		if langCode != "" {
			selectFields["names"] = bson.M{"$elemMatch": bson.M{"code": langCode}}
		} else {
			selectFields["names"] = 1
		}
	*/

	docList, total, err := svc.repo.FindStep(ctx, shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.ProductTypeInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc ProductTypeHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.ProductType) (common.BulkImport, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.ProductType](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.Code)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(ctx, shopID, "code", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.Code)
	}

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.ProductType, models.ProductTypeDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.ProductType) models.ProductTypeDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.ProductTypeDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.ProductType = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.ProductType, models.ProductTypeDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.ProductTypeDoc, error) {
			return svc.repo.FindByDocIndentityGuid(ctx, shopID, "code", guid)
		},
		func(doc models.ProductTypeDoc) bool {
			return doc.Code != ""
		},
		func(shopID string, authUsername string, data models.ProductType, doc models.ProductTypeDoc) error {

			doc.ProductType = data
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

func (svc ProductTypeHttpService) getDocIDKey(doc models.ProductType) string {
	return doc.Code
}

func (svc ProductTypeHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc ProductTypeHttpService) GetModuleName() string {
	return "productType"
}
