package inventory

import (
	"errors"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/utils"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IInventoryOptionMainService interface {
	CreateInventoryOptionMain(shopID string, authUsername string, invOpt models.InventoryOptionMain) (string, error)
	UpdateInventoryOptionMain(guid string, shopID string, authUsername string, invOpt models.InventoryOptionMain) error
	DeleteInventoryOptionMain(guid string, shopID string, username string) error
	InfoInventoryOptionMain(guid string, shopID string) (models.InventoryOptionMain, error)
	SearchInventoryOptionMain(shopID string, q string, page int, limit int) ([]models.InventoryOptionMain, paginate.PaginationData, error)
}

type InventoryOptionMainService struct {
	repo IInventoryOptionMainRepository
}

func NewInventoryOptionMainService(inventoryOptionRepository IInventoryOptionMainRepository) InventoryOptionMainService {
	return InventoryOptionMainService{
		repo: inventoryOptionRepository,
	}
}

func (svc InventoryOptionMainService) CreateInventoryOptionMain(shopID string, authUsername string, invOpt models.InventoryOptionMain) (string, error) {

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

func (svc InventoryOptionMainService) UpdateInventoryOptionMain(guid string, shopID string, authUsername string, invOpt models.InventoryOptionMain) error {

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

func (svc InventoryOptionMainService) DeleteInventoryOptionMain(guid string, shopID string, username string) error {

	err := svc.repo.Delete(guid, shopID, username)

	if err != nil {
		return err
	}

	return nil
}

func (svc InventoryOptionMainService) InfoInventoryOptionMain(guid string, shopID string) (models.InventoryOptionMain, error) {

	findDoc, err := svc.repo.FindByGuid(guid, shopID)

	if err != nil {
		return models.InventoryOptionMain{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.InventoryOptionMain{}, errors.New("document not found")
	}

	return findDoc, nil
}

func (svc InventoryOptionMainService) SearchInventoryOptionMain(shopID string, q string, page int, limit int) ([]models.InventoryOptionMain, paginate.PaginationData, error) {
	docList, pagination, err := svc.repo.FindPage(shopID, q, page, limit)

	if err != nil {
		return []models.InventoryOptionMain{}, pagination, err
	}

	return docList, pagination, nil
}
