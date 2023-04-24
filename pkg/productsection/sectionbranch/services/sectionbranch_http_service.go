package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/productsection/sectionbranch/models"
	"smlcloudplatform/pkg/productsection/sectionbranch/repositories"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type ISectionBranchHttpService interface {
	CreateSectionBranch(shopID string, authUsername string, doc models.SectionBranch) (string, error)
	UpdateSectionBranch(shopID string, guid string, authUsername string, doc models.SectionBranch) error
	DeleteSectionBranch(shopID string, guid string, authUsername string) error
	DeleteSectionBranchByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoSectionBranch(shopID string, guid string) (models.SectionBranchInfo, error)
	InfoSectionBranchByBranchCode(shopID string, branchcode string) (models.SectionBranchInfo, error)
	SearchSectionBranch(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SectionBranchInfo, mongopagination.PaginationData, error)
	SearchSectionBranchStep(shopID string, langBranchCode string, pageableStep micromodels.PageableStep) ([]models.SectionBranchInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.SectionBranch) (common.BulkImport, error)

	GetModuleName() string
}

type SectionBranchHttpService struct {
	repo repositories.ISectionBranchRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.SectionBranchActivity, models.SectionBranchDeleteActivity]
}

func NewSectionBranchHttpService(repo repositories.ISectionBranchRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *SectionBranchHttpService {

	insSvc := &SectionBranchHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.SectionBranchActivity, models.SectionBranchDeleteActivity](repo)

	return insSvc
}

func (svc SectionBranchHttpService) CreateSectionBranch(shopID string, authUsername string, doc models.SectionBranch) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "branchcode", doc.BranchCode)

	if err != nil {
		return "", err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", errors.New("branch code is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.SectionBranchDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.SectionBranch = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc SectionBranchHttpService) UpdateSectionBranch(shopID string, guid string, authUsername string, doc models.SectionBranch) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	findDoc.SectionBranch = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc SectionBranchHttpService) DeleteSectionBranch(shopID string, guid string, authUsername string) error {

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

func (svc SectionBranchHttpService) DeleteSectionBranchByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc SectionBranchHttpService) InfoSectionBranch(shopID string, guid string) (models.SectionBranchInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.SectionBranchInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.SectionBranchInfo{}, errors.New("document not found")
	}

	return findDoc.SectionBranchInfo, nil
}

func (svc SectionBranchHttpService) InfoSectionBranchByBranchCode(shopID string, branchcode string) (models.SectionBranchInfo, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "branchcode", branchcode)

	if err != nil {
		return models.SectionBranchInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.SectionBranchInfo{}, errors.New("document not found")
	}

	return findDoc.SectionBranchInfo, nil
}

func (svc SectionBranchHttpService) SearchSectionBranch(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SectionBranchInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"branchcode",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.SectionBranchInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc SectionBranchHttpService) SearchSectionBranchStep(shopID string, langBranchCode string, pageableStep micromodels.PageableStep) ([]models.SectionBranchInfo, int, error) {
	searchInFields := []string{
		"branchcode",
	}

	selectFields := map[string]interface{}{}

	if langBranchCode != "" {
		selectFields["names"] = bson.M{"$elemMatch": bson.M{"branchcode": langBranchCode}}
	} else {
		selectFields["names"] = 1
	}

	docList, total, err := svc.repo.FindStep(shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.SectionBranchInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc SectionBranchHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.SectionBranch) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.SectionBranch](dataList, svc.getDocIDKey)

	itemBranchCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemBranchCodeGuidList = append(itemBranchCodeGuidList, doc.BranchCode)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "branchcode", itemBranchCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.BranchCode)
	}

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.SectionBranch, models.SectionBranchDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.SectionBranch) models.SectionBranchDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.SectionBranchDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.SectionBranch = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.SectionBranch, models.SectionBranchDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.SectionBranchDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "branchcode", guid)
		},
		func(doc models.SectionBranchDoc) bool {
			return doc.BranchCode != ""
		},
		func(shopID string, authUsername string, data models.SectionBranch, doc models.SectionBranchDoc) error {

			doc.SectionBranch = data
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
		createDataKey = append(createDataKey, doc.BranchCode)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.BranchCode)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {

		updateDataKey = append(updateDataKey, doc.BranchCode)
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

func (svc SectionBranchHttpService) getDocIDKey(doc models.SectionBranch) string {
	return doc.BranchCode
}

func (svc SectionBranchHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc SectionBranchHttpService) GetModuleName() string {
	return "productSectionBranch"
}
