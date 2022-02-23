package inventoryservice

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IInventoryOptionRepository interface {
	Create(doc models.InventoryOption) (string, error)
	Update(guid string, doc models.InventoryOption) error
	Delete(guid string) error
	FindByGuid(guid string, merchantId string) (models.InventoryOption, error)
	FindPage(merchantId string, q string, page int, limit int) ([]models.InventoryOption, paginate.PaginationData, error)
}

type InventoryOptionRepository struct {
	pst microservice.IPersisterMongo
}

func NewInventoryOptionRepository(pst microservice.IPersisterMongo) IInventoryOptionRepository {
	return &InventoryOptionRepository{
		pst: pst,
	}
}

func (repo InventoryOptionRepository) Create(doc models.InventoryOption) (string, error) {
	idx, err := repo.pst.Create(&models.InventoryOption{}, doc)

	if err != nil {
		return "", err
	}

	return idx.Hex(), nil
}

func (repo InventoryOptionRepository) Update(guid string, doc models.InventoryOption) error {
	err := repo.pst.UpdateOne(&models.InventoryOption{}, "guidFixed", guid, doc)

	if err != nil {
		return err
	}

	return nil
}

func (repo InventoryOptionRepository) Delete(guid string) error {
	err := repo.pst.SoftDeleteByID(&models.InventoryOption{}, guid)

	if err != nil {
		return err
	}

	return nil
}

func (repo InventoryOptionRepository) FindByGuid(guid string, merchantId string) (models.InventoryOption, error) {

	doc := &models.InventoryOption{}
	err := repo.pst.FindOne(&models.InventoryOption{}, bson.M{"guidFixed": guid, "merchantId": merchantId, "deleted": false}, doc)

	if err != nil {
		return models.InventoryOption{}, err
	}

	return *doc, nil
}

func (repo InventoryOptionRepository) FindPage(merchantId string, q string, page int, limit int) ([]models.InventoryOption, paginate.PaginationData, error) {

	docList := []models.InventoryOption{}
	pagination, err := repo.pst.FindPage(&models.InventoryOption{}, limit, page, bson.M{
		"merchantId": merchantId,
		"deleted":    false,
		"$or": []interface{}{
			bson.M{"guidFixed": q},
			bson.M{"optionName1": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &docList)

	if err != nil {
		return []models.InventoryOption{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}
