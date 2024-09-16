package datatransfer

import (
	"context"
	"smlcloudplatform/internal/organization/businesstype/models"
	"smlcloudplatform/internal/repositories"
	"smlcloudplatform/pkg/microservice"
	micromodels "smlcloudplatform/pkg/microservice/models"
	msModels "smlcloudplatform/pkg/microservice/models"

	organizationBusinessTypeRepository "smlcloudplatform/internal/organization/businesstype/repositories"

	"github.com/userplant/mongopagination"
)

type OrganizationBusinessTypeDataTransfer struct {
	transferConnection IDataTransferConnection
}

type IOrganizationBusinessTypeDataTransferRepository interface {
	FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.BusinessTypeDoc, mongopagination.PaginationData, error)
}

type OrganizationBusinessTypeDataTransferRepository struct {
	pst microservice.IPersisterMongo
	repositories.SearchRepository[models.BusinessTypeDoc]
}

func NewOrganizationBusinessTypeDataTransferRepository(mongodbPersister microservice.IPersisterMongo) IOrganizationBusinessTypeDataTransferRepository {
	repo := &OrganizationBusinessTypeDataTransferRepository{
		pst: mongodbPersister,
	}

	repo.SearchRepository = repositories.NewSearchRepository[models.BusinessTypeDoc](mongodbPersister)
	return repo
}

func (repo OrganizationBusinessTypeDataTransferRepository) FindPage(ctx context.Context, shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.BusinessTypeDoc, mongopagination.PaginationData, error) {

	results, pagination, err := repo.SearchRepository.FindPage(ctx, shopID, searchInFields, pageable)

	if err != nil {
		return nil, mongopagination.PaginationData{}, err
	}

	return results, pagination, nil
}

func NewOrganizationBusinessTypeDataTransfer(transferConnection IDataTransferConnection) IDataTransfer {
	return &OrganizationBusinessTypeDataTransfer{
		transferConnection: transferConnection,
	}
}

func (dt *OrganizationBusinessTypeDataTransfer) StartTransfer(ctx context.Context, shopID string) error {

	sourceRepository := NewOrganizationBusinessTypeDataTransferRepository(dt.transferConnection.GetSourceConnection())
	targetRepository := organizationBusinessTypeRepository.NewBusinessTypeRepository(dt.transferConnection.GetTargetConnection())

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
