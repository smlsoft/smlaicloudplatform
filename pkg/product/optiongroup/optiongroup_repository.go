package optiongroup

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/product/optiongroup/models"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IOptionGroupRepository interface {
	Count(shopID string) (int, error)
	Create(category models.InventoryOptionGroup) (string, error)
	Update(shopID string, guid string, category models.InventoryOptionGroup) error
	Delete(shopID string, guid string, username string) error
	FindByGuid(shopID string, guid string) (models.InventoryOptionGroup, error)
	FindPage(shopID string, pageable micromodels.Pageable) ([]models.InventoryOptionGroup, mongopagination.PaginationData, error)
}

type OptionGroupRepository struct {
	pst microservice.IPersisterMongo
}

func NewOptionGroupRepository(pst microservice.IPersisterMongo) OptionGroupRepository {
	return OptionGroupRepository{
		pst: pst,
	}
}

func (repo OptionGroupRepository) Count(shopID string) (int, error) {

	count, err := repo.pst.Count(&models.InventoryOptionGroup{}, bson.M{"shopid": shopID})

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo OptionGroupRepository) Create(category models.InventoryOptionGroup) (string, error) {
	idx, err := repo.pst.Create(&models.InventoryOptionGroup{}, category)

	if err != nil {
		return "", err
	}

	return idx.Hex(), nil
}

func (repo OptionGroupRepository) Update(shopID string, guid string, category models.InventoryOptionGroup) error {
	filterDoc := map[string]interface{}{
		"shopid":    shopID,
		"guidfixed": guid,
	}
	err := repo.pst.UpdateOne(&models.InventoryOptionGroup{}, filterDoc, category)

	if err != nil {
		return err
	}

	return nil
}

func (repo OptionGroupRepository) Delete(shopID string, guid string, username string) error {
	err := repo.pst.SoftDelete(&models.InventoryOptionGroup{}, username, bson.M{"guidfixed": guid, "shopid": shopID})

	if err != nil {
		return err
	}

	return nil
}

func (repo OptionGroupRepository) FindByGuid(shopID string, guid string) (models.InventoryOptionGroup, error) {

	doc := &models.InventoryOptionGroup{}
	err := repo.pst.FindOne(&models.InventoryOptionGroup{}, bson.M{"guidfixed": guid, "shopid": shopID, "deletedat": bson.M{"$exists": false}}, doc)

	if err != nil {
		return models.InventoryOptionGroup{}, err
	}

	return *doc, nil
}

func (repo OptionGroupRepository) FindPage(shopID string, pageable micromodels.Pageable) ([]models.InventoryOptionGroup, mongopagination.PaginationData, error) {
	filterQueries := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"guidfixed": pageable.Query},
			bson.M{"optionName1": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + pageable.Query + ".*",
				Options: "",
			}}},
		},
	}

	docList := []models.InventoryOptionGroup{}
	pagination, err := repo.pst.FindPage(&models.InventoryOptionGroup{}, filterQueries, pageable, &docList)

	if err != nil {
		return []models.InventoryOptionGroup{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}
