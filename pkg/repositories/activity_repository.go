package repositories

import (
	"smlcloudplatform/internal/microservice"
	"time"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IActivityRepository[TCU any, TDEL any] interface {
	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]TDEL, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]TCU, mongopagination.PaginationData, error)
	FindDeletedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]TDEL, error)
	FindCreatedOrUpdatedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]TCU, error)
}
type ActivityRepository[TCU any, TDEL any] struct {
	pst microservice.IPersisterMongo
}

func NewActivityRepository[TCU any, TDEL any](pst microservice.IPersisterMongo) ActivityRepository[TCU, TDEL] {
	return ActivityRepository[TCU, TDEL]{
		pst: pst,
	}
}

func (repo ActivityRepository[TCU, TDEL]) FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]TDEL, mongopagination.PaginationData, error) {

	docList := []TDEL{}
	pagination, err := repo.pst.FindPage(new(TDEL), limit, page, bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$gte": lastUpdatedDate},
	}, &docList)

	if err != nil {
		return []TDEL{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo ActivityRepository[TCU, TDEL]) FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]TCU, mongopagination.PaginationData, error) {

	docList := []TCU{}
	pagination, err := repo.pst.FindPage(new(TCU), limit, page, bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$not": bson.M{"$gte": lastUpdatedDate}},
		"$or": []interface{}{
			bson.M{"createdat": bson.M{"$gte": lastUpdatedDate}},
			bson.M{"updatedat": bson.M{"$gte": lastUpdatedDate}},
		},
	}, &docList)

	if err != nil {
		return []TCU{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo ActivityRepository[TCU, TDEL]) FindDeletedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]TDEL, error) {

	docList := []TDEL{}

	filterQuery := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$gte": lastUpdatedDate},
	}

	tempOptions := &options.FindOptions{}
	tempOptions.SetSkip(int64(skip))
	tempOptions.SetLimit(int64(limit))

	err := repo.pst.Find(new(TDEL), filterQuery, &docList, tempOptions)

	if err != nil {
		return []TDEL{}, err
	}

	return docList, nil
}

func (repo ActivityRepository[TCU, TDEL]) FindCreatedOrUpdatedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]TCU, error) {

	docList := []TCU{}
	filterQuery := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$not": bson.M{"$gte": lastUpdatedDate}},
		"$or": []interface{}{
			bson.M{"createdat": bson.M{"$gte": lastUpdatedDate}},
			bson.M{"updatedat": bson.M{"$gte": lastUpdatedDate}},
		},
	}

	tempOptions := &options.FindOptions{}
	tempOptions.SetSkip(int64(skip))
	tempOptions.SetLimit(int64(limit))

	err := repo.pst.Find(new(TCU), filterQuery, &docList, tempOptions)

	if err != nil {
		return []TCU{}, err
	}

	return docList, nil
}
