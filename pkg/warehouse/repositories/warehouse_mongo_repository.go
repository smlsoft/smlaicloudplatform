package repositories

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/utils/mogoutil"
	"smlcloudplatform/pkg/warehouse/models"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
)

type IWarehouseRepository interface {
	Count(shopID string) (int, error)
	Create(doc models.WarehouseDoc) (string, error)
	CreateInBatch(docList []models.WarehouseDoc) error
	Update(shopID string, guid string, doc models.WarehouseDoc) error
	DeleteByGuidfixed(shopID string, guid string, username string) error
	Delete(shopID string, username string, filters map[string]interface{}) error
	FindPage(shopID string, searchInFields []string, pageable micromodels.Pageable) ([]models.WarehouseInfo, mongopagination.PaginationData, error)
	FindByGuid(shopID string, guid string) (models.WarehouseDoc, error)

	FindInItemGuid(shopID string, columnName string, itemGuidList []string) ([]models.WarehouseItemGuid, error)
	FindByDocIndentityGuid(shopID string, indentityField string, indentityValue interface{}) (models.WarehouseDoc, error)
	FindPageFilter(shopID string, filters map[string]interface{}, searchInFields []string, pageable micromodels.Pageable) ([]models.WarehouseInfo, mongopagination.PaginationData, error)
	FindStep(shopID string, filters map[string]interface{}, searchInFields []string, projects map[string]interface{}, pageableLimit micromodels.PageableStep) ([]models.WarehouseInfo, int, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.WarehouseDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.WarehouseActivity, mongopagination.PaginationData, error)
	FindDeletedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.WarehouseDeleteActivity, error)
	FindCreatedOrUpdatedStep(shopID string, lastUpdatedDate time.Time, filters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.WarehouseActivity, error)

	FindLocationPage(shopID string, pageable micromodels.Pageable) ([]models.LocationInfo, mongopagination.PaginationData, error)
	FindShelfPage(shopID string, pageable micromodels.Pageable) ([]models.ShelfInfo, mongopagination.PaginationData, error)
	FindWarehouseByLocation(shopID, warehouseCode, locationCode string) (models.WarehouseDoc, error)
	FindWarehouseByShelf(shopID, warehouseCode, locationCode, shelfCode string) (models.WarehouseDoc, error)

	Transaction(queryFunc func() error) error
}

type WarehouseRepository struct {
	pst microservice.IPersisterMongo
	repositories.CrudRepository[models.WarehouseDoc]
	repositories.SearchRepository[models.WarehouseInfo]
	repositories.GuidRepository[models.WarehouseItemGuid]
	repositories.ActivityRepository[models.WarehouseActivity, models.WarehouseDeleteActivity]
}

func NewWarehouseRepository(pst microservice.IPersisterMongo) *WarehouseRepository {

	insRepo := &WarehouseRepository{
		pst: pst,
	}

	insRepo.CrudRepository = repositories.NewCrudRepository[models.WarehouseDoc](pst)
	insRepo.SearchRepository = repositories.NewSearchRepository[models.WarehouseInfo](pst)
	insRepo.GuidRepository = repositories.NewGuidRepository[models.WarehouseItemGuid](pst)
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.WarehouseActivity, models.WarehouseDeleteActivity](pst)

	return insRepo
}

func (repo WarehouseRepository) Transaction(queryFunc func() error) error {
	return repo.pst.Transaction(queryFunc)
}

