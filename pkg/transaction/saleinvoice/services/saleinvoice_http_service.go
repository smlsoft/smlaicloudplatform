package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/services"
	trancache "smlcloudplatform/pkg/transaction/repositories"
	"smlcloudplatform/pkg/transaction/saleinvoice/models"
	"smlcloudplatform/pkg/transaction/saleinvoice/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"strconv"
	"strings"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type ISaleInvoiceHttpService interface {
	CreateSaleInvoice(shopID string, authUsername string, doc models.SaleInvoice) (string, string, error)
	UpdateSaleInvoice(shopID string, guid string, authUsername string, doc models.SaleInvoice) error
	DeleteSaleInvoice(shopID string, guid string, authUsername string) error
	DeleteSaleInvoiceByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoSaleInvoice(shopID string, guid string) (models.SaleInvoiceInfo, error)
	InfoSaleInvoiceByCode(shopID string, code string) (models.SaleInvoiceInfo, error)
	SearchSaleInvoice(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SaleInvoiceInfo, mongopagination.PaginationData, error)
	SearchSaleInvoiceStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SaleInvoiceInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.SaleInvoice) (common.BulkImport, error)

	GetModuleName() string
}

const (
	MODULE_NAME = "SI"
)

type SaleInvoiceHttpService struct {
	repo             repositories.ISaleInvoiceRepository
	repoCache        trancache.ICacheRepository
	cacheExpireDocNo time.Duration
	syncCacheRepo    mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.SaleInvoiceActivity, models.SaleInvoiceDeleteActivity]
}

func NewSaleInvoiceHttpService(repo repositories.ISaleInvoiceRepository, repoCache trancache.ICacheRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *SaleInvoiceHttpService {

	insSvc := &SaleInvoiceHttpService{
		repo:             repo,
		repoCache:        repoCache,
		syncCacheRepo:    syncCacheRepo,
		cacheExpireDocNo: time.Hour * 24,
	}

	insSvc.ActivityService = services.NewActivityService[models.SaleInvoiceActivity, models.SaleInvoiceDeleteActivity](repo)

	return insSvc
}

func (svc SaleInvoiceHttpService) getDocNoPrefix(docDate time.Time) string {
	docDateStr := docDate.Format("20060102")
	return fmt.Sprintf("%s%s", MODULE_NAME, docDateStr)
}

func (svc SaleInvoiceHttpService) generateNewDocNo(shopID, prefixDocNo string, docNumber int) (string, int, error) {
	prevoiusDocNumber, err := svc.repoCache.Get(shopID, prefixDocNo)

	if prevoiusDocNumber == 0 || err != nil {
		lastDoc, err := svc.repo.FindLastDocNo(shopID, prefixDocNo)

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

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "docno", newDocNo)

	if err != nil {
		return "", 0, err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", 0, errors.New("DocNo is exists")
	}

	return newDocNo, newDocNumber, nil
}

func (svc SaleInvoiceHttpService) CreateSaleInvoice(shopID string, authUsername string, doc models.SaleInvoice) (string, string, error) {

	docDate := doc.DocDatetime
	prefixDocNo := svc.getDocNoPrefix(docDate)

	newDocNo, newDocNumber, err := svc.generateNewDocNo(shopID, prefixDocNo, 1)

	if err != nil {
		return "", "", err
	}

	newGuidFixed := utils.NewGUID()

	docData := models.SaleInvoiceDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.SaleInvoice = doc

	docData.DocNo = newDocNo
	if doc.TaxDocNo == "" {
		docData.TaxDocNo = newDocNo
	}

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", "", err
	}

	go svc.repoCache.Save(shopID, prefixDocNo, newDocNumber, svc.cacheExpireDocNo)

	svc.saveMasterSync(shopID)

	return newGuidFixed, newDocNo, nil
}

func (svc SaleInvoiceHttpService) UpdateSaleInvoice(shopID string, guid string, authUsername string, doc models.SaleInvoice) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	findDoc.SaleInvoice = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc SaleInvoiceHttpService) DeleteSaleInvoice(shopID string, guid string, authUsername string) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	err = svc.repo.DeleteByGuidfixed(shopID, guid, authUsername)
	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc SaleInvoiceHttpService) DeleteSaleInvoiceByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc SaleInvoiceHttpService) InfoSaleInvoice(shopID string, guid string) (models.SaleInvoiceInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.SaleInvoiceInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.SaleInvoiceInfo{}, errors.New("document not found")
	}

	return findDoc.SaleInvoiceInfo, nil
}

func (svc SaleInvoiceHttpService) InfoSaleInvoiceByCode(shopID string, code string) (models.SaleInvoiceInfo, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "docno", code)

	if err != nil {
		return models.SaleInvoiceInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.SaleInvoiceInfo{}, errors.New("document not found")
	}

	return findDoc.SaleInvoiceInfo, nil
}

func (svc SaleInvoiceHttpService) SearchSaleInvoice(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SaleInvoiceInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"docno",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.SaleInvoiceInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc SaleInvoiceHttpService) SearchSaleInvoiceStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SaleInvoiceInfo, int, error) {
	searchInFields := []string{
		"docno",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, filters, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.SaleInvoiceInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc SaleInvoiceHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.SaleInvoice) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.SaleInvoice](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.DocNo)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "docno", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.DocNo)
	}

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.SaleInvoice, models.SaleInvoiceDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.SaleInvoice) models.SaleInvoiceDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.SaleInvoiceDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.SaleInvoice = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.SaleInvoice, models.SaleInvoiceDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.SaleInvoiceDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "docno", guid)
		},
		func(doc models.SaleInvoiceDoc) bool {
			return doc.DocNo != ""
		},
		func(shopID string, authUsername string, data models.SaleInvoice, doc models.SaleInvoiceDoc) error {

			doc.SaleInvoice = data
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

func (svc SaleInvoiceHttpService) getDocIDKey(doc models.SaleInvoice) string {
	return doc.DocNo
}

func (svc SaleInvoiceHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc SaleInvoiceHttpService) GetModuleName() string {
	return "saleInvoice"
}
