package services

import (
	"encoding/json"
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/transaction/documentformate/models"
	"smlcloudplatform/pkg/transaction/documentformate/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IDocumentFormateHttpService interface {
	CreateDocumentFormate(shopID string, authUsername string, doc models.DocumentFormate) (string, error)
	UpdateDocumentFormate(shopID string, guid string, authUsername string, doc models.DocumentFormate) error
	DeleteDocumentFormate(shopID string, guid string, authUsername string) error
	DeleteDocumentFormateByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoDocumentFormate(shopID string, guid string) (models.DocumentFormateInfo, error)
	InfoDocumentFormateByCode(shopID string, code string) (models.DocumentFormateInfo, error)
	SearchDocumentFormate(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.DocumentFormateInfo, mongopagination.PaginationData, error)
	SearchDocumentFormateStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.DocumentFormateInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.DocumentFormate) (common.BulkImport, error)
	GetModuleDefault() ([]map[string]interface{}, error)

	GetModuleName() string
}

type DocumentFormateHttpService struct {
	repo repositories.IDocumentFormateRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.DocumentFormateActivity, models.DocumentFormateDeleteActivity]
}

func NewDocumentFormateHttpService(repo repositories.IDocumentFormateRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *DocumentFormateHttpService {

	insSvc := &DocumentFormateHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.DocumentFormateActivity, models.DocumentFormateDeleteActivity](repo)

	return insSvc
}

func (svc DocumentFormateHttpService) CreateDocumentFormate(shopID string, authUsername string, doc models.DocumentFormate) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "doccode", doc.DocCode)

	if err != nil {
		return "", err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", errors.New("DocCode is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.DocumentFormateDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.DocumentFormate = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc DocumentFormateHttpService) UpdateDocumentFormate(shopID string, guid string, authUsername string, doc models.DocumentFormate) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	findDoc.DocumentFormate = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc DocumentFormateHttpService) DeleteDocumentFormate(shopID string, guid string, authUsername string) error {

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

func (svc DocumentFormateHttpService) DeleteDocumentFormateByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc DocumentFormateHttpService) InfoDocumentFormate(shopID string, guid string) (models.DocumentFormateInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.DocumentFormateInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.DocumentFormateInfo{}, errors.New("document not found")
	}

	return findDoc.DocumentFormateInfo, nil
}

func (svc DocumentFormateHttpService) InfoDocumentFormateByCode(shopID string, code string) (models.DocumentFormateInfo, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "doccode", code)

	if err != nil {
		return models.DocumentFormateInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.DocumentFormateInfo{}, errors.New("document not found")
	}

	return findDoc.DocumentFormateInfo, nil
}

func (svc DocumentFormateHttpService) SearchDocumentFormate(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.DocumentFormateInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"doccode",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.DocumentFormateInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc DocumentFormateHttpService) SearchDocumentFormateStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.DocumentFormateInfo, int, error) {
	searchInFields := []string{
		"doccode",
	}

	selectFields := map[string]interface{}{}

	/*
		if langCode != "" {
			selectFields["names"] = bson.M{"$elemMatch": bson.M{"code": langCode}}
		} else {
			selectFields["names"] = 1
		}
	*/

	docList, total, err := svc.repo.FindStep(shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.DocumentFormateInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc DocumentFormateHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.DocumentFormate) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.DocumentFormate](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.DocCode)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "doccode", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.DocCode)
	}

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.DocumentFormate, models.DocumentFormateDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.DocumentFormate) models.DocumentFormateDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.DocumentFormateDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.DocumentFormate = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.DocumentFormate, models.DocumentFormateDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.DocumentFormateDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "doccode", guid)
		},
		func(doc models.DocumentFormateDoc) bool {
			return doc.DocCode != ""
		},
		func(shopID string, authUsername string, data models.DocumentFormate, doc models.DocumentFormateDoc) error {

			doc.DocumentFormate = data
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
		createDataKey = append(createDataKey, doc.DocCode)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.DocCode)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {

		updateDataKey = append(updateDataKey, doc.DocCode)
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

func (svc DocumentFormateHttpService) getDocIDKey(doc models.DocumentFormate) string {
	return doc.DocCode
}

func (svc DocumentFormateHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc DocumentFormateHttpService) GetModuleName() string {
	return "documentFormate"
}

func (svc DocumentFormateHttpService) GetModuleDefault() ([]map[string]interface{}, error) {
	defaultModule := `
	[{"name":"Purchase","doccode":"PU","dateformate":"YYYYMMDD","docnumber":5},{"name":"Purchase Return","doccode":"PT","dateformate":"YYYYMMDD","docnumber":5},{"name":"SaleInvoice","doccode":"SI","dateformate":"YYYYMMDD","docnumber":5},{"name":"Sale Invoice Return","doccode":"ST","dateformate":"YYYYMMDD","docnumber":5},{"name":"Stock Adjustment","doccode":"AJ","dateformate":"YYYYMMDD","docnumber":5},{"name":"Stock Pickup Product","doccode":"IM","dateformate":"YYYYMMDD","docnumber":5},{"name":"Stock Receive Product","doccode":"IF","dateformate":"YYYYMMDD","docnumber":5},{"name":"Stock Return Product","doccode":"IR","dateformate":"YYYYMMDD","docnumber":5},{"name":"Stock Transfer","doccode":"TF","dateformate":"YYYYMMDD","docnumber":5}]
	`

	jsonData := []map[string]interface{}{}
	err := json.Unmarshal([]byte(defaultModule), &jsonData)

	if err != nil {
		return []map[string]interface{}{}, err
	}

	return jsonData, nil
}
