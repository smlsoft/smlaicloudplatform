package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/documentwarehouse/documentimage/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/utils/mogoutil"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IDocumentImageRepository interface {
	Create(doc models.DocumentImageDoc) (string, error)
	Update(shopID string, guid string, doc models.DocumentImageDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindOne(shopID string, filters map[string]interface{}) (models.DocumentImageDoc, error)
	FindByGuid(shopID string, guid string) (models.DocumentImageDoc, error)
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.DocumentImageInfo, mongopagination.PaginationData, error)

	UpdateDocumentImageStatus(shopID string, guid string, status int8) error
	UpdateDocumentImageStatusByDocumentRef(shopID string, docRef string, status int8) error
	SaveDocumentImageDocRefGroup(shopID string, docRef string, docImages []string) error
	ListDocumentImageGroup(shopID string, q string, page int, limit int) ([]models.DocumentImageGroup, mongopagination.PaginationData, error)
	GetDocumentImageGroup(shopID string, docRef string) (models.DocumentImageGroup, error)
}

type DocumentImageRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.DocumentImageDoc]
	repositories.SearchRepository[models.DocumentImageInfo]
}

func NewDocumentImageRepository(pst microservice.IPersisterMongo) DocumentImageRepository {
	insRepo := DocumentImageRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.DocumentImageDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.DocumentImageInfo](pst)

	return insRepo
}

func (repo DocumentImageRepository) UpdateDocumentImageStatus(shopID string, guid string, status int8) error {
	fillter := bson.M{
		"shopid":    shopID,
		"guidfixed": guid,
	}

	data := bson.M{
		"$set": bson.M{"status": status},
	}

	return repo.pst.Update(models.DocumentImageDoc{}, fillter, data)
}

func (repo DocumentImageRepository) UpdateDocumentImageStatusByDocumentRef(shopID string, docRef string, status int8) error {

	fillter := bson.M{
		"shopid":      shopID,
		"documentref": docRef,
	}

	data := bson.M{
		"$set": bson.M{"status": status},
	}

	return repo.pst.Update(models.DocumentImageDoc{}, fillter, data)
}

func (repo DocumentImageRepository) SaveDocumentImageDocRefGroup(shopID string, docRef string, docImages []string) error {

	fillter := bson.M{
		"shopid":    shopID,
		"guidfixed": bson.M{"$in": docImages},
	}

	data := bson.M{
		"$set": bson.M{"documentref": docRef},
	}

	return repo.pst.Update(models.DocumentImageDoc{}, fillter, data)
}

func (repo DocumentImageRepository) ListDocumentImageGroup(shopID string, q string, page int, limit int) ([]models.DocumentImageGroup, mongopagination.PaginationData, error) {

	shopQuery := bson.M{"$match": bson.M{
		"shopid": shopID,
		"documentref": bson.M{"$regex": primitive.Regex{
			Pattern: ".*" + q + ".*",
			Options: "",
		}},
		"name": bson.M{"$regex": primitive.Regex{
			Pattern: ".*" + q + ".*",
			Options: "",
		}},
	}}

	groupQuery := bson.M{"$group": bson.M{"_id": "$documentref", "documentimages": bson.M{"$push": bson.M{
		"guidfixed":  "$guidfixed",
		"name":       "$name",
		"imageuri":   "$imageuri",
		"docguidref": "$docguidref",
		"module":     "$module",
		"status":     "$status",
	}}}}

	projectQuery := bson.M{"$project": bson.M{"documentref": "$_id", "documentimages": 1}}

	aggData, err := repo.pst.AggregatePage(&models.DocumentImageGroup{}, limit, page, shopQuery, groupQuery, projectQuery)

	if err != nil {
		return nil, mongopagination.PaginationData{}, err
	}

	docList, err := mogoutil.AggregatePageDecode[models.DocumentImageGroup](aggData)

	if err != nil {
		return nil, mongopagination.PaginationData{}, err
	}

	return docList, aggData.Pagination, nil
}

func (repo DocumentImageRepository) GetDocumentImageGroup(shopID string, docRef string) (models.DocumentImageGroup, error) {

	shopQuery := bson.M{"$match": bson.M{"shopid": shopID, "documentref": docRef}}

	groupQuery := bson.M{"$group": bson.M{"_id": "$documentref", "documentimages": bson.M{"$push": bson.M{
		"guidfixed":  "$guidfixed",
		"name":       "$name",
		"imageuri":   "$imageuri",
		"docguidref": "$docguidref",
		"module":     "$module",
		"status":     "$status",
	}}}}

	projectQuery := bson.M{"$project": bson.M{"documentref": "$_id", "documentimages": 1}}

	results := []models.DocumentImageGroup{}
	err := repo.pst.Aggregate(&models.DocumentImageGroup{}, []interface{}{
		shopQuery,
		groupQuery,
		projectQuery,
	}, &results)

	if err != nil {
		return models.DocumentImageGroup{}, err
	}

	if len(results) < 1 {
		return models.DocumentImageGroup{}, nil
	}

	return results[0], nil
}
