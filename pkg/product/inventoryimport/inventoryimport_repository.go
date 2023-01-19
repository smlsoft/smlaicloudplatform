package inventoryimport

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/product/inventoryimport/models"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IInventoryImportRepository interface {
	CreateInBatch(inventories []models.InventoryImportDoc) error
	DeleteInBatch(shopID string, guidList []string) error
	DeleteInBatchCode(shopID string, codeList []string) error
	FindPage(shopID string, pageable micromodels.Pageable) ([]models.InventoryImportInfo, mongopagination.PaginationData, error)
}

type InventoryImportRepository struct {
	pst microservice.IPersisterMongo
}

func NewInventoryImportRepository(pst microservice.IPersisterMongo) InventoryImportRepository {
	return InventoryImportRepository{
		pst: pst,
	}
}

func (repo InventoryImportRepository) CreateInBatch(inventories []models.InventoryImportDoc) error {
	var tempList []interface{}

	for _, inv := range inventories {
		tempList = append(tempList, inv)
	}

	err := repo.pst.CreateInBatch(&models.InventoryImportDoc{}, tempList)

	if err != nil {
		return err
	}
	return nil
}

func (repo InventoryImportRepository) DeleteInBatch(shopID string, guidList []string) error {

	err := repo.pst.Delete(&models.InventoryImportDoc{}, bson.M{
		"shopid":    shopID,
		"guidfixed": bson.M{"$in": guidList},
	})

	if err != nil {
		return err
	}
	return nil
}

func (repo InventoryImportRepository) DeleteInBatchCode(shopID string, codeList []string) error {

	err := repo.pst.Delete(&models.InventoryImportDoc{}, bson.M{
		"shopid":   shopID,
		"itemcode": bson.M{"$in": codeList},
	})

	if err != nil {
		return err
	}
	return nil
}

func (repo InventoryImportRepository) FindPage(shopID string, pageable micromodels.Pageable) ([]models.InventoryImportInfo, mongopagination.PaginationData, error) {

	filterQueries := bson.M{
		"shopid": shopID,
	}

	docList := []models.InventoryImportInfo{}
	pagination, err := repo.pst.FindPage(&models.InventoryImportInfo{}, filterQueries, pageable, &docList)

	if err != nil {
		return []models.InventoryImportInfo{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}
