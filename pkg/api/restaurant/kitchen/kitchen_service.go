package kitchen

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

type IKitchenService interface {
	CreateKitchen(shopID string, authUsername string, doc restaurant.Kitchen) (string, error)
	UpdateKitchen(guid string, shopID string, authUsername string, doc restaurant.Kitchen) error
	DeleteKitchen(guid string, shopID string, authUsername string) error
	InfoKitchen(guid string, shopID string) (restaurant.KitchenInfo, error)
	SearchKitchen(shopID string, q string, page int, limit int) ([]restaurant.KitchenInfo, mongopagination.PaginationData, error)
}

type KitchenService struct {
	crudRepo   repositories.CrudRepository[restaurant.KitchenDoc]
	searchRepo repositories.SearchRepository[restaurant.KitchenInfo]
}

func NewKitchenService(crudRepo repositories.CrudRepository[restaurant.KitchenDoc], searchRepo repositories.SearchRepository[restaurant.KitchenInfo]) KitchenService {
	return KitchenService{
		crudRepo:   crudRepo,
		searchRepo: searchRepo,
	}
}

func (svc KitchenService) CreateKitchen(shopID string, authUsername string, doc restaurant.Kitchen) (string, error) {

	newGuidFixed := utils.NewGUID()

	docData := restaurant.KitchenDoc{}
	docData.ShopID = shopID
	docData.GuidFixed = newGuidFixed
	docData.Kitchen = doc

	docData.CreatedBy = authUsername
	docData.CreatedAt = time.Now()

	_, err := svc.crudRepo.Create(docData)

	if err != nil {
		return "", err
	}
	return newGuidFixed, nil
}

func (svc KitchenService) UpdateKitchen(guid string, shopID string, authUsername string, doc restaurant.Kitchen) error {

	findDoc, err := svc.crudRepo.FindByGuid(shopID, guid)

	if err != nil {
		return err
	}

	if findDoc.ID == primitive.NilObjectID {
		return errors.New("document not found")
	}

	findDoc.Kitchen = doc

	findDoc.UpdatedBy = authUsername
	findDoc.UpdatedAt = time.Now()

	err = svc.crudRepo.Update(guid, findDoc)

	if err != nil {
		return err
	}
	return nil
}

func (svc KitchenService) DeleteKitchen(guid string, shopID string, authUsername string) error {
	err := svc.crudRepo.Delete(shopID, guid, authUsername)

	if err != nil {
		return err
	}
	return nil
}

func (svc KitchenService) InfoKitchen(guid string, shopID string) (restaurant.KitchenInfo, error) {

	findDoc, err := svc.crudRepo.FindByGuid(shopID, guid)

	if err != nil {
		return restaurant.KitchenInfo{}, err
	}

	if findDoc.ID == primitive.NilObjectID {
		return restaurant.KitchenInfo{}, errors.New("document not found")
	}

	return findDoc.KitchenInfo, nil

}

func (svc KitchenService) SearchKitchen(shopID string, q string, page int, limit int) ([]restaurant.KitchenInfo, mongopagination.PaginationData, error) {
	searchCols := []string{
		"code",
	}

	for i := range [5]bool{} {
		searchCols = append(searchCols, fmt.Sprintf("name%d", (i+1)))
	}

	docList, pagination, err := svc.searchRepo.FindPage(shopID, searchCols, q, page, limit)

	if err != nil {
		return []restaurant.KitchenInfo{}, pagination, err
	}

	return docList, pagination, nil
}
