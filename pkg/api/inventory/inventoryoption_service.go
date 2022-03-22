package inventory

import (
	"errors"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IInventoryOptionService interface {
	CreateInventoryOption(shopID string, authUsername string, invOpt models.InventoryOption) (string, error)
	UpdateInventoryOption(guid string, shopID string, authUsername string, invOpt models.InventoryOption) error
	DeleteInventoryOption(guid string, shopID string, username string) error
	InfoInventoryOption(guid string, shopID string) (models.InventoryOption, error)
	SearchInventoryOption(shopID string, q string, page int, limit int) ([]models.InventoryOption, paginate.PaginationData, error)
}

type InventoryOptionService struct {
	repo IInventoryOptionRepository
}

func NewInventoryOptionService(inventoryOptionRepository IInventoryOptionRepository) InventoryOptionService {
	return InventoryOptionService{
		repo: inventoryOptionRepository,
	}
}

func (svc InventoryOptionService) CreateInventoryOption(shopID string, authUsername string, invOpt models.InventoryOption) (string, error) {

	newGuidFixed := utils.NewGUID()
	invOpt.ShopID = shopID
	invOpt.GuidFixed = newGuidFixed
	invOpt.CreatedBy = authUsername
	invOpt.CreatedAt = time.Now()

	_, err := svc.repo.Create(invOpt)

	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc InventoryOptionService) UpdateInventoryOption(guid string, shopID string, authUsername string, invOpt models.InventoryOption) error {

	findDoc, err := svc.repo.FindByGuid(guid, shopID)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.InventoryID = invOpt.InventoryID
	findDoc.OptionGroupID = invOpt.OptionGroupID
	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	if err != nil {
		return err
	}

	return nil
}

func (svc InventoryOptionService) DeleteInventoryOption(guid string, shopID string, username string) error {

	err := svc.repo.Delete(guid, shopID, username)

	if err != nil {
		return err
	}

	return nil
}

func (svc InventoryOptionService) InfoInventoryOption(guid string, shopID string) (models.InventoryOption, error) {

	findDoc, err := svc.repo.FindByGuid(guid, shopID)

	if err != nil {
		return models.InventoryOption{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.InventoryOption{}, errors.New("document not found")
	}

	return findDoc, nil
}

func (svc InventoryOptionService) SearchInventoryOption(shopID string, q string, page int, limit int) ([]models.InventoryOption, paginate.PaginationData, error) {
	docList, pagination, err := svc.repo.FindPage(shopID, q, page, limit)

	if err != nil {
		return []models.InventoryOption{}, pagination, err
	}

	return docList, pagination, nil
}
