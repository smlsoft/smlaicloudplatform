package repositories

import (
	"smlcloudplatform/internal/microservice"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SearchRepository[T any] struct {
	pst microservice.IPersisterMongo
}

func NewSearchRepository[T any](pst microservice.IPersisterMongo) SearchRepository[T] {
	return SearchRepository[T]{
		pst: pst,
	}
}

func (repo SearchRepository[T]) Find(shopID string, colNameSearch []string, q string) ([]T, error) {

	searchFilterList := []interface{}{}

	for _, colName := range colNameSearch {
		searchFilterList = append(searchFilterList, bson.M{colName: bson.M{"$regex": primitive.Regex{
			Pattern: ".*" + q + ".*",
			Options: "",
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

func (repo SearchRepository[T]) FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]T, mongopagination.PaginationData, error) {

	searchFilterList := []interface{}{}

	for _, colName := range colNameSearch {
		searchFilterList = append(searchFilterList, bson.M{colName: bson.M{"$regex": primitive.Regex{
			Pattern: ".*" + q + ".*",
			Options: "",
		}}})
	}

	docList := []T{}
	pagination, err := repo.pst.FindPage(new(T), limit, page, bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"$or":       searchFilterList,
	}, &docList)

	if err != nil {
		return []T{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo SearchRepository[T]) FindPageSort(shopID string, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]T, mongopagination.PaginationData, error) {

	searchFilterList := []interface{}{}

	for _, colName := range colNameSearch {
		searchFilterList = append(searchFilterList, bson.M{colName: bson.M{"$regex": primitive.Regex{
			Pattern: ".*" + q + ".*",
			Options: "",
		}}})
	}

	filters := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
	}

	if len(searchFilterList) > 0 {
		filters["$or"] = searchFilterList
	}

	docList := []T{}
	pagination, err := repo.pst.FindPageSort(new(T), limit, page, filters, sorts, &docList)

	if err != nil {
		return []T{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo SearchRepository[T]) FindPageFilterSort(shopID string, filters map[string]interface{}, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]T, mongopagination.PaginationData, error) {

	sorts["guidfixed"] = 1

	matchFilterList := []interface{}{}

	for key, value := range filters {
		matchFilterList = append(matchFilterList, bson.M{key: value})
	}

	searchFilterList := []interface{}{}

	for _, colName := range colNameSearch {
		searchFilterList = append(searchFilterList, bson.M{colName: bson.M{"$regex": primitive.Regex{
			Pattern: ".*" + q + ".*",
			Options: "",
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

	docList := []T{}
	pagination, err := repo.pst.FindPageSort(new(T), limit, page, queryFilters, sorts, &docList)

	if err != nil {
		return []T{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}
