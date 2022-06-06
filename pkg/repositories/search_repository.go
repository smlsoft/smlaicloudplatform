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

	docList := []T{}
	pagination, err := repo.pst.FindPageSort(new(T), limit, page, bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"$or":       searchFilterList,
	}, sorts, &docList)

	if err != nil {
		return []T{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}
