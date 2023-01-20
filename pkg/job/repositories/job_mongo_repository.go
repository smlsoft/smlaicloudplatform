package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/job/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/utils/mogoutil"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IJobRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.JobDoc) (string, error)
	CreateInBatch(docList []models.JobDoc) error
	Update(shopID string, guid string, doc models.JobDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.JobInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.JobDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.JobItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.JobDoc, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.JobInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, pageable micromodels.Pageable) ([]models.JobDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, pageable micromodels.Pageable) ([]models.JobActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, pageableStep micromodels.PageableStep) ([]models.JobDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, pageableStep micromodels.PageableStep) ([]models.JobActivity, error)

	FindPageJob(shopID string, module string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.JobInfo, mongopagination.PaginationData, error)
}

type JobRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.JobDoc]
	repositories.SearchRepository[models.JobInfo]
	repositories.GuidRepository[models.JobItemGuid]
	repositories.ActivityRepository[models.JobActivity, models.JobDeleteActivity]
}

func NewJobRepository(pst microservice.IPersisterMongo) *JobRepository {

	insRepo := &JobRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.JobDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.JobInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.JobItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.JobActivity, models.JobDeleteActivity](pst)

	return insRepo
}

func (repo *JobRepository) FindPageJob(shopID string, module string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.JobInfo, mongopagination.PaginationData, error) {

	matchFilterList := []interface{}{}

	for key, value := range filters {
		matchFilterList = append(matchFilterList, bson.M{key: value})
	}

	searchFilterList := []interface{}{}

	for _, colName := range searchInFields {
		searchFilterList = append(searchFilterList, bson.M{colName: bson.M{"$regex": primitive.Regex{
			Pattern: ".*" + pageable.Query + ".*",
			Options: "",
		}}})
	}

	queryFilters := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
	}

	if len(module) > 0 {
		queryFilters["module"] = module
	}

	if len(searchFilterList) > 0 {
		queryFilters["$or"] = searchFilterList
	}

	if len(matchFilterList) > 0 {
		queryFilters["$and"] = matchFilterList
	}

	if len(pageable.Sorts) > 0 {
		pageable.Sorts = append(pageable.Sorts, micromodels.KeyInt{Key: "guidfixed", Value: 1})
	}

	matchQuery := bson.M{
		"$match": queryFilters,
	}

	lookupQuery := bson.M{
		"$lookup": bson.M{
			"from": "documentImageGroups",
			"let":  bson.M{"shopid": "$shopid", "foreignField": "$foreignField"},
			"pipeline": []interface{}{
				bson.M{
					"$match": bson.M{
						"$expr": bson.M{"$and": []interface{}{
							bson.M{
								"$eq": []interface{}{"$shopid", "$$shopid"},
							},
							bson.M{
								"$eq": []interface{}{"$shopid", "$$shopid"},
							},
						}},
					},
				},
			},
			"as": "tempTotal",
		},
	}

	addFielsQuery := bson.M{
		"$addFields": bson.M{
			"total": bson.M{"$size": "$tempTotal"},
		},
	}

	projectQuery := bson.M{
		"$project": bson.M{
			"tempTotal": 0,
		},
	}

	sortFields := bson.D{}
	for _, sortTemp := range pageable.Sorts {
		tempSortVal := 1
		if sortTemp.Value < 1 {
			tempSortVal = -1
		}
		sortFields = append(sortFields, bson.E{Key: sortTemp.Key, Value: tempSortVal})
	}

	sortQuery := bson.M{
		"$sort": sortFields,
	}

	aggData, err := repo.pst.AggregatePage(models.JobInfo{}, pageable, matchQuery, lookupQuery, addFielsQuery, projectQuery, sortQuery)

	if err != nil {
		return []models.JobInfo{}, mongopagination.PaginationData{}, err
	}

	docList, err := mogoutil.AggregatePageDecode[models.JobInfo](aggData)

	if err != nil {
		return []models.JobInfo{}, mongopagination.PaginationData{}, err
	}

	if err != nil {
		return []models.JobInfo{}, mongopagination.PaginationData{}, err
	}

	return docList, aggData.Pagination, nil
}
