package inventoryimport

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/inventoryimport/models"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IInventoryOptionMainImportRepository interface {
	CreateInBatch(inventories []models.InventoryOptionMainImportDoc) error
	DeleteInBatch(shopID string, guidList []string) error
	DeleteInBatchCode(shopID string, codeList []string) error
	FindPage(shopID string, page int, limit int) ([]models.InventoryOptionMainImportInfo, mongopagination.PaginationData, error)
}

type InventoryOptionMainImportRepository struct {
	pst microservice.IPersisterMongo
}

func NewInventoryOptionMainImportRepository(pst microservice.IPersisterMongo) InventoryOptionMainImportRepository {
	return InventoryOptionMainImportRepository{
		pst: pst,
	}
}
func (repo InventoryOptionMainImportRepository) CreateInBatch(inventories []models.InventoryOptionMainImportDoc) error {
	var tempList []interface{}

	for _, inv := range inventories {
		tempList = append(tempList, inv)
	}

	err := repo.pst.CreateInBatch(&models.InventoryOptionMainImportDoc{}, tempList)

	if err != nil {
		return err
	}
	return nil
}

func (repo InventoryOptionMainImportRepository) DeleteInBatch(shopID string, guidList []string) error {

	err := repo.pst.Delete(&models.InventoryOptionMainImportDoc{}, bson.M{
		"shopid":    shopID,
		"guidfixed": bson.M{"$in": guidList},
	})

	if err != nil {
		return err
	}
	return nil
}

func (repo InventoryOptionMainImportRepository) DeleteInBatchCode(shopID string, codeList []string) error {

	err := repo.pst.Delete(&models.InventoryOptionMainImportDoc{}, bson.M{
		"shopid": shopID,
		"code":   bson.M{"$in": codeList},
	})

	if err != nil {
		return err
	}
	return nil
}

func (repo InventoryOptionMainImportRepository) FindPage(shopID string, page int, limit int) ([]models.InventoryOptionMainImportInfo, mongopagination.PaginationData, error) {

	docList := []models.InventoryOptionMainImportInfo{}
	pagination, err := repo.pst.FindPage(&models.InventoryOptionMainImportInfo{}, limit, page, bson.M{
		"shopid": shopID,
	}, &docList)

	if err != nil {
		return []models.InventoryOptionMainImportInfo{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}
