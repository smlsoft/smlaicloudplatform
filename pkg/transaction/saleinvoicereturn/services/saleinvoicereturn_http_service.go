package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/services"
	trancache "smlcloudplatform/pkg/transaction/repositories"
	"smlcloudplatform/pkg/transaction/saleinvoicereturn/models"
	"smlcloudplatform/pkg/transaction/saleinvoicereturn/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"strconv"
	"strings"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type ISaleInvoiceReturnHttpService interface {
	CreateSaleInvoiceReturn(shopID string, authUsername string, doc models.SaleInvoiceReturn) (string, string, error)
	UpdateSaleInvoiceReturn(shopID string, guid string, authUsername string, doc models.SaleInvoiceReturn) error
	DeleteSaleInvoiceReturn(shopID string, guid string, authUsername string) error
	DeleteSaleInvoiceReturnByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoSaleInvoiceReturn(shopID string, guid string) (models.SaleInvoiceReturnInfo, error)
	InfoSaleInvoiceReturnByCode(shopID string, code string) (models.SaleInvoiceReturnInfo, error)
	SearchSaleInvoiceReturn(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SaleInvoiceReturnInfo, mongopagination.PaginationData, error)
	SearchSaleInvoiceReturnStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SaleInvoiceReturnInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.SaleInvoiceReturn) (common.BulkImport, error)

	GetModuleName() string
}

const (
	MODULE_NAME = "ST"
)

type SaleInvoiceReturnHttpService struct {
	repo             repositories.ISaleInvoiceReturnRepository
	repoCache        trancache.ICacheRepository
	cacheExpireDocNo time.Duration
	syncCacheRepo    mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.SaleInvoiceReturnActivity, models.SaleInvoiceReturnDeleteActivity]
}

func NewSaleInvoiceReturnHttpService(repo repositories.ISaleInvoiceReturnRepository, repoCache trancache.ICacheRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *SaleInvoiceReturnHttpService {

	insSvc := &SaleInvoiceReturnHttpService{
		repo:             repo,
		repoCache:        repoCache,
		syncCacheRepo:    syncCacheRepo,
		cacheExpireDocNo: time.Hour * 24,
	}

	insSvc.ActivityService = services.NewActivityService[models.SaleInvoiceReturnActivity, models.SaleInvoiceReturnDeleteActivity](repo)

	return insSvc
}

func (svc SaleInvoiceReturnHttpService) getDocNoPrefix(docDate time.Time) string {
	docDateStr := docDate.Format("20060102")
	return fmt.Sprintf("%s%s", MODULE_NAME, docDateStr)
}

func (svc SaleInvoiceReturnHttpService) generateNewDocNo(shopID, prefixDocNo string, docNumber int) (string, int, error) {
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
func (svc SaleInvoiceReturnHttpService) CreateSaleInvoiceReturn(shopID string, authUsername string, doc models.SaleInvoiceReturn) (string, string, error) {
	timeNow := time.Now()
	prefixDocNo := svc.getDocNoPrefix(timeNow)

	newDocNo, newDocNumber, err := svc.generateNewDocNo(shopID, prefixDocNo, 1)

	if err != nil {
		return "", "", err
	}

	newGuidFixed := utils.NewGUID()

	docData := models.SaleInvoiceReturnDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.SaleInvoiceReturn = doc

	docData.DocNo = newDocNo
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

func (svc SaleInvoiceReturnHttpService) UpdateSaleInvoiceReturn(shopID string, guid string, authUsername string, doc models.SaleInvoiceReturn) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	findDoc.SaleInvoiceReturn = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc SaleInvoiceReturnHttpService) DeleteSaleInvoiceReturn(shopID string, guid string, authUsername string) error {

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

func (svc SaleInvoiceReturnHttpService) DeleteSaleInvoiceReturnByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc SaleInvoiceReturnHttpService) InfoSaleInvoiceReturn(shopID string, guid string) (models.SaleInvoiceReturnInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.SaleInvoiceReturnInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.SaleInvoiceReturnInfo{}, errors.New("document not found")
	}

	return findDoc.SaleInvoiceReturnInfo, nil
}

func (svc SaleInvoiceReturnHttpService) InfoSaleInvoiceReturnByCode(shopID string, code string) (models.SaleInvoiceReturnInfo, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "docno", code)

	if err != nil {
		return models.SaleInvoiceReturnInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.SaleInvoiceReturnInfo{}, errors.New("document not found")
	}

	return findDoc.SaleInvoiceReturnInfo, nil
}

func (svc SaleInvoiceReturnHttpService) SearchSaleInvoiceReturn(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SaleInvoiceReturnInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"docno",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.SaleInvoiceReturnInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc SaleInvoiceReturnHttpService) SearchSaleInvoiceReturnStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.SaleInvoiceReturnInfo, int, error) {
	searchInFields := []string{
		"docno",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, filters, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.SaleInvoiceReturnInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc SaleInvoiceReturnHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.SaleInvoiceReturn) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.SaleInvoiceReturn](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.SaleInvoiceReturn, models.SaleInvoiceReturnDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.SaleInvoiceReturn) models.SaleInvoiceReturnDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.SaleInvoiceReturnDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.SaleInvoiceReturn = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.SaleInvoiceReturn, models.SaleInvoiceReturnDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.SaleInvoiceReturnDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "docno", guid)
		},
		func(doc models.SaleInvoiceReturnDoc) bool {
			return doc.DocNo != ""
		},
		func(shopID string, authUsername string, data models.SaleInvoiceReturn, doc models.SaleInvoiceReturnDoc) error {

			doc.SaleInvoiceReturn = data
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

func (svc SaleInvoiceReturnHttpService) getDocIDKey(doc models.SaleInvoiceReturn) string {
	return doc.DocNo
}

func (svc SaleInvoiceReturnHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc SaleInvoiceReturnHttpService) GetModuleName() string {
	return "saleInvoiceReturn"
}
