package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/debtaccount/creditor/models"
	"smlcloudplatform/pkg/debtaccount/creditor/repositories"
	groupModels "smlcloudplatform/pkg/debtaccount/creditorgroup/models"
	groupRepositories "smlcloudplatform/pkg/debtaccount/creditorgroup/repositories"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/samber/lo"
	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ICreditorHttpService interface {
	CreateCreditor(shopID string, authUsername string, doc models.CreditorRequest) (string, error)
	UpdateCreditor(shopID string, guid string, authUsername string, doc models.CreditorRequest) error
	DeleteCreditor(shopID string, guid string, authUsername string) error
	DeleteCreditorByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoCreditor(shopID string, guid string) (models.CreditorInfo, error)
	InfoCreditorByCode(shopID string, code string) (models.CreditorInfo, error)
	SearchCreditor(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.CreditorInfo, mongopagination.PaginationData, error)
	SearchCreditorStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.CreditorInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.CreditorRequest) (common.BulkImport, error)

	GetModuleName() string
}

type CreditorHttpService struct {
	repo          repositories.ICreditorRepository
	repoGroup     groupRepositories.ICreditorGroupRepository
	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.CreditorActivity, models.CreditorDeleteActivity]
}

func NewCreditorHttpService(repo repositories.ICreditorRepository, repoGroup groupRepositories.ICreditorGroupRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *CreditorHttpService {

	insSvc := &CreditorHttpService{
		repo:          repo,
		repoGroup:     repoGroup,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.CreditorActivity, models.CreditorDeleteActivity](repo)

	return insSvc
}

func (svc CreditorHttpService) CreateCreditor(shopID string, authUsername string, doc models.CreditorRequest) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", doc.Code)

	if err != nil {
		return "", err
	}

	if findDoc.Code != "" {
		return "", errors.New("Code is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.CreditorDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Creditor = doc.Creditor
	docData.GroupGUIDs = &doc.Groups

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc CreditorHttpService) UpdateCreditor(shopID string, guid string, authUsername string, doc models.CreditorRequest) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.Creditor = doc.Creditor
	findDoc.GroupGUIDs = &doc.Groups

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc CreditorHttpService) DeleteCreditor(shopID string, guid string, authUsername string) error {

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

func (svc CreditorHttpService) DeleteCreditorByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc CreditorHttpService) InfoCreditor(shopID string, guid string) (models.CreditorInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.CreditorInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.CreditorInfo{}, errors.New("document not found")
	}

	findGroups, err := svc.repoGroup.FindByGuids(shopID, *findDoc.GroupGUIDs)

	if err != nil {
		return models.CreditorInfo{}, err
	}

	groupInfo := lo.Map[groupModels.CreditorGroupDoc, groupModels.CreditorGroupInfo](
		findGroups,
		func(docGroup groupModels.CreditorGroupDoc, idx int) groupModels.CreditorGroupInfo {
			return docGroup.CreditorGroupInfo
		})

	findDoc.CreditorInfo.Groups = &groupInfo

	return findDoc.CreditorInfo, nil

}

func (svc CreditorHttpService) InfoCreditorByCode(shopID string, code string) (models.CreditorInfo, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", code)

	if err != nil {
		return models.CreditorInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.CreditorInfo{}, errors.New("document not found")
	}

	findGroups, err := svc.repoGroup.FindByGuids(shopID, *findDoc.GroupGUIDs)

	if err != nil {
		return models.CreditorInfo{}, err
	}

	groupInfo := lo.Map[groupModels.CreditorGroupDoc, groupModels.CreditorGroupInfo](
		findGroups,
		func(docGroup groupModels.CreditorGroupDoc, idx int) groupModels.CreditorGroupInfo {
			return docGroup.CreditorGroupInfo
		})

	findDoc.CreditorInfo.Groups = &groupInfo

	return findDoc.CreditorInfo, nil

}

func (svc CreditorHttpService) SearchCreditor(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.CreditorInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"code",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.CreditorInfo{}, pagination, err
	}

	for idx, doc := range docList {
		if doc.GroupGUIDs != nil {
			findCustGroups, err := svc.repoGroup.FindByGuids(shopID, *doc.GroupGUIDs)
			if err != nil {
				return []models.CreditorInfo{}, pagination, err
			}

			custGroupInfo := lo.Map[groupModels.CreditorGroupDoc, groupModels.CreditorGroupInfo](
				findCustGroups,
				func(docGroup groupModels.CreditorGroupDoc, idx int) groupModels.CreditorGroupInfo {
					return docGroup.CreditorGroupInfo
				})

			docList[idx].Groups = &custGroupInfo
		}
	}

	return docList, pagination, nil
}

func (svc CreditorHttpService) SearchCreditorStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.CreditorInfo, int, error) {
	searchInFields := []string{
		"code",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.CreditorInfo{}, 0, err
	}

	for idx, doc := range docList {
		if doc.GroupGUIDs != nil {
			findCustGroups, err := svc.repoGroup.FindByGuids(shopID, *doc.GroupGUIDs)
			if err != nil {
				return []models.CreditorInfo{}, 0, err
			}

			custGroupInfo := lo.Map[groupModels.CreditorGroupDoc, groupModels.CreditorGroupInfo](
				findCustGroups,
				func(docGroup groupModels.CreditorGroupDoc, idx int) groupModels.CreditorGroupInfo {
					return docGroup.CreditorGroupInfo
				})

			docList[idx].Groups = &custGroupInfo
		}
	}

	return docList, total, nil
}

func (svc CreditorHttpService) SaveInBatch(shopID string, authUsername string, dataListParam []models.CreditorRequest) (common.BulkImport, error) {

	dataList := []models.Creditor{}
	for _, doc := range dataListParam {
		doc.GroupGUIDs = &doc.Groups
		dataList = append(dataList, doc.Creditor)
	}

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Creditor](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.Code)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "code", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.Code)
	}

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Creditor, models.CreditorDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Creditor) models.CreditorDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.CreditorDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Creditor = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Creditor, models.CreditorDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.CreditorDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "code", guid)
		},
		func(doc models.CreditorDoc) bool {
			return doc.Code != ""
		},
		func(shopID string, authUsername string, data models.Creditor, doc models.CreditorDoc) error {

			doc.Creditor = data
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
		createDataKey = append(createDataKey, doc.Code)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.Code)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {

		updateDataKey = append(updateDataKey, doc.Code)
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

func (svc CreditorHttpService) getDocIDKey(doc models.Creditor) string {
	return doc.Code
}

func (svc CreditorHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc CreditorHttpService) GetModuleName() string {
	return "creditor"
}
