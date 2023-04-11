package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"smlcloudplatform/pkg/warehouse/models"
	"smlcloudplatform/pkg/warehouse/repositories"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IWarehouseHttpService interface {
	CreateWarehouse(shopID string, authUsername string, doc models.Warehouse) (string, error)
	UpdateWarehouse(shopID string, guid string, authUsername string, doc models.Warehouse) error
	DeleteWarehouse(shopID string, guid string, authUsername string) error
	DeleteWarehouseByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoWarehouse(shopID string, guid string) (models.WarehouseInfo, error)
	InfoWarehouseByCode(shopID string, code string) (models.WarehouseInfo, error)
	SearchWarehouse(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.WarehouseInfo, mongopagination.PaginationData, error)
	SearchWarehouseStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.WarehouseInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Warehouse) (common.BulkImport, error)

	SearchLocation(shopID string, pageable micromodels.Pageable) ([]models.LocationInfo, mongopagination.PaginationData, error)
	SearchShelf(shopID string, pageable micromodels.Pageable) ([]models.ShelfInfo, mongopagination.PaginationData, error)

	InfoLocation(shopID, warehouseCode, locationCode string) (models.LocationInfo, error)
	CreateLocation(shopID, authUsername, warehouseCode string, doc models.LocationRequest) error
	UpdateLocation(shopID, authUsername, warehouseCode, locationCode string, doc models.LocationRequest) error

	InfoShelf(shopID, warehouseCode, locationCode, shelfCode string) (models.ShelfInfo, error)
	CreateShelf(shopID, authUsername, warehouseCode, locationCode string, doc models.ShelfRequest) error
	UpdateShelf(shopID, authUsername, warehouseCode, locationCode, shelfCode string, doc models.ShelfRequest) error

	GetModuleName() string
}

type WarehouseHttpService struct {
	repo repositories.IWarehouseRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.WarehouseActivity, models.WarehouseDeleteActivity]
}

func NewWarehouseHttpService(repo repositories.IWarehouseRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *WarehouseHttpService {

	insSvc := &WarehouseHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.WarehouseActivity, models.WarehouseDeleteActivity](repo)

	return insSvc
}

