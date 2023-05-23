package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/transaction/pay/models"
	"smlcloudplatform/pkg/transaction/pay/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IPayHttpService interface {
	CreatePay(shopID string, authUsername string, doc models.Pay) (string, error)
	UpdatePay(shopID string, guid string, authUsername string, doc models.Pay) error
	DeletePay(shopID string, guid string, authUsername string) error
	DeletePayByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoPay(shopID string, guid string) (models.PayInfo, error)
	InfoPayByCode(shopID string, code string) (models.PayInfo, error)
	SearchPay(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.PayInfo, mongopagination.PaginationData, error)
	SearchPayStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.PayInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Pay) (common.BulkImport, error)

	GetModuleName() string
}

type PayHttpService struct {
	repo repositories.IPayRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.PayActivity, models.PayDeleteActivity]
}

func NewPayHttpService(repo repositories.IPayRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *PayHttpService {

	insSvc := &PayHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.PayActivity, models.PayDeleteActivity](repo)

	return insSvc
}

func (svc PayHttpService) CreatePay(shopID string, authUsername string, doc models.Pay) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "docno", doc.DocNo)

	if err != nil {
		return "", err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", errors.New("DocNo is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.PayDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Pay = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc PayHttpService) UpdatePay(shopID string, guid string, authUsername string, doc models.Pay) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	findDoc.Pay = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc PayHttpService) DeletePay(shopID string, guid string, authUsername string) error {

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

func (svc PayHttpService) DeletePayByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc PayHttpService) InfoPay(shopID string, guid string) (models.PayInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.PayInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.PayInfo{}, errors.New("document not found")
	}

	return findDoc.PayInfo, nil
}

func (svc PayHttpService) InfoPayByCode(shopID string, code string) (models.PayInfo, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "docno", code)

	if err != nil {
		return models.PayInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.PayInfo{}, errors.New("document not found")
	}

	return findDoc.PayInfo, nil
}

func (svc PayHttpService) SearchPay(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.PayInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"docno",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.PayInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc PayHttpService) SearchPayStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.PayInfo, int, error) {
	searchInFields := []string{
		"docno",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, filters, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.PayInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc PayHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.Pay) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Pay](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Pay, models.PayDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Pay) models.PayDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.PayDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Pay = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Pay, models.PayDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.PayDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "docno", guid)
		},
		func(doc models.PayDoc) bool {
			return doc.DocNo != ""
		},
		func(shopID string, authUsername string, data models.Pay, doc models.PayDoc) error {

			doc.Pay = data
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

func (svc PayHttpService) getDocIDKey(doc models.Pay) string {
	return doc.DocNo
}

func (svc PayHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc PayHttpService) GetModuleName() string {
	return "pay"
}
