package services

import (
	"errors"
	"fmt"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/barcodemaster/models"
	"smlcloudplatform/pkg/product/barcodemaster/repositories"
	"smlcloudplatform/pkg/utils"
	"sync"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IBarcodeMasterService interface {
	IsExistsGuid(shopID string, guidFixed string) (bool, error)
	CreateInBatch(shopID string, authUsername string, inventories []models.BarcodeMaster) (models.BarcodeMasterBulkImport, error)
	CreateWithGuid(shopID string, authUsername string, guidFixed string, barcodemaster models.BarcodeMaster) (string, error)
	CreateBarcodeMaster(shopID string, authUsername string, barcodemaster models.BarcodeMaster) (string, string, error)
	UpdateBarcodeMaster(shopID string, guid string, authUsername string, barcodemaster models.BarcodeMaster) error
	DeleteBarcodeMaster(shopID string, guid string, username string) error
	InfoBarcodeMaster(shopID string, guid string) (models.BarcodeMasterInfo, error)
	InfoMongoBarcodeMaster(id string) (models.BarcodeMasterInfo, error)
	SearchBarcodeMaster(shopID string, q string, page int, limit int) ([]models.BarcodeMasterInfo, paginate.PaginationData, error)
	LastActivity(shopID string, lastUpdatedDate time.Time, page int, limit int) (common.LastActivity, paginate.PaginationData, error)
	UpdateProductCategory(shopID string, authUsername string, catId string, guid []string) error
}

type BarcodeMasterService struct {
	invRepo   repositories.IBarcodeMasterRepository
	invMqRepo repositories.IBarcodeMasterMQRepository
	cacheRepo mastersync.IMasterSyncCacheRepository
}

func NewBarcodeMasterService(barcodemasterRepo repositories.IBarcodeMasterRepository, barcodemasterMqRepo repositories.IBarcodeMasterMQRepository, cacheRepo mastersync.IMasterSyncCacheRepository) BarcodeMasterService {
	return BarcodeMasterService{
		invRepo:   barcodemasterRepo,
		invMqRepo: barcodemasterMqRepo,
		cacheRepo: cacheRepo,
	}
}

func (svc BarcodeMasterService) IsExistsGuid(shopID string, guidFixed string) (bool, error) {

	findDoc, err := svc.invRepo.FindByGuid(shopID, guidFixed)

	if err != nil {
		return false, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return false, nil
	}

	return true, nil
}

func (svc BarcodeMasterService) CreateInBatch(shopID string, authUsername string, inventories []models.BarcodeMaster) (models.BarcodeMasterBulkImport, error) {

	createDataList := []models.BarcodeMasterDoc{}
	duplicateDataList := []models.BarcodeMaster{}

	payloadBarcodeMasterList, payloadDuplicateBarcodeMasterList := filterDuplicateBarcodeMaster(inventories)

	itemCodeGuidList := []string{}
	for _, barcodemaster := range payloadBarcodeMasterList {
		itemCodeGuidList = append(itemCodeGuidList, barcodemaster.ItemGuid)
	}

	findItemGuid, err := svc.invRepo.FindByItemCodeGuid(shopID, itemCodeGuidList)

	if err != nil {
		return models.BarcodeMasterBulkImport{}, err
	}

	duplicateDataList, createDataList = preparePayloadDataBarcodeMaster(shopID, authUsername, findItemGuid, payloadBarcodeMasterList)

	updateSuccessDataList, updateFailDataList := updateOnDuplicateBarcodeMaster(shopID, authUsername, duplicateDataList, svc.invRepo)

	if len(createDataList) > 0 {
		err = svc.invRepo.CreateInBatch(createDataList)

		if err != nil {
			return models.BarcodeMasterBulkImport{}, err
		}
	}
	createDataKey := []string{}

	for _, inv := range createDataList {
		createDataKey = append(createDataKey, inv.ItemGuid)

		// reply kafka
		if svc.invMqRepo != nil {
			err = svc.invMqRepo.Create(inv.BarcodeMasterData)

			if err != nil {
				return models.BarcodeMasterBulkImport{}, err
			}
		}
	}

	payloadDuplicateDataKey := []string{}
	for _, inv := range payloadDuplicateBarcodeMasterList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, inv.ItemGuid)
	}

	updateDataKey := []string{}
	for _, inv := range updateSuccessDataList {
		updateDataKey = append(updateDataKey, inv.ItemGuid)
		// reply kafka
		if svc.invMqRepo != nil {
			err = svc.invMqRepo.Update(inv.BarcodeMasterData)

			if err != nil {
				return models.BarcodeMasterBulkImport{}, err
			}
		}
	}

	updateFailDataKey := []string{}
	for _, inv := range updateFailDataList {
		updateFailDataKey = append(updateFailDataKey, inv.ItemGuid)
	}

	svc.saveMasterSync(shopID)

	return models.BarcodeMasterBulkImport{
		Created:          createDataKey,
		Updated:          updateDataKey,
		UpdateFailed:     updateFailDataKey,
		PayloadDuplicate: payloadDuplicateDataKey,
	}, nil
}

