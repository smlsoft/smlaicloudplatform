package services

import (
	"errors"
	"regexp"
	"smlcloudplatform/pkg/smsreceive/smspatterns/models"
	"smlcloudplatform/pkg/smsreceive/smspatterns/repositories"
	"smlcloudplatform/pkg/utils"
	"strconv"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ISmsPatternsHttpService interface {
	CreateSmsPatterns(authUsername string, doc models.SmsPatterns) (string, error)
	UpdateSmsPatterns(guid string, authUsername string, doc models.SmsPatterns) error
	DeleteSmsPatterns(guid string) error
	InfoSmsPatterns(guid string) (models.SmsPatternsInfo, error)
	SearchSmsPatterns(q string, page int, limit int) ([]models.SmsPatternsInfo, mongopagination.PaginationData, error)
}

type SmsPatternsHttpService struct {
	repo repositories.ISmsPatternsRepository
}

func NewSmsPatternsHttpService(repo repositories.ISmsPatternsRepository) *SmsPatternsHttpService {

	return &SmsPatternsHttpService{
		repo: repo,
	}
}

func (svc SmsPatternsHttpService) CreateSmsPatterns(authUsername string, doc models.SmsPatterns) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := models.SmsPatternsDoc{}
	docData.GuidFixed = newGuidFixed
	docData.SmsPatterns = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err := svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc SmsPatternsHttpService) UpdateSmsPatterns(guid string, authUsername string, doc models.SmsPatterns) error {

	findDoc, err := svc.repo.FindByGuid(guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	tempCode := findDoc.Code

	findDoc.SmsPatterns = doc
	findDoc.Code = tempCode

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.UpdateByGuid(guid, findDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc SmsPatternsHttpService) DeleteSmsPatterns(guid string) error {

	findDoc, err := svc.repo.FindByGuid(guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	err = svc.repo.DeleteByGuid(guid)
	if err != nil {
		return err
	}

	return nil
}

func (svc SmsPatternsHttpService) InfoSmsPatterns(guid string) (models.SmsPatternsInfo, error) {

	findDoc, err := svc.repo.FindByGuid(guid)

	if err != nil {
		return models.SmsPatternsInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.SmsPatternsInfo{}, errors.New("document not found")
	}

	return findDoc.SmsPatternsInfo, nil

}

func (svc SmsPatternsHttpService) SearchSmsPatterns(q string, page int, limit int) ([]models.SmsPatternsInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"name",
		"address",
	}

	docList, pagination, err := svc.repo.FindPage(searchCols, q, page, limit)

	if err != nil {
		return []models.SmsPatternsInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func GetAmountFromPattern(pattern string, message string) (float64, error) {
	re := regexp.MustCompile(pattern)

	reVal := re.FindStringSubmatch(message)

	if len(reVal) > 1 {

		return strconv.ParseFloat(reVal[1], 64)

	}

	return 0.0, errors.New("message not match")
}
