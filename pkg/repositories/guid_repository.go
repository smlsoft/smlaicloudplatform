package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"smlcloudplatform/pkg/models/restaurant"

	"go.mongodb.org/mongo-driver/bson"
)

type IGuidRepo interface {
	models.CategoryItemGuid | restaurant.ShopZoneItemGuid | restaurant.ShopTableItemGuid | restaurant.PrinterTerminalItemGuid | restaurant.KitchenItemGuid
}

type GuidRepository[T IGuidRepo] struct {
	pst microservice.IPersisterMongo
}

func NewGuidRepository[T IGuidRepo](pst microservice.IPersisterMongo) GuidRepository[T] {
	return GuidRepository[T]{
		pst: pst,
	}
}

func (repo GuidRepository[T]) FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]T, error) {

	findDoc := []T{}
	err := repo.pst.Find(new(T), bson.M{"shopid": shopID, columnName: bson.M{"$in": itemGuidList}, "deletedat": bson.M{"$exists": false}}, &findDoc)

	if err != nil {
		return []T{}, err
	}
	return findDoc, nil
}
