package services

import (
	"errors"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/customershop/customer/models"
	"smlcloudplatform/pkg/customershop/customer/repositories"
	modelsCustomerGroup "smlcloudplatform/pkg/customershop/customergroup/models"
	repositoriesCustomerGroup "smlcloudplatform/pkg/customershop/customergroup/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/samber/lo"
	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ICustomerHttpService interface {
	CreateCustomer(shopID string, authUsername string, doc models.CustomerRequest) (string, error)
	UpdateCustomer(shopID string, guid string, authUsername string, doc models.CustomerRequest) error
	DeleteCustomer(shopID string, guid string, authUsername string) error
	DeleteCustomerByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoCustomer(shopID string, guid string) (models.CustomerInfo, error)
	SearchCustomer(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.CustomerInfo, mongopagination.PaginationData, error)
	SearchCustomerStep(shopID string, langCode string, pageable micromodels.PageableStep) ([]models.CustomerInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Customer) (common.BulkImport, error)
}

type CustomerHttpService struct {
	repo              repositories.ICustomerRepository
	repoCustomerGroup repositoriesCustomerGroup.ICustomerGroupRepository
}

func NewCustomerHttpService(repo repositories.ICustomerRepository, repoCustomerGroup repositoriesCustomerGroup.ICustomerGroupRepository) *CustomerHttpService {

	return &CustomerHttpService{
		repo:              repo,
		repoCustomerGroup: repoCustomerGroup,
	}
}

func (svc CustomerHttpService) getCustomerGroupByGUIDs(shopID string, groupGUIDs []string) ([]modelsCustomerGroup.CustomerGroupInfo, error) {
	docCustomerGroups, err := svc.repoCustomerGroup.FindByGuids(shopID, groupGUIDs)

	if err != nil {
		return []modelsCustomerGroup.CustomerGroupInfo{}, err
	}

	customerGroups := lo.Map[modelsCustomerGroup.CustomerGroupDoc, modelsCustomerGroup.CustomerGroupInfo](docCustomerGroups, func(docCustomerGroup modelsCustomerGroup.CustomerGroupDoc, idx int) modelsCustomerGroup.CustomerGroupInfo {
		return docCustomerGroup.CustomerGroupInfo
	})

	return customerGroups, nil
}

func (svc CustomerHttpService) CreateCustomer(shopID string, authUsername string, doc models.CustomerRequest) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", doc.Code)

	if err != nil {
		return "", err
	}

	if findDoc.Code != "" {
		return "", errors.New("code is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.CustomerDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Customer = doc.Customer

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	customerGroups := []modelsCustomerGroup.CustomerGroupInfo{}
	if doc.Groups != nil && len(doc.Groups) > 0 {
		customerGroupGUIDs := lo.Map[models.CustomerGroupRequest, string](
			doc.Groups,
			func(docGroup models.CustomerGroupRequest, idx int) string {
				return docGroup.GuidFixed
			})

		customerGroups, err = svc.getCustomerGroupByGUIDs(shopID, customerGroupGUIDs)
		if err != nil {
			return "", err
		}
	}

	docData.Groups = &customerGroups

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc CustomerHttpService) UpdateCustomer(shopID string, guid string, authUsername string, doc models.CustomerRequest) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.Customer = doc.Customer

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	customerGroups := []modelsCustomerGroup.CustomerGroupInfo{}
	if doc.Groups != nil && len(doc.Groups) > 0 {
		customerGroupGUIDs := lo.Map[models.CustomerGroupRequest, string](
			doc.Groups,
			func(docGroup models.CustomerGroupRequest, idx int) string {
				return docGroup.GuidFixed
			})
		customerGroups, err = svc.getCustomerGroupByGUIDs(shopID, customerGroupGUIDs)
		if err != nil {
			return err
		}
	}

	findDoc.Groups = &customerGroups

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc CustomerHttpService) DeleteCustomer(shopID string, guid string, authUsername string) error {

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

func (svc CustomerHttpService) InfoCustomer(shopID string, guid string) (models.CustomerInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.CustomerInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.CustomerInfo{}, errors.New("document not found")
	}

	return findDoc.CustomerInfo, nil

}

func (svc CustomerHttpService) SearchCustomer(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.CustomerInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"code",
		"names.name",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.CustomerInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc CustomerHttpService) SearchCustomerStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.CustomerInfo, int, error) {
	searchInFields := []string{
		"code",
		"names.name",
	}

	selectCols := []string{
		"guidfixed", "code", "personaltype", "images",
		"names", "addressforbilling", "addressforshipping",
		"taxid", "email",
	}

	selectFields := map[string]interface{}{}

	for _, col := range selectCols {
		selectFields[col] = 1
	}

	if langCode != "" {
		selectFields["names"] = bson.M{"$elemMatch": bson.M{"code": langCode}}
	} else {
		selectFields["names"] = 1
	}

	filterQuery := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, filterQuery, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.CustomerInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc CustomerHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.Customer) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Customer](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Customer, models.CustomerDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Customer) models.CustomerDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.CustomerDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Customer = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Customer, models.CustomerDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.CustomerDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "code", guid)
		},
		func(doc models.CustomerDoc) bool {
			return doc.Code != ""
		},
		func(shopID string, authUsername string, data models.Customer, doc models.CustomerDoc) error {

			doc.Customer = data
			doc.UpdatedBy = authUsername
			doc.UpdatedAt = time.Now()

			groupGUIDs := []string{}

			for _, tempGroup := range *doc.Groups {
				groupGUIDs = append(groupGUIDs, tempGroup.GuidFixed)
			}

			customerGroups := []modelsCustomerGroup.CustomerGroupInfo{}
			if doc.Groups != nil && len(*doc.Groups) > 0 {
				customerGroups, err = svc.getCustomerGroupByGUIDs(shopID, groupGUIDs)
				if err != nil {
					return err
				}
			}

			doc.Groups = &customerGroups

			err = svc.repo.Update(shopID, doc.GuidFixed, doc)
			if err != nil {
				return nil
			}
			return nil
		},
	)

	for docIdx, tempDoc := range createDataList {
		customerGroups := []modelsCustomerGroup.CustomerGroupInfo{}

		groupGUIDs := lo.Map[modelsCustomerGroup.CustomerGroupInfo, string](
			*tempDoc.Groups,
			func(doc modelsCustomerGroup.CustomerGroupInfo, idx int) string {
				return doc.GuidFixed
			})

		if tempDoc.Groups != nil && len(*tempDoc.Groups) > 0 {
			customerGroups, err = svc.getCustomerGroupByGUIDs(shopID, groupGUIDs)
			if err != nil {
				return common.BulkImport{}, err
			}
		}

		createDataList[docIdx].Groups = &customerGroups
	}

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

func (svc CustomerHttpService) DeleteCustomerByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc CustomerHttpService) getDocIDKey(doc models.Customer) string {
	return doc.Code
}