func filterDuplicateBarcodeMaster(inventories []models.BarcodeMaster) (itemTemp []models.BarcodeMaster, itemDuplicate []models.BarcodeMaster) {
	tempFilterDict := map[string]models.BarcodeMaster{}
	for _, barcodemaster := range inventories {
		if _, ok := tempFilterDict[barcodemaster.ItemGuid]; ok {
			itemDuplicate = append(itemDuplicate, barcodemaster)

		}
		tempFilterDict[barcodemaster.ItemGuid] = barcodemaster
	}

	for _, barcodemaster := range tempFilterDict {
		itemTemp = append(itemTemp, barcodemaster)
	}

	return itemTemp, itemDuplicate
}

func updateOnDuplicateBarcodeMaster(shopID string, authUsername string, duplicateDataList []models.BarcodeMaster, repo repositories.IBarcodeMasterRepository) ([]models.BarcodeMasterDoc, []models.BarcodeMaster) {
	updateSuccessDataList := []models.BarcodeMasterDoc{}
	updateFailDataList := []models.BarcodeMaster{}
	for _, inv := range duplicateDataList {
		findDoc, err := repo.FindByItemGuid(shopID, inv.ItemGuid)

		if err != nil || findDoc.ID == primitive.NilObjectID {
			updateFailDataList = append(updateFailDataList, inv)
			continue
		}

		findDoc.BarcodeMaster = inv

		findDoc.UpdatedBy = authUsername
		findDoc.UpdatedAt = time.Now()
		findDoc.LastUpdatedAt = time.Now()

		err = repo.Update(shopID, findDoc.GuidFixed, findDoc)

		if err != nil {
			updateFailDataList = append(updateFailDataList, inv)
			continue
		}

		updateSuccessDataList = append(updateSuccessDataList, findDoc)
	}
	return updateSuccessDataList, updateFailDataList
}

func preparePayloadDataBarcodeMaster(shopID string, authUsername string, findItemGuid []models.BarcodeMasterItemGuid, payloadBarcodeMasterList []models.BarcodeMaster) ([]models.BarcodeMaster, []models.BarcodeMasterDoc) {
	createDataList := []models.BarcodeMasterDoc{}
	duplicateDataList := []models.BarcodeMaster{}
	tempItemGuidDict := make(map[string]bool)

	for _, itemGuid := range findItemGuid {
		tempItemGuidDict[itemGuid.ItemGuid] = true
	}

	for _, barcodemaster := range payloadBarcodeMasterList {

		if _, ok := tempItemGuidDict[barcodemaster.ItemGuid]; ok {
			duplicateDataList = append(duplicateDataList, barcodemaster)
		} else {
			newGuid := utils.NewGUID()

			invDoc := models.BarcodeMasterDoc{}

			invDoc.GuidFixed = newGuid
			invDoc.ShopID = shopID
			invDoc.BarcodeMaster = barcodemaster

			invDoc.CreatedBy = authUsername
			invDoc.CreatedAt = time.Now()
			invDoc.LastUpdatedAt = time.Now()

			createDataList = append(createDataList, invDoc)
		}
	}
	return duplicateDataList, createDataList
}

