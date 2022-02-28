package merchant

import (
	"errors"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IMerchantService interface {
	CreateMerchant(username string, merchant models.Merchant) (string, error)
	UpdateMerchant(guid string, username string, merchant models.Merchant) error
	DeleteMerchant(guid string, username string) error
	InfoMerchant(guid string) (models.MerchantInfo, error)
	SearchMerchant(q string, page int, limit int) ([]models.MerchantInfo, paginate.PaginationData, error)
}

type MerchantService struct {
	repo             IMerchantRepository
	merchantUserRepo IMerchantUserRepository
}

func NewMerchantService(repo IMerchantRepository, merchantUserRepo IMerchantUserRepository) IMerchantService {
	return &MerchantService{
		repo:             repo,
		merchantUserRepo: merchantUserRepo,
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

	svc.merchantUserRepo.Save(merchantId, username, models.ROLE_OWNER)

	return merchantId, nil
}

func (svc *MerchantService) UpdateMerchant(guid string, username string, merchant models.Merchant) error {

	findMerchant, err := svc.repo.FindByGuid(guid)

	if err != nil {
		return err
	}

	if findMerchant.Id == primitive.NilObjectID {
		return errors.New("merchant not found")
	}

	findMerchant.Name1 = merchant.Name1
	findMerchant.UpdatedBy = username
	findMerchant.UpdatedAt = time.Now()

	err = svc.repo.Update(guid, findMerchant)

	if err != nil {
		return err
	}

	return nil
}

func (svc *MerchantService) DeleteMerchant(guid string, username string) error {

	err := svc.repo.Delete(guid)

	if err != nil {
		return err
	}
	return nil
}

func (svc *MerchantService) InfoMerchant(guid string) (models.MerchantInfo, error) {
	findMerchant, err := svc.repo.FindByGuid(guid)

	if err != nil {
		return models.MerchantInfo{}, err
	}

	return models.MerchantInfo{
		Id:        findMerchant.Id,
		GuidFixed: findMerchant.GuidFixed,
		Name1:     findMerchant.Name1,
	}, nil
}

func (svc *MerchantService) SearchMerchant(q string, page int, limit int) ([]models.MerchantInfo, paginate.PaginationData, error) {
	merchantList, pagination, err := svc.repo.FindPage(q, page, limit)

	if err != nil {
		return merchantList, pagination, err
	}

	return merchantList, pagination, nil
}
