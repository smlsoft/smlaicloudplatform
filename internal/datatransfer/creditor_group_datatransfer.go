package datatransfer

import (
	"context"
	"smlaicloudplatform/internal/debtaccount/creditorgroup/models"
	creditorGroupRepositories "smlaicloudplatform/internal/debtaccount/creditorgroup/repositories"
	"smlaicloudplatform/internal/repositories"
	"smlaicloudplatform/pkg/microservice"
	msModels "smlaicloudplatform/pkg/microservice/models"

	"github.com/smlsoft/mongopagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreditorGroupDataTransfer struct {
	transferConnection IDataTransferConnection
}

type ICreditorGroupDataTransferRepository interface {
	FindPage(ctx context.Context, shopID string, searchInFields []string, pageable msModels.Pageable) ([]models.CreditorGroupDoc, mongopagination.PaginationData, error)
}

type CreditorGroupDataTransferRepository struct {
	pst microservice.IPersisterMongo
	repositories.SearchRepository[models.CreditorGroupDoc]
}

func NewCreditorGroupDataTransferRepository(mongodbPersister microservice.IPersisterMongo) ICreditorGroupDataTransferRepository {
	repo := &CreditorGroupDataTransferRepository{
		pst: mongodbPersister,
	}

	repo.SearchRepository = repositories.NewSearchRepository[models.CreditorGroupDoc](mongodbPersister)
	return repo
}

func (repo CreditorGroupDataTransferRepository) FindPage(ctx context.Context, shopID string, searchInFields []string, pageable msModels.Pageable) ([]models.CreditorGroupDoc, mongopagination.PaginationData, error) {

	results, pagination, err := repo.SearchRepository.FindPage(ctx, shopID, searchInFields, pageable)

	if err != nil {
		return nil, mongopagination.PaginationData{}, err
	}

	return results, pagination, nil
}

func NewCreditorGroupDataTransfer(transferConnection IDataTransferConnection) IDataTransfer {

	return &CreditorGroupDataTransfer{
		transferConnection: transferConnection,
	}
}

func (pdt *CreditorGroupDataTransfer) StartTransfer(ctx context.Context, shopID string, targetShopID string) error {

	sourceRepository := NewCreditorGroupDataTransferRepository(pdt.transferConnection.GetSourceConnection())
	targetRepository := creditorGroupRepositories.NewCreditorGroupRepository(pdt.transferConnection.GetTargetConnection())

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

			if targetShopID != "" {
				for i := range docs {
					docs[i].ShopID = targetShopID
					docs[i].ID = primitive.NewObjectID()
				}
			}

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
