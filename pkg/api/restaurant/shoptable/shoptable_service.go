package shoptable

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

type IShopTableService interface {
	CreateShopTable(shopID string, authUsername string, doc restaurant.ShopTable) (string, error)
	UpdateShopTable(guid string, shopID string, authUsername string, doc restaurant.ShopTable) error
	DeleteShopTable(guid string, shopID string, authUsername string) error
	InfoShopTable(guid string, shopID string) (restaurant.ShopTableInfo, error)
	SearchShopTable(shopID string, q string, page int, limit int) ([]restaurant.ShopTableInfo, mongopagination.PaginationData, error)
}

type ShopTableService struct {
	crudRepo   repositories.CrudRepository[restaurant.ShopTableDoc]
	searchRepo repositories.SearchRepository[restaurant.ShopTableInfo]
}

func NewShopTableService(crudRepo repositories.CrudRepository[restaurant.ShopTableDoc], searchRepo repositories.SearchRepository[restaurant.ShopTableInfo]) ShopTableService {
	return ShopTableService{
		crudRepo:   crudRepo,
		searchRepo: searchRepo,
	}
}

func (svc ShopTableService) CreateShopTable(shopID string, authUsername string, doc restaurant.ShopTable) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := restaurant.ShopTableDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.ShopTable = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err := svc.crudRepo.Create(docData)

	if err != nil {
		return "", err
	}
	return newGuidFixed, nil
}

func (svc ShopTableService) UpdateShopTable(guid string, shopID string, authUsername string, doc restaurant.ShopTable) error {

	findDoc, err := svc.crudRepo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.ShopTable = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.crudRepo.Update(guid, findDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc ShopTableService) DeleteShopTable(guid string, shopID string, authUsername string) error {
	err := svc.crudRepo.Delete(shopID, guid, authUsername)

	if err != nil {
		return err
	}
	return nil
}

func (svc ShopTableService) InfoShopTable(guid string, shopID string) (restaurant.ShopTableInfo, error) {

	findDoc, err := svc.crudRepo.FindByGuid(shopID, guid)

	if err != nil {
		return restaurant.ShopTableInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return restaurant.ShopTableInfo{}, errors.New("document not found")
	}

	return findDoc.ShopTableInfo, nil

}

func (svc ShopTableService) SearchShopTable(shopID string, q string, page int, limit int) ([]restaurant.ShopTableInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"code",
	}

	for i := range [5]bool{} {
		searchCols = append(searchCols, fmt.Sprintf("name%d", (i+1)))
	}

	docList, pagination, err := svc.searchRepo.FindPage(shopID, searchCols, q, page, limit)

	if err != nil {
		return []restaurant.ShopTableInfo{}, pagination, err
	}

	return docList, pagination, nil
}
