package printer

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/restaurant/printer/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IPrinterService interface {
	CreatePrinter(shopID string, authUsername string, doc models.Printer) (string, error)
	UpdatePrinter(shopID string, guid string, authUsername string, doc models.Printer) error
	DeletePrinter(shopID string, guid string, authUsername string) error
	InfoPrinter(shopID string, guid string) (models.PrinterInfo, error)
	SearchPrinter(shopID string, pageable micromodels.Pageable) ([]models.PrinterInfo, mongopagination.PaginationData, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Printer) (common.BulkImport, error)
	SearchPrinterStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.PrinterInfo, int, error)

	GetModuleName() string
}

type PrinterService struct {
	repo          PrinterRepository
	syncCacheRepo mastersync.IMasterSyncCacheRepository

	services.ActivityService[models.PrinterActivity, models.PrinterDeleteActivity]
}

func NewPrinterService(repo PrinterRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) PrinterService {

	insSvc := PrinterService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.PrinterActivity, models.PrinterDeleteActivity](repo)
	return insSvc
}

func (svc PrinterService) CreatePrinter(shopID string, authUsername string, doc models.Printer) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := models.PrinterDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Printer = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err := svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc PrinterService) UpdatePrinter(shopID string, guid string, authUsername string, doc models.Printer) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.Printer = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc PrinterService) DeletePrinter(shopID string, guid string, authUsername string) error {
	err := svc.repo.DeleteByGuidfixed(shopID, guid, authUsername)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc PrinterService) InfoPrinter(shopID string, guid string) (models.PrinterInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.PrinterInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.PrinterInfo{}, errors.New("document not found")
	}

	return findDoc.PrinterInfo, nil

}

func (svc PrinterService) SearchPrinter(shopID string, pageable micromodels.Pageable) ([]models.PrinterInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"code",
		"name1",
		"name2",
		"name3",
		"name4",
		"name5",
	}

	for i := range [5]bool{} {
		searchInFields = append(searchInFields, fmt.Sprintf("name%d", (i+1)))
	}

	docList, pagination, err := svc.repo.FindPage(shopID, searchInFields, pageable)

	if err != nil {
		return []models.PrinterInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc PrinterService) SearchPrinterStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.PrinterInfo, int, error) {
	searchInFields := []string{
		"code",
		"name1",
		"name2",
		"name3",
		"name4",
		"name5",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.PrinterInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc PrinterService) SaveInBatch(shopID string, authUsername string, dataList []models.Printer) (common.BulkImport, error) {

	// createDataList := []models.PrinterDoc{}
	// duplicateDataList := []models.Printer{}

	payloadCategoryList, payloadDuplicateCategoryList := importdata.FilterDuplicate[models.Printer](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadCategoryList {
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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Printer, models.PrinterDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadCategoryList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Printer) models.PrinterDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.PrinterDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Printer = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			dataDoc.LastUpdatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Printer, models.PrinterDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.PrinterDoc, error) {
			return svc.repo.FindByGuid(shopID, guid)
		},
		func(doc models.PrinterDoc) bool {
			return false
		},
		func(shopID string, authUsername string, data models.Printer, doc models.PrinterDoc) error {

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
	for _, doc := range payloadDuplicateCategoryList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.Code)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {
		updateDataKey = append(updateDataKey, doc.Code)
	}

	updateFailDataKey := []string{}
	for _, doc := range updateFailDataList {
		updateFailDataKey = append(updateFailDataKey, doc.Code)
	}

	svc.saveMasterSync(shopID)

	return common.BulkImport{
		Created:          createDataKey,
		Updated:          updateDataKey,
		UpdateFailed:     updateFailDataKey,
		PayloadDuplicate: payloadDuplicateDataKey,
	}, nil
}

func (svc PrinterService) getDocIDKey(doc models.Printer) string {
	return doc.Code
}

func (svc PrinterService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc PrinterService) GetModuleName() string {
	return "printer"
}
