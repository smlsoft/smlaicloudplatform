package inventoryimport

import (
	"context"
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/product/inventoryimport/models"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type ICategoryImportRepository interface {
	CreateInBatch(ctx context.Context, docList []models.CategoryImportDoc) error
	DeleteInBatch(ctx context.Context, shopID string, guidList []string) error
	DeleteInBatchCode(ctx context.Context, shopID string, codeList []string) error
	FindPage(ctx context.Context, shopID string, pageable micromodels.Pageable) ([]models.CategoryImportInfo, mongopagination.PaginationData, error)
}

type CategoryImportRepository struct {
	pst microservice.IPersisterMongo
}

func NewCategoryImportRepository(pst microservice.IPersisterMongo) CategoryImportRepository {
	return CategoryImportRepository{
		pst: pst,
	}
}
func (repo CategoryImportRepository) CreateInBatch(ctx context.Context, inventories []models.CategoryImportDoc) error {
	var tempList []interface{}

	for _, inv := range inventories {
		tempList = append(tempList, inv)
	}

	err := repo.pst.CreateInBatch(ctx, &models.CategoryImportDoc{}, tempList)

	if err != nil {
		return err
	}
	return nil
}

func (repo CategoryImportRepository) DeleteInBatch(ctx context.Context, shopID string, guidList []string) error {

	err := repo.pst.Delete(ctx, &models.CategoryImportDoc{}, bson.M{
		"shopid":    shopID,
		"guidfixed": bson.M{"$in": guidList},
	})

	if err != nil {
		return err
	}
	return nil
}

func (repo CategoryImportRepository) DeleteInBatchCode(ctx context.Context, shopID string, codeList []string) error {

	err := repo.pst.Delete(ctx, &models.CategoryImportDoc{}, bson.M{
		"shopid": shopID,
		"code":   bson.M{"$in": codeList},
	})

	if err != nil {
		return err
	}
	return nil
}

func (repo CategoryImportRepository) FindPage(ctx context.Context, shopID string, pageable micromodels.Pageable) ([]models.CategoryImportInfo, mongopagination.PaginationData, error) {
	filterQueries := bson.M{
		"shopid": shopID,
	}

	docList := []models.CategoryImportInfo{}
	pagination, err := repo.pst.FindPage(ctx, &models.CategoryImportInfo{}, filterQueries, pageable, &docList)

	if err != nil {
		return []models.CategoryImportInfo{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}
