package services

import (
	"context"
	"errors"
	"fmt"
	mastersync "smlaicloudplatform/internal/mastersync/repositories"
	common "smlaicloudplatform/internal/models"
	businessTypeRepositories "smlaicloudplatform/internal/organization/businesstype/repositories"
	"smlaicloudplatform/internal/organization/company/models"
	"smlaicloudplatform/internal/organization/company/repositories"
	deparmentRepositories "smlaicloudplatform/internal/organization/department/repositories"
	"smlaicloudplatform/internal/services"
	"smlaicloudplatform/internal/utils"
	"smlaicloudplatform/internal/utils/importdata"
	micromodels "smlaicloudplatform/pkg/microservice/models"
	"time"

	"github.com/smlsoft/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type ICompanyHttpService interface {
	CreateCompany(shopID string, authUsername string, doc models.Company) (string, error)
	UpdateCompany(shopID string, guid string, authUsername string, doc models.Company) error
	DeleteCompany(shopID string, guid string, authUsername string) error
	DeleteCompanyByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoCompany(shopID string, guid string) (models.CompanyInfoResponse, error)
	InfoCompanyByCode(shopID string, code string) (models.CompanyInfoResponse, error)
	SearchCompany(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.CompanyInfoResponse, mongopagination.PaginationData, error)
	SearchCompanyStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.CompanyInfoResponse, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Company) (common.BulkImport, error)

	GetModuleName() string
}

type CompanyHttpService struct {
	repo             repositories.ICompanyRepository
	repoDepartment   deparmentRepositories.IDepartmentRepository
	repoBusinessType businessTypeRepositories.IBusinessTypeRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.CompanyActivity, models.CompanyDeleteActivity]
	contextTimeout time.Duration
}

func NewCompanyHttpService(repo repositories.ICompanyRepository, repoDepartment deparmentRepositories.IDepartmentRepository, repoBusinessType businessTypeRepositories.IBusinessTypeRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *CompanyHttpService {

	contextTimeout := time.Duration(15) * time.Second

	insSvc := &CompanyHttpService{
		repo:             repo,
		repoDepartment:   repoDepartment,
		repoBusinessType: repoBusinessType,
		syncCacheRepo:    syncCacheRepo,
		contextTimeout:   contextTimeout,
	}

	insSvc.ActivityService = services.NewActivityService[models.CompanyActivity, models.CompanyDeleteActivity](repo)

	return insSvc
}

func (svc CompanyHttpService) getContextTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), svc.contextTimeout)
}

