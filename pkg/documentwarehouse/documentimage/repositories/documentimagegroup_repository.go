package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/documentwarehouse/documentimage/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/utils/mogoutil"

	micromodels "smlcloudplatform/internal/microservice/models"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IDocumentImageGroupRepository interface {
	CountByTask(shopID string, taskGUID string) (int, error)
	CountRejectByTask(shopID string, taskGUID string) (int, error)
	Create(doc models.DocumentImageGroupDoc) (string, error)
	CreateInBatch(doc []models.DocumentImageGroupDoc) error
	Update(shopID string, guid string, doc models.DocumentImageGroupDoc) error
	UpdateXOrder(shopID string, taskGUID string, GUID string, xorder uint) error
	DeleteByGuidfixed(shopID string, guid string) error
	DeleteByDocumentImageGUIDsHasOne(shopID string, imageGUIDs []string) error
	DeleteByGUIDsIsDocumentImageEmpty(shopID string, GUIDs []string) error
	RemoveDocumentImageByDocumentImageGUIDs(shopID string, imageGUIDs []string) error
	DeleteByDocumentImageGUIDsHasOneWithoutDocumentImageGroupGUID(shopID string, withoutGUID string, imageGUIDs []string) error
	DeleteByGUIDsIsDocumentImageEmptyWithoutDocumentImageGroupGUID(shopID string, withoutGUID string, GUIDs []string) error
	DeleteByGUIDIsDocumentImageEmpty(shopID string, imageGroupGUID string) error
	RemoveDocumentImageByDocumentImageGUIDsWithoutDocumentImageGroupGUID(shopID string, withoutGUID string, imageGUIDs []string) error
	FindOne(shopID string, filters interface{}) (models.DocumentImageGroupDoc, error)
	FindByGuid(shopID string, guid string) (models.DocumentImageGroupDoc, error)

	FindStatusByDocumentImageGroupTask(shopID string, taskGUID string) ([]models.DocumentImageGroupStatus, error)
	FindLastOneByTask(shopID string, taskGUID string) (models.DocumentImageGroupDoc, error)
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.DocumentImageGroupInfo, mongopagination.PaginationData, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.DocumentImageGroupInfo, mongopagination.PaginationData, error)
	FindByTaskGUID(shopID string, taskGUID string) ([]models.DocumentImageGroupDoc, error)

	UpdateTaskIsCompletedByTaskGUID(shopID string, taskGUID string, isCompleted bool) error
	FindOneByReference(shopID string, reference models.Reference) (models.DocumentImageGroupDoc, error)
	FindOneByDocumentImageGUID(shopID string, documentImageGUID string) (models.DocumentImageGroupDoc, error)
	FindByDocumentImageGUIDs(shopID string, documentImageGUIDs []string) ([]models.DocumentImageGroupInfo, error)
	FindByReference(shopID string, reference models.Reference) ([]models.DocumentImageGroupDoc, error)
	FindByReferenceDocNo(shopID string, docNo string) ([]models.DocumentImageGroupDoc, error)
	FindPageImageGroup(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.DocumentImageGroupInfo, mongopagination.PaginationData, error)
	Transaction(fnc func() error) error

	FindOneByDocumentImageGUIDAll(documentImageGUID string) (models.DocumentImageGroupDoc, error)
	UpdateStatusByTask(shopID string, taskGUID string, status int8) error
}

type DocumentImageGroupRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.DocumentImageGroupDoc]
	repositories.SearchRepository[models.DocumentImageGroupInfo]
}

func NewDocumentImageGroupRepository(pst microservice.IPersisterMongo) DocumentImageGroupRepository {
	insRepo := DocumentImageGroupRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.DocumentImageGroupDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.DocumentImageGroupInfo](pst)

	return insRepo
}

func (repo DocumentImageGroupRepository) Transaction(fnc func() error) error {
	return repo.pst.Transaction(fnc)
}

func (repo DocumentImageGroupRepository) UpdateStatusByTask(shopID string, taskGUID string, status int8) error {
	return repo.pst.Update(models.DocumentImageGroupDoc{}, bson.M{"shopid": shopID, "taskguid": taskGUID}, bson.M{"$set": bson.M{"status": status}})
}

