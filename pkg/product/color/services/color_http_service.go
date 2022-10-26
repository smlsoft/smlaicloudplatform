package services

import (
	"errors"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/color/models"
	"smlcloudplatform/pkg/product/color/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IColorHttpService interface {
	CreateColor(shopID string, authUsername string, doc models.Color) (string, error)
	UpdateColor(shopID string, guid string, authUsername string, doc models.Color) error
	DeleteColor(shopID string, guid string, authUsername string) error
	InfoColor(shopID string, guid string) (models.ColorInfo, error)
	SearchColor(shopID string, q string, page int, limit int, sort map[string]int) ([]models.ColorInfo, mongopagination.PaginationData, error)
	SearchColorLimit(shopID string, langCode string, q string, skip int, limit int, sort map[string]int) ([]models.ColorInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Color) (common.BulkImport, error)
}

type ColorHttpService struct {
	repo repositories.IColorRepository
}

func NewColorHttpService(repo repositories.IColorRepository) *ColorHttpService {

	return &ColorHttpService{
		repo: repo,
	}
}

func (svc ColorHttpService) CreateColor(shopID string, authUsername string, doc models.Color) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", doc.Code)

	if err != nil {
		return "", err
	}

	if findDoc.Code != "" {
		return "", errors.New("Code is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.ColorDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Color = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc ColorHttpService) UpdateColor(shopID string, guid string, authUsername string, doc models.Color) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.Color = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc ColorHttpService) DeleteColor(shopID string, guid string, authUsername string) error {

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

func (svc ColorHttpService) InfoColor(shopID string, guid string) (models.ColorInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.ColorInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.ColorInfo{}, errors.New("document not found")
	}

	return findDoc.ColorInfo, nil

}

func (svc ColorHttpService) SearchColor(shopID string, q string, page int, limit int, sort map[string]int) ([]models.ColorInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"guidfixed",
		"code",
	}

	docList, pagination, err := svc.repo.FindPageSort(shopID, searchCols, q, page, limit, sort)

	if err != nil {
		return []models.ColorInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc ColorHttpService) SearchColorLimit(shopID string, langCode string, q string, skip int, limit int, sort map[string]int) ([]models.ColorInfo, int, error) {
	searchCols := []string{
		"guidfixed",
		"code",
	}

	projectQuery := map[string]interface{}{
		"guidfixed":      1,
		"code":           1,
		"colorselect":    1,
		"colorsystem":    1,
		"colorhex":       1,
		"colorselecthex": 1,
		"colorsystemhex": 1,
	}

	if langCode != "" {
		projectQuery["names"] = bson.M{"$elemMatch": bson.M{"code": langCode}}
	} else {
		projectQuery["names"] = 1
	}

	docList, total, err := svc.repo.FindLimit(shopID, searchCols, q, skip, limit, sort, projectQuery)

	if err != nil {
		return []models.ColorInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc ColorHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.Color) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Color](dataList, svc.getDocIDKey)

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

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Color, models.ColorDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Color) models.ColorDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.ColorDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Color = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Color, models.ColorDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.ColorDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "code", guid)
		},
		func(doc models.ColorDoc) bool {
			return doc.Code != ""
		},
		func(shopID string, authUsername string, data models.Color, doc models.ColorDoc) error {

			doc.Color = data
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

func (svc ColorHttpService) getDocIDKey(doc models.Color) string {
	return doc.Code
}
