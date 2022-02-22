package merchantservice

import (
	"errors"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"time"
)

type IMerchantService interface {
	CreateMerchant(username string, merchant models.Merchant) (string, error)
	UpdateMerchant(guid string, username string, merchant models.Merchant) error
	DeleteMerchant(guid string, username string) error
}

type MerchantService struct {
	repo IMerchantRepository
}

func NewMerchantService(repo IMerchantRepository) IMerchantService {
	return &MerchantService{
		repo: repo,
	}
}

func (svc *MerchantService) CreateMerchant(username string, merchant models.Merchant) (string, error) {

	merchantId := utils.NewGUID()
	merchant.GuidFixed = merchantId
	merchant.CreatedBy = username
	merchant.CreatedAt = time.Now()

	_, err := svc.repo.Create(merchant)

	if err != nil {
		return "", err
	}

	return merchantId, nil
}

func (svc *MerchantService) UpdateMerchant(guid string, username string, merchant models.Merchant) error {

	findMerchant, err := svc.repo.FindByGuid(guid)

	if err != nil {
		return err
	}

	// *** warning feature change to check by role owner
	if len(findMerchant.CreatedBy) < 1 {
		return errors.New("username invalid")
	}

	findMerchant.Name1 = merchant.Name1
	findMerchant.UpdatedBy = username
	findMerchant.UpdatedAt = time.Now()

	err = svc.repo.Update(guid, merchant)

	if err != nil {
		return err
	}

	return nil
}

func (svc *MerchantService) DeleteMerchant(guid string, username string) error {
	findMerchant, err := svc.repo.FindByGuid(guid)

	if err != nil {
		return err
	}

	// *** warning feature change to check by role owner
	if len(findMerchant.CreatedBy) < 1 {
		return errors.New("username invalid")
	}

	err = svc.repo.Delete(guid)

	if err != nil {
		return err
	}
	return nil
}