func (repo DocumentImageGroupRepository) FindStatusByDocumentImageGroupTask(shopID string, taskGUID string) ([]models.DocumentImageGroupStatus, error) {

	filters := bson.M{
		"shopid":    shopID,
		"taskguid":  taskGUID,
		"deletedat": bson.M{"$exists": false},
	}
	docList := []models.DocumentImageGroupStatus{}
	err := repo.pst.Find(models.DocumentImageGroupDoc{}, filters, &docList)

	if err != nil {
		return []models.DocumentImageGroupStatus{}, err
	}

	return docList, nil
}

func (repo DocumentImageGroupRepository) CountByTask(shopID string, taskGUID string) (int, error) {

	filters := bson.M{
		"shopid":    shopID,
		"taskguid":  taskGUID,
		"deletedat": bson.M{"$exists": false},
	}

	return repo.pst.Count(models.DocumentImageGroupDoc{}, filters)
}

func (repo DocumentImageGroupRepository) CountRejectByTask(shopID string, taskGUID string) (int, error) {

	filters := bson.M{
		"shopid":   shopID,
		"taskguid": taskGUID,
		"$or": []interface{}{
			bson.M{"status": models.IMAGE_REJECT},
			bson.M{"status": models.IMAGE_REJECT_KEYING},
		},
		"deletedat": bson.M{"$exists": false},
	}

	return repo.pst.Count(models.DocumentImageGroupDoc{}, filters)
}

func (repo DocumentImageGroupRepository) UpdateTaskIsCompletedByTaskGUID(shopID string, taskGUID string, isCompleted bool) error {

	filters := bson.M{
		"shopid":    shopID,
		"taskguid":  taskGUID,
		"deletedat": bson.M{"$exists": false},
	}

	err := repo.pst.Update(models.DocumentImageGroupDoc{}, filters, bson.M{"$set": bson.M{"iscompleted": isCompleted}})

	if err != nil {
		return err
	}

	return nil
}

func (repo DocumentImageGroupRepository) UpdateXOrder(shopID string, taskGUID string, GUID string, xorder uint) error {

	filters := bson.M{
		"shopid":    shopID,
		"taskguid":  taskGUID,
		"guidfixed": GUID,
		"deletedat": bson.M{"$exists": false},
	}

	err := repo.pst.UpdateOne(models.DocumentImageGroupDoc{}, filters, bson.M{"xorder": xorder})

	if err != nil {
		return err
	}

	return nil
}

func (repo DocumentImageGroupRepository) FindLastOneByTask(shopID string, taskGUID string) (models.DocumentImageGroupDoc, error) {

	results := []models.DocumentImageGroupDoc{}
	err := repo.pst.Aggregate(models.DocumentImageGroupDoc{}, []interface{}{
		bson.M{"$match": bson.M{
			"shopid":    shopID,
			"taskguid":  taskGUID,
			"deletedat": bson.M{"$exists": false},
		}},
		bson.M{"$sort": bson.M{"xorder": -1}},
		bson.M{"$limit": 1},
	}, &results)

	if err != nil {
		return models.DocumentImageGroupDoc{}, err
	}

	if len(results) < 1 {
		return models.DocumentImageGroupDoc{}, nil
	}

	return results[0], nil

}

func (repo DocumentImageGroupRepository) FindOneByReference(shopID string, reference models.Reference) (models.DocumentImageGroupDoc, error) {

	results := []models.DocumentImageGroupDoc{}
	err := repo.pst.Aggregate(models.DocumentImageGroupDoc{}, []interface{}{
		bson.M{"$match": bson.M{
			"shopid":            shopID,
			"references.module": reference.Module,
			"references.docno":  reference.DocNo,
			"deletedat":         bson.M{"$exists": false},
		}},
		bson.M{"$limit": 1},
	}, &results)

	if err != nil {
		return models.DocumentImageGroupDoc{}, err
	}

	if len(results) < 1 {
		return models.DocumentImageGroupDoc{}, nil
	}

	return results[0], nil

}

