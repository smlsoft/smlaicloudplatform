package repositories

import (
	"smlcloudplatform/internal/microservice"
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
	FindPage(shopID string, colNameSearch []string, q string, page int, limit int) ([]models.FileFolderInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.FileFolderDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.FileFolderItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.FileFolderDoc, error)
	FindPageSort(shopID string, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]models.FileFolderInfo, mongopagination.PaginationData, error)
	FindLimit(shopID string, filters map[string]interface{}, colNameSearch []string, q string, skip int, limit int, sorts map[string]int, projects map[string]interface{}) ([]models.FileFolderInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.FileFolderDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.FileFolderActivity, mongopagination.PaginationData, error)
	FindDeletedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.FileFolderDeleteActivity, error)
	FindCreatedOrUpdatedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.FileFolderActivity, error)

	FindPageFileFolder(shopID string, module string, filters map[string]interface{}, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]models.FileFolderInfo, mongopagination.PaginationData, error)
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

func (repo *FileFolderRepository) FindPageFileFolder(shopID string, module string, filters map[string]interface{}, colNameSearch []string, q string, page int, limit int, sorts map[string]int) ([]models.FileFolderInfo, mongopagination.PaginationData, error) {

	matchFilterList := []interface{}{}

	for key, value := range filters {
		matchFilterList = append(matchFilterList, bson.M{key: value})
	}

	searchFilterList := []interface{}{}

	for _, colName := range colNameSearch {
		searchFilterList = append(searchFilterList, bson.M{colName: bson.M{"$regex": primitive.Regex{
			Pattern: ".*" + q + ".*",
			Options: "",
		}}})
	}

	queryFilters := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
	}

	if len(module) > 0 {
		queryFilters["module"] = module
	} else {
		// searchFilterList = append(searchFilterList, bson.M{"module": ""})
		// searchFilterList = append(searchFilterList, bson.M{"module": bson.M{"$exists": false}})
	}

	if len(searchFilterList) > 0 {
		queryFilters["$or"] = searchFilterList
	}

	if len(matchFilterList) > 0 {
		queryFilters["$and"] = matchFilterList
	}

	if len(sorts) > 0 {
		sorts["guidfixed"] = 1
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

	sortTemp := bson.M{}
	for sortKey, sortVal := range sorts {
		tempSortVal := 1
		if sorts[sortKey] < sortVal {
			tempSortVal = -1
		}
		sortTemp[sortKey] = tempSortVal
	}

	sortTemp["guidfixed"] = 1

	sortQuery := bson.M{
		"$sort": sortTemp,
	}

	aggData, err := repo.pst.AggregatePage(models.FileFolderInfo{}, limit, page, matchQuery, lookupQuery, addFielsQuery, projectQuery, sortQuery)

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
