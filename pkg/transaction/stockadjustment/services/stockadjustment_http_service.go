package services

import (
	"context"
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/services"
	trancache "smlcloudplatform/pkg/transaction/repositories"
	"smlcloudplatform/pkg/transaction/stockadjustment/models"
	"smlcloudplatform/pkg/transaction/stockadjustment/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"strconv"
	"strings"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IStockAdjustmentHttpService interface {
	CreateStockAdjustment(shopID string, authUsername string, doc models.StockAdjustment) (string, string, error)
	UpdateStockAdjustment(shopID string, guid string, authUsername string, doc models.StockAdjustment) error
	DeleteStockAdjustment(shopID string, guid string, authUsername string) error
	DeleteStockAdjustmentByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoStockAdjustment(shopID string, guid string) (models.StockAdjustmentInfo, error)
	InfoStockAdjustmentByCode(shopID string, code string) (models.StockAdjustmentInfo, error)
	SearchStockAdjustment(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.StockAdjustmentInfo, mongopagination.PaginationData, error)
	SearchStockAdjustmentStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.StockAdjustmentInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.StockAdjustment) (common.BulkImport, error)

	GetModuleName() string
}

const (
	MODULE_NAME = "AJ"
)

type StockAdjustmentHttpService struct {
	repoMq           repositories.IStockAdjustmentMessageQueueRepository
	repo             repositories.IStockAdjustmentRepository
	repoCache        trancache.ICacheRepository
	cacheExpireDocNo time.Duration
	syncCacheRepo    mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.StockAdjustmentActivity, models.StockAdjustmentDeleteActivity]
	contextTimeout time.Duration
}

func NewStockAdjustmentHttpService(
	repo repositories.IStockAdjustmentRepository,
	repoCache trancache.ICacheRepository,
	repoMq repositories.IStockAdjustmentMessageQueueRepository,
	syncCacheRepo mastersync.IMasterSyncCacheRepository,
) *StockAdjustmentHttpService {

	contextTimeout := time.Duration(15) * time.Second

	insSvc := &StockAdjustmentHttpService{
		repo:             repo,
		repoMq:           repoMq,
		repoCache:        repoCache,
		syncCacheRepo:    syncCacheRepo,
		cacheExpireDocNo: time.Hour * 24,
		contextTimeout:   contextTimeout,
	}

	insSvc.ActivityService = services.NewActivityService[models.StockAdjustmentActivity, models.StockAdjustmentDeleteActivity](repo)

	return insSvc
}

func (svc StockAdjustmentHttpService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc StockAdjustmentHttpService) getDocNoPrefix(docDate time.Time) string {
	docDateStr := docDate.Format("20060102")
	return fmt.Sprintf("%s%s", MODULE_NAME, docDateStr)
}

func (svc StockAdjustmentHttpService) generateNewDocNo(ctx context.Context, shopID, prefixDocNo string, docNumber int) (string, int, error) {
	prevoiusDocNumber, err := svc.repoCache.Get(shopID, prefixDocNo)

	if prevoiusDocNumber == 0 || err != nil {
		lastDoc, err := svc.repo.FindLastDocNo(ctx, shopID, prefixDocNo)

		if err != nil {
			return "", 0, err
		}

		if len(lastDoc.DocNo) > 0 {
			rawNumber := strings.Replace(lastDoc.DocNo, prefixDocNo, "", -1)
			prevoiusDocNumber, err = strconv.Atoi(rawNumber)

			if err != nil {
				prevoiusDocNumber = 0
			}
		}

	}

	newDocNumber := prevoiusDocNumber + 1
	newDocNo := fmt.Sprintf("%s%05d", prefixDocNo, newDocNumber)

	findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "docno", newDocNo)

	if err != nil {
		return "", 0, err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", 0, errors.New("DocNo is exists")
	}

	return newDocNo, newDocNumber, nil
}

