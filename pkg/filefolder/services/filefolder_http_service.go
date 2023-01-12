package services

import (
	"errors"
	"fmt"
	"smlcloudplatform/pkg/filefolder/models"
	"smlcloudplatform/pkg/filefolder/repositories"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IFileFolderHttpService interface {
	CreateFileFolder(shopID string, authUsername string, doc models.FileFolder) (string, error)
	UpdateFileFolder(shopID string, guid string, authUsername string, doc models.FileFolder) error
	DeleteFileFolder(shopID string, guid string, authUsername string) error
	DeleteFileFolderByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoFileFolder(shopID string, guid string) (models.FileFolderInfo, error)
	SearchFileFolder(shopID string, q string, page int, limit int, sort map[string]int) ([]models.FileFolderInfo, mongopagination.PaginationData, error)
	SearchFileFolderStep(shopID string, langCode string, q string, skip int, limit int, sort map[string]int) ([]models.FileFolderInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.FileFolder) (common.BulkImport, error)

	GetModuleName() string
}

type FileFolderHttpService struct {
	repo repositories.IFileFolderRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.FileFolderActivity, models.FileFolderDeleteActivity]
}

func NewFileFolderHttpService(repo repositories.IFileFolderRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *FileFolderHttpService {

	insSvc := &FileFolderHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.FileFolderActivity, models.FileFolderDeleteActivity](repo)

	return insSvc
}

func (svc FileFolderHttpService) CreateFileFolder(shopID string, authUsername string, doc models.FileFolder) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "name", doc.Name)

	if err != nil {
		return "", err
	}

	if findDoc.Name != "" {
		return "", errors.New("Name is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.FileFolderDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.FileFolder = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc FileFolderHttpService) UpdateFileFolder(shopID string, guid string, authUsername string, doc models.FileFolder) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.FileFolder = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc FileFolderHttpService) DeleteFileFolder(shopID string, guid string, authUsername string) error {

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

func (svc FileFolderHttpService) DeleteFileFolderByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc FileFolderHttpService) InfoFileFolder(shopID string, guid string) (models.FileFolderInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.FileFolderInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.FileFolderInfo{}, errors.New("document not found")
	}

	return findDoc.FileFolderInfo, nil

}

func (svc FileFolderHttpService) SearchFileFolder(shopID string, q string, page int, limit int, sort map[string]int) ([]models.FileFolderInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"name",
	}

	docList, pagination, err := svc.repo.FindPageSort(shopID, searchCols, q, page, limit, sort)

	if err != nil {
		return []models.FileFolderInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc FileFolderHttpService) SearchFileFolderStep(shopID string, langCode string, q string, skip int, limit int, sort map[string]int) ([]models.FileFolderInfo, int, error) {
	searchCols := []string{
		"name",
	}

	projectQuery := map[string]interface{}{}

	docList, total, err := svc.repo.FindLimit(shopID, map[string]interface{}{}, searchCols, q, skip, limit, sort, projectQuery)

	if err != nil {
		return []models.FileFolderInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc FileFolderHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.FileFolder) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.FileFolder](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.Name)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "name", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.Name)
	}

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.FileFolder, models.FileFolderDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.FileFolder) models.FileFolderDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.FileFolderDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.FileFolder = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.FileFolder, models.FileFolderDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.FileFolderDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "name", guid)
		},
		func(doc models.FileFolderDoc) bool {
			return doc.Name != ""
		},
		func(shopID string, authUsername string, data models.FileFolder, doc models.FileFolderDoc) error {

			doc.FileFolder = data
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
		createDataKey = append(createDataKey, doc.Name)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.Name)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {

		updateDataKey = append(updateDataKey, doc.Name)
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

func (svc FileFolderHttpService) getDocIDKey(doc models.FileFolder) string {
	return doc.Name
}

func (svc FileFolderHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc FileFolderHttpService) GetModuleName() string {
	return "fileFolder"
}
