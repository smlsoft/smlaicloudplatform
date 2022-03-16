package inventory

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IOptionGroupRepository interface {
	Count(shopId string) (int, error)
	Create(category models.InventoryOptionGroup) (string, error)
	Update(guid string, category models.InventoryOptionGroup) error
	Delete(guid string, shopId string) error
	FindByGuid(guid string, shopId string) (models.InventoryOptionGroup, error)
	FindPage(shopId string, q string, page int, limit int) ([]models.InventoryOptionGroup, paginate.PaginationData, error)
}

type OptionGroupRepository struct {
	pst microservice.IPersisterMongo
}

func NewOptionGroupRepository(pst microservice.IPersisterMongo) IOptionGroupRepository {
	return &OptionGroupRepository{
		pst: pst,
	}
}

func (repo *OptionGroupRepository) Count(shopId string) (int, error) {

	count, err := repo.pst.Count(&models.InventoryOptionGroup{}, bson.M{"shopId": shopId})

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *OptionGroupRepository) Create(category models.InventoryOptionGroup) (string, error) {
	idx, err := repo.pst.Create(&models.InventoryOptionGroup{}, category)

	if err != nil {
		return "", err
	}

	return idx.Hex(), nil
}

func (repo *OptionGroupRepository) Update(guid string, category models.InventoryOptionGroup) error {
	err := repo.pst.UpdateOne(&models.InventoryOptionGroup{}, "guidFixed", guid, category)

	if err != nil {
		return err
	}

	return nil
}

func (repo *OptionGroupRepository) Delete(guid string, shopId string) error {
	err := repo.pst.SoftDelete(&models.InventoryOptionGroup{}, bson.M{"guidFixed": guid, "shopId": shopId})

	if err != nil {
		return err
	}

	return nil
}

func (repo *OptionGroupRepository) FindByGuid(guid string, shopId string) (models.InventoryOptionGroup, error) {

	doc := &models.InventoryOptionGroup{}
	err := repo.pst.FindOne(&models.InventoryOptionGroup{}, bson.M{"guidFixed": guid, "shopId": shopId, "deleted": false}, doc)

	if err != nil {
		return models.InventoryOptionGroup{}, err
	}

	return *doc, nil
}

func (repo *OptionGroupRepository) FindPage(shopId string, q string, page int, limit int) ([]models.InventoryOptionGroup, paginate.PaginationData, error) {

	docList := []models.InventoryOptionGroup{}
	pagination, err := repo.pst.FindPage(&models.InventoryOptionGroup{}, limit, page, bson.M{
		"shopId":  shopId,
		"deleted": false,
		"$or": []interface{}{
			bson.M{"guidFixed": q},
			bson.M{"optionName1": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &docList)

	if err != nil {
		return []models.InventoryOptionGroup{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}
