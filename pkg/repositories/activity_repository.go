package repositories

import (
	"smlcloudplatform/internal/microservice"
	"time"

	micromodels "smlcloudplatform/internal/microservice/models"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IActivityRepository[TCU any, TDEL any] interface {
	FindDeletedPage(shopID string, lastUpdatedDate time.Time, pageable micromodels.Pageable) ([]TDEL, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, pageable micromodels.Pageable) ([]TCU, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, pageableStep micromodels.PageableStep) ([]TDEL, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, pageableStep micromodels.PageableStep) ([]TCU, error)
}
type ActivityRepository[TCU any, TDEL any] struct {
	pst microservice.IPersisterMongo
}

func NewActivityRepository[TCU any, TDEL any](pst microservice.IPersisterMongo) ActivityRepository[TCU, TDEL] {
	return ActivityRepository[TCU, TDEL]{
		pst: pst,
	}
}

func (repo ActivityRepository[TCU, TDEL]) FindDeletedPage(shopID string, lastUpdatedDate time.Time, pageable micromodels.Pageable) ([]TDEL, mongopagination.PaginationData, error) {

	filterQueries := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$gte": lastUpdatedDate},
	}

	docList := []TDEL{}
	pagination, err := repo.pst.FindPage(new(TDEL), filterQueries, pageable, &docList)

	if err != nil {
		return []TDEL{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo ActivityRepository[TCU, TDEL]) FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, pageable micromodels.Pageable) ([]TCU, mongopagination.PaginationData, error) {

	filterQueries := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$not": bson.M{"$gte": lastUpdatedDate}},
		"$or": []interface{}{
			bson.M{"createdat": bson.M{"$gte": lastUpdatedDate}},
			bson.M{"updatedat": bson.M{"$gte": lastUpdatedDate}},
		},
	}

	docList := []TCU{}
	pagination, err := repo.pst.FindPage(new(TCU), filterQueries, pageable, &docList)

	if err != nil {
		return []TCU{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo ActivityRepository[TCU, TDEL]) FindDeletedStep(shopID string, lastUpdatedDate time.Time, pageableStep micromodels.PageableStep) ([]TDEL, error) {

	docList := []TDEL{}

	filterQuery := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$gte": lastUpdatedDate},
	}

	tempOptions := &options.FindOptions{}
	tempOptions.SetSkip(int64(pageableStep.Skip))
	tempOptions.SetLimit(int64(pageableStep.Limit))

	err := repo.pst.Find(new(TDEL), filterQuery, &docList, tempOptions)

	if err != nil {
		return []TDEL{}, err
	}

	return docList, nil
}

func (repo ActivityRepository[TCU, TDEL]) FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, pageableStep micromodels.PageableStep) ([]TCU, error) {

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
	tempOptions.SetSkip(int64(pageableStep.Skip))
	tempOptions.SetLimit(int64(pageableStep.Limit))

	err := repo.pst.Find(new(TCU), filterQuery, &docList, tempOptions)

	if err != nil {
		return []TCU{}, err
	}

	return docList, nil
}
