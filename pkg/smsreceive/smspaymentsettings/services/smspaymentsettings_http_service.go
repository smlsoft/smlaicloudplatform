package services

import (
	"errors"
	smspatternsrepo "smlcloudplatform/pkg/smsreceive/smspatterns/repositories"
	"smlcloudplatform/pkg/smsreceive/smspaymentsettings/models"
	"smlcloudplatform/pkg/smsreceive/smspaymentsettings/repositories"
	"smlcloudplatform/pkg/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ISmsPaymentSettingsHttpService interface {
	SaveSmsPaymentSettings(shopID string, authUsername string, doc models.SmsPaymentSettings) error
	InfoSmsPaymentSettings(shopID string) (models.SmsPaymentSettingsInfo, error)
}

type SmsPaymentSettingsHttpService struct {
	repo        repositories.SmsPaymentSettingsRepository
	repoPattern smspatternsrepo.ISmsPatternsRepository
}

func NewSmsPaymentSettingsHttpService(repo repositories.SmsPaymentSettingsRepository, repoPattern smspatternsrepo.ISmsPatternsRepository) SmsPaymentSettingsHttpService {

	return SmsPaymentSettingsHttpService{
		repo:        repo,
		repoPattern: repoPattern,
	}
}

func (svc SmsPaymentSettingsHttpService) SaveSmsPaymentSettings(shopID string, authUsername string, doc models.SmsPaymentSettings) error {

	findPattern, err := svc.repoPattern.FindByCode(doc.PatternCode)

	if err != nil {
		return err
	}

	if len(findPattern.Code) < 1 {
		return errors.New("pattern code not found")
	}

	findDoc, err := svc.repo.FindOne(shopID, map[string]interface{}{})

	if err != nil {
		return err
	}

	isExitsSetting, err := svc.isExistsPaymentSettings(findDoc)

	if err != nil {
		return err
	}

	if isExitsSetting {
		return svc.updateSmsPaymentSettings(shopID, findDoc.GuidFixed, authUsername, doc)
	} else {
		return svc.createSmsPaymentSettings(shopID, authUsername, doc)
	}

}

func (svc SmsPaymentSettingsHttpService) isExistsPaymentSettings(findDoc models.SmsPaymentSettingsDoc) (bool, error) {

	if len(findDoc.ShopID) > 0 {
		return true, nil
	}

	return false, nil
}

func (svc SmsPaymentSettingsHttpService) createSmsPaymentSettings(shopID string, authUsername string, doc models.SmsPaymentSettings) error {

	newGuidFixed := utils.NewGUID()

	docData := models.SmsPaymentSettingsDoc{}
	docData.SmsPaymentSettings = doc

	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err := svc.repo.Create(docData)

	if err != nil {
		return err
	}

	return nil
}

func (svc SmsPaymentSettingsHttpService) updateSmsPaymentSettings(shopID string, guid string, authUsername string, doc models.SmsPaymentSettings) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.SmsPaymentSettings = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc SmsPaymentSettingsHttpService) InfoSmsPaymentSettings(shopID string) (models.SmsPaymentSettingsInfo, error) {

	findDoc, err := svc.repo.FindOne(shopID, map[string]interface{}{})

	if err != nil {
		return models.SmsPaymentSettingsInfo{}, err
	}

	return findDoc.SmsPaymentSettingsInfo, nil

}
