package category

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ICategoryRepository interface {
	Count(shopID string) (int, error)
	Create(category models.CategoryDoc) (string, error)
	Update(guid string, category models.CategoryDoc) error
	Delete(guid string, username string, shopID string) error
	FindByGuid(guid string, shopID string) (models.CategoryDoc, error)
	FindPage(shopID string, q string, page int, limit int) ([]models.CategoryInfo, paginate.PaginationData, error)
}

type CategoryRepository struct {
	pst microservice.IPersisterMongo
}

func NewCategoryRepository(pst microservice.IPersisterMongo) CategoryRepository {
	return CategoryRepository{
		pst: pst,
	}
}

func (repo CategoryRepository) Count(shopID string) (int, error) {

	count, err := repo.pst.Count(&models.CategoryDoc{}, bson.M{"shopID": shopID})

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo CategoryRepository) Create(category models.CategoryDoc) (string, error) {
	idx, err := repo.pst.Create(&models.CategoryDoc{}, category)

	if err != nil {
		return "", err
	}

	return idx.Hex(), nil
}

func (repo CategoryRepository) Update(guid string, category models.CategoryDoc) error {
	err := repo.pst.UpdateOne(&models.CategoryDoc{}, "guidFixed", guid, category)

	if err != nil {
		return err
	}

	return nil
}

func (repo CategoryRepository) Delete(guid string, shopID string, username string) error {
	err := repo.pst.SoftDelete(&models.CategoryDoc{}, username, bson.M{"guidFixed": guid, "shopID": shopID})

	if err != nil {
		return err
	}

	return nil
}

func (repo CategoryRepository) FindByGuid(guid string, shopID string) (models.CategoryDoc, error) {

	doc := &models.CategoryDoc{}

	err := repo.pst.FindOne(&models.CategoryInfo{}, bson.M{"guidFixed": guid, "shopID": shopID, "deletedAt": bson.M{"$exists": false}}, doc)

	if err != nil {
		return models.CategoryDoc{}, err
	}

	return *doc, nil
}

func (repo CategoryRepository) FindPage(shopID string, q string, page int, limit int) ([]models.CategoryInfo, paginate.PaginationData, error) {

	docList := []models.CategoryInfo{}
	pagination, err := repo.pst.FindPage(&models.CategoryInfo{}, limit, page, bson.M{
		"shopID":    shopID,
		"deletedAt": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"guidFixed": q},
			bson.M{"name1": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &docList)

	if err != nil {
		return []models.CategoryInfo{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}
