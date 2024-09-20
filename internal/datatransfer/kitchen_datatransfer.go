package datatransfer

import (
	"context"
	"smlcloudplatform/internal/repositories"
	"smlcloudplatform/internal/restaurant/kitchen"
	"smlcloudplatform/internal/restaurant/kitchen/models"
	"smlcloudplatform/pkg/microservice"
	msModels "smlcloudplatform/pkg/microservice/models"

	"github.com/smlsoft/mongopagination"
)

type RestaurantKitchenDataTransfer struct {
	transferConnection IDataTransferConnection
}

type IRestaurantKitchenDataTransferRepository interface {
	FindPage(ctx context.Context, shopID string, searchInFields []string, pageable msModels.Pageable) ([]models.KitchenDoc, mongopagination.PaginationData, error)
}

type RestaurantKitchenDataTransferRepository struct {
	pst microservice.IPersisterMongo
	repositories.SearchRepository[models.KitchenDoc]
}

func NewRestaurantKitchenDataTransferRepository(mongodbPersister microservice.IPersisterMongo) IRestaurantKitchenDataTransferRepository {

	repo := &RestaurantKitchenDataTransferRepository{
		pst: mongodbPersister,
	}

	repo.SearchRepository = repositories.NewSearchRepository[models.KitchenDoc](mongodbPersister)
	return repo
}

func (repo RestaurantKitchenDataTransferRepository) FindPage(ctx context.Context, shopID string, searchInFields []string, pageable msModels.Pageable) ([]models.KitchenDoc, mongopagination.PaginationData, error) {

	results, pagination, err := repo.SearchRepository.FindPage(ctx, shopID, searchInFields, pageable)

	if err != nil {
		return nil, mongopagination.PaginationData{}, err
	}

	return results, pagination, nil
}

func NewRestaurantKitchenDataTransfer(transferConnection IDataTransferConnection) IDataTransfer {
	return &RestaurantKitchenDataTransfer{
		transferConnection: transferConnection,
	}
}

func (pbd *RestaurantKitchenDataTransfer) StartTransfer(ctx context.Context, shopID string) error {

	sourceRestaurantKitchenRepository := NewRestaurantKitchenDataTransferRepository(pbd.transferConnection.GetSourceConnection())
	targetRestaurantKitchenRepository := kitchen.NewKitchenRepository(pbd.transferConnection.GetTargetConnection())

	pageRequest := msModels.Pageable{
		Limit: 100,
		Page:  1,
	}

	for {
		kitchenDocs, pages, err := sourceRestaurantKitchenRepository.FindPage(ctx, shopID, []string{}, pageRequest)
		if err != nil {
			return err
		}

		if len(kitchenDocs) > 0 {
			err = targetRestaurantKitchenRepository.CreateInBatch(ctx, kitchenDocs)
			if err != nil {
				return err
			}
		}

		if pages.TotalPage > int64(pageRequest.Page) {
			pageRequest.Page++
		} else {
			break
		}
	}

	return nil
}
