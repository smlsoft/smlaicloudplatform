package services

import (
	"context"
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	"smlcloudplatform/pkg/services"
	trancache "smlcloudplatform/pkg/transaction/repositories"
	"smlcloudplatform/pkg/transaction/stockbalancedetail/models"
	"smlcloudplatform/pkg/transaction/stockbalancedetail/repositories"
	"smlcloudplatform/pkg/utils"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IStockBalanceDetailHttpService interface {
	CreateStockBalanceDetail(shopID string, authUsername string, doc []models.StockBalanceDetail) error
	UpdateStockBalanceDetail(shopID string, guid string, authUsername string, doc models.StockBalanceDetail) error
	DeleteStockBalanceDetail(shopID string, guid string, authUsername string) error
	DeleteStockBalanceDetailByGUIDs(shopID string, authUsername string, GUIDs []string) error
	DeleteStockBalanceDetailByDocNo(shopID string, authUsername string, docNo string) error
	InfoStockBalanceDetail(shopID string, guid string) (models.StockBalanceDetailInfo, error)
	InfoStockBalanceDetailByCode(shopID string, code string) (models.StockBalanceDetailInfo, error)
	SearchStockBalanceDetail(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.StockBalanceDetailInfo, mongopagination.PaginationData, error)
	SearchStockBalanceDetailStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.StockBalanceDetailInfo, int, error)

	GetModuleName() string
}

type StockBalanceDetailHttpService struct {
	repoMq           repositories.IStockBalanceDetailMessageQueueRepository
	repo             repositories.IStockBalanceDetailRepository
	repoCache        trancache.ICacheRepository
	cacheExpireDocNo time.Duration
	syncCacheRepo    mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.StockBalanceDetailActivity, models.StockBalanceDetailDeleteActivity]
	contextTimeout time.Duration
}

func NewStockBalanceDetailHttpService(
	repo repositories.IStockBalanceDetailRepository,
	repoCache trancache.ICacheRepository,
	repoMq repositories.IStockBalanceDetailMessageQueueRepository,
	syncCacheRepo mastersync.IMasterSyncCacheRepository,
) *StockBalanceDetailHttpService {

	contextTimeout := time.Duration(15) * time.Second

	insSvc := &StockBalanceDetailHttpService{
		repoMq:           repoMq,
		repo:             repo,
		repoCache:        repoCache,
		syncCacheRepo:    syncCacheRepo,
		cacheExpireDocNo: time.Hour * 24,
		contextTimeout:   contextTimeout,
	}

	insSvc.ActivityService = services.NewActivityService[models.StockBalanceDetailActivity, models.StockBalanceDetailDeleteActivity](repo)

	return insSvc
}

func (svc StockBalanceDetailHttpService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc StockBalanceDetailHttpService) CreateStockBalanceDetail(shopID string, authUsername string, docs []models.StockBalanceDetail) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	prepareDocs := []models.StockBalanceDetailDoc{}

	for _, doc := range docs {
		newGuidFixed := utils.NewGUID()

		docData := models.StockBalanceDetailDoc{}
		docData.ShopID = shopID
		docData.GuidFixed = newGuidFixed
		docData.StockBalanceDetail = doc

		docData.CreatedBy = authUsername
		docData.CreatedAt = time.Now()

		prepareDocs = append(prepareDocs, docData)
	}

	err := svc.repo.CreateInBatch(ctx, prepareDocs)

	if err != nil {
		return err
	}

	go func() {
		svc.repoMq.CreateInBatch(prepareDocs)
		svc.saveMasterSync(shopID)
	}()

	return nil
}

func (svc StockBalanceDetailHttpService) UpdateStockBalanceDetail(shopID string, guid string, authUsername string, doc models.StockBalanceDetail) error {

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
	docData.StockBalanceDetail = doc

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

func (svc StockBalanceDetailHttpService) DeleteStockBalanceDetail(shopID string, guid string, authUsername string) error {

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

func (svc StockBalanceDetailHttpService) DeleteStockBalanceDetailByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	// prepare item for message queue
	docs, err := svc.repo.FindByGuids(ctx, shopID, GUIDs)
	if err != nil {
		return err
	}

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err = svc.repo.Delete(ctx, shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	func() {
		svc.repoMq.DeleteInBatch(docs)
		svc.saveMasterSync(shopID)
	}()

	return nil
}

func (svc StockBalanceDetailHttpService) DeleteStockBalanceDetailByDocNo(shopID string, authUsername string, docNo string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	docs, err := svc.repo.FindByDocIndentityGuids(ctx, shopID, "docno", docNo)
	if err != nil {
		return err
	}

	deleteFilterQuery := map[string]interface{}{
		"docno": docNo,
	}

	err = svc.repo.Delete(ctx, shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	func() {
		svc.repoMq.DeleteInBatch(docs)
		svc.saveMasterSync(shopID)
	}()

	return nil
}

func (svc StockBalanceDetailHttpService) InfoStockBalanceDetail(shopID string, guid string) (models.StockBalanceDetailInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return models.StockBalanceDetailInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.StockBalanceDetailInfo{}, errors.New("document not found")
	}

	return findDoc.StockBalanceDetailInfo, nil
}

func (svc StockBalanceDetailHttpService) InfoStockBalanceDetailByCode(shopID string, code string) (models.StockBalanceDetailInfo, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "docno", code)

	if err != nil {
		return models.StockBalanceDetailInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.StockBalanceDetailInfo{}, errors.New("document not found")
	}

	return findDoc.StockBalanceDetailInfo, nil
}

func (svc StockBalanceDetailHttpService) SearchStockBalanceDetail(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.StockBalanceDetailInfo, mongopagination.PaginationData, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"docno",
	}

	docList, pagination, err := svc.repo.FindPageFilter(ctx, shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.StockBalanceDetailInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc StockBalanceDetailHttpService) SearchStockBalanceDetailStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.StockBalanceDetailInfo, int, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"docno",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(ctx, shopID, filters, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.StockBalanceDetailInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc StockBalanceDetailHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc StockBalanceDetailHttpService) GetModuleName() string {
	return "stockBalanceDetail"
}
