package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/internal/microservice/models"
	"strings"

	"github.com/userplant/mongopagination"
	m "github.com/veer66/mapkha"
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

	filterQuery := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
	}

	searchFilterQuery := repo.CreateTextFilter(searchInFields, q)

	if len(searchFilterQuery) > 0 {
		if filterQuery["$or"] == nil {
			filterQuery["$or"] = searchFilterQuery
		} else {
			filterQuery["$or"] = append(filterQuery["$or"].([]interface{}), searchFilterQuery...)
		}
	}

	docList := []T{}
	err := repo.pst.Find(new(T), filterQuery, &docList)

	if err != nil {
		return []T{}, err
	}

	return docList, nil
}

func (repo SearchRepository[T]) FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableStep models.PageableStep) ([]T, int, error) {

	filterQuery := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
	}

	matchFilterList := []interface{}{}

	for key, value := range filters {
		matchFilterList = append(matchFilterList, bson.M{key: value})
	}

	if len(matchFilterList) > 0 {
		filterQuery["$and"] = matchFilterList
	}

	searchFilterQuery := repo.CreateTextFilter(searchInFields, pageableStep.Query)

	if len(searchFilterQuery) > 0 {
		if filterQuery["$or"] == nil {
			filterQuery["$or"] = searchFilterQuery
		} else {
			filterQuery["$or"] = append(filterQuery["$or"].([]interface{}), searchFilterQuery...)
		}
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

	filterQuery := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
	}

	searchFilterQuery := repo.CreateTextFilter(searchInFields, pageable.Query)

	if len(searchFilterQuery) > 0 {
		if filterQuery["$or"] == nil {
			filterQuery["$or"] = searchFilterQuery
		} else {
			filterQuery["$or"] = append(filterQuery["$or"].([]interface{}), searchFilterQuery...)
		}
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

	searchFilterQuery := repo.CreateTextFilter(searchInFields, pageable.Query)

	// matchFilterList = append(matchFilterList, searchFilterQuery...)

	queryFilters := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
	}

	if len(matchFilterList) > 0 {
		queryFilters["$and"] = matchFilterList
	}

	if len(pageable.Sorts) < 1 {
		pageable.Sorts = append(pageable.Sorts, models.KeyInt{Key: "createdat", Value: 1})
	}

	if len(searchFilterQuery) > 0 {
		if queryFilters["$or"] == nil {
			queryFilters["$or"] = searchFilterQuery
		} else {
			queryFilters["$or"] = append(queryFilters["$or"].([]interface{}), searchFilterQuery...)
		}
	}

	docList := []T{}
	pagination, err := repo.pst.FindPage(new(T), queryFilters, pageable, &docList)

	if err != nil {
		return []T{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}

var wordCut *m.Wordcut

func (repo SearchRepository[T]) getTokenizer() (*m.Wordcut, error) {

	if wordCut == nil {
		dict, err := m.LoadDefaultDict()
		if err != nil {
			return nil, err
		}

		wordCut = m.NewWordcut(dict)
	}

	return wordCut, nil
}

func (repo SearchRepository[T]) CreateTextFilter(searchFields []string, query string) []interface{} {
	searchTerms := repo.extractSearchTerms(query)
	fieldFilters := repo.generateFieldFilters(searchFields, searchTerms)

	return fieldFilters
}

func (repo SearchRepository[T]) extractSearchTerms(query string) []string {
	trimmedQuery := strings.Trim(query, " ")
	splitBySpace := strings.Split(trimmedQuery, " ")

	searchTerms := []string{}
	tokenizer, err := repo.getTokenizer()

	if err != nil {
		return searchTerms
	}

	for _, term := range splitBySpace {
		if len(term) > 0 {
			searchTerms = append(searchTerms, tokenizer.Segment(term)...)
		}
	}

	return searchTerms
}

func (repo SearchRepository[T]) generateFieldFilters(searchFields []string, searchTerms []string) []interface{} {
	fieldFilters := []interface{}{}

	for _, field := range searchFields {
		termFilters := []interface{}{}

		for _, searchTerm := range searchTerms {
			termFilters = append(termFilters, bson.M{
				field: primitive.Regex{
					Pattern: searchTerm,
				},
			})
		}

		if len(termFilters) > 0 {
			fieldFilters = append(fieldFilters, bson.M{"$and": termFilters})
		}

	}

	return fieldFilters
}
