package inventoryimport

import (
	"context"
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/product/inventoryimport/models"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IInventoryOptionMainImportRepository interface {
	CreateInBatch(ctx context.Context, docs []models.InventoryOptionMainImportDoc) error
	DeleteInBatch(ctx context.Context, shopID string, guidList []string) error
	DeleteInBatchCode(ctx context.Context, shopID string, codeList []string) error
	FindPage(ctx context.Context, shopID string, pageable micromodels.Pageable) ([]models.InventoryOptionMainImportInfo, mongopagination.PaginationData, error)
}

type InventoryOptionMainImportRepository struct {
	pst microservice.IPersisterMongo
}

func NewInventoryOptionMainImportRepository(pst microservice.IPersisterMongo) InventoryOptionMainImportRepository {
	return InventoryOptionMainImportRepository{
		pst: pst,
	}
}
func (repo InventoryOptionMainImportRepository) CreateInBatch(ctx context.Context, docs []models.InventoryOptionMainImportDoc) error {
	var tempList []interface{}

	for _, inv := range docs {
		tempList = append(tempList, inv)
	}

	err := repo.pst.CreateInBatch(ctx, &models.InventoryOptionMainImportDoc{}, tempList)

	if err != nil {
		return err
	}
	return nil
}

func (repo InventoryOptionMainImportRepository) DeleteInBatch(ctx context.Context, shopID string, guidList []string) error {

	err := repo.pst.Delete(ctx, &models.InventoryOptionMainImportDoc{}, bson.M{
		"shopid":    shopID,
		"guidfixed": bson.M{"$in": guidList},
	})

	if err != nil {
		return err
	}
	return nil
}

func (repo InventoryOptionMainImportRepository) DeleteInBatchCode(ctx context.Context, shopID string, codeList []string) error {

	err := repo.pst.Delete(ctx, &models.InventoryOptionMainImportDoc{}, bson.M{
		"shopid": shopID,
		"code":   bson.M{"$in": codeList},
	})

	if err != nil {
		return err
	}
	return nil
}

func (repo InventoryOptionMainImportRepository) FindPage(ctx context.Context, shopID string, pageable micromodels.Pageable) ([]models.InventoryOptionMainImportInfo, mongopagination.PaginationData, error) {

	filterQueries := bson.M{
		"shopid": shopID,
	}

	docList := []models.InventoryOptionMainImportInfo{}
	pagination, err := repo.pst.FindPage(ctx, &models.InventoryOptionMainImportInfo{}, filterQueries, pageable, &docList)

	if err != nil {
		return []models.InventoryOptionMainImportInfo{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}