func (repo DocumentImageGroupRepository) FindOneByDocumentImageGUIDAll(documentImageGUID string) (models.DocumentImageGroupDoc, error) {

	matchQuery := bson.M{"$match": bson.M{
		"imagereferences.documentimageguid": documentImageGUID,
	}}

	results := []models.DocumentImageGroupDoc{}
	err := repo.pst.Aggregate(models.DocumentImageGroupDoc{}, []interface{}{
		matchQuery,
		bson.M{"$limit": 1},
	}, &results)

	if err != nil {
		return models.DocumentImageGroupDoc{}, err
	}

	if len(results) < 1 {
		return models.DocumentImageGroupDoc{}, nil
	}

	return results[0], nil
}

func (repo DocumentImageGroupRepository) FindOneByDocumentImageGUID(shopID string, documentImageGUID string) (models.DocumentImageGroupDoc, error) {

	matchQuery := bson.M{"$match": bson.M{
		"shopid":                            shopID,
		"imagereferences":                   bson.M{"$exists": true},
		"imagereferences.documentimageguid": documentImageGUID,
		"deletedat":                         bson.M{"$exists": false},
	}}

	results := []models.DocumentImageGroupDoc{}
	err := repo.pst.Aggregate(models.DocumentImageGroupDoc{}, []interface{}{
		matchQuery,
		bson.M{"$limit": 1},
	}, &results)

	if err != nil {
		return models.DocumentImageGroupDoc{}, err
	}

	if len(results) < 1 {
		return models.DocumentImageGroupDoc{}, nil
	}

	return results[0], nil
}

func (repo DocumentImageGroupRepository) FindWithoutGUIDByDocumentImageGUIDs(shopID string, documentImageGroupGUID string, documentImageGUIDs []string) ([]models.DocumentImageGroupInfo, error) {

	matchQuery := bson.M{"$match": bson.M{
		"shopid":                            shopID,
		"guidfixed":                         bson.M{"$ne": documentImageGroupGUID},
		"imagereferences":                   bson.M{"$exists": true},
		"imagereferences.documentimageguid": bson.M{"$in": documentImageGUIDs},
		"deletedat":                         bson.M{"$exists": false},
	}}

	results := []models.DocumentImageGroupInfo{}
	err := repo.pst.Aggregate(models.DocumentImageGroupDoc{}, []interface{}{
		matchQuery,
	}, &results)

	if err != nil {
		return []models.DocumentImageGroupInfo{}, err
	}

	return results, nil
}

func (repo DocumentImageGroupRepository) FindByDocumentImageGUIDs(shopID string, documentImageGUIDs []string) ([]models.DocumentImageGroupInfo, error) {

	matchQuery := bson.M{"$match": bson.M{
		"shopid":                            shopID,
		"imagereferences":                   bson.M{"$exists": true},
		"imagereferences.documentimageguid": bson.M{"$in": documentImageGUIDs},
		"deletedat":                         bson.M{"$exists": false},
	}}

	results := []models.DocumentImageGroupInfo{}
	err := repo.pst.Aggregate(models.DocumentImageGroupDoc{}, []interface{}{
		matchQuery,
	}, &results)

	if err != nil {
		return []models.DocumentImageGroupInfo{}, err
	}

	return results, nil
}

func (repo DocumentImageGroupRepository) FindByReferenceDocNo(shopID string, docNo string) ([]models.DocumentImageGroupDoc, error) {
	docList := []models.DocumentImageGroupDoc{}
	err := repo.pst.Find(models.DocumentImageGroupDoc{}, bson.M{
		"references.docno": docNo,
		"deletedat":        bson.M{"$exists": false},
	}, &docList)

	if err != nil {
		return []models.DocumentImageGroupDoc{}, err
	}

	return docList, nil
}

func (repo DocumentImageGroupRepository) FindByTaskGUID(shopID string, taskGUID string) ([]models.DocumentImageGroupDoc, error) {
	docList := []models.DocumentImageGroupDoc{}
	err := repo.pst.Find(models.DocumentImageGroupDoc{}, bson.M{
		"taskguid":  taskGUID,
		"deletedat": bson.M{"$exists": false},
	}, &docList)

	if err != nil {
		return []models.DocumentImageGroupDoc{}, err
	}

	return docList, nil
}

