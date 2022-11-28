package services

import (
	"errors"
	"smlcloudplatform/pkg/customershop/customergroup/models"
	"smlcloudplatform/pkg/customershop/customergroup/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ICustomerGroupHttpService interface {
	CreateCustomerGroup(shopID string, authUsername string, doc models.CustomerGroup) (string, error)
	UpdateCustomerGroup(shopID string, guid string, authUsername string, doc models.CustomerGroup) error
	DeleteCustomerGroup(shopID string, guid string, authUsername string) error
	InfoCustomerGroup(shopID string, guid string) (models.CustomerGroupInfo, error)
	SearchCustomerGroup(shopID string, q string, page int, limit int, sort map[string]int) ([]models.CustomerGroupInfo, mongopagination.PaginationData, error)
	SearchCustomerGroupStep(shopID string, langCode string, q string, skip int, limit int, sort map[string]int) ([]models.CustomerGroupInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.CustomerGroup) (common.BulkImport, error)
}

type CustomerGroupHttpService struct {
	repo repositories.ICustomerGroupRepository
}

func NewCustomerGroupHttpService(repo repositories.ICustomerGroupRepository) *CustomerGroupHttpService {

	return &CustomerGroupHttpService{
		repo: repo,
	}
}

func (svc CustomerGroupHttpService) CreateCustomerGroup(shopID string, authUsername string, doc models.CustomerGroup) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "guid", doc.GUID)

	if err != nil {
		return "", err
	}

	if findDoc.GUID != "" {
		return "", errors.New("GUID is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.CustomerGroupDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.CustomerGroup = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc CustomerGroupHttpService) UpdateCustomerGroup(shopID string, guid string, authUsername string, doc models.CustomerGroup) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.CustomerGroup = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc CustomerGroupHttpService) DeleteCustomerGroup(shopID string, guid string, authUsername string) error {

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

func (svc CustomerGroupHttpService) InfoCustomerGroup(shopID string, guid string) (models.CustomerGroupInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.CustomerGroupInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.CustomerGroupInfo{}, errors.New("document not found")
	}

	return findDoc.CustomerGroupInfo, nil

}

func (svc CustomerGroupHttpService) SearchCustomerGroup(shopID string, q string, page int, limit int, sort map[string]int) ([]models.CustomerGroupInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"guidfixed",
		"guid",
	}

	docList, pagination, err := svc.repo.FindPageSort(shopID, searchCols, q, page, limit, sort)

	if err != nil {
		return []models.CustomerGroupInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc CustomerGroupHttpService) SearchCustomerGroupStep(shopID string, langCode string, q string, skip int, limit int, sort map[string]int) ([]models.CustomerGroupInfo, int, error) {
	searchCols := []string{
		"guidfixed",
		"guid",
	}

	projectQuery := map[string]interface{}{
		"guidfixed":    1,
		"guid":         1,
		"customercode": 1,
		"names":        1,
	}

	if langCode != "" {
		projectQuery["names"] = bson.M{"$elemMatch": bson.M{"code": langCode}}
	} else {
		projectQuery["names"] = 1
	}

	docList, total, err := svc.repo.FindLimit(shopID, searchCols, q, skip, limit, sort, projectQuery)

	if err != nil {
		return []models.CustomerGroupInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc CustomerGroupHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.CustomerGroup) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.CustomerGroup](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.GUID)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "guid", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.GUID)
	}

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.CustomerGroup, models.CustomerGroupDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.CustomerGroup) models.CustomerGroupDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.CustomerGroupDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.CustomerGroup = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.CustomerGroup, models.CustomerGroupDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.CustomerGroupDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "guid", guid)
		},
		func(doc models.CustomerGroupDoc) bool {
			return doc.GUID != ""
		},
		func(shopID string, authUsername string, data models.CustomerGroup, doc models.CustomerGroupDoc) error {

			doc.CustomerGroup = data
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
		createDataKey = append(createDataKey, doc.GUID)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.GUID)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {

		updateDataKey = append(updateDataKey, doc.GUID)
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

func (svc CustomerGroupHttpService) getDocIDKey(doc models.CustomerGroup) string {
	return doc.GUID
}
