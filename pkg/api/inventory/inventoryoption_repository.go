package inventory

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IInventoryOptionMainRepository interface {
	Create(doc models.InventoryOptionMainDoc) (string, error)
	Update(guid string, doc models.InventoryOptionMainDoc) error
	Delete(shopID string, guid string, username string) error
	FindByGuid(shopID string, guid string) (models.InventoryOptionMainDoc, error)
	FindPage(shopID string, q string, page int, limit int) ([]models.InventoryOptionMainInfo, paginate.PaginationData, error)
}

type InventoryOptionMainRepository struct {
	pst microservice.IPersisterMongo
}

func NewInventoryOptionMainRepository(pst microservice.IPersisterMongo) InventoryOptionMainRepository {
	return InventoryOptionMainRepository{
		pst: pst,
	}
}

func (repo InventoryOptionMainRepository) Create(doc models.InventoryOptionMainDoc) (string, error) {
	idx, err := repo.pst.Create(&models.InventoryOptionMainDoc{}, doc)

	if err != nil {
		return "", err
	}

	return idx.Hex(), nil
}

func (repo InventoryOptionMainRepository) Update(guid string, doc models.InventoryOptionMainDoc) error {
	err := repo.pst.UpdateOne(&models.InventoryOptionMainDoc{}, "guidfixed", guid, doc)

	if err != nil {
		return err
	}

	return nil
}

func (repo InventoryOptionMainRepository) Delete(shopID string, guid string, username string) error {
	err := repo.pst.SoftDelete(&models.InventoryOptionMainDoc{}, username, bson.M{"guidfixed": guid, "shopid": shopID})

	if err != nil {
		return err
	}

	return nil
}

func (repo InventoryOptionMainRepository) FindByGuid(shopID string, guid string) (models.InventoryOptionMainDoc, error) {

	doc := &models.InventoryOptionMainDoc{}
	err := repo.pst.FindOne(&models.InventoryOptionMainDoc{}, bson.M{"guidfixed": guid, "shopid": shopID, "deletedat": bson.M{"$exists": false}}, doc)

	if err != nil {
		return models.InventoryOptionMainDoc{}, err
	}

	return *doc, nil
}

func (repo InventoryOptionMainRepository) FindPage(shopID string, q string, page int, limit int) ([]models.InventoryOptionMainInfo, paginate.PaginationData, error) {

	docList := []models.InventoryOptionMainInfo{}
	pagination, err := repo.pst.FindPage(&models.InventoryOptionMainInfo{}, limit, page, bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"guidfixed": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
			bson.M{"inventoryID": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
			bson.M{"optionGroupID": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &docList)

	if err != nil {
		return []models.InventoryOptionMainInfo{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}
