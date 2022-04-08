package inventory

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IInventoryOptionMainRepository interface {
	Create(doc models.InventoryOptionMain) (string, error)
	Update(guid string, doc models.InventoryOptionMain) error
	Delete(guid string, shopID string, username string) error
	FindByGuid(guid string, shopID string) (models.InventoryOptionMain, error)
	FindPage(shopID string, q string, page int, limit int) ([]models.InventoryOptionMain, paginate.PaginationData, error)
}

type InventoryOptionMainRepository struct {
	pst microservice.IPersisterMongo
}

func NewInventoryOptionMainRepository(pst microservice.IPersisterMongo) InventoryOptionMainRepository {
	return InventoryOptionMainRepository{
		pst: pst,
	}
}

func (repo InventoryOptionMainRepository) Create(doc models.InventoryOptionMain) (string, error) {
	idx, err := repo.pst.Create(&models.InventoryOptionMain{}, doc)

	if err != nil {
		return "", err
	}

	return idx.Hex(), nil
}

func (repo InventoryOptionMainRepository) Update(guid string, doc models.InventoryOptionMain) error {
	err := repo.pst.UpdateOne(&models.InventoryOptionMain{}, "guidfixed", guid, doc)

	if err != nil {
		return err
	}

	return nil
}

func (repo InventoryOptionMainRepository) Delete(guid string, shopID string, username string) error {
	err := repo.pst.SoftDelete(&models.InventoryOptionMain{}, username, bson.M{"guidfixed": guid, "shopid": shopID})

	if err != nil {
		return err
	}

	return nil
}

func (repo InventoryOptionMainRepository) FindByGuid(guid string, shopID string) (models.InventoryOptionMain, error) {

	doc := &models.InventoryOptionMain{}
	err := repo.pst.FindOne(&models.InventoryOptionMain{}, bson.M{"guidfixed": guid, "shopid": shopID, "deletedat": bson.M{"$exists": false}}, doc)

	if err != nil {
		return models.InventoryOptionMain{}, err
	}

	return *doc, nil
}

func (repo InventoryOptionMainRepository) FindPage(shopID string, q string, page int, limit int) ([]models.InventoryOptionMain, paginate.PaginationData, error) {

	docList := []models.InventoryOptionMain{}
	pagination, err := repo.pst.FindPage(&models.InventoryOptionMain{}, limit, page, bson.M{
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
		return []models.InventoryOptionMain{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}