func (repo DocumentImageGroupRepository) FindByReference(shopID string, reference models.Reference) ([]models.DocumentImageGroupDoc, error) {
	docList := []models.DocumentImageGroupDoc{}
	err := repo.pst.Find(models.DocumentImageGroupDoc{}, bson.M{
		"references.module": reference.Module,
		"references.docno":  reference.DocNo,
		"deletedat":         bson.M{"$exists": false},
	}, &docList)

	if err != nil {
		return []models.DocumentImageGroupDoc{}, err
	}

	return docList, nil
}

func (repo DocumentImageGroupRepository) DeleteByGuidfixed(shopID string, guid string) error {
	return repo.pst.Delete(models.DocumentImageGroupDoc{}, bson.M{
		"shopid":    shopID,
		"guidfixed": guid,
	})
}

func (repo DocumentImageGroupRepository) DeleteByGUIDsIsDocumentImageEmpty(shopID string, GUIDs []string) error {
	return repo.pst.Delete(models.DocumentImageGroupDoc{}, bson.M{
		"shopid":          shopID,
		"imagereferences": bson.M{"$exists": true, "$size": 0},
		"$or": []interface{}{
			bson.M{"references": bson.M{"$exists": false}},
			bson.M{"references": bson.M{"$size": 0}},
			bson.M{"references": nil},
		},
	})
}

func (repo DocumentImageGroupRepository) DeleteByDocumentImageGUIDsHasOne(shopID string, imageGUIDs []string) error {
	return repo.pst.Delete(models.DocumentImageGroupDoc{}, bson.M{
		"shopid":                            shopID,
		"imagereferences":                   bson.M{"$exists": true, "$size": 1},
		"imagereferences.documentimageguid": bson.M{"$in": imageGUIDs},
		"$or": []interface{}{
			bson.M{"references": bson.M{"$exists": false}},
			bson.M{"references": bson.M{"$size": 0}},
		},
	})
}

func (repo DocumentImageGroupRepository) RemoveDocumentImageByDocumentImageGUIDs(shopID string, imageGUIDs []string) error {

	filterQuery := bson.M{
		"shopid":                            shopID,
		"imagereferences.documentimageguid": bson.M{"$in": imageGUIDs},
		"imagereferences":                   bson.M{"$exists": true},
		"$or": []interface{}{
			bson.M{"references": bson.M{"$exists": false}},
			bson.M{"references": bson.M{"$size": 0}},
			bson.M{"references": nil},
		},
	}

	removeQuery := bson.M{
		"$pull": bson.M{"imagereferences": bson.M{"documentimageguid": bson.M{"$in": imageGUIDs}}},
	}

	return repo.pst.Update(models.DocumentImageGroupDoc{}, filterQuery, removeQuery)
}

func (repo DocumentImageGroupRepository) DeleteByGUIDsIsDocumentImageEmptyWithoutDocumentImageGroupGUID(shopID string, withoutGUID string, GUIDs []string) error {

	filterQuery := bson.D{
		{Key: "shopid", Value: shopID},
		{Key: "guidfixed", Value: bson.M{"$ne": withoutGUID, "$in": GUIDs}},
		// {Key: "guidfixed", Value: bson.M{"$in": GUIDs}},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "imagereferences", Value: bson.D{{Key: "$exists", Value: false}}}},
			bson.D{{Key: "imagereferences", Value: bson.D{{Key: "$size", Value: 0}}}},
		}},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "references", Value: bson.D{{Key: "$exists", Value: false}}}},
			bson.D{{Key: "references", Value: bson.D{{Key: "$size", Value: 0}}}},
		}},
	}

	return repo.pst.Delete(models.DocumentImageGroupDoc{}, filterQuery)
}