func (svc WarehouseHttpService) CreateWarehouse(shopID string, authUsername string, doc models.Warehouse) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", doc.Code)

	if err != nil {
		return "", err
	}

	if findDoc.Code != "" {
		return "", errors.New("Code is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.WarehouseDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Warehouse = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc WarehouseHttpService) UpdateWarehouse(shopID string, guid string, authUsername string, doc models.Warehouse) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.Warehouse = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc WarehouseHttpService) CreateLocation(shopID, authUsername, warehouseCode string, doc models.LocationRequest) error {

	findDoc, err := svc.repo.FindWarehouseByLocation(shopID, warehouseCode, doc.Code)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	locations := findDoc.Warehouse.Location

	for _, location := range *locations {
		if location.Code == doc.Code {
			return errors.New("location code is exists")
		}
	}

	*findDoc.Location = append(*findDoc.Location, models.Location{
		Code:  doc.Code,
		Names: doc.Names,
	})

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, findDoc.GuidFixed, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc WarehouseHttpService) UpdateLocation(shopID, authUsername, warehouseCode, locationCode string, doc models.LocationRequest) error {

	updateDoc := models.WarehouseDoc{}
	removeDoc := models.WarehouseDoc{}

	findDoc, err := svc.repo.FindWarehouseByLocation(shopID, warehouseCode, locationCode)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	if warehouseCode != doc.WarehouseCode {

		findDocWarehouse, err := svc.repo.FindByDocIndentityGuid(shopID, "code", doc.WarehouseCode)

		if err != nil {
			return err
		}

		if len(findDocWarehouse.GuidFixed) < 1 {
			return errors.New("document not found")
		}

		updateDoc = findDocWarehouse

		// clear doc
		removeDoc = findDoc
		tempLocation := []models.Location{}

		for _, location := range *removeDoc.Location {
			if location.Code != locationCode {
				tempLocation = append(tempLocation, location)
			}
		}

		removeDoc.Location = &tempLocation
	} else {
		updateDoc = findDoc
	}

	if warehouseCode == doc.WarehouseCode {
		if locationCode != doc.Code {
			locations := updateDoc.Warehouse.Location

			for _, location := range *locations {
				if location.Code == doc.Code {
					return errors.New("location code is exists")
				}
			}

			*updateDoc.Location = append(*updateDoc.Location, models.Location{
				Code:  doc.Code,
				Names: doc.Names,
				Shelf: &doc.Shelf,
			})

		} else {
			for i, location := range *updateDoc.Location {
				if location.Code == doc.Code {
					location.Names = doc.Names
					(*updateDoc.Location)[i] = location
				}
			}
		}
	} else {
		for _, location := range *updateDoc.Location {
			if location.Code == doc.Code {
				return errors.New("location code is exists")
			}
		}

		*updateDoc.Location = append(*updateDoc.Location, models.Location{
			Code:  doc.Code,
			Names: doc.Names,
			Shelf: &doc.Shelf,
		})
	}

	err = svc.repo.Transaction(func() error {
		updateDoc.UpdatedBy = authUsername
		updateDoc.UpdatedAt = time.Now()

		err = svc.repo.Update(shopID, updateDoc.GuidFixed, updateDoc)

		if err != nil {
			return err
		}

		if len(removeDoc.GuidFixed) > 0 {
			removeDoc.UpdatedBy = authUsername
			removeDoc.UpdatedAt = time.Now()

			err = svc.repo.Update(shopID, removeDoc.GuidFixed, removeDoc)

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc WarehouseHttpService) UpdateLocation2(shopID, authUsername, warehouseCode, locationCode string, doc models.LocationRequest) error {
	// Retrieve the original warehouse containing the location
	foundDoc, err := svc.repo.FindWarehouseByLocation(shopID, warehouseCode, locationCode)
	if err != nil {
		return err
	}

	if len(foundDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	// If the warehouse code is different, find the new warehouse to move the location to
	var targetDoc models.WarehouseDoc
	if warehouseCode != doc.WarehouseCode {
		targetDoc, err = svc.repo.FindByDocIndentityGuid(shopID, "code", doc.WarehouseCode)
		if err != nil {
			return err
		}

		if len(targetDoc.GuidFixed) < 1 {
			return errors.New("document not found")
		}
	} else {
		targetDoc = foundDoc
	}

	// Check if the new location code already exists in the target warehouse
	if warehouseCode != doc.WarehouseCode {
		for _, location := range *targetDoc.Location {
			if location.Code == doc.Code {
				return errors.New("location code is exists")
			}
		}
	}

	// Update the location in the target warehouse and remove it from the original warehouse if necessary
	err = svc.updateLocationInWarehouse(shopID, authUsername, foundDoc, targetDoc, doc, warehouseCode, locationCode)
	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)
	return nil
}

func (svc WarehouseHttpService) updateLocationInWarehouse(shopID, authUsername string, foundDoc, targetDoc models.WarehouseDoc, doc models.LocationRequest, warehouseCode, locationCode string) error {
	// Create the updated location object
	updatedLocation := models.Location{
		Code:  doc.Code,
		Names: doc.Names,
		Shelf: &doc.Shelf,
	}

	// Update the target warehouse's location list
	*targetDoc.Location = append(*targetDoc.Location, updatedLocation)

	// Remove the location from the original warehouse if it's different from the target warehouse
	if warehouseCode != doc.WarehouseCode {
		*foundDoc.Location = removeLocationFromWarehouse(*foundDoc.Location, locationCode)
	}

	// Perform database updates in a transaction
	err := svc.repo.Transaction(func() error {
		// Update the target warehouse
		targetDoc.UpdatedBy = authUsername
		targetDoc.UpdatedAt = time.Now()
		err := svc.repo.Update(shopID, targetDoc.GuidFixed, targetDoc)
		if err != nil {
			return err
		}

		// Update the original warehouse if necessary
		if warehouseCode != doc.WarehouseCode {
			foundDoc.UpdatedBy = authUsername
			foundDoc.UpdatedAt = time.Now()
			err = svc.repo.Update(shopID, foundDoc.GuidFixed, foundDoc)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func removeLocationFromWarehouse(locations []models.Location, locationCode string) []models.Location {
	updatedLocations := []models.Location{}

	for _, location := range locations {
		if location.Code != locationCode {
			updatedLocations = append(updatedLocations, location)
		}
	}

	return updatedLocations
}

func (svc WarehouseHttpService) CreateShelf(shopID, authUsername, warehouseCode, locationCode string, doc models.ShelfRequest) error {
	findDoc, err := svc.repo.FindWarehouseByShelf(shopID, warehouseCode, locationCode, doc.Code)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	locations := findDoc.Warehouse.Location

	for indexLocation, location := range *locations {

		if location.Code == locationCode {
			shelves := location.Shelf

			for _, shelf := range *shelves {
				if shelf.Code == doc.Code {
					return errors.New("shelf Code is exists")
				}
			}

			tempLocation := (*findDoc.Location)[indexLocation]
			tempShelf := *tempLocation.Shelf

			tempShelf = append(tempShelf, models.Shelf{
				Code: doc.Code,
				Name: doc.Name,
			})

			(*findDoc.Location)[indexLocation].Shelf = &tempShelf

			break
		}
	}

	return nil
}

func (svc WarehouseHttpService) UpdateShelf(shopID, authUsername, warehouseCode, locationCode, shelfCode string, doc models.ShelfRequest) error {

	updateDoc := models.WarehouseDoc{}
	removeDoc := models.WarehouseDoc{}

	findDoc, err := svc.repo.FindWarehouseByShelf(shopID, warehouseCode, locationCode, shelfCode)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	if warehouseCode != doc.WarehouseCode {

		findDocWarehouse, err := svc.repo.FindByDocIndentityGuid(shopID, "code", doc.WarehouseCode)

		if err != nil {
			return err
		}

		if len(findDocWarehouse.GuidFixed) < 1 {
			return errors.New("document not found")
		}

		updateDoc = findDocWarehouse

		// clear doc
		removeDoc = findDoc
		tempRemoveShelf := []models.Shelf{}

		for i, location := range *removeDoc.Location {
			if location.Code == locationCode {
				for _, shelf := range *location.Shelf {
					if shelf.Code != shelfCode {
						tempRemoveShelf = append(tempRemoveShelf, shelf)
					}
				}

				(*removeDoc.Location)[i].Shelf = &tempRemoveShelf
			}
		}

	} else {
		updateDoc = findDoc
	}

	if warehouseCode == doc.WarehouseCode {
		if shelfCode != doc.Code {

			locations := updateDoc.Warehouse.Location

			for indexLocation, location := range *locations {

				if location.Code == locationCode {
					shelves := location.Shelf

					for _, shelf := range *shelves {
						if shelf.Code == doc.Code {
							return errors.New("shelf code is exists")
						}
					}

					tempLocation := (*updateDoc.Location)[indexLocation]
					tempShelf := *tempLocation.Shelf

					newShelf := models.Shelf{
						Code: doc.Code,
						Name: doc.Name,
					}

					tempShelf = append(tempShelf, newShelf)

					(*updateDoc.Location)[indexLocation].Shelf = &tempShelf

				}
			}

		} else {
			locations := updateDoc.Warehouse.Location

			for indexLocation, location := range *locations {
				shelves := *location.Shelf
				for indexShelf, shelf := range shelves {
					if shelf.Code == doc.Code {
						tempLocation := (*updateDoc.Location)[indexLocation]
						tempShelf := *tempLocation.Shelf
						tempShelf[indexShelf].Name = doc.Name
						(*updateDoc.Location)[indexLocation].Shelf = &tempShelf
					}
				}
			}

		}
	} else {

		locations := updateDoc.Warehouse.Location

		for indexLocation, location := range *locations {
			if location.Shelf == nil {
				(*updateDoc.Location)[indexLocation].Shelf = &[]models.Shelf{}
			} else {

				shelves := *location.Shelf

				for _, shelf := range shelves {
					if shelf.Code == doc.Code {
						return errors.New("shelf code is exists")
					}
				}
			}

			tempShelfs := (*updateDoc.Location)[indexLocation].Shelf
			*(*updateDoc.Location)[indexLocation].Shelf = append(*tempShelfs, models.Shelf{
				Code: doc.Code,
				Name: doc.Name,
			})
		}
	}

	err = svc.repo.Transaction(func() error {
		updateDoc.UpdatedBy = authUsername
		updateDoc.UpdatedAt = time.Now()

		err = svc.repo.Update(shopID, updateDoc.GuidFixed, updateDoc)

		if err != nil {
			return err
		}

		if len(removeDoc.GuidFixed) > 0 {
			removeDoc.UpdatedBy = authUsername
			removeDoc.UpdatedAt = time.Now()

			err = svc.repo.Update(shopID, removeDoc.GuidFixed, removeDoc)

			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (svc WarehouseHttpService) DeleteWarehouse(shopID string, guid string, authUsername string) error {

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

func (svc WarehouseHttpService) DeleteWarehouseByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc WarehouseHttpService) InfoWarehouse(shopID string, guid string) (models.WarehouseInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.WarehouseInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.WarehouseInfo{}, errors.New("document not found")
	}

	return findDoc.WarehouseInfo, nil
}

func (svc WarehouseHttpService) InfoWarehouseByCode(shopID string, code string) (models.WarehouseInfo, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", code)

	if err != nil {
		return models.WarehouseInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.WarehouseInfo{}, errors.New("document not found")
	}

	return findDoc.WarehouseInfo, nil
}

func (svc WarehouseHttpService) InfoLocation(shopID, warehouseCode, locationCode string) (models.LocationInfo, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", warehouseCode)

	if err != nil {
		return models.LocationInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.LocationInfo{}, errors.New("document not found")
	}

	locationInfo := models.LocationInfo{}

	locationInfo.GuidFixed = findDoc.GuidFixed
	locationInfo.WarehouseCode = findDoc.Code
	locationInfo.WarehouseNames = findDoc.Names

	for _, location := range *findDoc.Location {
		if location.Code == locationCode {
			locationInfo.LocationCode = location.Code
			locationInfo.LocationNames = location.Names
			locationInfo.Shelf = *location.Shelf
			break
		}
	}

	return locationInfo, nil
}

func (svc WarehouseHttpService) InfoShelf(shopID, warehouseCode, locationCode, shelfCode string) (models.ShelfInfo, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", warehouseCode)

	if err != nil {
		return models.ShelfInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.ShelfInfo{}, errors.New("document not found")
	}

	shelfInfo := models.ShelfInfo{}

	shelfInfo.GuidFixed = findDoc.GuidFixed
	shelfInfo.WarehouseCode = findDoc.Code
	shelfInfo.WarehouseNames = findDoc.Names

	for _, location := range *findDoc.Location {
		if location.Code == locationCode {
			shelfInfo.LocationCode = location.Code
			shelfInfo.LocationNames = location.Names
			for _, shelf := range *location.Shelf {
				if shelf.Code == shelfCode {
					shelfInfo.ShelfCode = shelf.Code
					shelfInfo.ShelfName = shelf.Name
					break
				}
			}
			break
		}
	}

	return shelfInfo, nil
}

func (svc WarehouseHttpService) SearchWarehouse(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.WarehouseInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"code",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.WarehouseInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc WarehouseHttpService) SearchWarehouseStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.WarehouseInfo, int, error) {
	searchInFields := []string{
		"code",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.WarehouseInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc WarehouseHttpService) SearchLocation(shopID string, pageable micromodels.Pageable) ([]models.LocationInfo, mongopagination.PaginationData, error) {

	docList, pagination, err := svc.repo.FindLocationPage(shopID, pageable)

	if err != nil {
		return []models.LocationInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc WarehouseHttpService) SearchShelf(shopID string, pageable micromodels.Pageable) ([]models.ShelfInfo, mongopagination.PaginationData, error) {

	docList, pagination, err := svc.repo.FindShelfPage(shopID, pageable)

	if err != nil {
		return []models.ShelfInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc WarehouseHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.Warehouse) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Warehouse](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Warehouse, models.WarehouseDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Warehouse) models.WarehouseDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.WarehouseDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Warehouse = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Warehouse, models.WarehouseDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.WarehouseDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "code", guid)
		},
		func(doc models.WarehouseDoc) bool {
			return doc.Code != ""
		},
		func(shopID string, authUsername string, data models.Warehouse, doc models.WarehouseDoc) error {

			doc.Warehouse = data
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

func (svc WarehouseHttpService) getDocIDKey(doc models.Warehouse) string {
	return doc.Code
}

func (svc WarehouseHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc WarehouseHttpService) GetModuleName() string {
	return "warehouse"
}
