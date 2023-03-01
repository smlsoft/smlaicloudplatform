package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/internal/microservice/models"
	"strings"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SearchRepository[T any] struct {
	pst microservice.IPersisterMongo
}

func NewSearchRepository[T any](pst microservice.IPersisterMongo) SearchRepository[T] {
	return SearchRepository[T]{
		pst: pst,
	}
}

func (repo SearchRepository[T]) Find(shopID string, searchInFields []string, q string) ([]T, error) {

	searchFilterList := []interface{}{}

	for _, colName := range searchInFields {
		searchFilterList = append(searchFilterList, bson.M{colName: bson.M{"$regex": primitive.Regex{
			Pattern: ".*" + q + ".*",
		}}})
	}

	docList := []T{}
	err := repo.pst.Find(new(T), bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"$or":       searchFilterList,
	}, &docList)

	if err != nil {
		return []T{}, err
	}

	return docList, nil
}

func (repo SearchRepository[T]) FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableStep models.PageableStep) ([]T, int, error) {

	searchFilterList := []interface{}{}

	for _, colName := range searchInFields {
		searchFilterList = append(searchFilterList, bson.M{colName: bson.M{"$regex": primitive.Regex{
			Pattern: ".*" + pageableStep.Query + ".*",
			Options: "i",
		}}})
	}

	filterQuery := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"$or":       searchFilterList,
	}

	matchFilterList := []interface{}{}

	for key, value := range filters {
		matchFilterList = append(matchFilterList, bson.M{key: value})
	}

	if len(matchFilterList) > 0 {
		filterQuery["$and"] = matchFilterList
	}

	tempSkip := int64(pageableStep.Skip)
	tempLimit := int64(pageableStep.Limit)

	tempOptions := &options.FindOptions{}
	tempOptions.SetSkip(tempSkip)
	tempOptions.SetLimit(tempLimit)

	projectOptions := bson.M{}

	for key, val := range projects {
		projectOptions[key] = val
	}

	tempOptions.SetProjection(projectOptions)

	for _, pageSort := range pageableStep.Sorts {
		tempOptions.SetSort(bson.M{pageSort.Key: pageSort.Value})
	}

	if len(pageableStep.Sorts) < 1 {
		tempOptions.SetSort(bson.M{"createdat": 1})
	}

	docList := []T{}
	err := repo.pst.Find(new(T), filterQuery, &docList, tempOptions)

	if err != nil {
		return []T{}, 0, err
	}

	count, err := repo.pst.Count(new(T), filterQuery)

	if err != nil {
		return []T{}, 0, err
	}

	return docList, count, nil
}

func (repo SearchRepository[T]) FindPage(shopID string, searchInFields []string, pageable models.Pageable) ([]T, mongopagination.PaginationData, error) {

	searchFilterList := []interface{}{}

	for _, colName := range searchInFields {
		searchFilterList = append(searchFilterList, bson.M{colName: bson.M{"$regex": primitive.Regex{
			Pattern: ".*" + pageable.Query + ".*",
			Options: "i",
		}}})
	}

	filterQuery := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
	}

	if len(searchFilterList) > 0 {
		filterQuery["$or"] = searchFilterList
	}

	if len(pageable.Sorts) < 1 {
		pageable.Sorts = append(pageable.Sorts, models.KeyInt{Key: "createdat", Value: 1})
	}

	docList := []T{}
	pagination, err := repo.pst.FindPage(new(T), filterQuery, pageable, &docList)

	if err != nil {
		return []T{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo SearchRepository[T]) FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable models.Pageable) ([]T, mongopagination.PaginationData, error) {

	matchFilterList := []interface{}{}

	for key, value := range filters {
		matchFilterList = append(matchFilterList, bson.M{key: value})
	}

	searchFilterList := []interface{}{}

	if len(pageable.Query) > 0 {
		for _, colName := range searchInFields {
			searchFilterList = append(searchFilterList, bson.M{colName: bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + pageable.Query + ".*",
				Options: "i",
			}}})
		}
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

	if len(pageable.Sorts) < 1 {
		pageable.Sorts = append(pageable.Sorts, models.KeyInt{Key: "createdat", Value: 1})
	}

	docList := []T{}
	pagination, err := repo.pst.FindPage(new(T), queryFilters, pageable, &docList)

	if err != nil {
		return []T{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (SearchRepository[T]) SearchTextFilter(searchInFields []string, q string) primitive.M {
	prepareText := strings.Trim(q, " ")
	textSearchList := strings.Split(prepareText, " ")

	print(textSearchList)

	colFilter := []interface{}{}
	for _, col := range searchInFields {
		fieldFilter := []interface{}{}
		for _, textSearch := range textSearchList {
			fieldFilter = append(fieldFilter, bson.M{
				col: primitive.Regex{
					Pattern: ".*" + textSearch + ".*",
					Options: "i",
				},
			})
		}

		colFilter = append(colFilter, bson.M{"$or": fieldFilter})

	}

	return bson.M{
		"$and": colFilter,
	}
}