func (svc BarcodeMasterService) CreateWithGuid(shopID string, authUsername string, guidFixed string, barcodemaster models.BarcodeMaster) (string, error) {

	reqBarcodes := []string{}

	for _, barcode := range *barcodemaster.Barcodes {
		reqBarcodes = append(reqBarcodes, barcode.Barcode)
	}

	findDocBarcodes, err := svc.invRepo.FindByBarcodes(shopID, reqBarcodes)

	if err != nil {
		return "", err
	}

	if len(findDocBarcodes) > 0 {
		tempBarcode := *findDocBarcodes[0].Barcodes
		return "", fmt.Errorf("barcode '%s' is exists", tempBarcode[0].Barcode)
	}

	newGuid := guidFixed

	invDoc := models.BarcodeMasterDoc{}

	invDoc.GuidFixed = newGuid
	invDoc.ShopID = shopID
	invDoc.BarcodeMaster = barcodemaster

	invDoc.CreatedBy = authUsername
	invDoc.CreatedAt = time.Now()
	invDoc.LastUpdatedAt = time.Now()

	mongoIdx, err := svc.invRepo.Create(invDoc)

	if err != nil {
		return "", err
	}

	if svc.invMqRepo != nil {
		err = svc.invMqRepo.Create(invDoc.BarcodeMasterData)

		if err != nil {
			return "", err
		}
	}

	svc.saveMasterSync(shopID)

	return mongoIdx, nil
}

// func (svc BarcodeMasterService) CreateBulk(shopID string, authUsername string, guidFixed string, inventories []models.BarcodeMaster) (error) {

// 	newGuid := guidFixed

// 	invDocList := make([]models.BarcodeMasterDoc, len(inventories))

// 	for index, inv := range inventories {
// 		invDoc := models.BarcodeMasterDoc{}

// 		invDoc.GuidFixed = newGuid
// 		invDoc.ShopID = shopID
// 		invDoc.BarcodeMaster = inv

// 		invDoc.CreatedBy = authUsername
// 		invDoc.CreatedAt = time.Now()

// 		invDocList[index] = invDoc
// 	}

// 	mongoIdx, err := svc.invRepo.Create(invDoc)

// 	if err != nil {
// 		return "", err
// 	}

// 	if svc.invMqRepo != nil {
// 		err = svc.invMqRepo.Create(invDoc.BarcodeMasterData)

// 		if err != nil {
// 			return "", err
// 		}
// 	}

// 	return mongoIdx, nil
// }

func (svc BarcodeMasterService) CreateBarcodeMaster(shopID string, authUsername string, barcodemaster models.BarcodeMaster) (string, string, error) {

	newGuid := utils.NewGUID()
	mongoIdx, err := svc.CreateWithGuid(shopID, authUsername, newGuid, barcodemaster)
	return mongoIdx, newGuid, err
}

