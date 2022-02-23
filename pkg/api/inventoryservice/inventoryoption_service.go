package inventoryservice

import (
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
)

type IInventoryOptionService interface {
	CreateInventoryOption(merchantId string, authUsername string, invOpt models.InventoryOption) (string, error)
	UpdateInventoryOption(guid string, merchantId string, authUsername string, invOpt models.InventoryOption) error
	DeleteInventoryOption(guid string, merchantId string) error
	InfoInventoryOption(guid string, merchantId string) (models.InventoryOption, error)
	SearchInventoryOption(merchantId string, q string, page int, limit int) ([]models.InventoryOption, paginate.PaginationData, error)
}

type InventoryOptionService struct {
	repo IInventoryOptionRepository
}

func NewInventoryOptionService(inventoryOptionRepository IInventoryOptionRepository) IInventoryOptionService {
	return &InventoryOptionService{
		repo: inventoryOptionRepository,
	}
}

func (svc *InventoryOptionService) CreateInventoryOption(merchantId string, authUsername string, invOpt models.InventoryOption) (string, error) {

	newGuidFixed := utils.NewGUID()
	invOpt.MerchantId = merchantId
	invOpt.GuidFixed = newGuidFixed
	invOpt.Deleted = false
	invOpt.CreatedBy = authUsername
	invOpt.CreatedAt = time.Now()

	_, err := svc.repo.Create(invOpt)

	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc *InventoryOptionService) UpdateInventoryOption(guid string, merchantId string, authUsername string, invOpt models.InventoryOption) error {

	findDoc, err := svc.repo.FindByGuid(guid, merchantId)

	if err != nil {
		return err
	}

	findDoc.InventoryId = invOpt.InventoryId
	findDoc.OptionGroupId = invOpt.OptionGroupId
	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	if err != nil {
		return err
	}

	return nil
}

func (svc *InventoryOptionService) DeleteInventoryOption(guid string, merchantId string) error {

	err := svc.repo.Delete(guid, merchantId)

	if err != nil {
		return err
	}

	return nil
}

func (svc *InventoryOptionService) InfoInventoryOption(guid string, merchantId string) (models.InventoryOption, error) {

	findDoc, err := svc.repo.FindByGuid(guid, merchantId)

	if err != nil {
		return models.InventoryOption{}, err
	}

	return findDoc, nil
}

func (svc *InventoryOptionService) SearchInventoryOption(merchantId string, q string, page int, limit int) ([]models.InventoryOption, paginate.PaginationData, error) {
	docList, pagination, err := svc.repo.FindPage(merchantId, q, page, limit)

	if err != nil {
		return []models.InventoryOption{}, pagination, err
	}

	return docList, pagination, nil
}
