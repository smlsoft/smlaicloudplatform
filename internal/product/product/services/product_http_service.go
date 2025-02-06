package services

import (
	"context"
	"errors"
	"fmt"
	mastersync "smlaicloudplatform/internal/mastersync/repositories"
	common "smlaicloudplatform/internal/models"
	"smlaicloudplatform/internal/product/product/models"
	"smlaicloudplatform/internal/product/product/repositories"
	"smlaicloudplatform/internal/services"
	"smlaicloudplatform/internal/utils"
	"smlaicloudplatform/internal/utils/importdata"
	micromodels "smlaicloudplatform/pkg/microservice/models"
	"time"

	"github.com/smlsoft/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IProductHttpService interface {
	CreateProduct(shopID string, authUsername string, doc models.Product) (string, error)
	UpdateProduct(shopID string, guid string, authUsername string, doc models.Product) error
	DeleteProduct(shopID string, guid string, authUsername string) error
	DeleteProductByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoProduct(shopID string, guid string) (models.ProductInfo, error)
	InfoProductByCode(shopID, productCode string) (models.ProductInfo, error)
	SearchProduct(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.ProductInfo, mongopagination.PaginationData, error)
	SearchProductStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.ProductInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Product) (common.BulkImport, error)

	GetModuleName() string
}

type ProductHttpService struct {
	repo repositories.IProductRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.ProductActivity, models.ProductDeleteActivity]
	contextTimeout time.Duration
}

func NewProductHttpService(repo repositories.IProductRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *ProductHttpService {

	contextTimeout := time.Duration(15) * time.Second

	insSvc := &ProductHttpService{
		repo:           repo,
		syncCacheRepo:  syncCacheRepo,
		contextTimeout: contextTimeout,
	}

	insSvc.ActivityService = services.NewActivityService[models.ProductActivity, models.ProductDeleteActivity](repo)

	return insSvc
}

func (svc ProductHttpService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc ProductHttpService) CreateProduct(shopID string, authUsername string, doc models.Product) (string, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "code", doc.Code)

	if err != nil {
		return "", err
	}

	if findDoc.Code != "" {
		return "", errors.New("Code is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.ProductDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Product = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(ctx, docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc ProductHttpService) UpdateProduct(shopID string, guid string, authUsername string, doc models.Product) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	findDoc.Product = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(ctx, shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc ProductHttpService) DeleteProduct(shopID string, guid string, authUsername string) error {

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

func (svc ProductHttpService) DeleteProductByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(ctx, shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc ProductHttpService) InfoProduct(shopID string, guid string) (models.ProductInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return models.ProductInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.ProductInfo{}, errors.New("document not found")
	}

	return findDoc.ProductInfo, nil
}

func (svc ProductHttpService) InfoProductByCode(shopID, productCode string) (models.ProductInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindOneByCode(ctx, shopID, productCode)

	if err != nil {
		return models.ProductInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.ProductInfo{}, errors.New("document not found")
	}

	return findDoc.ProductInfo, nil
}

func (svc ProductHttpService) SearchProduct(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.ProductInfo, mongopagination.PaginationData, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"code",
	}

	docList, pagination, err := svc.repo.FindPageFilter(ctx, shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.ProductInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc ProductHttpService) SearchProductStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.ProductInfo, int, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"code",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(ctx, shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.ProductInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc ProductHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.Product) (common.BulkImport, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Product](dataList, svc.getDocIDKey)

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
			return svc.repo.FindByDocIndentityGuid(ctx, shopID, "code", guid)
		},
		func(doc models.ProductDoc) bool {
			return doc.Code != ""
		},
		func(shopID string, authUsername string, data models.Product, doc models.ProductDoc) error {

			doc.Product = data
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

func (svc ProductHttpService) getDocIDKey(doc models.Product) string {
	return doc.Code
}

func (svc ProductHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc ProductHttpService) GetModuleName() string {
	return "product"
}
