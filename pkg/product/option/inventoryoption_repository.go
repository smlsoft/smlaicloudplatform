package option

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/product/option/models"
	"smlcloudplatform/pkg/repositories"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IOptionRepository interface {
	Create(doc models.InventoryOptionMainDoc) (string, error)
	Update(shopID string, guid string, doc models.InventoryOptionMainDoc) error
	Delete(shopID string, guid string, username string) error
	FindByGuid(shopID string, guid string) (models.InventoryOptionMainDoc, error)
	FindPage(shopID string, pageable micromodels.Pageable) ([]models.InventoryOptionMainInfo, mongopagination.PaginationData, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.InventoryOptionMainDoc, error)
}

type OptionRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.InventoryOptionMainDoc]
}

func NewOptionRepository(pst microservice.IPersisterMongo) *OptionRepository {
	insRepo := &OptionRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.InventoryOptionMainDoc](pst)
	return insRepo
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

func (repo OptionRepository) FindPage(shopID string, pageable micromodels.Pageable) ([]models.InventoryOptionMainInfo, mongopagination.PaginationData, error) {

	filterQueries := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"guidfixed": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + pageable.Query + ".*",
				Options: "",
			}}},
			bson.M{"code": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + pageable.Query + ".*",
				Options: "",
			}}},
			bson.M{"name1": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + pageable.Query + ".*",
				Options: "",
			}}},
		},
	}

	docList := []models.InventoryOptionMainInfo{}
	pagination, err := repo.pst.FindPage(&models.InventoryOptionMainInfo{}, filterQueries, pageable, &docList)

	if err != nil {
		return []models.InventoryOptionMainInfo{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}
