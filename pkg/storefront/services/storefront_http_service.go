package services

import (
	"errors"
	"smlcloudplatform/pkg/storefront/models"
	"smlcloudplatform/pkg/storefront/repositories"
	"smlcloudplatform/pkg/utils"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IStorefrontHttpService interface {
	CreateStorefront(shopID string, authUsername string, doc models.Storefront) (string, error)
	UpdateStorefront(shopID string, guid string, authUsername string, doc models.Storefront) error
	DeleteStorefront(shopID string, guid string, authUsername string) error
	InfoStorefront(shopID string, guid string) (models.StorefrontInfo, error)
	SearchStorefront(shopID string, q string, page int, limit int, sort map[string]int) ([]models.StorefrontInfo, mongopagination.PaginationData, error)
}

type StorefrontHttpService struct {
	repo repositories.IStorefrontRepository
}

func NewStorefrontHttpService(repo repositories.IStorefrontRepository) *StorefrontHttpService {

	return &StorefrontHttpService{
		repo: repo,
	}
}

func (svc StorefrontHttpService) CreateStorefront(shopID string, authUsername string, doc models.Storefront) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := models.StorefrontDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Storefront = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err := svc.repo.Create(docData)

	if err != nil {
		return "", err
	}

	return newGuidFixed, nil
}

func (svc StorefrontHttpService) UpdateStorefront(shopID string, guid string, authUsername string, doc models.Storefront) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.Storefront = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}

	return nil
}

func (svc StorefrontHttpService) DeleteStorefront(shopID string, guid string, authUsername string) error {

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

func (svc StorefrontHttpService) InfoStorefront(shopID string, guid string) (models.StorefrontInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.StorefrontInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.StorefrontInfo{}, errors.New("document not found")
	}

	return findDoc.StorefrontInfo, nil

}

func (svc StorefrontHttpService) SearchStorefront(shopID string, q string, page int, limit int, sort map[string]int) ([]models.StorefrontInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"guidfixed",
		"code",
	}

	docList, pagination, err := svc.repo.FindPageSort(shopID, searchCols, q, page, limit, sort)

	if err != nil {
		return []models.StorefrontInfo{}, pagination, err
	}

	return docList, pagination, nil
}
