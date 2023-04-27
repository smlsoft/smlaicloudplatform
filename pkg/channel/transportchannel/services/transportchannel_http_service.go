package services

import (
	"errors"
	"fmt"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/channel/transportchannel/models"
	"smlcloudplatform/pkg/channel/transportchannel/repositories"
	mastersync "smlcloudplatform/pkg/mastersync/repositories"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/services"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type ITransportChannelHttpService interface {
	CreateTransportChannel(shopID string, authUsername string, doc models.TransportChannel) (string, error)
	UpdateTransportChannel(shopID string, guid string, authUsername string, doc models.TransportChannel) error
	DeleteTransportChannel(shopID string, guid string, authUsername string) error
	DeleteTransportChannelByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoTransportChannel(shopID string, guid string) (models.TransportChannelInfo, error)
	InfoTransportChannelByCode(shopID string, code string) (models.TransportChannelInfo, error)
	SearchTransportChannel(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.TransportChannelInfo, mongopagination.PaginationData, error)
	SearchTransportChannelStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.TransportChannelInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.TransportChannel) (common.BulkImport, error)

	GetModuleName() string
}

type TransportChannelHttpService struct {
	repo repositories.ITransportChannelRepository

	syncCacheRepo mastersync.IMasterSyncCacheRepository
	services.ActivityService[models.TransportChannelActivity, models.TransportChannelDeleteActivity]
}

func NewTransportChannelHttpService(repo repositories.ITransportChannelRepository, syncCacheRepo mastersync.IMasterSyncCacheRepository) *TransportChannelHttpService {

	insSvc := &TransportChannelHttpService{
		repo:          repo,
		syncCacheRepo: syncCacheRepo,
	}

	insSvc.ActivityService = services.NewActivityService[models.TransportChannelActivity, models.TransportChannelDeleteActivity](repo)

	return insSvc
}

func (svc TransportChannelHttpService) CreateTransportChannel(shopID string, authUsername string, doc models.TransportChannel) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", doc.Code)

	if err != nil {
		return "", err
	}

	if len(findDoc.GuidFixed) > 0 {
		return "", errors.New("Code is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.TransportChannelDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.TransportChannel = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	svc.saveMasterSync(shopID)

	return newGuidFixed, nil
}

func (svc TransportChannelHttpService) UpdateTransportChannel(shopID string, guid string, authUsername string, doc models.TransportChannel) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	findDoc.TransportChannel = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc TransportChannelHttpService) DeleteTransportChannel(shopID string, guid string, authUsername string) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if len(findDoc.GuidFixed) < 1 {
		return errors.New("document not found")
	}

	err = svc.repo.DeleteByGuidfixed(shopID, guid, authUsername)
	if err != nil {
		return err
	}

	svc.saveMasterSync(shopID)

	return nil
}

func (svc TransportChannelHttpService) DeleteTransportChannelByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc TransportChannelHttpService) InfoTransportChannel(shopID string, guid string) (models.TransportChannelInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.TransportChannelInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.TransportChannelInfo{}, errors.New("document not found")
	}

	return findDoc.TransportChannelInfo, nil
}

func (svc TransportChannelHttpService) InfoTransportChannelByCode(shopID string, code string) (models.TransportChannelInfo, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", code)

	if err != nil {
		return models.TransportChannelInfo{}, err
	}

	if len(findDoc.GuidFixed) < 1 {
		return models.TransportChannelInfo{}, errors.New("document not found")
	}

	return findDoc.TransportChannelInfo, nil
}

func (svc TransportChannelHttpService) SearchTransportChannel(shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.TransportChannelInfo, mongopagination.PaginationData, error) {
	searchInFields := []string{
		"code",
		"names.name",
	}

	docList, pagination, err := svc.repo.FindPageFilter(shopID, filters, searchInFields, pageable)

	if err != nil {
		return []models.TransportChannelInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc TransportChannelHttpService) SearchTransportChannelStep(shopID string, langCode string, pageableStep micromodels.PageableStep) ([]models.TransportChannelInfo, int, error) {
	searchInFields := []string{
		"code",
		"names.name",
	}

	selectFields := map[string]interface{}{}

	docList, total, err := svc.repo.FindStep(shopID, map[string]interface{}{}, searchInFields, selectFields, pageableStep)

	if err != nil {
		return []models.TransportChannelInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc TransportChannelHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.TransportChannel) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.TransportChannel](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.TransportChannel, models.TransportChannelDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.TransportChannel) models.TransportChannelDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.TransportChannelDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.TransportChannel = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.TransportChannel, models.TransportChannelDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.TransportChannelDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "code", guid)
		},
		func(doc models.TransportChannelDoc) bool {
			return doc.Code != ""
		},
		func(shopID string, authUsername string, data models.TransportChannel, doc models.TransportChannelDoc) error {

			doc.TransportChannel = data
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

func (svc TransportChannelHttpService) getDocIDKey(doc models.TransportChannel) string {
	return doc.Code
}

func (svc TransportChannelHttpService) saveMasterSync(shopID string) {
	if svc.syncCacheRepo != nil {
		err := svc.syncCacheRepo.Save(shopID, svc.GetModuleName())

		if err != nil {
			fmt.Printf("save %s cache error :: %s", svc.GetModuleName(), err.Error())
		}
	}
}

func (svc TransportChannelHttpService) GetModuleName() string {
	return "transportChannel"
}
