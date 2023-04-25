package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/productsection/sectionbusinesstype/models"
	"smlcloudplatform/pkg/productsection/sectionbusinesstype/repositories"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type ISectionBusinessTypeHttpService interface {
	SaveSectionBusinessType(shopID string, authUsername string, doc models.SectionBusinessType) (string, error)
	DeleteSectionBusinessType(shopID string, guid string, authUsername string) error
	DeleteSectionBusinessTypeByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoSectionBusinessType(shopID string, guid string) (models.SectionBusinessTypeInfo, error)
	SearchSectionBusinessType(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SectionBusinessTypeInfo, mongopagination.PaginationData, error)
	SearchSectionBusinessTypeStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.SectionBusinessTypeInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.SectionBusinessType) (common.BulkImport, error)

	GetModuleName() string
}

type SectionBusinessTypeHttpService struct {
	repo repositories.ISectionBusinessTypeRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.SectionBusinessTypeActivity, models.SectionBusinessTypeDeleteActivity]
}

func NewSectionBusinessTypeHttpService(repo repositories.ISectionBusinessTypeRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *SectionBusinessTypeHttpService {

	insSvc := &SectionBusinessTypeHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.SectionBusinessTypeActivity, models.SectionBusinessTypeDeleteActivity](repo)

	return insSvc
}

func (svc SectionBusinessTypeHttpService) SaveSectionBusinessType(shopID string, authUsername string, doc models.SectionBusinessType) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "businesstypecode", doc.BusinessTypeCode)

	if err != nil {
		return "", err
	}

	guidFixed := ""
	if len(findDoc.GuidFixed) < 1 {
		guidFixed, err = svc.create(findDoc, shopID, authUsername, doc)
	} else {
		err = svc.update(findDoc, shopID, authUsername, doc)
		guidFixed = findDoc.GuidFixed
	}

	if err != nil {
		return "", err
	}

	return guidFixed, nil
}

func (svc SectionBusinessTypeHttpService) create(findDoc models.SectionBusinessTypeDoc, shopID string, authUsername string, doc models.SectionBusinessType) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := models.SectionBusinessTypeDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.SectionBusinessType = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err := svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc SectionBusinessTypeHttpService) update(findDoc models.SectionBusinessTypeDoc, shopID string, authUsername string, doc models.SectionBusinessType) error {

	findDoc.SectionBusinessType = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err := svc.repo.Update(shopID, findDoc.GuidFixed, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc SectionBusinessTypeHttpService) DeleteSectionBusinessType(shopID string, guid string, authUsername string) error {

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

func (svc SectionBusinessTypeHttpService) DeleteSectionBusinessTypeByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc SectionBusinessTypeHttpService) InfoSectionBusinessType(shopID string, guid string) (models.SectionBusinessTypeInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.SectionBusinessTypeInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.SectionBusinessTypeInfo{}, errors.New("document not found")
	}

	return findDoc.SectionBusinessTypeInfo, nil
}

func (svc SectionBusinessTypeHttpService) SearchSectionBusinessType(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.SectionBusinessTypeInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"businesstypecode",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.SectionBusinessTypeInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc SectionBusinessTypeHttpService) SearchSectionBusinessTypeStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.SectionBusinessTypeInfo, int, error) {
	searchInFields := []string{
		"businesstypecode",
	}

	selectFields := map[string]interface{}{}

	if langCode != "" {
		selectFields["names"] = bson.M{"$elemMatch": bson.M{"code": langCode}}
	} else {
		selectFields["names"] = 1
	}

	docList, total, err := svc.repo.FindStep(shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.SectionBusinessTypeInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc SectionBusinessTypeHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.SectionBusinessType) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.SectionBusinessType](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.BusinessTypeCode)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "businesstypecode", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.BusinessTypeCode)
	}

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.SectionBusinessType, models.SectionBusinessTypeDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.SectionBusinessType) models.SectionBusinessTypeDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.SectionBusinessTypeDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.SectionBusinessType = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.SectionBusinessType, models.SectionBusinessTypeDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.SectionBusinessTypeDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "businesstypecode", guid)
		},
		func(doc models.SectionBusinessTypeDoc) bool {
			return doc.BusinessTypeCode != ""
		},
		func(shopID string, authUsername string, data models.SectionBusinessType, doc models.SectionBusinessTypeDoc) error {

			doc.SectionBusinessType = data
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
		createDataKey = append(createDataKey, doc.BusinessTypeCode)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.BusinessTypeCode)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {

		updateDataKey = append(updateDataKey, doc.BusinessTypeCode)
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

func (svc SectionBusinessTypeHttpService) getDocIDKey(doc models.SectionBusinessType) string {
	return doc.BusinessTypeCode
}

func (svc SectionBusinessTypeHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc SectionBusinessTypeHttpService) GetModuleName() string {
	return "sectionBusinessType"
}
