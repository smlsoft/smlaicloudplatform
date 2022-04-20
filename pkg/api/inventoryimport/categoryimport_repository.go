package inventoryimport

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
)

type ICategoryImportRepository interface {
	CreateInBatch(inventories []models.CategoryImportDoc) error
	DeleteInBatch(shopID string, guidList []string) error
	FindPage(shopID string, page int, limit int) ([]models.CategoryImportInfo, paginate.PaginationData, error)
}

type CategoryImportRepository struct {
	pst microservice.IPersisterMongo
}

func NewCategoryImportRepository(pst microservice.IPersisterMongo) CategoryImportRepository {
	return CategoryImportRepository{
		pst: pst,
	}
}
func (repo CategoryImportRepository) CreateInBatch(inventories []models.CategoryImportDoc) error {
	var tempList []interface{}

	for _, inv := range inventories {
		tempList = append(tempList, inv)
	}

	err := repo.pst.CreateInBatch(&models.CategoryImportDoc{}, tempList)

	if err != nil {
		return err
	}
	return nil
}

func (repo CategoryImportRepository) DeleteInBatch(shopID string, guidList []string) error {

	err := repo.pst.Delete(&models.CategoryImportDoc{}, bson.M{
		"shopid":    shopID,
		"guidfixed": bson.M{"$in": guidList},
	})

	if err != nil {
		return err
	}
	return nil
}

func (repo CategoryImportRepository) FindPage(shopID string, page int, limit int) ([]models.CategoryImportInfo, paginate.PaginationData, error) {

	docList := []models.CategoryImportInfo{}
	pagination, err := repo.pst.FindPage(&models.CategoryImportInfo{}, limit, page, bson.M{
		"shopid": shopID,
	}, &docList)

	if err != nil {
		return []models.CategoryImportInfo{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}
