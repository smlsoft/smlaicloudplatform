package services

import (
	"errors"
	"smlcloudplatform/pkg/utils"
	"smlcloudplatform/pkg/vfgl/accountperiodmaster/models"
	"smlcloudplatform/pkg/vfgl/accountperiodmaster/repositories"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAccountPeriodMasterHttpService interface {
	CreateAccountPeriodMaster(shopID string, authUsername string, doc models.AccountPeriodMaster) (string, error)
	UpdateAccountPeriodMaster(shopID string, guid string, authUsername string, doc models.AccountPeriodMaster) error
	DeleteAccountPeriodMaster(shopID string, guid string, authUsername string) error
	DeleteAccountPeriodMasterByGUIDs(shopID string, authUsername string, GUIDs []string) error
	InfoAccountPeriodMaster(shopID string, guid string) (models.AccountPeriodMasterInfo, error)
	SearchAccountPeriodMaster(shopID string, q string, page int, limit int, sort map[string]int) ([]models.AccountPeriodMasterInfo, mongopagination.PaginationData, error)
	SearchAccountPeriodMasterStep(shopID string, langCode string, q string, skip int, limit int, sort map[string]int) ([]models.AccountPeriodMasterInfo, int, error)
}

type AccountPeriodMasterHttpService struct {
	repo repositories.IAccountPeriodMasterRepository
}

func NewAccountPeriodMasterHttpService(repo repositories.IAccountPeriodMasterRepository) *AccountPeriodMasterHttpService {

	return &AccountPeriodMasterHttpService{
		repo: repo,
	}
}

func (svc AccountPeriodMasterHttpService) CreateAccountPeriodMaster(shopID string, authUsername string, doc models.AccountPeriodMaster) (string, error) {

	findDocExists, err := svc.repo.FindByDateRange(shopID, doc.StartDate, doc.EndDate)

	if err != nil {
		return "", err
	}

	if len(findDocExists.GuidFixed) > 0 {
		return "", errors.New("date range already exists")
	}

	newGuidFixed := utils.NewGUID()

	docData := models.AccountPeriodMasterDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.AccountPeriodMaster = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc AccountPeriodMasterHttpService) UpdateAccountPeriodMaster(shopID string, guid string, authUsername string, doc models.AccountPeriodMaster) error {

	findDocExists, err := svc.repo.FindByDateRange(shopID, doc.StartDate, doc.EndDate)

	if err != nil {
		return err
	}

	if len(findDocExists.GuidFixed) > 0 && findDocExists.GuidFixed != guid {
		return errors.New("date range already exists")
	}

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.AccountPeriodMaster = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc AccountPeriodMasterHttpService) DeleteAccountPeriodMaster(shopID string, guid string, authUsername string) error {

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

func (svc AccountPeriodMasterHttpService) DeleteAccountPeriodMasterByGUIDs(shopID string, authUsername string, GUIDs []string) error {

	deleteFilterQuery := map[string]interface{}{
		"guidfixed": bson.M{"$in": GUIDs},
	}

	err := svc.repo.Delete(shopID, authUsername, deleteFilterQuery)
	if err != nil {
		return err
	}

	return nil
}

func (svc AccountPeriodMasterHttpService) InfoAccountPeriodMaster(shopID string, guid string) (models.AccountPeriodMasterInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.AccountPeriodMasterInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.AccountPeriodMasterInfo{}, errors.New("document not found")
	}

	return findDoc.AccountPeriodMasterInfo, nil

}

func (svc AccountPeriodMasterHttpService) SearchAccountPeriodMaster(shopID string, q string, page int, limit int, sort map[string]int) ([]models.AccountPeriodMasterInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"guidfixed",
		"docno",
	}

	docList, pagination, err := svc.repo.FindPageSort(shopID, searchCols, q, page, limit, sort)

	if err != nil {
		return []models.AccountPeriodMasterInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc AccountPeriodMasterHttpService) SearchAccountPeriodMasterStep(shopID string, langCode string, q string, skip int, limit int, sort map[string]int) ([]models.AccountPeriodMasterInfo, int, error) {
	searchCols := []string{
		"guidfixed",
		"docno",
	}

	projectQuery := map[string]interface{}{
		"guidfixed": 1,
		"docno":     1,
	}

	if langCode != "" {
		projectQuery["names"] = bson.M{"$elemMatch": bson.M{"code": langCode}}
	} else {
		projectQuery["names"] = 1
	}

	docList, total, err := svc.repo.FindLimit(shopID, searchCols, q, skip, limit, sort, projectQuery)

	if err != nil {
		return []models.AccountPeriodMasterInfo{}, 0, err
	}

	return docList, total, nil
}