func (repo WarehouseRepository) FindLocationPage(shopID string, pageable micromodels.Pageable) ([]models.LocationInfo, mongopagination.PaginationData, error) {

	criteria := []interface{}{}

	mainQuery := bson.M{
		"$match": bson.M{
			"shopid":    shopID,
			"deletedat": bson.M{"$exists": false},
		},
	}
	criteria = append(criteria, mainQuery)

	searchFilterQuery := repo.CreateTextFilter([]string{"location.code", "location.code.names.name"}, pageable.Query)

	if len(searchFilterQuery) > 0 {
		searchQuery := bson.M{"$match": bson.M{"$or": searchFilterQuery}}
		criteria = append(criteria, searchQuery)
	}

	unwindQuery := bson.M{"$unwind": "$location"}
	criteria = append(criteria, unwindQuery)

	projectQuery := bson.M{"$project": bson.M{
		"guidfixed":      "$guidfixed",
		"warehousecode":  "$code",
		"warehousenames": "$names",
		"locationcode":   "$location.code",
		"locationnames":  "$location.names",
		"shelf":          1,
	}}
	criteria = append(criteria, projectQuery)

	aggData, err := repo.pst.AggregatePage(models.LocationInfo{}, pageable, criteria...)

	if err != nil {
		return []models.LocationInfo{}, mongopagination.PaginationData{}, err
	}

	docList, err := mogoutil.AggregatePageDecode[models.LocationInfo](aggData)

	if err != nil {
		return []models.LocationInfo{}, mongopagination.PaginationData{}, err
	}

	return docList, aggData.Pagination, nil
}

func (repo WarehouseRepository) FindShelfPage(shopID string, pageable micromodels.Pageable) ([]models.ShelfInfo, mongopagination.PaginationData, error) {

	criteria := []interface{}{}

	mainQuery := bson.M{
		"$match": bson.M{
			"shopid":    shopID,
			"deletedat": bson.M{"$exists": false},
		},
	}
	criteria = append(criteria, mainQuery)

	searchFilterQuery := repo.CreateTextFilter([]string{"location.shelf.code", "location.shelf.name"}, pageable.Query)

	if len(searchFilterQuery) > 0 {
		searchQuery := bson.M{"$match": bson.M{"$or": searchFilterQuery}}
		criteria = append(criteria, searchQuery)
	}

	unwindQueryLevel1 := bson.M{"$unwind": "$location"}
	criteria = append(criteria, unwindQueryLevel1)

	unwindQueryLevel2 := bson.M{"$unwind": "$location.shelf"}
	criteria = append(criteria, unwindQueryLevel2)

	projectQuery := bson.M{"$project": bson.M{
		"guidfixed":      "$guidfixed",
		"warehousecode":  "$code",
		"warehousenames": "$names",
		"locationcode":   "$location.code",
		"locationnames":  "$location.names",
		"shelfcode":      "$location.shelf.code",
		"shelfname":      "$location.shelf.name",
	}}

	criteria = append(criteria, projectQuery)

	aggData, err := repo.pst.AggregatePage(models.ShelfInfo{}, pageable, criteria...)

	if err != nil {
		return []models.ShelfInfo{}, mongopagination.PaginationData{}, err
	}

	docList, err := mogoutil.AggregatePageDecode[models.ShelfInfo](aggData)

	if err != nil {
		return []models.ShelfInfo{}, mongopagination.PaginationData{}, err
	}

	return docList, aggData.Pagination, nil
}

func (repo WarehouseRepository) FindWarehouseByLocation(shopID, warehouseCode, locationCode string) (models.WarehouseDoc, error) {

	filters := bson.M{
		"shopid":        shopID,
		"code":          warehouseCode,
		"location.code": locationCode,
	}

	doc := models.WarehouseDoc{}
	err := repo.pst.FindOne(models.WarehouseDoc{}, filters, &doc)

	if err != nil {
		return models.WarehouseDoc{}, err
	}

	return doc, nil
}

func (repo WarehouseRepository) FindWarehouseByShelf(shopID, warehouseCode, locationCode, shelfCode string) (models.WarehouseDoc, error) {

	filters := bson.M{
		"shopid":              shopID,
		"code":                warehouseCode,
		"location.code":       locationCode,
		"location.shelf.code": shelfCode,
	}

	doc := models.WarehouseDoc{}
	err := repo.pst.FindOne(models.WarehouseDoc{}, filters, &doc)

	if err != nil {
		return models.WarehouseDoc{}, err
	}

	return doc, nil
}