func (svc StockAdjustmentHttpService) CreateStockAdjustment(shopID string, authUsername string, doc models.StockAdjustment) (string, string, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	docDate := doc.DocDatetime
	prefixDocNo := svc.getDocNoPrefix(docDate)

	newDocNo, newDocNumber, err := svc.generateNewDocNo(ctx, shopID, prefixDocNo, 1)

	if err != nil {
		return "", "", err
	}

	newGuidFixed := utils.NewGUID()

	docData := models.StockAdjustmentDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.StockAdjustment = doc

	docData.DocNo = newDocNo
	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(ctx, docData)

	if err != nil {
		return "", "", err
	}

	go svc.repoCache.Save(shopID, prefixDocNo, newDocNumber, svc.cacheExpireDocNo)

	go func() {
		svc.repoMq.Create(docData)
		svc.repoCache.Save(shopID, prefixDocNo, newDocNumber, svc.cacheExpireDocNo)
		svc.saveMasterSync(shopID)
	}()

	return newGuidFixed, newDocNo, nil
}

func (svc StockAdjustmentHttpService) UpdateStockAdjustment(shopID string, guid string, authUsername string, doc models.StockAdjustment) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	docData := findDoc
	docData.StockAdjustment = doc

	docData.DocNo = findDoc.DocNo
	docData.UpdatedBy = authUsername
	docData.UpdatedAt = time.Now()

	err = svc.repo.Update(ctx, shopID, guid, docData)

	if err != nil {
		return err
	}

	func() {
		svc.repoMq.Update(docData)
		svc.saveMasterSync(shopID)
	}()

	return nil
}

func (svc StockAdjustmentHttpService) DeleteStockAdjustment(shopID string, guid string, authUsername string) error {

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

	func() {
		svc.repoMq.Delete(findDoc)
		svc.saveMasterSync(shopID)
	}()

	return nil
}

func (svc StockAdjustmentHttpService) DeleteStockAdjustmentByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(ctx, shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	func() {
		docs, _ := svc.repo.FindByGuids(ctx, shopID, GUIDs)
		svc.repoMq.DeleteInBatch(docs)
		svc.saveMasterSync(shopID)
	}()

	return nil
}

func (svc StockAdjustmentHttpService) InfoStockAdjustment(shopID string, guid string) (models.StockAdjustmentInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return models.StockAdjustmentInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.StockAdjustmentInfo{}, errors.New("document not found")
	}

	return findDoc.StockAdjustmentInfo, nil
}

func (svc StockAdjustmentHttpService) InfoStockAdjustmentByCode(shopID string, code string) (models.StockAdjustmentInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "docno", code)

	if err != nil {
		return models.StockAdjustmentInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.StockAdjustmentInfo{}, errors.New("document not found")
	}

	return findDoc.StockAdjustmentInfo, nil
}

func (svc StockAdjustmentHttpService) SearchStockAdjustment(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.StockAdjustmentInfo, mongopagination.PaginationData, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"docno",
	}

	docList, pagination, err := svc.repo.FindPageFilter(ctx, shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.StockAdjustmentInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc StockAdjustmentHttpService) SearchStockAdjustmentStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.StockAdjustmentInfo, int, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"docno",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(ctx, shopID, filters, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.StockAdjustmentInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc StockAdjustmentHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.StockAdjustment) (common.BulkImport, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.StockAdjustment](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.StockAdjustment, models.StockAdjustmentDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.StockAdjustment) models.StockAdjustmentDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.StockAdjustmentDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.StockAdjustment = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.StockAdjustment, models.StockAdjustmentDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.StockAdjustmentDoc, error) {
			return svc.repo.FindByDocIndentityGuid(ctx, shopID, "docno", guid)
		},
		func(doc models.StockAdjustmentDoc) bool {
			return doc.DocNo != ""
		},
		func(shopID string, authUsername string, data models.StockAdjustment, doc models.StockAdjustmentDoc) error {

			doc.StockAdjustment = data
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

	svc.saveMasterSync(shopID)

	return common.BulkImport{
		Created:          createDataKey,
		Updated:          updateDataKey,
		UpdateFailed:     updateFailDataKey,
		PayloadDuplicate: payloadDuplicateDataKey,
	}, nil
}

func (svc StockAdjustmentHttpService) getDocIDKey(doc models.StockAdjustment) string {
	return doc.DocNo
}

func (svc StockAdjustmentHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc StockAdjustmentHttpService) GetModuleName() string {
	return "stockAdjustment"
}
