package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/debtaccount/debtor/models"
	"smlcloudplatform/pkg/debtaccount/debtor/repositories"
	groupModels "smlcloudplatform/pkg/debtaccount/debtorgroup/models"
	groupRepositories "smlcloudplatform/pkg/debtaccount/debtorgroup/repositories"
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

type IDebtorHttpService interface {
	CreateDebtor(shopID string, authUsername string, doc models.DebtorRequest) (string, error)
	UpdateDebtor(shopID string, guid string, authUsername string, doc models.DebtorRequest) error
	DeleteDebtor(shopID string, guid string, authUsername string) error
	DeleteDebtorByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoDebtor(shopID string, guid string) (models.DebtorInfo, error)
	InfoDebtorByCode(shopID string, code string) (models.DebtorInfo, error)
	SearchDebtor(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.DebtorInfo, mongopagination.PaginationData, error)
	SearchDebtorStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.DebtorInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.DebtorRequest) (common.BulkImport, error)

	GetModuleName() string
}

type DebtorHttpService struct {
	repo      repositories.IDebtorRepository
	repoGroup groupRepositories.IDebtorGroupRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.DebtorActivity, models.DebtorDeleteActivity]
}

func NewDebtorHttpService(repo repositories.IDebtorRepository, repoGroup groupRepositories.IDebtorGroupRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *DebtorHttpService {

	insSvc := &DebtorHttpService{
		repo:          repo,
		repoGroup:     repoGroup,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.DebtorActivity, models.DebtorDeleteActivity](repo)

	return insSvc
}

func (svc DebtorHttpService) CreateDebtor(shopID string, authUsername string, doc models.DebtorRequest) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", doc.Code)

	if err != nil {
		return "", err
	}

	if findDoc.Code != "" {
		return "", errors.New("Code is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.DebtorDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Debtor = doc.Debtor
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

func (svc DebtorHttpService) UpdateDebtor(shopID string, guid string, authUsername string, doc models.DebtorRequest) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.Debtor = doc.Debtor
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

func (svc DebtorHttpService) DeleteDebtor(shopID string, guid string, authUsername string) error {

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

func (svc DebtorHttpService) DeleteDebtorByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc DebtorHttpService) InfoDebtor(shopID string, guid string) (models.DebtorInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.DebtorInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.DebtorInfo{}, errors.New("document not found")
	}

	findGroups, err := svc.repoGroup.FindByGuids(shopID, *findDoc.GroupGUIDs)

	if err != nil {
		return models.DebtorInfo{}, err
	}

	custGroupInfo := lo.Map[groupModels.DebtorGroupDoc, groupModels.DebtorGroupInfo](
		findGroups,
		func(docGroup groupModels.DebtorGroupDoc, idx int) groupModels.DebtorGroupInfo {
			return docGroup.DebtorGroupInfo
		})

	findDoc.DebtorInfo.Groups = &custGroupInfo

	return findDoc.DebtorInfo, nil

}

func (svc DebtorHttpService) InfoDebtorByCode(shopID string, code string) (models.DebtorInfo, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", code)

	if err != nil {
		return models.DebtorInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.DebtorInfo{}, errors.New("document not found")
	}

	findGroups, err := svc.repoGroup.FindByGuids(shopID, *findDoc.GroupGUIDs)

	if err != nil {
		return models.DebtorInfo{}, err
	}

	custGroupInfo := lo.Map[groupModels.DebtorGroupDoc, groupModels.DebtorGroupInfo](
		findGroups,
		func(docGroup groupModels.DebtorGroupDoc, idx int) groupModels.DebtorGroupInfo {
			return docGroup.DebtorGroupInfo
		})

	findDoc.DebtorInfo.Groups = &custGroupInfo

	return findDoc.DebtorInfo, nil

}

func (svc DebtorHttpService) SearchDebtor(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.DebtorInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"code",
		"names.name",
		"groups",
		"fundcode",
		"addressforbilling.address.0",
		"addressforbilling.phoneprimary",
		"addressforbilling.phonesecondary",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.DebtorInfo{}, pagination, err
	}

	for idx, doc := range docList {
		if doc.GroupGUIDs != nil {
			findCustGroups, err := svc.repoGroup.FindByGuids(shopID, *doc.GroupGUIDs)
			if err != nil {
				return []models.DebtorInfo{}, pagination, err
			}

			custGroupInfo := lo.Map[groupModels.DebtorGroupDoc, groupModels.DebtorGroupInfo](
				findCustGroups,
				func(docGroup groupModels.DebtorGroupDoc, idx int) groupModels.DebtorGroupInfo {
					return docGroup.DebtorGroupInfo
				})

			docList[idx].Groups = &custGroupInfo
		}
	}

	return docList, pagination, nil
}

func (svc DebtorHttpService) SearchDebtorStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.DebtorInfo, int, error) {
	searchInFields := []string{
		"code",
		"names.name",
		"groups",
		"fundcode",
		"addressforbilling.address.0",
		"addressforbilling.phoneprimary",
		"addressforbilling.phonesecondary",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, filters, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.DebtorInfo{}, 0, err
	}

	for idx, doc := range docList {
		if doc.GroupGUIDs != nil {
			findCustGroups, err := svc.repoGroup.FindByGuids(shopID, *doc.GroupGUIDs)
			if err != nil {
				return []models.DebtorInfo{}, 0, err
			}

			custGroupInfo := lo.Map[groupModels.DebtorGroupDoc, groupModels.DebtorGroupInfo](
				findCustGroups,
				func(docGroup groupModels.DebtorGroupDoc, idx int) groupModels.DebtorGroupInfo {
					return docGroup.DebtorGroupInfo
				})

			docList[idx].Groups = &custGroupInfo
		}
	}

	return docList, total, nil
}

func (svc DebtorHttpService) SaveInBatch(shopID string, authUsername string, dataListReq []models.DebtorRequest) (common.BulkImport, error) {

	dataList := []models.Debtor{}
	for _, doc := range dataListReq {
		doc.GroupGUIDs = &doc.Groups
		dataList = append(dataList, doc.Debtor)
	}

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Debtor](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Debtor, models.DebtorDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Debtor) models.DebtorDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.DebtorDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Debtor = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Debtor, models.DebtorDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.DebtorDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "code", guid)
		},
		func(doc models.DebtorDoc) bool {
			return doc.Code != ""
		},
		func(shopID string, authUsername string, data models.Debtor, doc models.DebtorDoc) error {

			doc.Debtor = data
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

func (svc DebtorHttpService) getDocIDKey(doc models.Debtor) string {
	return doc.Code
}

func (svc DebtorHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc DebtorHttpService) GetModuleName() string {
	return "debtor"
}
