package services

import (
	"errors"
	"fmt"
	"smlcloudplatform/pkg/shopdesign/zonedesign/models"
	"smlcloudplatform/pkg/shopdesign/zonedesign/repositories"
	"smlcloudplatform/pkg/utils"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IZoneDesignService interface {
	CreateZoneDesign(shopID string, authUsername string, category models.ZoneDesign) (string, error)
	UpdateZoneDesign(shopID string, guid string, authUsername string, category models.ZoneDesign) error
	DeleteZoneDesign(shopID string, guid string, authUsername string) error
	InfoZoneDesign(shopID string, guid string) (models.ZoneDesignInfo, error)
	SearchZoneDesign(shopID string, q string, page int, limit int) ([]models.ZoneDesignInfo, mongopagination.PaginationData, error)
}

type ZoneDesignService struct {
	repo repositories.IZoneDesignRepository
}

func NewZoneDesignService(repo repositories.IZoneDesignRepository) ZoneDesignService {
	return ZoneDesignService{
		repo: repo,
	}
}

func (svc ZoneDesignService) CreateZoneDesign(shopID string, authUsername string, doc models.ZoneDesign) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := models.ZoneDesignDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.ZoneDesign = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err := svc.repo.Create(docData)

	if err != nil {
		return "", err
	}
	return newGuidFixed, nil
}

func (svc ZoneDesignService) UpdateZoneDesign(shopID string, guid string, authUsername string, category models.ZoneDesign) error {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.ZoneDesign = category

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.repo.Update(shopID, guid, findDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc ZoneDesignService) DeleteZoneDesign(shopID string, guid string, authUsername string) error {
	err := svc.repo.DeleteByGuidfixed(shopID, guid, authUsername)

	if err != nil {
		return err
	}
	return nil
}

func (svc ZoneDesignService) InfoZoneDesign(shopID string, guid string) (models.ZoneDesignInfo, error) {

	findDoc, err := svc.repo.FindByGuid(shopID, guid)

	if err != nil {
		return models.ZoneDesignInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return models.ZoneDesignInfo{}, errors.New("document not found")
	}

	return findDoc.ZoneDesignInfo, nil

}

func (svc ZoneDesignService) SearchZoneDesign(shopID string, q string, page int, limit int) ([]models.ZoneDesignInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"guidfixed",
		"guidfixed",
	}

	for i := range [5]bool{} {
		searchCols = append(searchCols, fmt.Sprintf("name%d", (i+1)))
	}

	docList, pagination, err := svc.repo.FindPage(shopID, searchCols, q, page, limit)

	if err != nil {
		return []models.ZoneDesignInfo{}, pagination, err
	}

	return docList, pagination, nil
}
