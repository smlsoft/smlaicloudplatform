package option

import (
	"errors"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/product/option/models"
	"smlcloudplatform/pkg/utils"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IOptionService interface {
	CreateOption(shopID string, authUsername string, invOpt models.InventoryOptionMain) (string, error)
	UpdateOption(shopID string, guid string, authUsername string, invOpt models.InventoryOptionMain) error
	DeleteOption(shopID string, guid string, username string) error
	InfoOption(shopID string, guid string) (models.InventoryOptionMainInfo, error)
	InfoWTFArray(shopID string, codes []string) ([]interface{}, error)
	SearchOption(shopID string, pageable micromodels.Pageable) ([]models.InventoryOptionMainInfo, mongopagination.PaginationData, error)
}

type OptionService struct {
	repo IOptionRepository
}

func NewOptionService(inventoryOptionRepository IOptionRepository) OptionService {
	return OptionService{
		repo: inventoryOptionRepository,
	}
}

func (svc OptionService) CreateOption(shopID string, authUsername string, invOpt models.InventoryOptionMain) (string, error) {

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

func (svc OptionService) UpdateOption(shopID string, guid string, authUsername string, invOpt models.InventoryOptionMain) error {

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

func (svc OptionService) DeleteOption(shopID string, guid string, username string) error {

	err := svc.repo.Delete(shopID, guid, username)

	if err != nil {
		return err
	}

	return nil
}

func (svc OptionService) InfoOption(shopID string, guid string) (models.InventoryOptionMainInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.InventoryOptionMainInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.InventoryOptionMainInfo{}, errors.New("document not found")
	}

	return findDoc.InventoryOptionMainInfo, nil
}

func (svc OptionService) InfoWTFArray(shopID string, codes []string) ([]interface{}, error) {
	docList := []interface{}{}

	for _, code := range codes {
		findDoc, err := svc.repo.FindByDocIndentityGuid(shopID, "code", code)
		if err != nil || findDoc.ID == primitive.NilObjectID {
			// add item empty
			docList = append(docList, nil)
		} else {
			docList = append(docList, findDoc.InventoryOptionMainInfo)
		}
	}

	return docList, nil
}

func (svc OptionService) SearchOption(shopID string, pageable micromodels.Pageable) ([]models.InventoryOptionMainInfo, mongopagination.PaginationData, error) {
	docList, pagination, err := svc.repo.FindPage(shopID, pageable)

	if err != nil {
		return []models.InventoryOptionMainInfo{}, pagination, err
	}

	return docList, pagination, nil
}
