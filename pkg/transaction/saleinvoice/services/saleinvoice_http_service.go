package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/transaction/saleinvoice/models"
	"smlcloudplatform/pkg/transaction/saleinvoice/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type ISaleInvoiceHttpService interface {
	CreateSaleInvoice(shopID string, authUsername string, doc models.SaleInvoice) (string, error)
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

type SaleInvoiceHttpService struct {
	repo repositories.ISaleInvoiceRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.SaleInvoiceActivity, models.SaleInvoiceDeleteActivity]
}

func NewSaleInvoiceHttpService(repo repositories.ISaleInvoiceRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *SaleInvoiceHttpService {

	insSvc := &SaleInvoiceHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.SaleInvoiceActivity, models.SaleInvoiceDeleteActivity](repo)

	return insSvc
}

func (svc SaleInvoiceHttpService) CreateSaleInvoice(shopID string, authUsername string, doc models.SaleInvoice) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "docno", doc.Docno)

	if err != nil {
		return "", err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", errors.New("Docno is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.SaleInvoiceDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.SaleInvoice = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
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
		itemCodeGuidList = append(itemCodeGuidList, doc.Docno)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "docno", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.Docno)
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
			return doc.Docno != ""
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
		createDataKey = append(createDataKey, doc.Docno)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.Docno)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {

		updateDataKey = append(updateDataKey, doc.Docno)
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
	return doc.Docno
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
