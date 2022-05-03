package shopzone

import (
	"errors"
	"fmt"
	"smlcloudplatform/pkg/models/restaurant"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/utils"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IShopZoneService interface {
	CreateShopZone(shopID string, authUsername string, doc restaurant.ShopZone) (string, error)
	UpdateShopZone(guid string, shopID string, authUsername string, doc restaurant.ShopZone) error
	DeleteShopZone(guid string, shopID string, authUsername string) error
	InfoShopZone(guid string, shopID string) (restaurant.ShopZoneInfo, error)
	SearchShopZone(shopID string, q string, page int, limit int) ([]restaurant.ShopZoneInfo, mongopagination.PaginationData, error)
}

type ShopZoneService struct {
	crudRepo   repositories.CrudRepository[restaurant.ShopZoneDoc]
	searchRepo repositories.SearchRepository[restaurant.ShopZoneInfo]
}

func NewShopZoneService(crudRepo repositories.CrudRepository[restaurant.ShopZoneDoc], searchRepo repositories.SearchRepository[restaurant.ShopZoneInfo]) ShopZoneService {
	return ShopZoneService{
		crudRepo:   crudRepo,
		searchRepo: searchRepo,
	}
}

func (svc ShopZoneService) CreateShopZone(shopID string, authUsername string, doc restaurant.ShopZone) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := restaurant.ShopZoneDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.ShopZone = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err := svc.crudRepo.Create(docData)

	if err != nil {
		return "", err
	}
	return newGuidFixed, nil
}

func (svc ShopZoneService) UpdateShopZone(guid string, shopID string, authUsername string, doc restaurant.ShopZone) error {

	findDoc, err := svc.crudRepo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.ShopZone = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.crudRepo.Update(guid, findDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc ShopZoneService) DeleteShopZone(guid string, shopID string, authUsername string) error {
	err := svc.crudRepo.Delete(shopID, guid, authUsername)

	if err != nil {
		return err
	}
	return nil
}

func (svc ShopZoneService) InfoShopZone(guid string, shopID string) (restaurant.ShopZoneInfo, error) {

	findDoc, err := svc.crudRepo.FindByGuid(shopID, guid)

	if err != nil {
		return restaurant.ShopZoneInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return restaurant.ShopZoneInfo{}, errors.New("document not found")
	}

	return findDoc.ShopZoneInfo, nil

}

func (svc ShopZoneService) SearchShopZone(shopID string, q string, page int, limit int) ([]restaurant.ShopZoneInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"code",
	}

	for i := range [5]bool{} {
		searchCols = append(searchCols, fmt.Sprintf("name%d", (i+1)))
	}

	docList, pagination, err := svc.searchRepo.FindPage(shopID, searchCols, q, page, limit)

	if err != nil {
		return []restaurant.ShopZoneInfo{}, pagination, err
	}

	return docList, pagination, nil
}
