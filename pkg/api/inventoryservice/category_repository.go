package inventoryservice

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ICategoryRepository interface{}

type CategoryRepository struct {
	pst microservice.IPersisterMongo
}

func NewCategoryRepository(pst microservice.IPersisterMongo) ICategoryRepository {
	return &CategoryRepository{
		pst: pst,
	}
}

func (repo *CategoryRepository) Count(merchantId string) (int, error) {

	count, err := repo.pst.Count(&models.Category{}, bson.M{"merchantId": merchantId})

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *CategoryRepository) Create(category models.Category) (string, error) {
	idx, err := repo.pst.Create(&models.Category{}, category)

	if err != nil {
		return "", err
	}

	return idx.Hex(), nil
}

func (repo *CategoryRepository) Update(guid string, category models.Category) error {
	err := repo.pst.UpdateOne(&models.Category{}, "guidFixed", guid, category)

	if err != nil {
		return err
	}

	return nil
}

func (repo *CategoryRepository) Delete(guid string) error {
	err := repo.pst.SoftDeleteByID(&models.Category{}, guid)

	if err != nil {
		return err
	}

	return nil
}

func (repo *CategoryRepository) FindByGuid(guid string, merchantId string) (models.Category, error) {

	doc := &models.Category{}
	err := repo.pst.FindOne(&models.Category{}, bson.M{"guidFixed": guid, "merchantId": merchantId, "deleted": false}, doc)

	if err != nil {
		return models.Category{}, err
	}

	return *doc, nil
}

func (repo *CategoryRepository) FindPage(merchantId string, q string, page int, limit int) ([]models.Category, paginate.PaginationData, error) {

	docList := []models.Category{}
	pagination, err := repo.pst.FindPage(&models.Category{}, limit, page, bson.M{
		"merchantId": merchantId,
		"deleted":    false,
		"$or": []interface{}{
			bson.M{"guidFixed": q},
			bson.M{"name1": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &docList)

	if err != nil {
		return []models.Category{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}