func (repo DocumentImageGroupRepository) DeleteByGUIDIsDocumentImageEmpty(shopID string, imageGroupGUID string) error {

	filterQuery := bson.D{
		{Key: "shopid", Value: shopID},
		{Key: "guidfixed", Value: imageGroupGUID},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "imagereferences", Value: bson.D{{Key: "$exists", Value: false}}}},
			bson.D{{Key: "imagereferences", Value: bson.D{{Key: "$size", Value: 0}}}},
		}},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "references", Value: bson.D{{Key: "$exists", Value: false}}}},
			bson.D{{Key: "references", Value: bson.D{{Key: "$size", Value: 0}}}},
		}},
	}

	return repo.pst.Delete(models.DocumentImageGroupDoc{}, filterQuery)
}

func (repo DocumentImageGroupRepository) DeleteByDocumentImageGUIDsHasOneWithoutDocumentImageGroupGUID(shopID string, withoutGUID string, imageGUIDs []string) error {
	return repo.pst.Delete(models.DocumentImageGroupDoc{}, bson.M{
		"shopid":                            shopID,
		"guidfixed":                         bson.M{"$ne": withoutGUID},
		"imagereferences":                   bson.M{"$exists": true, "$size": 1},
		"imagereferences.documentimageguid": bson.M{"$in": imageGUIDs},
		"$or": []interface{}{
			bson.M{"references": bson.M{"$exists": false}},
			bson.M{"references": bson.M{"$size": 0}},
		},
	})
}

func (repo DocumentImageGroupRepository) RemoveDocumentImageByDocumentImageGUIDsWithoutDocumentImageGroupGUID(shopID string, withoutGUID string, imageGUIDs []string) error {

	filterQuery := bson.M{
		"shopid":                            shopID,
		"guidfixed":                         bson.M{"$ne": withoutGUID},
		"imagereferences.documentimageguid": bson.M{"$in": imageGUIDs},
		"imagereferences":                   bson.M{"$exists": true},
		"$or": []interface{}{
			bson.M{"references": bson.M{"$exists": false}},
			bson.M{"references": bson.M{"$size": 0}},
		},
	}

	removeQuery := bson.M{
		"$pull": bson.M{"imagereferences": bson.M{"documentimageguid": bson.M{"$in": imageGUIDs}}},
	}

	return repo.pst.Update(models.DocumentImageGroupDoc{}, filterQuery, removeQuery)
}

func (repo DocumentImageGroupRepository) FindPageImageGroup(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.DocumentImageGroupInfo, mongopagination.PaginationData, error) {

	matchFilterList := []interface{}{}

	for key, value := range filters {
		matchFilterList = append(matchFilterList, bson.M{key: value})
	}

	searchFilterList := []interface{}{}

	for _, colName := range searchInFields {
		searchFilterList = append(searchFilterList, bson.M{colName: bson.M{"$regex": primitive.Regex{
			Pattern: ".*" + pageable.Query + ".*",
			Options: "i",
		}}})
	}

	queryFilters := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
	}

	if len(searchFilterList) > 0 {
		queryFilters["$or"] = searchFilterList
	}

	if len(matchFilterList) > 0 {
		queryFilters["$and"] = matchFilterList
	}

	pageable.Sorts = append(pageable.Sorts, micromodels.KeyInt{Key: "guidfixed", Value: 1})

	matchQuery := bson.M{
		"$match": queryFilters,
	}

	sortFields := bson.D{}
	for _, sortTemp := range pageable.Sorts {
		tempSortVal := 1
		if sortTemp.Value < 1 {
			tempSortVal = -1
		}
		sortFields = append(sortFields, bson.E{Key: sortTemp.Key, Value: tempSortVal})
	}

	sortQuery := bson.M{
		"$sort": sortFields,
	}

	aggData, err := repo.pst.AggregatePage(models.DocumentImageGroupInfo{}, pageable, matchQuery, sortQuery)

	if err != nil {
		return []models.DocumentImageGroupInfo{}, mongopagination.PaginationData{}, err
	}

	docList, err := mogoutil.AggregatePageDecode[models.DocumentImageGroupInfo](aggData)

	if err != nil {
		return []models.DocumentImageGroupInfo{}, mongopagination.PaginationData{}, err
	}

	if err != nil {
		return []models.DocumentImageGroupInfo{}, mongopagination.PaginationData{}, err
	}

	return docList, aggData.Pagination, nil
}
