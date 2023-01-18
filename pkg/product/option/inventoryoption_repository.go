package option

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/option/models"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IOptionRepository interface {
	Create(doc models.InventoryOptionMainDoc) (string, error)
	Update(shopID string, guid string, doc models.InventoryOptionMainDoc) error
	Delete(shopID string, guid string, username string) error
	FindByGuid(shopID string, guid string) (models.InventoryOptionMainDoc, error)
	FindPage(shopID string, q string, page int, limit int) ([]models.InventoryOptionMainInfo, mongopagination.PaginationData, error)
}

type OptionRepository struct {
	pst microservice.IPersisterMongo
}

func NewOptionRepository(pst microservice.IPersisterMongo) *OptionRepository {
	return &OptionRepository{
		pst: pst,
	}
}

func (repo OptionRepository) Create(doc models.InventoryOptionMainDoc) (string, error) {
	idx, err := repo.pst.Create(&models.InventoryOptionMainDoc{}, doc)

	if err != nil {
		return "", err
	}

	return idx.Hex(), nil
}

func (repo OptionRepository) Update(shopID string, guid string, doc models.InventoryOptionMainDoc) error {

	filterDoc := map[string]interface{}{
		"shopid":    shopID,
		"guidfixed": guid,
	}

	err := repo.pst.UpdateOne(&models.InventoryOptionMainDoc{}, filterDoc, doc)

	if err != nil {
		return err
	}

	return nil
}

func (repo OptionRepository) Delete(shopID string, guid string, username string) error {
	err := repo.pst.SoftDelete(&models.InventoryOptionMainDoc{}, username, bson.M{"guidfixed": guid, "shopid": shopID})

	if err != nil {
		return err
	}

	return nil
}

func (repo OptionRepository) FindByGuid(shopID string, guid string) (models.InventoryOptionMainDoc, error) {

	doc := &models.InventoryOptionMainDoc{}
	err := repo.pst.FindOne(&models.InventoryOptionMainDoc{}, bson.M{"guidfixed": guid, "shopid": shopID, "deletedat": bson.M{"$exists": false}}, doc)

	if err != nil {
		return models.InventoryOptionMainDoc{}, err
	}

	return *doc, nil
}

func (repo OptionRepository) FindPage(shopID string, q string, page int, limit int) ([]models.InventoryOptionMainInfo, mongopagination.PaginationData, error) {

	docList := []models.InventoryOptionMainInfo{}
	pagination, err := repo.pst.FindPage(&models.InventoryOptionMainInfo{}, limit, page, bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"guidfixed": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
			bson.M{"code": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
			bson.M{"name1": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &docList)

	if err != nil {
		return []models.InventoryOptionMainInfo{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}
