package services

import (
	"errors"
	common "smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/product/unit/models"
	"smlcloudplatform/pkg/product/unit/repositories"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/utils/importdata"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IUnitHttpService interface {
	CreateUnit(shopID string, authUsername string, doc models.Unit) (string, error)
	UpdateUnit(shopID string, guid string, authUsername string, doc models.Unit) error
	UpdateFieldUnit(shopID string, guid string, authUsername string, doc models.Unit) error
	DeleteUnit(shopID string, guid string, authUsername string) error
	InfoUnit(shopID string, guid string) (models.UnitInfo, error)
	SearchUnit(shopID string, q string, page int, limit int, sort map[string]int) ([]models.UnitInfo, mongopagination.PaginationData, error)
	SearchUnitLimit(shopID string, langCode string, q string, skip int, limit int, sort map[string]int) ([]models.UnitInfo, int, error)
	SaveInBatch(shopID string, authUsername string, dataList []models.Unit) (common.BulkImport, error)
}

type UnitHttpService struct {
	repo repositories.IUnitRepository
}

func NewUnitHttpService(repo repositories.IUnitRepository) *UnitHttpService {

	return &UnitHttpService{
		repo: repo,
	}
}

func (svc UnitHttpService) CreateUnit(shopID string, authUsername string, doc models.Unit) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "unitcode", doc.UnitCode)

	if err != nil {
		return "", err
	}

	if findDoc.UnitCode != "" {
		return "", errors.New("unit code is exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.UnitDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Unit = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc UnitHttpService) UpdateUnit(shopID string, guid string, authUsername string, doc models.Unit) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	tempCode := findDoc.UnitCode

	findDoc.Unit = doc

	//
	findDoc.UnitCode = tempCode

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc UnitHttpService) UpdateFieldUnit(shopID string, guid string, authUsername string, doc models.Unit) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	temp := map[string]common.NameX{}

	for _, v := range *findDoc.Names {
		temp[*v.Code] = v
	}

	for _, v := range *doc.Names {
		temp[*v.Code] = v
	}

	tempNames := []common.NameX{}

	for _, v := range temp {
		tempNames = append(tempNames, v)
	}

	findDoc.Unit.Names = &tempNames

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc UnitHttpService) DeleteUnit(shopID string, guid string, authUsername string) error {

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

func (svc UnitHttpService) InfoUnit(shopID string, guid string) (models.UnitInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.UnitInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.UnitInfo{}, errors.New("document not found")
	}

	return findDoc.UnitInfo, nil

}

func (svc UnitHttpService) SearchUnit(shopID string, q string, page int, limit int, sort map[string]int) ([]models.UnitInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"guidfixed",
		"unitcode",
	}

	docList, pagination, err := svc.repo.FindPageSort(shopID, searchCols, q, page, limit, sort)

	if err != nil {
		return []models.UnitInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc UnitHttpService) SearchUnitLimit(shopID string, langCode string, q string, skip int, limit int, sort map[string]int) ([]models.UnitInfo, int, error) {
	searchCols := []string{
		"guidfixed",
		"unitcode",
	}

	projectQuery := map[string]interface{}{
		"guidfixed": 1,
		"unitcode":  1,
	}

	if langCode != "" {
		projectQuery["names"] = bson.M{"$elemMatch": bson.M{"code": langCode}}
	} else {
		projectQuery["names"] = 1
	}

	docList, total, err := svc.repo.FindLimit(shopID, searchCols, q, skip, limit, sort, map[string]interface{}{})

	if err != nil {
		return []models.UnitInfo{}, 0, err
	}

	return docList, total, nil
}

func (svc UnitHttpService) SaveInBatch(shopID string, authUsername string, dataList []models.Unit) (common.BulkImport, error) {

	payloadList, payloadDuplicateList := importdata.FilterDuplicate[models.Unit](dataList, svc.getDocIDKey)

	itemCodeGuidList := []string{}
	for _, doc := range payloadList {
		itemCodeGuidList = append(itemCodeGuidList, doc.UnitCode)
	}

	findItemGuid, err := svc.repo.FindInItemGuid(shopID, "unitcode", itemCodeGuidList)

	if err != nil {
		return common.BulkImport{}, err
	}

	foundItemGuidList := []string{}
	for _, doc := range findItemGuid {
		foundItemGuidList = append(foundItemGuidList, doc.DocNo)
	}

	duplicateDataList, createDataList := importdata.PreparePayloadData[models.Unit, models.UnitDoc](
		shopID,
		authUsername,
		foundItemGuidList,
		payloadList,
		svc.getDocIDKey,
		func(shopID string, authUsername string, doc models.Unit) models.UnitDoc {
			newGuid := utils.NewGUID()

			dataDoc := models.UnitDoc{}

			dataDoc.GuidFixed = newGuid
			dataDoc.ShopID = shopID
			dataDoc.Unit = doc

			currentTime := time.Now()
			dataDoc.CreatedBy = authUsername
			dataDoc.CreatedAt = currentTime
			return dataDoc
		},
	)

	updateSuccessDataList, updateFailDataList := importdata.UpdateOnDuplicate[models.Unit, models.UnitDoc](
		shopID,
		authUsername,
		duplicateDataList,
		svc.getDocIDKey,
		func(shopID string, guid string) (models.UnitDoc, error) {
			return svc.repo.FindByDocIndentityGuid(shopID, "unitcode", guid)
		},
		func(doc models.UnitDoc) bool {
			return doc.UnitCode != ""
		},
		func(shopID string, authUsername string, data models.Unit, doc models.UnitDoc) error {

			doc.Unit = data
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
		createDataKey = append(createDataKey, doc.UnitCode)
	}

	payloadDuplicateDataKey := []string{}
	for _, doc := range payloadDuplicateList {
		payloadDuplicateDataKey = append(payloadDuplicateDataKey, doc.UnitCode)
	}

	updateDataKey := []string{}
	for _, doc := range updateSuccessDataList {

		updateDataKey = append(updateDataKey, doc.UnitCode)
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

func (svc UnitHttpService) getDocIDKey(doc models.Unit) string {
	return doc.UnitCode
}
