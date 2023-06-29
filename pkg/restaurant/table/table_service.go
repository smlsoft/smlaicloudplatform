package table

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/restaurant/table/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ITableService interface {
	CreateTable(shopID string, authUsername string, doc models.Table) (string, error)
	UpdateTable(shopID string, guid string, authUsername string, doc models.Table) error
	DeleteTable(shopID string, guid string, authUsername string) error
	InfoTable(shopID string, guid string) (models.TableInfo, error)
	SearchTable(shopID string, pageable micromodels.Pageable) ([]models.TableInfo, mongopagination.PaginationData, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Table) (common.BulkImport, error)

	GetModuleName() string
}

type TableService struct {
	repo          ITableRepository
	syncCacheRepo mastersync.IMasterSyncCacheRepository

	services.ActivityService[models.TableActivity, models.TableDeleteActivity]
}

func NewTableService(repo TableRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *TableService {
	insSvc := TableService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.TableActivity, models.TableDeleteActivity](repo)
	return &insSvc
}

func (svc TableService) CreateTable(shopID string, authUsername string, doc models.Table) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := models.TableDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Table = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err := svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc TableService) UpdateTable(shopID string, guid string, authUsername string, doc models.Table) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.Table = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc TableService) DeleteTable(shopID string, guid string, authUsername string) error {
	err := svc.repo.DeleteByGuidfixed(shopID, guid, authUsername)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc TableService) InfoTable(shopID string, guid string) (models.TableInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.TableInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.TableInfo{}, errors.New("document not found")
	}

	return findDoc.TableInfo, nil

}

func (svc TableService) SearchTable(shopID string, pageable micromodels.Pageable) ([]models.TableInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"code",
		"names.name",
	}

	for i := range [5]bool{} {
		searchInFields = append(searchInFields, fmt.Sprintf("name%d", (i+1)))
	}

	docList, pagination, err := svc.repo.FindPage(shopID, searchInFields, pageable)

	if err != nil {
		return []models.TableInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc TableService) SaveInBatch(shopID string, authUsername string, dataList []models.Table) (common.BulkImport, error) {

	// createDataList := []models.TableDoc{}
	// duplicateDataList := []models.Table{}

	payloadCategoryList, payloadDuplicateCategoryList := importdata.FilterDuplicate[models.Table](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadCategoryList {
		itemCodeGuidList = append(itemCodeGuidList, doc.Number)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "code", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.Code)
	}

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Table, models.TableDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadCategoryList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Table) models.TableDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.TableDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Table = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			dataDoc.LastUpdatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Table, models.TableDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.TableDoc, error) {
			return svc.repo.FindByGuid(shopID, guid)
		},
		func(doc models.TableDoc) bool {
			return false
		},
		func(shopID string, authUsername string, data models.Table, doc models.TableDoc) error {

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
		createDataKey = append(createDataKey, doc.Number)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateCategoryList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.Number)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {
		updateDataKey = append(updateDataKey, doc.Number)
	}

	updateFailDataKey := []string{}
	for _, doc := range updateFailDataList {
		updateFailDataKey = append(updateFailDataKey, doc.Number)
	}

	svc.saveMasterSync(shopID)

	return common.BulkImport{
		Created:          createDataKey,
		Updated:          updateDataKey,
		UpdateFailed:     updateFailDataKey,
		PayloadDuplicate: payloadDuplicateDataKey,
	}, nil
}

func (svc TableService) getDocIDKey(doc models.Table) string {
	return doc.Number
}

func (svc TableService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc TableService) GetModuleName() string {
	return "restauranttable"
}
