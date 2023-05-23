package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/transaction/paid/models"
	"smlcloudplatform/pkg/transaction/paid/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IPaidHttpService interface {
	CreatePaid(shopID string, authUsername string, doc models.Paid) (string, error)
	UpdatePaid(shopID string, guid string, authUsername string, doc models.Paid) error
	DeletePaid(shopID string, guid string, authUsername string) error
	DeletePaidByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoPaid(shopID string, guid string) (models.PaidInfo, error)
	InfoPaidByCode(shopID string, code string) (models.PaidInfo, error)
	SearchPaid(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.PaidInfo, mongopagination.PaginationData, error)
	SearchPaidStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.PaidInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Paid) (common.BulkImport, error)

	GetModuleName() string
}

type PaidHttpService struct {
	repo repositories.IPaidRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.PaidActivity, models.PaidDeleteActivity]
}

func NewPaidHttpService(repo repositories.IPaidRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *PaidHttpService {

	insSvc := &PaidHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.PaidActivity, models.PaidDeleteActivity](repo)

	return insSvc
}

func (svc PaidHttpService) CreatePaid(shopID string, authUsername string, doc models.Paid) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "docno", doc.DocNo)

	if err != nil {
		return "", err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", errors.New("DocNo is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.PaidDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Paid = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc PaidHttpService) UpdatePaid(shopID string, guid string, authUsername string, doc models.Paid) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	findDoc.Paid = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc PaidHttpService) DeletePaid(shopID string, guid string, authUsername string) error {

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

func (svc PaidHttpService) DeletePaidByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc PaidHttpService) InfoPaid(shopID string, guid string) (models.PaidInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.PaidInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.PaidInfo{}, errors.New("document not found")
	}

	return findDoc.PaidInfo, nil
}

func (svc PaidHttpService) InfoPaidByCode(shopID string, code string) (models.PaidInfo, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "docno", code)

	if err != nil {
		return models.PaidInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.PaidInfo{}, errors.New("document not found")
	}

	return findDoc.PaidInfo, nil
}

func (svc PaidHttpService) SearchPaid(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.PaidInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"docno",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.PaidInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc PaidHttpService) SearchPaidStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.PaidInfo, int, error) {
	searchInFields := []string{
		"docno",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, filters, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.PaidInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc PaidHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.Paid) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Paid](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Paid, models.PaidDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Paid) models.PaidDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.PaidDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Paid = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Paid, models.PaidDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.PaidDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "docno", guid)
		},
		func(doc models.PaidDoc) bool {
			return doc.DocNo != ""
		},
		func(shopID string, authUsername string, data models.Paid, doc models.PaidDoc) error {

			doc.Paid = data
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

func (svc PaidHttpService) getDocIDKey(doc models.Paid) string {
	return doc.DocNo
}

func (svc PaidHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc PaidHttpService) GetModuleName() string {
	return "paid"
}
