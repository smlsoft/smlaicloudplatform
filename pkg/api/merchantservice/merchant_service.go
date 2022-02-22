package merchantservice

import (
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"time"
)

type IMerchantService interface {
	CreateMerchant(username string, merchant models.Merchant) (string, error)
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
