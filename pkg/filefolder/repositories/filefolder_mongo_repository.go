package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/filefolder/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/utils/mogoutil"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IFileFolderRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.FileFolderDoc) (string, error)
	CreateInBatch(docList []models.FileFolderDoc) error
	Update(shopID string, guid string, doc models.FileFolderDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.FileFolderInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.FileFolderDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.FileFolderItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.FileFolderDoc, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.FileFolderInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, pageable micromodels.Pageable) ([]models.FileFolderDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, pageable micromodels.Pageable) ([]models.FileFolderActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, pageableStep micromodels.PageableStep) ([]models.FileFolderDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, pageableStep micromodels.PageableStep) ([]models.FileFolderActivity, error)

	FindPageFileFolder(shopID string, module string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.FileFolderInfo, mongopagination.PaginationData, error)
}

type FileFolderRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.FileFolderDoc]
	repositories.SearchRepository[models.FileFolderInfo]
	repositories.GuidRepository[models.FileFolderItemGuid]
	repositories.ActivityRepository[models.FileFolderActivity, models.FileFolderDeleteActivity]
}

func NewFileFolderRepository(pst microservice.IPersisterMongo) *FileFolderRepository {

	insRepo := &FileFolderRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.FileFolderDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.FileFolderInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.FileFolderItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.FileFolderActivity, models.FileFolderDeleteActivity](pst)

	return insRepo
}

func (repo *FileFolderRepository) FindPageFileFolder(shopID string, module string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.FileFolderInfo, mongopagination.PaginationData, error) {

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
			"as": "tempDocs",
		},
	}

	addFielsQuery := bson.M{
		"$addFields": bson.M{
			"total": bson.M{"$size": "$tempDocs"},
		},
	}

	projectQuery := bson.M{
		"$project": bson.M{
			"tempDocs": 0,
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

	aggData, err := repo.pst.AggregatePage(models.FileFolderInfo{}, pageable, matchQuery, lookupQuery, addFielsQuery, projectQuery, sortQuery)

	if err != nil {
		return []models.FileFolderInfo{}, mongopagination.PaginationData{}, err
	}

	docList, err := mogoutil.AggregatePageDecode[models.FileFolderInfo](aggData)

	if err != nil {
		return []models.FileFolderInfo{}, mongopagination.PaginationData{}, err
	}

	if err != nil {
		return []models.FileFolderInfo{}, mongopagination.PaginationData{}, err
	}

	return docList, aggData.Pagination, nil
}
