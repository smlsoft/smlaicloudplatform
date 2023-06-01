package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/task/models"
	"smlcloudplatform/pkg/utils/mogoutil"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ITaskRepository interface {
	FindTaskChild(shopID string, rejectFromTaskGUID string) (models.TaskChild, error)
	FindLastTaskByCode(shopID string, codeFormat string) (models.TaskDoc, error)
	FindOneTaskByCode(shopID string, taskCode string) (models.TaskInfo, error)
	Count(shopID string) (int, error)
	CountTaskParent(shopID string, taskGUID string) (int, error)
	Create(doc models.TaskDoc) (string, error)
	CreateInBatch(docList []models.TaskDoc) error
	Update(shopID string, guid string, doc models.TaskDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.TaskInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.TaskDoc, error)

	UpdateTotalDocumentImageGroup(shopID string, taskGUID string, total int, countStatus []models.TotalStatus) error
	UpdateTotalRejectDocumentImageGroup(shopID string, taskGUID string, total int) error

	FindPageByTaskReject(shopID string, module string, taskGUID string) ([]models.TaskInfo, error)
	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.TaskItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.TaskDoc, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.TaskInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.TaskDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.TaskActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.TaskDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.TaskActivity, error)

	FindPageTask(shopID string, module string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.TaskInfo, mongopagination.PaginationData, error)
}

type TaskRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.TaskDoc]
	repositories.SearchRepository[models.TaskInfo]
	repositories.GuidRepository[models.TaskItemGuid]
	repositories.ActivityRepository[models.TaskActivity, models.TaskDeleteActivity]
}

func NewTaskRepository(pst microservice.IPersisterMongo) *TaskRepository {

	insRepo := TaskRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.TaskDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.TaskInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.TaskItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.TaskActivity, models.TaskDeleteActivity](pst)

	return &insRepo
}

func (repo *TaskRepository) FindTaskChild(shopID string, rejectFromTaskGUID string) (models.TaskChild, error) {

	queryFilters := bson.M{
		"shopid":             shopID,
		"deletedat":          bson.M{"$exists": false},
		"rejectfromtaskguid": rejectFromTaskGUID,
	}

	opts := &options.FindOneOptions{}

	opts.SetSort(bson.M{"code": 1})

	findDoc := models.TaskChild{}
	err := repo.pst.FindOne(models.TaskChild{}, queryFilters, &findDoc, opts)

	if err != nil {
		return models.TaskChild{}, err
	}

	return findDoc, nil
}

func (repo *TaskRepository) FindLastTaskByCode(shopID string, codeFormat string) (models.TaskDoc, error) {

	queryFilters := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"code": bson.M{"$regex": primitive.Regex{
			Pattern: codeFormat + ".*",
			Options: "i",
		}},
	}

	opts := &options.FindOneOptions{}

	opts.SetSort(bson.M{"code": -1})

	findDoc := new(models.TaskDoc)
	err := repo.pst.FindOne(models.TaskDoc{}, queryFilters, &findDoc, opts)
	// err := repo.pst.FindOne(models.TaskDoc{}, bson.M{"shopid": shopID}, &findDoc)

	if err != nil {
		return models.TaskDoc{}, err
	}

	return *findDoc, nil
}

func (repo *TaskRepository) CountTaskParent(shopID string, taskGUID string) (int, error) {

	queryFilters := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
	}

	queryFilters["parentguidfixed"] = taskGUID

	count, err := repo.pst.Count(models.TaskInfo{}, queryFilters)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *TaskRepository) UpdateTotalDocumentImageGroup(shopID string, taskGUID string, totalDoc int, totalDocStatus []models.TotalStatus) error {

	queryFilters := bson.M{
		"shopid":    shopID,
		"guidfixed": taskGUID,
		"deletedat": bson.M{"$exists": false},
	}

	queryFilters["guidfixed"] = taskGUID

	err := repo.pst.UpdateOne(models.TaskDocumentTotal{}, queryFilters, models.TaskDocumentTotal{TotalDocument: totalDoc, TotalDocumentStatus: &totalDocStatus})

	if err != nil {
		return err
	}

	return nil
}

func (repo *TaskRepository) UpdateTotalRejectDocumentImageGroup(shopID string, taskGUID string, total int) error {

	queryFilters := bson.M{
		"shopid":    shopID,
		"guidfixed": taskGUID,
		"deletedat": bson.M{"$exists": false},
	}

	queryFilters["guidfixed"] = taskGUID

	err := repo.pst.UpdateOne(models.TaskTotalReject{}, queryFilters, models.TaskTotalReject{ToTalReject: total})

	if err != nil {
		return err
	}

	return nil
}

func (repo *TaskRepository) FindPageByTaskReject(shopID string, module string, taskGUID string) ([]models.TaskInfo, error) {

	queryFilters := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
	}

	if len(module) > 0 {
		queryFilters["module"] = module
	}

	queryFilters["parentguidfixed"] = taskGUID

	docList := []models.TaskInfo{}

	err := repo.pst.Find(models.TaskInfo{}, queryFilters, &docList)

	if err != nil {
		return []models.TaskInfo{}, err
	}

	return docList, nil
}

func (repo *TaskRepository) FindOneTaskByCode(shopID string, taskCode string) (models.TaskInfo, error) {

	queryFilters := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"code":      taskCode,
	}

	findDoc := models.TaskInfo{}

	err := repo.pst.FindOne(models.TaskInfo{}, queryFilters, &findDoc)

	if err != nil {
		return models.TaskInfo{}, err
	}

	return findDoc, nil
}

func (repo *TaskRepository) FindPageTask(shopID string, module string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.TaskInfo, mongopagination.PaginationData, error) {

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

	pageable.Sorts = append(pageable.Sorts, micromodels.KeyInt{Key: "guidfixed", Value: 1})

	matchQuery := bson.M{
		"$match": queryFilters,
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

	aggData, err := repo.pst.AggregatePage(models.TaskInfo{}, pageable, matchQuery, sortQuery)

	if err != nil {
		return []models.TaskInfo{}, mongopagination.PaginationData{}, err
	}

	docList, err := mogoutil.AggregatePageDecode[models.TaskInfo](aggData)

	if err != nil {
		return []models.TaskInfo{}, mongopagination.PaginationData{}, err
	}

	if err != nil {
		return []models.TaskInfo{}, mongopagination.PaginationData{}, err
	}

	if docList == nil {
		docList = []models.TaskInfo{}
	}

	return docList, aggData.Pagination, nil
}
