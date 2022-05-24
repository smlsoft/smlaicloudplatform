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
	UpdateInventoryOptionMain(shopID string, guid string, authUsername string, invOpt models.InventoryOptionMain) error
	DeleteInventoryOptionMain(shopID string, guid string, username string) error
	InfoInventoryOptionMain(shopID string, guid string) (models.InventoryOptionMainInfo, error)
	SearchInventoryOptionMain(shopID string, q string, page int, limit int) ([]models.InventoryOptionMainInfo, paginate.PaginationData, error)
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

	invOptDoc := models.InventoryOptionMainDoc{}
	invOptDoc.ShopID = shopID
	invOptDoc.GuidFixed = newGuidFixed

	invOptDoc.InventoryOptionMain = invOpt

	invOptDoc.CreatedBy = authUsername
	invOptDoc.CreatedAt = time.Now()

	if invOptDoc.InventoryOptionMain.Choices == nil {
		invOptDoc.InventoryOptionMain.Choices = &[]models.Choice{}
	}

	_, err := svc.repo.Create(invOptDoc)

	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc InventoryOptionMainService) UpdateInventoryOptionMain(shopID string, guid string, authUsername string, invOpt models.InventoryOptionMain) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.InventoryOptionMain = invOpt
	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	if findDoc.InventoryOptionMain.Choices == nil {
		findDoc.InventoryOptionMain.Choices = &[]models.Choice{}
	}

	svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc InventoryOptionMainService) DeleteInventoryOptionMain(shopID string, guid string, username string) error {

	err := svc.repo.Delete(shopID, guid, username)

	if err != nil {
		return err
	}

	return nil
}

func (svc InventoryOptionMainService) InfoInventoryOptionMain(shopID string, guid string) (models.InventoryOptionMainInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.InventoryOptionMainInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.InventoryOptionMainInfo{}, errors.New("document not found")
	}

	return findDoc.InventoryOptionMainInfo, nil
}

func (svc InventoryOptionMainService) SearchInventoryOptionMain(shopID string, q string, page int, limit int) ([]models.InventoryOptionMainInfo, paginate.PaginationData, error) {
	docList, pagination, err := svc.repo.FindPage(shopID, q, page, limit)

	if err != nil {
		return []models.InventoryOptionMainInfo{}, pagination, err
	}

	return docList, pagination, nil
}
