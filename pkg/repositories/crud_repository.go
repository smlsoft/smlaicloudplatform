package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models/restaurant"
	"smlcloudplatform/pkg/models/vfgl"

	"go.mongodb.org/mongo-driver/bson"
)

type ICrudRepo interface {
	restaurant.ShopZoneDoc | restaurant.ShopTableDoc | restaurant.PrinterTerminalDoc | restaurant.KitchenDoc | vfgl.JournalDoc
}
type CrudRepository[T ICrudRepo] struct {
	pst microservice.IPersisterMongo
}

func NewCrudRepository[T ICrudRepo](pst microservice.IPersisterMongo) CrudRepository[T] {
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

func (repo CrudRepository[T]) Update(guid string, doc T) error {
	err := repo.pst.UpdateOne(new(T), "guidfixed", guid, doc)

	if err != nil {
		return err
	}

	return nil
}

func (repo CrudRepository[T]) Delete(shopID string, guid string, username string) error {
	err := repo.pst.SoftDelete(new(T), username, bson.M{"guidfixed": guid, "shopid": shopID})

	if err != nil {
		return err
	}

	return nil
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
