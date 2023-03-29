package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/debtaccount/creditorgroup/models"
	"smlcloudplatform/pkg/debtaccount/creditorgroup/repositories"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ICreditorGroupHttpService interface {
	CreateCreditorGroup(shopID string, authUsername string, doc models.CreditorGroup) (string, error)
	UpdateCreditorGroup(shopID string, guid string, authUsername string, doc models.CreditorGroup) error
	DeleteCreditorGroup(shopID string, guid string, authUsername string) error
	DeleteCreditorGroupByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoCreditorGroup(shopID string, guid string) (models.CreditorGroupInfo, error)
	SearchCreditorGroup(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.CreditorGroupInfo, mongopagination.PaginationData, error)
	SearchCreditorGroupStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.CreditorGroupInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.CreditorGroup) (common.BulkImport, error)

	GetModuleName() string
}

type CreditorGroupHttpService struct {
	repo repositories.ICreditorGroupRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.CreditorGroupActivity, models.CreditorGroupDeleteActivity]
}

func NewCreditorGroupHttpService(repo repositories.ICreditorGroupRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *CreditorGroupHttpService {

	insSvc := &CreditorGroupHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.CreditorGroupActivity, models.CreditorGroupDeleteActivity](repo)

	return insSvc
}

func (svc CreditorGroupHttpService) CreateCreditorGroup(shopID string, authUsername string, doc models.CreditorGroup) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "groupcode", doc.GroupCode)

	if err != nil {
		return "", err
	}

	if findDoc.GroupCode != "" {
		return "", errors.New("GroupCode is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.CreditorGroupDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.CreditorGroup = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc CreditorGroupHttpService) UpdateCreditorGroup(shopID string, guid string, authUsername string, doc models.CreditorGroup) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.CreditorGroup = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc CreditorGroupHttpService) DeleteCreditorGroup(shopID string, guid string, authUsername string) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	err = svc.repo.DeleteByGuidfixed(shopID, guid, authUsername)
	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc CreditorGroupHttpService) DeleteCreditorGroupByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc CreditorGroupHttpService) InfoCreditorGroup(shopID string, guid string) (models.CreditorGroupInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.CreditorGroupInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.CreditorGroupInfo{}, errors.New("document not found")
	}

	return findDoc.CreditorGroupInfo, nil

}

func (svc CreditorGroupHttpService) SearchCreditorGroup(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.CreditorGroupInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"groupcode",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.CreditorGroupInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc CreditorGroupHttpService) SearchCreditorGroupStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.CreditorGroupInfo, int, error) {
	searchInFields := []string{
		"groupcode",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.CreditorGroupInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc CreditorGroupHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.CreditorGroup) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.CreditorGroup](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.GroupCode)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "groupcode", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.GroupCode)
	}

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.CreditorGroup, models.CreditorGroupDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.CreditorGroup) models.CreditorGroupDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.CreditorGroupDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.CreditorGroup = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.CreditorGroup, models.CreditorGroupDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.CreditorGroupDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "groupcode", guid)
		},
		func(doc models.CreditorGroupDoc) bool {
			return doc.GroupCode != ""
		},
		func(shopID string, authUsername string, data models.CreditorGroup, doc models.CreditorGroupDoc) error {

			doc.CreditorGroup = data
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
		createDataKey = append(createDataKey, doc.GroupCode)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.GroupCode)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {

		updateDataKey = append(updateDataKey, doc.GroupCode)
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

func (svc CreditorGroupHttpService) getDocIDKey(doc models.CreditorGroup) string {
	return doc.GroupCode
}

func (svc CreditorGroupHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc CreditorGroupHttpService) GetModuleName() string {
	return "creditorGroup"
}
