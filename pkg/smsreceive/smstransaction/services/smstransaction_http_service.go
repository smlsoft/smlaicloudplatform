package services

import (
	"errors"
	"regexp"
	"smlcloudplatform/pkg/smsreceive/smstransaction/models"
	"smlcloudplatform/pkg/smsreceive/smstransaction/repositories"
	"strconv"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ISmsTransactionHttpService interface {
	CreateSmsTransaction(shopID string, authUsername string, doc models.SmsTransaction) (string, error)
	UpdateSmsTransaction(guid string, shopID string, authUsername string, doc models.SmsTransaction) error
	DeleteSmsTransaction(guid string, shopID string, authUsername string) error
	InfoSmsTransaction(guid string, shopID string) (models.SmsTransactionInfo, error)
	SearchSmsTransaction(shopID string, q string, page int, limit int, sort map[string]int) ([]models.SmsTransactionInfo, mongopagination.PaginationData, error)
	CheckSMS(shopID string, amount float64, startTime time.Time, endTime time.Time) (float64, error)
}

type SmsTransactionHttpService struct {
	repo    repositories.ISmsTransactionRepository
	genGUID func() string
	timeNow func() time.Time
}

func NewSmsTransactionHttpService(repo repositories.ISmsTransactionRepository, genGUID func() string, timeNow func() time.Time) *SmsTransactionHttpService {

	return &SmsTransactionHttpService{
		repo:    repo,
		genGUID: genGUID,
		timeNow: timeNow,
	}
}

func (svc SmsTransactionHttpService) CreateSmsTransaction(shopID string, authUsername string, doc models.SmsTransaction) (string, error) {

	findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "transid", doc.TransId)

	if err != nil {
		return "", err
	}

	if findDoc.TransId != "" {
		return "", errors.New("TransId is exists")
	}

	newGuidFixed := svc.genGUID()

	docData := models.SmsTransactionDoc{}
	docData.ShopID = shopID
	docData.SmsTransaction = doc

	docData.GuidFixed = newGuidFixed
	docData.TransId = svc.genGUID()

	docData.CreatedBy = authUsername
	docData.CreatedAt = svc.timeNow()

	_, err = svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc SmsTransactionHttpService) UpdateSmsTransaction(guid string, shopID string, authUsername string, doc models.SmsTransaction) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	tempTransId := findDoc.TransId

	findDoc.SmsTransaction = doc
	findDoc.TransId = tempTransId

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = svc.timeNow()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc SmsTransactionHttpService) DeleteSmsTransaction(guid string, shopID string, authUsername string) error {

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

func (svc SmsTransactionHttpService) InfoSmsTransaction(guid string, shopID string) (models.SmsTransactionInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.SmsTransactionInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.SmsTransactionInfo{}, errors.New("document not found")
	}

	return findDoc.SmsTransactionInfo, nil

}

func (svc SmsTransactionHttpService) SearchSmsTransaction(shopID string, q string, page int, limit int, sort map[string]int) ([]models.SmsTransactionInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"guidfixed",
		"transid",
	}

	docList, pagination, err := svc.repo.FindPageSort(shopID, searchCols, q, page, limit, sort)

	if err != nil {
		return []models.SmsTransactionInfo{}, pagination, err
	}

	return docList, pagination, nil
}

func (svc SmsTransactionHttpService) CheckSMS(shopID string, amount float64, startTime time.Time, endTime time.Time) (float64, error) {

	// timeNowUnix := time.Unix(svc.timeNow().Unix(), 0)

	// startTime := timeNowUnix.Add(time.Duration(-3) * time.Minute).Unix()
	// endTime := timeNowUnix.Add(time.Duration(10) * time.Minute).Unix()

	addressKey := "kbank"

	smsList, err := svc.repo.FindFilterSms(shopID, addressKey, startTime, endTime)
	if err != nil {
		return 0, err
	}

	var amountVal float64 = 0

	for _, smsMessage := range smsList {
		// msg1 := "12/04/63 09:25 บชX231148X รับโอนจากX815923X 1170.00บ คงเหลือ 2160.29บ"

		re := regexp.MustCompile(`[0-9]{2}\/[0-9]{2}\/[0-9]{2} [0-9]{2}:[0-9]{2} บชX[0-9].*X (?P<Amount>[0-9].*)บ คงเหลือ [0-9].*บ`)

		reVal := re.FindStringSubmatch(smsMessage.Body)

		if len(reVal) > 1 {

			amountVal, err = strconv.ParseFloat(reVal[1], 64)

			if err != nil {
				return 0, err
			}

			if amountVal > 0 {
				break
			}
		}
	}

	return amountVal, nil
}
