package chartofaccountadmin

import (
	"context"
	chartOfAccountModels "smlcloudplatform/internal/vfgl/chartofaccount/models"
	"smlcloudplatform/pkg/microservice"

	"go.mongodb.org/mongo-driver/bson"
)

type IChartOfAccountAdminRepository interface {
	FindChartOfAccountDocByShopID(ctx context.Context, shopID string) ([]chartOfAccountModels.ChartOfAccountDoc, error)
}

type ChartOfAccountAdminRepository struct {
	pst microservice.IPersisterMongo
}

func NewChartOfAccountAdminRepository(pst microservice.IPersisterMongo) IChartOfAccountAdminRepository {
	return &ChartOfAccountAdminRepository{
		pst: pst,
	}
}

func (r ChartOfAccountAdminRepository) FindChartOfAccountDocByShopID(ctx context.Context, shopID string) ([]chartOfAccountModels.ChartOfAccountDoc, error) {

	docList := []chartOfAccountModels.ChartOfAccountDoc{}

	err := r.pst.Find(ctx, &chartOfAccountModels.ChartOfAccountDoc{}, bson.M{"shopid": shopID}, &docList)
	if err != nil {
		return nil, err
	}

	return docList, nil
}
