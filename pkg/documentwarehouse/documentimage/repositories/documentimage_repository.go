package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/documentwarehouse/documentimage/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IDocumentImageRepository interface {
	Create(doc models.DocumentImageDoc) (string, error)
	CreateInBatch(doc []models.DocumentImageDoc) error
	Update(shopID string, guid string, doc models.DocumentImageDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindOne(shopID string, filters interface{}) (models.DocumentImageDoc, error)
	FindByReferenceDocNo(shopID string, docNo string) ([]models.DocumentImageDoc, error)
	FindByReference(shopID string, reference models.Reference) ([]models.DocumentImageDoc, error)
	FindByGuid(shopID string, guid string) (models.DocumentImageDoc, error)
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.DocumentImageInfo, mongopagination.PaginationData, error)
	FindPageFilterSort(shopID string, filters map[string]interface{}, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]models.DocumentImageInfo, mongopagination.PaginationData, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.DocumentImageItemGuid, error)
	FindInGUIDs(shopID string, docImageGUIDs []string) ([]models.DocumentImageDoc, error)

	FindAll() ([]models.DocumentImageDoc, error)
	UpdateAll(doc models.DocumentImageDoc) error

	// UpdateDocumentImageStatus(shopID string, guid string, docnoGUIDRef string, status int8) error
	// UpdateDocumentImageStatusByDocumentRef(shopID string, docRef string, docnoGUIDRef string, status int8) error
	// SaveDocumentImageDocRefGroup(shopID string, docRef string, docImages []models.DocumentImageGroup) error
	// ListDocumentImageGroup(shopID string, filters map[string]interface{}, q string, page int, limit int) ([]models.DocumentImageGroup, mongopagination.PaginationData, error)
	// GetDocumentImageGroup(shopID string, docRef string) (models.DocumentImageGroup, error)
}

type DocumentImageRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.DocumentImageDoc]
	repositories.SearchRepository[models.DocumentImageInfo]
	repositories.GuidRepository[models.DocumentImageItemGuid]
}

func NewDocumentImageRepository(pst microservice.IPersisterMongo) DocumentImageRepository {
	insRepo := DocumentImageRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.DocumentImageDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.DocumentImageInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.DocumentImageItemGuid](pst)

	return insRepo
}

func (repo DocumentImageRepository) FindInGUIDs(shopID string, docImageGUIDs []string) ([]models.DocumentImageDoc, error) {
	docList := []models.DocumentImageDoc{}
	err := repo.pst.Find(models.DocumentImageDoc{}, bson.M{
		"guidfixed": bson.M{
			"$in": docImageGUIDs,
		},
		"deletedat": bson.M{"$exists": false},
	}, &docList)

	if err != nil {
		return []models.DocumentImageDoc{}, err
	}

	return docList, nil
}

func (repo DocumentImageRepository) FindByReferenceDocNo(shopID string, docNo string) ([]models.DocumentImageDoc, error) {
	docList := []models.DocumentImageDoc{}
	err := repo.pst.Find(models.DocumentImageDoc{}, bson.M{
		"references.docno": docNo,
		"deletedat":        bson.M{"$exists": false},
	}, &docList)

	if err != nil {
		return []models.DocumentImageDoc{}, err
	}

	return docList, nil
}

func (repo DocumentImageRepository) FindByReference(shopID string, reference models.Reference) ([]models.DocumentImageDoc, error) {
	docList := []models.DocumentImageDoc{}
	err := repo.pst.Find(models.DocumentImageDoc{}, bson.M{
		"references.module": reference.Module,
		"references.docno":  reference.DocNo,
		"deletedat":         bson.M{"$exists": false},
	}, &docList)

	if err != nil {
		return []models.DocumentImageDoc{}, err
	}

	return docList, nil
}

func (repo DocumentImageRepository) UpdateReject(shopID string, authUsername string, updatedAt time.Time, docImageGUID string, isReject bool) error {
	fillter := bson.M{
		"shopid":    shopID,
		"guidfixed": docImageGUID,
	}

	data := bson.M{
		"$set": bson.M{"isreject": isReject, "updatedby": authUsername, "updatedat": updatedAt},
	}

	return repo.pst.Update(models.DocumentImageDoc{}, fillter, data)
}

