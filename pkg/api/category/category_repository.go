package category

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
)

type ICategoryRepository interface {
	Count(shopID string) (int, error)
	Create(category models.CategoryDoc) (string, error)
	CreateInBatch(inventories []models.CategoryDoc) error
	Update(shopID string, guid string, category models.CategoryDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindOne(shopID string, filters map[string]interface{}) (models.CategoryDoc, error)
	FindByGuid(shopID string, guid string) (models.CategoryDoc, error)
	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.CategoryItemGuid, error)
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.CategoryInfo, mongopagination.PaginationData, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.CategoryDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.CategoryActivity, mongopagination.PaginationData, error)
}

type CategoryRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.CategoryDoc]
	repositories.SearchRepository[models.CategoryInfo]
	repositories.GuidRepository[models.CategoryItemGuid]
	repositories.ActivityRepository[models.CategoryActivity, models.CategoryDeleteActivity]
}

func NewCategoryRepository(pst microservice.IPersisterMongo) CategoryRepository {
	insRepo := CategoryRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.CategoryDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.CategoryInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.CategoryItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.CategoryActivity, models.CategoryDeleteActivity](pst)

	return insRepo

}

// func (repo CategoryRepository) Count(shopID string) (int, error) {

// 	count, err := repo.pst.Count(&models.CategoryDoc{}, bson.M{"shopid": shopID})

// 	if err != nil {
// 		return 0, err
// 	}
// 	return count, nil
// }

// func (repo CategoryRepository) Create(category models.CategoryDoc) (string, error) {
// 	idx, err := repo.pst.Create(&models.CategoryDoc{}, category)

// 	if err != nil {
// 		return "", err
// 	}

// 	return idx.Hex(), nil
// }

// func (repo CategoryRepository) CreateInBatch(docList []models.CategoryDoc) error {
// 	var tempList []interface{}

// 	for _, inv := range docList {
// 		tempList = append(tempList, inv)
// 	}

// 	err := repo.pst.CreateInBatch(&models.CategoryDoc{}, tempList)

// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (repo CategoryRepository) Update(shopID string, guid string, category models.CategoryDoc) error {

// 	filterDoc := map[string]interface{}{
// 		"shopid":    shopID,
// 		"guidfixed": guid,
// 	}

// 	err := repo.pst.UpdateOne(&models.CategoryDoc{}, filterDoc, category)

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (repo CategoryRepository) Delete(shopID string, guid string, username string) error {
// 	err := repo.pst.SoftDelete(&models.CategoryDoc{}, username, bson.M{"guidfixed": guid, "shopid": shopID})

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (repo CategoryRepository) FindByGuid(shopID string, guid string) (models.CategoryDoc, error) {

// 	doc := &models.CategoryDoc{}

// 	err := repo.pst.FindOne(&models.CategoryInfo{}, bson.M{"guidfixed": guid, "shopid": shopID, "deletedat": bson.M{"$exists": false}}, doc)

// 	if err != nil {
// 		return models.CategoryDoc{}, err
// 	}

// 	return *doc, nil
// }

// func (repo CategoryRepository) FindByCategoryGuid(shopID string, categoryguid string) (models.CategoryDoc, error) {

// 	doc := &models.CategoryDoc{}

// 	err := repo.pst.FindOne(&models.CategoryInfo{}, bson.M{"shopid": shopID, "categoryguid": categoryguid, "deletedat": bson.M{"$exists": false}}, doc)

// 	if err != nil {
// 		return models.CategoryDoc{}, err
// 	}

// 	return *doc, nil
// }

// func (repo CategoryRepository) FindByCategoryGuidList(shopID string, categoryGuidList []string) ([]models.CategoryItemGuid, error) {

// 	findDoc := []models.CategoryItemGuid{}
// 	err := repo.pst.Find(&models.CategoryItemGuid{}, bson.M{"shopid": shopID, "categoryguid": bson.M{"$in": categoryGuidList}, "deletedat": bson.M{"$exists": false}}, &findDoc)

// 	if err != nil {
// 		return []models.CategoryItemGuid{}, err
// 	}
// 	return findDoc, nil
// }

// func (repo CategoryRepository) FindPage(shopID string, q string, page int, limit int) ([]models.CategoryInfo, paginate.PaginationData, error) {

// 	docList := []models.CategoryInfo{}
// 	pagination, err := repo.pst.FindPage(&models.CategoryInfo{}, limit, page, bson.M{
// 		"shopid":    shopID,
// 		"deletedat": bson.M{"$exists": false},
// 		"$or": []interface{}{
// 			bson.M{"guidfixed": q},
// 			bson.M{"name1": bson.M{"$regex": primitive.Regex{
// 				Pattern: ".*" + q + ".*",
// 				Options: "",
// 			}}},
// 		},
// 	}, &docList)

// 	if err != nil {
// 		return []models.CategoryInfo{}, paginate.PaginationData{}, err
// 	}

// 	return docList, pagination, nil
// }

// func (repo CategoryRepository) FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.CategoryDeleteActivity, paginate.PaginationData, error) {

// 	docList := []models.CategoryDeleteActivity{}
// 	pagination, err := repo.pst.FindPage(&models.CategoryDeleteActivity{}, limit, page, bson.M{
// 		"shopid":    shopID,
// 		"deletedat": bson.M{"$gte": lastUpdatedDate},
// 	}, &docList)

// 	if err != nil {
// 		return []models.CategoryDeleteActivity{}, paginate.PaginationData{}, err
// 	}

// 	return docList, pagination, nil
// }

// func (repo CategoryRepository) FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.CategoryActivity, paginate.PaginationData, error) {

// 	docList := []models.CategoryActivity{}
// 	pagination, err := repo.pst.FindPage(&models.CategoryActivity{}, limit, page, bson.M{
// 		"shopid":    shopID,
// 		"deletedat": bson.M{"$not": bson.M{"$gte": lastUpdatedDate}},
// 		"$or": []interface{}{
// 			bson.M{"createdat": bson.M{"$gte": lastUpdatedDate}},
// 			bson.M{"updatedat": bson.M{"$gte": lastUpdatedDate}},
// 		},
// 	}, &docList)

// 	if err != nil {
// 		return []models.CategoryActivity{}, paginate.PaginationData{}, err
// 	}

// 	return docList, pagination, nil
// }
