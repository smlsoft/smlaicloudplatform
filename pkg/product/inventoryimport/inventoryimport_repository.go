package inventoryimport

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/inventoryimport/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IInventoryImportRepository interface {
	CreateInBatch(inventories []models.InventoryImportDoc) error
	DeleteInBatch(shopID string, guidList []string) error
	DeleteInBatchCode(shopID string, codeList []string) error
	FindPage(shopID string, page int, limit int) ([]models.InventoryImportInfo, paginate.PaginationData, error)
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

func (repo InventoryImportRepository) FindPage(shopID string, page int, limit int) ([]models.InventoryImportInfo, paginate.PaginationData, error) {

	docList := []models.InventoryImportInfo{}
	pagination, err := repo.pst.FindPage(&models.InventoryImportInfo{}, limit, page, bson.M{
		"shopid": shopID,
	}, &docList)

	if err != nil {
		return []models.InventoryImportInfo{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}
