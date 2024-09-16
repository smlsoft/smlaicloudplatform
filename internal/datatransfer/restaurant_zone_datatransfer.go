package datatransfer

import (
	"context"
	"smlcloudplatform/internal/repositories"
	restaurantZoneRepository "smlcloudplatform/internal/restaurant/zone"
	"smlcloudplatform/internal/restaurant/zone/models"
	"smlcloudplatform/pkg/microservice"
	micromodels "smlcloudplatform/pkg/microservice/models"
	msModels "smlcloudplatform/pkg/microservice/models"

	"github.com/userplant/mongopagination"
)

type RestaurantZoneDataTransfer struct {
	transferConnection IDataTransferConnection
}

type IRestaurantZoneDataTransferRepository interface {
	FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.ZoneDoc, mongopagination.PaginationData, error)
}

type RestaurantZoneDataTransferRepository struct {
	pst microservice.IPersisterMongo
	repositories.SearchRepository[models.ZoneDoc]
}

func NewRestaurantZoneDataTransferRepository(mongodbPersister microservice.IPersisterMongo) IRestaurantZoneDataTransferRepository {
	repo := &RestaurantZoneDataTransferRepository{
		pst: mongodbPersister,
	}

	repo.SearchRepository = repositories.NewSearchRepository[models.ZoneDoc](mongodbPersister)
	return repo
}

func (repo RestaurantZoneDataTransferRepository) FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.ZoneDoc, mongopagination.PaginationData, error) {

	results, pagination, err := repo.SearchRepository.FindPage(ctx, shopID, searchInFields, pageable)

	if err != nil {
		return nil, mongopagination.PaginationData{}, err
	}

	return results, pagination, nil
}

func NewRestaurantZoneDataTransfer(transferConnection IDataTransferConnection) IDataTransfer {
	return &RestaurantZoneDataTransfer{
		transferConnection: transferConnection,
	}
}

func (pdt *RestaurantZoneDataTransfer) StartTransfer(ctx context.Context, shopID string) error {

	sourceRepository := NewRestaurantZoneDataTransferRepository(pdt.transferConnection.GetSourceConnection())
	targetRepository := restaurantZoneRepository.NewZoneRepository(pdt.transferConnection.GetTargetConnection())

	pageRequest := msModels.Pageable{
		Limit: 100,
		Page:  1,
	}

	for {

		docs, pages, err := sourceRepository.FindPage(ctx, shopID, nil, pageRequest)
		if err != nil {
			return err
		}

		if len(docs) > 0 {
			err = targetRepository.CreateInBatch(ctx, docs)
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
