package category

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ICategoryRepository interface {
	Count(shopID string) (int, error)
	Create(category models.CategoryDoc) (string, error)
	CreateInBatch(inventories []models.CategoryDoc) error
	Update(guid string, category models.CategoryDoc) error
	Delete(shopID string, guid string, username string) error
	FindByGuid(shopID string, guid string) (models.CategoryDoc, error)
	FindByCategoryGuid(shopID string, guid string) (models.CategoryDoc, error)
	FindByCategoryGuidList(shopID string, categoryGuidList []string) ([]models.CategoryItemCategoryGuid, error)
	FindPage(shopID string, q string, page int, limit int) ([]models.CategoryInfo, paginate.PaginationData, error)
	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.CategoryDeleteActivity, paginate.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.CategoryActivity, paginate.PaginationData, error)
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

	count, err := repo.pst.Count(&models.CategoryDoc{}, bson.M{"shopid": shopID})

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

func (repo CategoryRepository) CreateInBatch(inventories []models.CategoryDoc) error {
	var tempList []interface{}

	for _, inv := range inventories {
		tempList = append(tempList, inv)
	}

	err := repo.pst.CreateInBatch(&models.CategoryDoc{}, tempList)

	if err != nil {
		return err
	}
	return nil
}

func (repo CategoryRepository) Update(guid string, category models.CategoryDoc) error {
	err := repo.pst.UpdateOne(&models.CategoryDoc{}, "guidfixed", guid, category)

	if err != nil {
		return err
	}

	return nil
}

func (repo CategoryRepository) Delete(shopID string, guid string, username string) error {
	err := repo.pst.SoftDelete(&models.CategoryDoc{}, username, bson.M{"guidfixed": guid, "shopid": shopID})

	if err != nil {
		return err
	}

	return nil
}

func (repo CategoryRepository) FindByGuid(shopID string, guid string) (models.CategoryDoc, error) {

	doc := &models.CategoryDoc{}

	err := repo.pst.FindOne(&models.CategoryInfo{}, bson.M{"guidfixed": guid, "shopid": shopID, "deletedat": bson.M{"$exists": false}}, doc)

	if err != nil {
		return models.CategoryDoc{}, err
	}

	return *doc, nil
}

func (repo CategoryRepository) FindByCategoryGuid(shopID string, categoryguid string) (models.CategoryDoc, error) {

	doc := &models.CategoryDoc{}

	err := repo.pst.FindOne(&models.CategoryInfo{}, bson.M{"shopid": shopID, "categoryguid": categoryguid, "deletedat": bson.M{"$exists": false}}, doc)

	if err != nil {
		return models.CategoryDoc{}, err
	}

	return *doc, nil
}

func (repo CategoryRepository) FindByCategoryGuidList(shopID string, categoryGuidList []string) ([]models.CategoryItemCategoryGuid, error) {

	findDoc := []models.CategoryItemCategoryGuid{}
	err := repo.pst.Find(&models.CategoryItemCategoryGuid{}, bson.M{"shopid": shopID, "categoryguid": bson.M{"$in": categoryGuidList}, "deletedat": bson.M{"$exists": false}}, &findDoc)

	if err != nil {
		return []models.CategoryItemCategoryGuid{}, err
	}
	return findDoc, nil
}

func (repo CategoryRepository) FindPage(shopID string, q string, page int, limit int) ([]models.CategoryInfo, paginate.PaginationData, error) {

	docList := []models.CategoryInfo{}
	pagination, err := repo.pst.FindPage(&models.CategoryInfo{}, limit, page, bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"guidfixed": q},
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

func (repo CategoryRepository) FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.CategoryDeleteActivity, paginate.PaginationData, error) {

	docList := []models.CategoryDeleteActivity{}
	pagination, err := repo.pst.FindPage(&models.CategoryDeleteActivity{}, limit, page, bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$gte": lastUpdatedDate},
	}, &docList)

	if err != nil {
		return []models.CategoryDeleteActivity{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo CategoryRepository) FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.CategoryActivity, paginate.PaginationData, error) {

	docList := []models.CategoryActivity{}
	pagination, err := repo.pst.FindPage(&models.CategoryActivity{}, limit, page, bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$not": bson.M{"$gte": lastUpdatedDate}},
		"$or": []interface{}{
			bson.M{"createdat": bson.M{"$gte": lastUpdatedDate}},
			bson.M{"updatedat": bson.M{"$gte": lastUpdatedDate}},
		},
	}, &docList)

	if err != nil {
		return []models.CategoryActivity{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}
