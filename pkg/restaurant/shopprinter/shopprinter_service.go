package shopprinter

import (
	"errors"
	"fmt"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/restaurant/shopprinter/models"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"sync"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IShopPrinterService interface {
	CreateShopPrinter(shopID string, authUsername string, doc models.PrinterTerminal) (string, error)
	UpdateShopPrinter(shopID string, guid string, authUsername string, doc models.PrinterTerminal) error
	DeleteShopPrinter(shopID string, guid string, authUsername string) error
	InfoShopPrinter(shopID string, guid string) (models.PrinterTerminalInfo, error)
	SearchShopPrinter(shopID string, q string, page int, limit int) ([]models.PrinterTerminalInfo, mongopagination.PaginationData, error)
	LastActivity(shopID string, lastUpdatedDate time.Time, page int, limit int) (common.LastActivity, mongopagination.PaginationData, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.PrinterTerminal) (common.BulkImport, error)
}

type ShopPrinterService struct {
	repo      ShopPrinterRepository
	cacheRepo mastersync.IMasterSyncCacheRepository
}

func NewShopPrinterService(repo ShopPrinterRepository, cacheRepo mastersync.IMasterSyncCacheRepository) ShopPrinterService {

	return ShopPrinterService{
		repo:      repo,
		cacheRepo: cacheRepo,
	}
}

func (svc ShopPrinterService) CreateShopPrinter(shopID string, authUsername string, doc models.PrinterTerminal) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := models.PrinterTerminalDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.PrinterTerminal = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err := svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc ShopPrinterService) UpdateShopPrinter(shopID string, guid string, authUsername string, doc models.PrinterTerminal) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.PrinterTerminal = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc ShopPrinterService) DeleteShopPrinter(shopID string, guid string, authUsername string) error {
	err := svc.repo.DeleteByGuidfixed(shopID, guid, authUsername)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc ShopPrinterService) InfoShopPrinter(shopID string, guid string) (models.PrinterTerminalInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.PrinterTerminalInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.PrinterTerminalInfo{}, errors.New("document not found")
	}

	return findDoc.PrinterTerminalInfo, nil

}

func (svc ShopPrinterService) SearchShopPrinter(shopID string, q string, page int, limit int) ([]models.PrinterTerminalInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"guidfixed",
		"code",
	}

	for i := range [5]bool{} {
		searchCols = append(searchCols, fmt.Sprintf("name%d", (i+1)))
	}

	docList, pagination, err := svc.repo.FindPage(shopID, searchCols, q, page, limit)

	if err != nil {
		return []models.PrinterTerminalInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc ShopPrinterService) LastActivity(shopID string, lastUpdatedDate time.Time, page int, limit int) (common.LastActivity, mongopagination.PaginationData, error) {
	var wg sync.WaitGroup

	wg.Add(1)
	var deleteDocList []models.PrinterTerminalDeleteActivity
	var pagination1 mongopagination.PaginationData
	var err1 error

	go func() {
		deleteDocList, pagination1, err1 = svc.repo.FindDeletedPage(shopID, lastUpdatedDate, page, limit)
		wg.Done()
	}()

	wg.Add(1)
	var createAndUpdateDocList []models.PrinterTerminalActivity
	var pagination2 mongopagination.PaginationData
	var err2 error

	go func() {
		createAndUpdateDocList, pagination2, err2 = svc.repo.FindCreatedOrUpdatedPage(shopID, lastUpdatedDate, page, limit)
		wg.Done()
	}()

	wg.Wait()

	if err1 != nil {
		return common.LastActivity{}, pagination1, err1
	}

	if err2 != nil {
		return common.LastActivity{}, pagination2, err2
	}

	lastActivity := common.LastActivity{}

	lastActivity.Remove = &deleteDocList
	lastActivity.New = &createAndUpdateDocList

	pagination := pagination1

	if pagination.Total < pagination2.Total {
		pagination = pagination2
	}

	return lastActivity, pagination, nil
}

func (svc ShopPrinterService) SaveInBatch(shopID string, authUsername string, dataList []models.PrinterTerminal) (common.BulkImport, error) {

	// createDataList := []models.PrinterTerminalDoc{}
	// duplicateDataList := []models.PrinterTerminal{}

	payloadCategoryList, payloadDuplicateCategoryList := importdata.FilterDuplicate[models.PrinterTerminal](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.PrinterTerminal, models.PrinterTerminalDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadCategoryList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.PrinterTerminal) models.PrinterTerminalDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.PrinterTerminalDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.PrinterTerminal = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			dataDoc.LastUpdatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.PrinterTerminal, models.PrinterTerminalDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.PrinterTerminalDoc, error) {
			return svc.repo.FindByGuid(shopID, guid)
		},
		func(doc models.PrinterTerminalDoc) bool {
			return false
		},
		func(shopID string, authUsername string, data models.PrinterTerminal, doc models.PrinterTerminalDoc) error {

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

func (svc ShopPrinterService) getDocIDKey(doc models.PrinterTerminal) string {
	return doc.Code
}

func (svc ShopPrinterService) saveMasterSync(shopID string) {
	err := svc.cacheRepo.Save(shopID)

	if err != nil {
		fmt.Println("save shop printer master cache error :: " + err.Error())
	}
}
