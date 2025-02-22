package services

import (
	"context"
	"errors"
	"fmt"
	"smlaicloudplatform/internal/logger"
	mastersync "smlaicloudplatform/internal/mastersync/repositories"
	common "smlaicloudplatform/internal/models"
	"smlaicloudplatform/internal/pos/shift/models"
	"smlaicloudplatform/internal/pos/shift/repositories"
	"smlaicloudplatform/internal/services"
	"smlaicloudplatform/internal/utils"
	"smlaicloudplatform/internal/utils/importdata"
	micromodels "smlaicloudplatform/pkg/microservice/models"
	"time"

	"github.com/smlsoft/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IShiftHttpService interface {
	CreateShift(shopID string, authUsername string, doc models.Shift) (string, error)
	UpdateShift(shopID string, guid string, authUsername string, doc models.Shift) error
	DeleteShift(shopID string, guid string, authUsername string) error
	DeleteShiftByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoShift(shopID string, guid string) (models.ShiftInfo, error)
	InfoShiftByCode(shopID string, code string) (models.ShiftInfo, error)
	SearchShift(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.ShiftInfo, mongopagination.PaginationData, error)
	SearchShiftStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ShiftInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Shift) (common.BulkImport, error)

	GetModuleName() string
}

type ShiftHttpService struct {
	repo          repositories.IShiftRepository
	repoMq        repositories.IShiftMessageQueueRepository
	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.ShiftActivity, models.ShiftDeleteActivity]
	contextTimeout time.Duration
}

func NewShiftHttpService(
	repo repositories.IShiftRepository,
	repoMq repositories.IShiftMessageQueueRepository,
	syncCacheRepo mastersync.IMasterSyncCacheRepository,

	contextTimeout time.Duration,
) *ShiftHttpService {

	insSvc := &ShiftHttpService{
		repo:          repo,
		repoMq:        repoMq,
		syncCacheRepo: syncCacheRepo,

		contextTimeout: contextTimeout,
	}

	insSvc.ActivityService = services.NewActivityService[models.ShiftActivity, models.ShiftDeleteActivity](repo)

	return insSvc
}

func (svc ShiftHttpService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc ShiftHttpService) CreateShift(shopID string, authUsername string, doc models.Shift) (string, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "docno", doc.DocNo)

	if err != nil {
		return "", err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", errors.New("DocNo is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.ShiftDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Shift = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(ctx, docData)

	if err != nil {
		return "", err
	}

	go func() {
		svc.saveMasterSync(shopID)
		err = svc.repoMq.Create(docData)
		if err != nil {
			logger.GetLogger().Errorf("Create shift message queue error :: %s", err.Error())
		}
	}()

	return newGuidFixed, nil
}

func (svc ShiftHttpService) UpdateShift(shopID string, guid string, authUsername string, doc models.Shift) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	findDoc.Shift = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(ctx, shopID, guid, findDoc)

	if err != nil {
		return err
	}

	go func() {
		svc.saveMasterSync(shopID)
		err = svc.repoMq.Update(findDoc)
		if err != nil {
			logger.GetLogger().Errorf("Update shift message queue error :: %s", err.Error())
		}
	}()

	return nil
}

func (svc ShiftHttpService) DeleteShift(shopID string, guid string, authUsername string) error {

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

	go func() {
		svc.saveMasterSync(shopID)
		err = svc.repoMq.Delete(findDoc)
		if err != nil {
			logger.GetLogger().Errorf("Delete creditor message queue error :: %s", err.Error())
		}
	}()

	return nil
}

func (svc ShiftHttpService) DeleteShiftByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}
	findDocs, err := svc.repo.FindByGuids(ctx, shopID, GUIDs)

	if err != nil {
		return err
	}

	err = svc.repo.Delete(ctx, shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}
	go func() {
		svc.saveMasterSync(shopID)
		err = svc.repoMq.DeleteInBatch(findDocs)
		if err != nil {
			logger.GetLogger().Errorf("Delete creditor message queue error :: %s", err.Error())
		}
	}()

	return nil
}

func (svc ShiftHttpService) InfoShift(shopID string, guid string) (models.ShiftInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return models.ShiftInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.ShiftInfo{}, errors.New("document not found")
	}

	return findDoc.ShiftInfo, nil
}

func (svc ShiftHttpService) InfoShiftByCode(shopID string, code string) (models.ShiftInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "docno", code)

	if err != nil {
		return models.ShiftInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.ShiftInfo{}, errors.New("document not found")
	}

	return findDoc.ShiftInfo, nil
}

func (svc ShiftHttpService) SearchShift(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.ShiftInfo, mongopagination.PaginationData, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"docno",
		"username",
		"remark",
	}

	docList, pagination, err := svc.repo.FindPageFilter(ctx, shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.ShiftInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc ShiftHttpService) SearchShiftStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.ShiftInfo, int, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"docno",
		"username",
		"remark",
	}

	selectFields := map[string]interface{}{}

	/*
		if langCode != "" {
			selectFields["names"] = bson.M{"$elemMatch": bson.M{"code": langCode}}
		} else {
			selectFields["names"] = 1
		}
	*/

	docList, total, err := svc.repo.FindStep(ctx, shopID, filters, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.ShiftInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc ShiftHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.Shift) (common.BulkImport, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Shift](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.DocNo)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(ctx, shopID, "docno", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.DocNo)
	}

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Shift, models.ShiftDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Shift) models.ShiftDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.ShiftDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Shift = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Shift, models.ShiftDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.ShiftDoc, error) {
			return svc.repo.FindByDocIndentityGuid(ctx, shopID, "docno", guid)
		},
		func(doc models.ShiftDoc) bool {
			return doc.DocNo != ""
		},
		func(shopID string, authUsername string, data models.Shift, doc models.ShiftDoc) error {

			doc.Shift = data
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
		createDataKey = append(createDataKey, doc.DocNo)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.DocNo)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {

		updateDataKey = append(updateDataKey, doc.DocNo)
	}

	updateFailDataKey := []string{}
	for _, doc := range updateFailDataList {
		updateFailDataKey = append(updateFailDataKey, svc.getDocIDKey(doc))
	}

	go func() {
		svc.saveMasterSync(shopID)
		err = svc.repoMq.CreateInBatch(createDataList)
		if err != nil {
			logger.GetLogger().Errorf("Create shift message queue error :: %s", err.Error())
		}
		svc.repoMq.UpdateInBatch(updateSuccessDataList)

		if err != nil {
			logger.GetLogger().Errorf("Update shift message queue error :: %s", err.Error())
		}
	}()

	return common.BulkImport{
		Created:          createDataKey,
		Updated:          updateDataKey,
		UpdateFailed:     updateFailDataKey,
		PayloadDuplicate: payloadDuplicateDataKey,
	}, nil
}

func (svc ShiftHttpService) getDocIDKey(doc models.Shift) string {
	return doc.DocNo
}

func (svc ShiftHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc ShiftHttpService) GetModuleName() string {
	return "shift"
}
