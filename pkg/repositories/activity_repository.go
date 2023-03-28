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
	// InitialActivityRepository(pst microservice.IPersisterMongo)
	FindDeletedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]TDEL, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]TCU, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]TDEL, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]TCU, error)
}
type ActivityRepository[TCU any, TDEL any] struct {
	pst microservice.IPersisterMongo
}

func NewActivityRepository[TCU any, TDEL any](pst microservice.IPersisterMongo) ActivityRepository[TCU, TDEL] {
	return ActivityRepository[TCU, TDEL]{
		pst: pst,
	}
}

func (repo *ActivityRepository[TCU, TDEL]) InitialActivityRepository(pst microservice.IPersisterMongo) {
	repo.pst = pst
}

func (repo ActivityRepository[TCU, TDEL]) FindDeletedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]TDEL, mongopagination.PaginationData, error) {

	filterQueries := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$gte": lastUpdatedDate},
	}

	extraFilterQueries := repo.generateExtraFilters(extraFilters)
	if len(extraFilterQueries) > 0 {
		filterQueries["$and"] = extraFilterQueries
	}

	docList := []TDEL{}
	pagination, err := repo.pst.FindPage(new(TDEL), filterQueries, pageable, &docList)

	if err != nil {
		return []TDEL{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo ActivityRepository[TCU, TDEL]) FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]TCU, mongopagination.PaginationData, error) {

	filterQueries := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$not": bson.M{"$gte": lastUpdatedDate}},
		"$or": []interface{}{
			bson.M{"createdat": bson.M{"$gte": lastUpdatedDate}},
			bson.M{"updatedat": bson.M{"$gte": lastUpdatedDate}},
		},
	}

	extraFilterQueries := repo.generateExtraFilters(extraFilters)
	if len(extraFilterQueries) > 0 {
		filterQueries["$and"] = extraFilterQueries
	}

	docList := []TCU{}
	pagination, err := repo.pst.FindPage(new(TCU), filterQueries, pageable, &docList)

	if err != nil {
		return []TCU{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo ActivityRepository[TCU, TDEL]) FindDeletedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]TDEL, error) {

	docList := []TDEL{}

	filterQueries := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$gte": lastUpdatedDate},
	}

	extraFilterQueries := repo.generateExtraFilters(extraFilters)
	if len(extraFilterQueries) > 0 {
		filterQueries["$and"] = extraFilterQueries
	}

	tempOptions := &options.FindOptions{}
	tempOptions.SetSkip(int64(pageableStep.Skip))
	tempOptions.SetLimit(int64(pageableStep.Limit))

	err := repo.pst.Find(new(TDEL), filterQueries, &docList, tempOptions)

	if err != nil {
		return []TDEL{}, err
	}

	return docList, nil
}

func (repo ActivityRepository[TCU, TDEL]) FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]TCU, error) {

	docList := []TCU{}
	filterQueries := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$not": bson.M{"$gte": lastUpdatedDate}},
		"$or": []interface{}{
			bson.M{"createdat": bson.M{"$gte": lastUpdatedDate}},
			bson.M{"updatedat": bson.M{"$gte": lastUpdatedDate}},
		},
	}

	extraFilterQueries := repo.generateExtraFilters(extraFilters)
	if len(extraFilterQueries) > 0 {
		filterQueries["$and"] = extraFilterQueries
	}

	tempOptions := &options.FindOptions{}
	tempOptions.SetSkip(int64(pageableStep.Skip))
	tempOptions.SetLimit(int64(pageableStep.Limit))

	err := repo.pst.Find(new(TCU), filterQueries, &docList, tempOptions)

	if err != nil {
		return []TCU{}, err
	}

	return docList, nil
}

func (repo ActivityRepository[TCU, TDEL]) generateExtraFilters(extraFilters map[string]interface{}) []interface{} {
	if len(extraFilters) > 0 {
		tempExtraFilters := []interface{}{}
		for key, value := range extraFilters {
			tempExtraFilters = append(tempExtraFilters, bson.M{key: value})
		}

		return tempExtraFilters
	}

	return []interface{}{}
}
