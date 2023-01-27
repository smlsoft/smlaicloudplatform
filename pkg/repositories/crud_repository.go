package repositories

import (
	"errors"
	"smlcloudplatform/internal/microservice"

	"go.mongodb.org/mongo-driver/bson"
)

type ICRUDRepository[T any] interface {
	Count(shopID string) (int, error)
	Create(doc T) (string, error)
	CreateInBatch(docList []T) error
	Update(shopID string, guid string, doc T) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	FindOne(shopID string, filters interface{}) (T, error)
	FindByGuid(shopID string, guid string) (T, error)
	FindByGuids(shopID string, guids []string) ([]T, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (T, error)
}

type CrudRepository[T any] struct {
	pst microservice.IPersisterMongo
}

func NewCrudRepository[T any](pst microservice.IPersisterMongo) CrudRepository[T] {
	return CrudRepository[T]{
		pst: pst,
	}
}

func (repo CrudRepository[T]) Count(shopID string) (int, error) {

	count, err := repo.pst.Count(new(T), bson.M{"shopid": shopID})

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo CrudRepository[T]) Create(doc T) (string, error) {
	idx, err := repo.pst.Create(new(T), doc)

	if err != nil {
		return "", err
	}

	return idx.Hex(), nil
}

func (repo CrudRepository[T]) CreateInBatch(docList []T) error {
	var tempList []interface{}

	for _, inv := range docList {
		tempList = append(tempList, inv)
	}

	err := repo.pst.CreateInBatch(new(T), tempList)

	if err != nil {
		return err
	}
	return nil
}

func (repo CrudRepository[T]) Update(shopID string, guid string, doc T) error {
	filterDoc := map[string]interface{}{
		"shopid":    shopID,
		"guidfixed": guid,
	}

	err := repo.pst.UpdateOne(new(T), filterDoc, doc)

	if err != nil {
		return err
	}

	return nil
}

func (repo CrudRepository[T]) Delete(shopID string, username string, filters map[string]interface{}) error {

	filterQuery := bson.M{}

	for col, val := range filters {
		filterQuery[col] = val
	}

	filterQuery["shopid"] = shopID

	err := repo.pst.SoftDelete(new(T), username, filterQuery)

	if err != nil {
		return err
	}

	return nil
}

func (repo CrudRepository[T]) DeleteByGuidfixed(shopID string, guid string, username string) error {
	err := repo.pst.SoftDelete(new(T), username, bson.M{"guidfixed": guid, "shopid": shopID})

	if err != nil {
		return err
	}

	return nil
}

func (repo CrudRepository[T]) FindOne(shopID string, filters interface{}) (T, error) {

	var filterQuery interface{}

	switch filters.(type) {
	case bson.M:
		tempFilterQuery := filters.(bson.M)
		tempFilterQuery["shopid"] = shopID
		tempFilterQuery["deletedat"] = bson.M{"$exists": false}
		filterQuery = tempFilterQuery
	case bson.D:
		tempFilterQuery := filters.(bson.D)
		tempFilterQuery = append(tempFilterQuery, bson.E{"shopid", shopID})
		tempFilterQuery = append(tempFilterQuery, bson.E{"deletedat", bson.D{{"$exists", false}}})

		filterQuery = tempFilterQuery
	default:
		return *new(T), errors.New("invalid query filter type")
	}

	doc := new(T)

	err := repo.pst.FindOne(new(T), filterQuery, doc)

	if err != nil {
		return *new(T), err
	}

	return *doc, nil
}

func (repo CrudRepository[T]) FindByGuid(shopID string, guid string) (T, error) {

	doc := new(T)

	err := repo.pst.FindOne(new(T), bson.M{"guidfixed": guid, "shopid": shopID, "deletedat": bson.M{"$exists": false}}, doc)

	if err != nil {
		return *new(T), err
	}

	return *doc, nil
}

func (repo CrudRepository[T]) FindByGuids(shopID string, guids []string) ([]T, error) {

	doc := new([]T)

	err := repo.pst.Find(new(T), bson.M{"guidfixed": bson.M{"$in": guids}, "shopid": shopID, "deletedat": bson.M{"$exists": false}}, doc)

	if err != nil {
		return *new([]T), err
	}

	return *doc, nil
}

func (repo CrudRepository[T]) FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (T, error) {

	doc := new(T)

	err := repo.pst.FindOne(new(T), bson.M{"shopid": shopID, "deletedat": bson.M{"$exists": false}, indentityField: indentityValue}, doc)

	if err != nil {
		return *new(T), err
	}

	return *doc, nil
}
