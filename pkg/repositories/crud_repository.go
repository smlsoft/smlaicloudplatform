package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models/restaurant"

	"go.mongodb.org/mongo-driver/bson"
)

type ICrudRepo interface {
	restaurant.ShopZoneDoc
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

	count, err := repo.pst.Count(&T{}, bson.M{"shopid": shopID})

	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo CrudRepository[T]) Create(doc T) (string, error) {
	idx, err := repo.pst.Create(&T{}, doc)

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

	err := repo.pst.CreateInBatch(&T{}, tempList)

	if err != nil {
		return err
	}
	return nil
}

func (repo CrudRepository[T]) Update(guid string, doc T) error {
	err := repo.pst.UpdateOne(&T{}, "guidfixed", guid, doc)

	if err != nil {
		return err
	}

	return nil
}

func (repo CrudRepository[T]) Delete(shopID string, guid string, username string) error {
	err := repo.pst.SoftDelete(&T{}, username, bson.M{"guidfixed": guid, "shopid": shopID})

	if err != nil {
		return err
	}

	return nil
}

func (repo CrudRepository[T]) FindByGuid(shopID string, guid string) (T, error) {

	doc := &T{}

	err := repo.pst.FindOne(&T{}, bson.M{"guidfixed": guid, "shopid": shopID, "deletedat": bson.M{"$exists": false}}, doc)

	if err != nil {
		return T{}, err
	}

	return *doc, nil
}
