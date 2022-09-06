package services

import (
	"errors"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"smlcloudplatform/pkg/warehouse/models"
	"smlcloudplatform/pkg/warehouse/repositories"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IWarehouseHttpService interface {
	CreateWarehouse(shopID string, authUsername string, doc models.Warehouse) (string, error)
	UpdateWarehouse(shopID string, guid string, authUsername string, doc models.Warehouse) error
	DeleteWarehouse(shopID string, guid string, authUsername string) error
	InfoWarehouse(shopID string, guid string) (models.WarehouseInfo, error)
	SearchWarehouse(shopID string, q string, page int, limit int, sort map[string]int) ([]models.WarehouseInfo, mongopagination.PaginationData, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Warehouse) (common.BulkImport, error)
}

type WarehouseHttpService struct {
	repo repositories.IWarehouseRepository
}

func NewWarehouseHttpService(repo repositories.IWarehouseRepository) *WarehouseHttpService {

	return &WarehouseHttpService{
		repo: repo,
	}
}

func (svc WarehouseHttpService) CreateWarehouse(shopID string, authUsername string, doc models.Warehouse) (string, error) {

	if svc.isDuplicatLocation(*doc.Locations) {
		return "", errors.New("location code is duplicated")
	}

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", doc.Code)

	if err != nil {
		return "", err
	}

	if findDoc.Code != "" {
		return "", errors.New("code is exists")
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

	return newGuidFixed, nil
}

func (svc WarehouseHttpService) UpdateWarehouse(shopID string, guid string, authUsername string, doc models.Warehouse) error {
	if svc.isDuplicatLocation(*doc.Locations) {
		return errors.New("location code is duplicated")
	}

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

	return nil
}

func (svc WarehouseHttpService) isDuplicatLocation(locations []models.Location) bool {
	locationKey := map[string]struct{}{}
	for _, loction := range locations {
		if _, ok := locationKey[loction.Code]; ok {
			return true
		}
		locationKey[loction.Code] = struct{}{}
	}

	return false
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

func (svc WarehouseHttpService) SearchWarehouse(shopID string, q string, page int, limit int, sort map[string]int) ([]models.WarehouseInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"guidfixed",
		"code",
	}

	docList, pagination, err := svc.repo.FindPageSort(shopID, searchCols, q, page, limit, sort)

	if err != nil {
		return []models.WarehouseInfo{}, pagination, err
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
