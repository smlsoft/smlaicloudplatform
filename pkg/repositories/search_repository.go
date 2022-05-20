package repositories

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models/restaurant"
	"smlcloudplatform/pkg/models/vfgl"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ISearchRepo interface {
	restaurant.ShopZoneInfo | restaurant.ShopTableInfo | restaurant.PrinterTerminalInfo | restaurant.KitchenInfo | vfgl.JournalInfo
}

type SearchRepository[T ISearchRepo] struct {
	pst microservice.IPersisterMongo
}

func NewSearchRepository[T ISearchRepo](pst microservice.IPersisterMongo) SearchRepository[T] {
	return SearchRepository[T]{
		pst: pst,
	}
}

func (repo SearchRepository[T]) FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]T, mongopagination.PaginationData, error) {

	searchFilterList := []interface{}{
		bson.M{"guidfixed": q},
	}

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