func (svc BarcodeMasterService) UpdateBarcodeMaster(shopID string, guid string, authUsername string, barcodemaster models.BarcodeMaster) error {

	findDoc, err := svc.invRepo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	if len(*barcodemaster.Barcodes) > 0 {
		reqBarcodes := []string{}
		idxReqBarcode := map[string]struct{}{}

		for _, barcode := range *barcodemaster.Barcodes {
			reqBarcodes = append(reqBarcodes, barcode.Barcode)
			idxReqBarcode[barcode.Barcode] = struct{}{}
		}

		findDocBarcodes, err := svc.invRepo.FindByBarcodes(shopID, reqBarcodes)

		if err != nil {
			return err
		}

		for _, doc := range findDocBarcodes {
			for _, barcode := range *doc.Barcodes {

				if doc.GuidFixed != findDoc.GuidFixed {
					if _, ok := idxReqBarcode[barcode.Barcode]; ok {
						return fmt.Errorf("barcode %s is exists", barcode.Barcode)
					}
				}
			}
		}
	}

	findDoc.BarcodeMaster = barcodemaster

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()
	findDoc.LastUpdatedAt = time.Now()

	err = svc.invRepo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	err = svc.invMqRepo.Update(findDoc.BarcodeMasterData)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc BarcodeMasterService) DeleteBarcodeMaster(shopID string, guid string, username string) error {

	err := svc.invRepo.Delete(shopID, guid, username)

	if err != nil {
		return err
	}

	docIndentity := common.Identity{
		ShopID:    shopID,
		GuidFixed: guid,
	}

	err = svc.invMqRepo.Delete(docIndentity)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc BarcodeMasterService) InfoMongoBarcodeMaster(id string) (models.BarcodeMasterInfo, error) {
	start := time.Now()

	idx, err := primitive.ObjectIDFromHex(id)
	findDoc, err := svc.invRepo.FindByID(idx)

	if err != nil && err.Error() != "mongo: no documents in result" {
		return models.BarcodeMasterInfo{}, err
	}

	elapsed := time.Since(start)
	fmt.Printf("mongo :: pure id :: %s\n", elapsed)

	return findDoc.BarcodeMasterInfo, nil
}

func (svc BarcodeMasterService) InfoBarcodeMaster(shopID string, guid string) (models.BarcodeMasterInfo, error) {

	findDoc, err := svc.invRepo.FindByGuid(shopID, guid)

	if err != nil && err.Error() != "mongo: no documents in result" {
		return models.BarcodeMasterInfo{}, err
	}

	// if findDoc.ID == primitive.NilObjectID {
	// 	return models.BarcodeMasterInfo{}, nil
	// }

	return findDoc.BarcodeMasterInfo, nil
}

func (svc BarcodeMasterService) SearchBarcodeMaster(shopID string, q string, page int, limit int) ([]models.BarcodeMasterInfo, paginate.PaginationData, error) {
	docList, pagination, err := svc.invRepo.FindPage(shopID, q, page, limit)

	if err != nil {
		return []models.BarcodeMasterInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc BarcodeMasterService) LastActivity(shopID string, lastUpdatedDate time.Time, page int, limit int) (common.LastActivity, paginate.PaginationData, error) {

	var wg sync.WaitGroup

	wg.Add(1)
	var deleteDocList []models.BarcodeMasterDeleteActivity
	var pagination1 paginate.PaginationData
	var err1 error

	go func() {
		deleteDocList, pagination1, err1 = svc.invRepo.FindDeletedPage(shopID, lastUpdatedDate, page, limit)
		wg.Done()
	}()

	wg.Add(1)
	var createAndUpdateDocList []models.BarcodeMasterActivity
	var pagination2 paginate.PaginationData
	var err2 error

	go func() {
		createAndUpdateDocList, pagination2, err2 = svc.invRepo.FindCreatedOrUpdatedPage(shopID, lastUpdatedDate, page, limit)
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

func (svc BarcodeMasterService) UpdateProductCategory(shopID string, authUsername string, catId string, guids []string) error {

	for _, guid := range guids {

		findDoc, err := svc.invRepo.FindByItemGuid(shopID, guid)

		if err != nil {
			return err
		}

		if findDoc.ID == primitive.NilObjectID {
			return errors.New("document not found")
		}

		findDoc.CategoryGuid = catId
		findDoc.UpdatedBy = authUsername
		findDoc.UpdatedAt = time.Now()

		err = svc.invRepo.Update(shopID, findDoc.GuidFixed, findDoc)

		if err != nil {
			return err
		}

		err = svc.invMqRepo.Update(findDoc.BarcodeMasterData)

		if err != nil {
			return err
		}
	}

	if len(guids) > 0 {
		svc.saveMasterSync(shopID)
	}

	return nil
}

func (svc BarcodeMasterService) saveMasterSync(shopID string) {
	err := svc.cacheRepo.Save(shopID)

	if err != nil {
		fmt.Println("save barcodemaster master cache error :: " + err.Error())
	}
}
