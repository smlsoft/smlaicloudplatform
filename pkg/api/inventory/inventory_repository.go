package inventory

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IInventoryRepository interface {
	Create(inventory models.Inventory) (string, error)
	Update(guid string, inventory models.Inventory) error
	Delete(guid string, shopId string) error
	FindByGuid(guid string, shopId string) (models.Inventory, error)
	FindPage(shopId string, q string, page int, limit int) ([]models.Inventory, paginate.PaginationData, error)
}

type InventoryRepository struct {
	pst microservice.IPersisterMongo
}

func NewInventoryRepository(pst microservice.IPersisterMongo) IInventoryRepository {
	return &InventoryRepository{
		pst: pst,
	}
}

func (repo *InventoryRepository) Create(inventory models.Inventory) (string, error) {
	idx, err := repo.pst.Create(&models.Inventory{}, inventory)

	if err != nil {
		return "", err
	}
	return idx.Hex(), nil
}

func (repo *InventoryRepository) Update(guid string, inventory models.Inventory) error {

	err := repo.pst.UpdateOne(&models.Inventory{}, "guidFixed", guid, inventory)

	if err != nil {
		return err
	}

	return nil
}

func (repo *InventoryRepository) Delete(guid string, shopId string) error {

	err := repo.pst.SoftDelete(&models.Inventory{}, bson.M{"guidFixed": guid, "shopId": shopId})

	if err != nil {
		return err
	}
	return nil
}

func (repo *InventoryRepository) FindByGuid(guid string, shopId string) (models.Inventory, error) {

	findDoc := &models.Inventory{}
	err := repo.pst.FindOne(&models.Inventory{}, bson.M{"shopId": shopId, "guidFixed": guid, "deleted": false}, findDoc)

	if err != nil {
		return models.Inventory{}, err
	}
	return *findDoc, nil
}

func (repo *InventoryRepository) FindPage(shopId string, q string, page int, limit int) ([]models.Inventory, paginate.PaginationData, error) {

	docList := []models.Inventory{}
	pagination, err := repo.pst.FindPage(&models.Inventory{}, limit, page, bson.M{
		"shopId":  shopId,
		"deleted": false,
		"$or": []interface{}{
			bson.M{"guidFixed": q},
			bson.M{"name1": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &docList)

	if err != nil {
		return []models.Inventory{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}