func (svc CompanyHttpService) CreateCompany(shopID string, authUsername string, doc models.Company) (string, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindOneFilter(
		ctx,
		shopID,
		map[string]interface{}{
			"code": doc.Code,
		},
	)
	if err != nil {
		return "", err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", errors.New("company is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.CompanyDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Company = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(ctx, docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc CompanyHttpService) UpdateCompany(shopID string, guid string, authUsername string, doc models.Company) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	docData := findDoc

	docData.Company = doc

	docData.GuidFixed = findDoc.GuidFixed
	docData.UpdatedBy = authUsername
	docData.UpdatedAt = time.Now()

	err = svc.repo.Update(ctx, shopID, guid, docData)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc CompanyHttpService) DeleteCompany(shopID string, guid string, authUsername string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	err = svc.repo.DeleteByGuidfixed(ctx, shopID, guid, authUsername)
	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc CompanyHttpService) DeleteCompanyByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(ctx, shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc CompanyHttpService) InfoCompany(shopID string, guid string) (models.CompanyInfoResponse, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByGuid(ctx, shopID, guid)

	if err != nil {
		return models.CompanyInfoResponse{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.CompanyInfoResponse{}, errors.New("document not found")
	}

	resultDoc, err := svc.mapCompanyInfo(findDoc.CompanyInfo, shopID)

	if err != nil {
		return models.CompanyInfoResponse{}, err
	}

	return resultDoc, nil
}

func (svc CompanyHttpService) mapCompanyInfo(findInfo models.CompanyInfo, shopID string) (models.CompanyInfoResponse, error) {

	resultDoc := models.CompanyInfoResponse{}
	resultDoc.CompanyInfo = findInfo

	return resultDoc, nil
}

func (svc CompanyHttpService) InfoCompanyByCode(shopID string, code string) (models.CompanyInfoResponse, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	findDoc, err := svc.repo.FindByDocIndentityGuid(ctx, shopID, "code", code)

	if err != nil {
		return models.CompanyInfoResponse{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.CompanyInfoResponse{}, errors.New("document not found")
	}

	resultDoc, err := svc.mapCompanyInfo(findDoc.CompanyInfo, shopID)

	if err != nil {
		return models.CompanyInfoResponse{}, err
	}

	return resultDoc, nil
}

func (svc CompanyHttpService) SearchCompany(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.CompanyInfoResponse, mongopagination.PaginationData, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"code",
		"names.name",
	}

	docList, pagination, err := svc.repo.FindPageFilter(ctx, shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.CompanyInfoResponse{}, pagination, err
	}

	resultDocs := []models.CompanyInfoResponse{}
	for _, docInfo := range docList {

		resultDoc, err := svc.mapCompanyInfo(docInfo, shopID)

		if err != nil {
			return []models.CompanyInfoResponse{}, pagination, err
		}

		resultDocs = append(resultDocs, resultDoc)
	}

	return resultDocs, pagination, nil
}

func (svc CompanyHttpService) SearchCompanyStep(shopID string, langCode string, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.CompanyInfoResponse, int, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	searchInFields := []string{
		"code",
		"names.name",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(ctx, shopID, filters, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.CompanyInfoResponse{}, 0, err
	}

	resultDocs := []models.CompanyInfoResponse{}
	for _, docInfo := range docList {

		resultDoc, err := svc.mapCompanyInfo(docInfo, shopID)

		if err != nil {
			return []models.CompanyInfoResponse{}, 0, err
		}

		resultDocs = append(resultDocs, resultDoc)
	}

	return resultDocs, total, nil
}

func (svc CompanyHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.Company) (common.BulkImport, error) {

	ctx, ctxCancel := svc.getContextTimeout()
	defer ctxCancel()

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Company](dataList, svc.getDocIDKey)

	itemCodeGuidList := []interface{}{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.Code)
	}

	findItemGuid, err := svc.repo.FindInItemGuids(ctx, shopID, "code", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.Code)
	}

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Company, models.CompanyDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Company) models.CompanyDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.CompanyDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Company = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Company, models.CompanyDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.CompanyDoc, error) {
			return svc.repo.FindByDocIndentityGuid(ctx, shopID, "code", guid)
		},
		func(doc models.CompanyDoc) bool {
			return doc.Code != ""
		},
		func(shopID string, authUsername string, data models.Company, doc models.CompanyDoc) error {

			doc.Company = data
			doc.UpdatedBy = authUsername
			doc.UpdatedAt = time.Now()

			err = svc.repo.Update(ctx, shopID, doc.GuidFixed, doc)
			if err != nil {
				return nil
			}
			return nil
		},
	)

	if len(createDataList) > 0 {
		err = svc.repo.CreateInBatch(ctx, createDataList)

		if err != nil {
			return common.BulkImport{}, err
		}

	}

	createDataKey := []interface{}{}

	for _, doc := range createDataList {
		createDataKey = append(createDataKey, doc.Code)
	}

	payloadDuplicateDataKey := []interface{}{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.Code)
	}

	updateDataKey := []interface{}{}
	for _, doc := range updateSuccessDataList {

		updateDataKey = append(updateDataKey, doc.Code)
	}

	updateFailDataKey := []interface{}{}
	for _, doc := range updateFailDataList {
		updateFailDataKey = append(updateFailDataKey, svc.getDocIDKey(doc))
	}

	svc.saveMasterSync(shopID)

	tempCreateDataKey := svc.toSliceString(createDataKey)
	tempUpdateDataKey := svc.toSliceString(updateDataKey)
	tempUpdateFailDataKey := svc.toSliceString(updateFailDataKey)
	tempPayloadDuplicateDataKey := svc.toSliceString(payloadDuplicateDataKey)

	return common.BulkImport{
		Created:          tempCreateDataKey,
		Updated:          tempUpdateDataKey,
		UpdateFailed:     tempUpdateFailDataKey,
		PayloadDuplicate: tempPayloadDuplicateDataKey,
	}, nil
}

func (svc CompanyHttpService) toSliceString(data []interface{}) []string {
	tempData := make([]string, len(data))
	for i, v := range data {
		tempData[i] = v.(string)
	}

	return tempData
}

func (svc CompanyHttpService) getDocIDKey(doc models.Company) string {
	return doc.Code
}

func (svc CompanyHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc CompanyHttpService) GetModuleName() string {
	return "company"
}
