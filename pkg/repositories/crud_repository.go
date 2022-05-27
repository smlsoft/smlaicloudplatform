package repositories

import (
	"smlcloudplatform/internal/microservice"

	"go.mongodb.org/mongo-driver/bson"
)

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

func (repo CrudRepository[T]) FindOne(shopID string, filters map[string]interface{}) (T, error) {

	doc := new(T)

	filterQuery := bson.M{}

	for col, val := range filters {
		filterQuery[col] = val
	}

	filterQuery["shopid"] = shopID
	filterQuery["deletedat"] = bson.M{"$exists": false}

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

func (repo CrudRepository[T]) FindByDocIndentiryGuid(shopID string, indentityField string, guid string) (T, error) {

	doc := new(T)

	err := repo.pst.FindOne(new(T), bson.M{indentityField: guid, "shopid": shopID, "deletedat": bson.M{"$exists": false}}, doc)

	if err != nil {
		return *new(T), err
	}

	return *doc, nil
}
