package datatransfer

import (
	"context"
	"smlcloudplatform/internal/repositories"
	restaurantTableRepository "smlcloudplatform/internal/restaurant/table"
	"smlcloudplatform/internal/restaurant/table/models"
	"smlcloudplatform/pkg/microservice"
	micromodels "smlcloudplatform/pkg/microservice/models"
	msModels "smlcloudplatform/pkg/microservice/models"

	"github.com/userplant/mongopagination"
)

type RestaurantTableDataTransfer struct {
	transferConnection IDataTransferConnection
}

type IRestaurantTableDataTransferRepository interface {
	FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.TableDoc, mongopagination.PaginationData, error)
}

type RestaurantTableDataTransferRepository struct {
	pst microservice.IPersisterMongo
	repositories.SearchRepository[models.TableDoc]
}

func NewRestaurantTableDataTransferRepository(mongodbPersister microservice.IPersisterMongo) IRestaurantTableDataTransferRepository {
	repo := &RestaurantTableDataTransferRepository{
		pst: mongodbPersister,
	}

	repo.SearchRepository = repositories.NewSearchRepository[models.TableDoc](mongodbPersister)
	return repo
}

func (repo RestaurantTableDataTransferRepository) FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.TableDoc, mongopagination.PaginationData, error) {

	results, pagination, err := repo.SearchRepository.FindPage(ctx, shopID, searchInFields, pageable)

	if err != nil {
		return nil, mongopagination.PaginationData{}, err
	}

	return results, pagination, nil
}

func NewRestaurantTableDataTransfer(transferConnection IDataTransferConnection) IDataTransfer {
	return &RestaurantTableDataTransfer{
		transferConnection: transferConnection,
	}
}

func (pdt *RestaurantTableDataTransfer) StartTransfer(ctx context.Context, shopID string) error {

	sourceRepository := NewRestaurantTableDataTransferRepository(pdt.transferConnection.GetSourceConnection())
	targetRepository := restaurantTableRepository.NewTableRepository(pdt.transferConnection.GetTargetConnection())

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
