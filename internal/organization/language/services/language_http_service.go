package services

import (
	"context"
	"errors"
	"fmt"
	mastersync "smlaicloudplatform/internal/mastersync/repositories"
	common "smlaicloudplatform/internal/models"
	"smlaicloudplatform/internal/organization/language/models"
	"smlaicloudplatform/internal/organization/language/repositories"
	"smlaicloudplatform/internal/services"
	"smlaicloudplatform/internal/utils"
	"smlaicloudplatform/internal/utils/importdata"
	micromodels "smlaicloudplatform/pkg/microservice/models"
	"time"

	"github.com/smlsoft/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type ILanguageHttpService interface {
	CreateLanguage(shopID string, authUsername string, doc models.Language) (string, error)
	UpdateLanguage(shopID string, guid string, authUsername string, doc models.Language) error
	DeleteLanguage(shopID string, guid string, authUsername string) error
	DeleteLanguageByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoLanguage(shopID string, guid string) (models.LanguageInfo, error)
	InfoLanguageByCode(shopID, languageCode string) (models.LanguageInfo, error)
	SearchLanguage(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.LanguageInfo, mongopagination.PaginationData, error)
	SearchLanguageStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.LanguageInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Language) (common.BulkImport, error)

	GetModuleName() string
}

type LanguageHttpService struct {
	repo repositories.ILanguageRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.LanguageActivity, models.LanguageDeleteActivity]
	contextTimeout time.Duration
}

func NewLanguageHttpService(repo repositories.ILanguageRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *LanguageHttpService {

	contextTimeout := time.Duration(15) * time.Second

	insSvc := &LanguageHttpService{
		repo:           repo,
		syncCacheRepo:  syncCacheRepo,
		contextTimeout: contextTimeout,
	}

	insSvc.ActivityService = services.NewActivityService[models.LanguageActivity, models.LanguageDeleteActivity](repo)

	return insSvc
}

func (svc LanguageHttpService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc LanguageHttpService) CreateLanguage(shopID string, authUsername string, doc models.Language) (string, error) {

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

	docData := models.LanguageDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Language = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(ctx, docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc LanguageHttpService) UpdateLanguage(shopID string, guid string, authUsername string, doc models.Language) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	findDoc.Language = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(ctx, shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc LanguageHttpService) DeleteLanguage(shopID string, guid string, authUsername string) error {

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

func (svc LanguageHttpService) DeleteLanguageByGUIDs(shopID string, authUsername string, GUIDs []string) error {

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

func (svc LanguageHttpService) InfoLanguage(shopID string, guid string) (models.LanguageInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return models.LanguageInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.LanguageInfo{}, errors.New("document not found")
	}

	return findDoc.LanguageInfo, nil
}

func (svc LanguageHttpService) InfoLanguageByCode(shopID, languageCode string) (models.LanguageInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindOneByCode(ctx, shopID, languageCode)

	if err != nil {
		return models.LanguageInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.LanguageInfo{}, errors.New("document not found")
	}

	return findDoc.LanguageInfo, nil
}

func (svc LanguageHttpService) SearchLanguage(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.LanguageInfo, mongopagination.PaginationData, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"code",
	}

	docList, pagination, err := svc.repo.FindPageFilter(ctx, shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.LanguageInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc LanguageHttpService) SearchLanguageStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.LanguageInfo, int, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"code",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(ctx, shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.LanguageInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc LanguageHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.Language) (common.BulkImport, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Language](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Language, models.LanguageDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Language) models.LanguageDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.LanguageDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Language = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Language, models.LanguageDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.LanguageDoc, error) {
			return svc.repo.FindByDocIndentityGuid(ctx, shopID, "code", guid)
		},
		func(doc models.LanguageDoc) bool {
			return doc.Code != ""
		},
		func(shopID string, authUsername string, data models.Language, doc models.LanguageDoc) error {

			doc.Language = data
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

func (svc LanguageHttpService) getDocIDKey(doc models.Language) string {
	return doc.Code
}

func (svc LanguageHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc LanguageHttpService) GetModuleName() string {
	return "language"
}