func (repo DocumentImageRepository) FindAll() ([]models.DocumentImageDoc, error) {
	docList := []models.DocumentImageDoc{}
	err := repo.pst.Find(models.DocumentImageDoc{}, bson.M{}, &docList)

	if err != nil {
		return []models.DocumentImageDoc{}, err
	}

	return docList, nil
}

func (repo DocumentImageRepository) UpdateAll(doc models.DocumentImageDoc) error {
	filterDoc := map[string]interface{}{
		"guidfixed": doc.GuidFixed,
	}

	err := repo.pst.UpdateOne(models.DocumentImageDoc{}, filterDoc, doc)

	if err != nil {
		return err
	}

	return nil
}

/*
func (repo DocumentImageRepository) UpdateDocumentImageStatus(shopID string, guid string, docnoGUIDRef string, status int8) error {
	fillter := bson.M{
		"shopid":    shopID,
		"guidfixed": guid,
	}

	data := bson.M{
		"$set": bson.M{"status": status},
	}

	return repo.pst.Update(models.DocumentImageDoc{}, fillter, data)
}

func (repo DocumentImageRepository) UpdateDocumentImageStatusByDocumentRef(shopID string, docRef string, docnoGUIDRef string, status int8) error {

	fillter := bson.M{
		"shopid":      shopID,
		"documentref": docRef,
	}

	data := bson.M{
		"$set": bson.M{
			"docguidref": docnoGUIDRef,
			"status":     status,
		},
	}

	return repo.pst.Update(models.DocumentImageDoc{}, fillter, data)
}

func (repo DocumentImageRepository) SaveDocumentImageDocRefGroup(shopID string, docRef string, docImages []models.DocumentImageGroup) error {

	fillter := bson.M{
		"shopid":    shopID,
		"guidfixed": bson.M{"$in": docImages},
	}

	data := bson.M{
		"$set": bson.M{
			"documentref": docRef},
	}

	return repo.pst.Update(models.DocumentImageDoc{}, fillter, data)
}

func (repo DocumentImageRepository) ListDocumentImageGroup(shopID string, filters map[string]interface{}, q string, page int, limit int) ([]models.DocumentImageGroup, mongopagination.PaginationData, error) {

	searchFilter := bson.M{
		"shopid": shopID,
	}

	if len(filters) > 0 {
		for key, val := range filters {
			searchFilter[key] = val
		}
	}

	if len(strings.TrimSpace(q)) > 0 {

		searchFields := []string{"documentref", "name"}

		tempFilters := []interface{}{}

		for _, fieldName := range searchFields {
			tempFilters = append(tempFilters, bson.M{fieldName: bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}})
		}

		searchFilter["$or"] = tempFilters

	}

	searchQuery := bson.M{"$match": searchFilter}

	groupQuery := bson.M{"$group": bson.M{"_id": "$documentref", "documentimages": bson.M{"$push": bson.M{
		"guidfixed":  "$guidfixed",
		"name":       "$name",
		"imageuri":   "$imageuri",
		"docguidref": "$docguidref",
		"module":     "$module",
		"status":     "$status",
		"uploadedby": "$uploadedby",
		"uploadedat": "$uploadedat",
	}}}}

	projectQuery := bson.M{"$project": bson.M{"documentref": "$_id", "documentimages": 1}}

	aggData, err := repo.pst.AggregatePage(&models.DocumentImageGroup{}, limit, page, searchQuery, groupQuery, projectQuery)

	if err != nil {
		return nil, mongopagination.PaginationData{}, err
	}

	docList, err := mogoutil.AggregatePageDecode[models.DocumentImageGroup](aggData)

	if err != nil {
		return []models.DocumentImageGroup{}, mongopagination.PaginationData{}, err
	}

	if len(docList) < 1 {
		return []models.DocumentImageGroup{}, mongopagination.PaginationData{}, nil
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
		"uploadedby": "$uploadedby",
		"uploadedat": "$uploadedat",
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
*/
