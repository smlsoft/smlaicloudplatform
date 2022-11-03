package repositories

import (
	"errors"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/inventory/models"
	"smlcloudplatform/pkg/repositories"
	"time"

	"smlcloudplatform/pkg/utils/mogoutil"

	mongopagination "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IInventoryRepository interface {
	CreateInBatch(inventories []models.InventoryDoc) error
	Create(inventory models.InventoryDoc) (string, error)
	Update(shopID string, guid string, inventory models.InventoryDoc) error
	Delete(shopID string, guid string, username string) error
	FindByItemCodeGuid(shopID string, itemCodeGuidList []string) ([]models.InventoryItemGuid, error)
	FindByID(id primitive.ObjectID) (models.InventoryDoc, error)
	FindByGuid(shopID string, guid string) (models.InventoryDoc, error)
	FindByItemCode(shopID string, itemCode string) (models.InventoryDoc, error)
	FindPage(shopID string, filters map[string]interface{}, q string, page int, limit int) ([]models.InventoryInfo, mongopagination.PaginationData, error)
	FindByItemGuid(shopId string, itemguid string) (models.InventoryDoc, error)
	FindByItemGuidList(shopID string, guidList []string) ([]models.InventoryDoc, error)
	FindByItemBarcode(shopId string, barcode string) (models.InventoryDoc, error)
	FindByBarcodes(shopID string, barcodes []string) ([]models.InventoryDoc, error)

	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.InventoryDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.InventoryActivity, mongopagination.PaginationData, error)
	FindDeletedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.InventoryDeleteActivity, error)
	FindCreatedOrUpdatedOffset(shopID string, lastUpdatedDate time.Time, skip int, limit int) ([]models.InventoryActivity, error)
}

type InventoryRepository struct {
	pst microservice.IPersisterMongo
	repositories.ActivityRepository[models.InventoryActivity, models.InventoryDeleteActivity]
}

func NewInventoryRepository(pst microservice.IPersisterMongo) InventoryRepository {

	insRepo := InventoryRepository{
		pst: pst,
	}
	insRepo.ActivityRepository = repositories.NewActivityRepository[models.InventoryActivity, models.InventoryDeleteActivity](pst)
	return insRepo
}

func (repo InventoryRepository) CreateInBatch(inventories []models.InventoryDoc) error {
	var tempList []interface{}

	for _, inv := range inventories {
		tempList = append(tempList, inv)
	}

	err := repo.pst.CreateInBatch(&models.InventoryDoc{}, tempList)

	if err != nil {
		return err
	}
	return nil
}

func (repo InventoryRepository) Create(inventory models.InventoryDoc) (string, error) {
	idx, err := repo.pst.Create(&models.InventoryDoc{}, inventory)

	if err != nil {
		return "", err
	}
	return idx.Hex(), nil
}

func (repo InventoryRepository) Update(shopID string, guid string, inventory models.InventoryDoc) error {

	filterDoc := map[string]interface{}{
		"shopid":    shopID,
		"guidfixed": guid,
		"deletedat": bson.M{"$exists": false},
	}

	err := repo.pst.UpdateOne(&models.InventoryDoc{}, filterDoc, inventory)

	if err != nil {
		return err
	}

	return nil
}

func (repo InventoryRepository) Delete(shopID string, guid string, username string) error {

	err := repo.pst.SoftDeleteLastUpdate(&models.InventoryDoc{}, username, bson.M{"guidfixed": guid, "shopid": shopID})

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	return nil
}

func (repo InventoryRepository) FindByBarcodes(shopID string, barcodes []string) ([]models.InventoryDoc, error) {

	findDoc := []models.InventoryDoc{}

	filters := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"barcodes":  bson.M{"$exists": true},
		"$and": []interface{}{
			bson.M{"$or": []interface{}{bson.M{
				"barcodes.barcode": bson.M{"$in": barcodes},
				"barcode":          bson.M{"$in": barcodes},
			}}},
		},
	}

	err := repo.pst.Find(&models.InventoryDoc{}, filters, &findDoc)

	if err != nil {
		return []models.InventoryDoc{}, err
	}
	return findDoc, nil
}

func (repo InventoryRepository) FindByItemCodeGuid(shopID string, itemCodeGuidList []string) ([]models.InventoryItemGuid, error) {

	findDoc := []models.InventoryItemGuid{}
	err := repo.pst.Find(&models.InventoryItemGuid{}, bson.M{"shopid": shopID, "itemguid": bson.M{"$in": itemCodeGuidList}, "deletedat": bson.M{"$exists": false}}, &findDoc)

	if err != nil {
		return []models.InventoryItemGuid{}, err
	}
	return findDoc, nil
}

func (repo InventoryRepository) FindByID(id primitive.ObjectID) (models.InventoryDoc, error) {

	findDocList := []models.InventoryDoc{}

	pipeline := []bson.M{
		{"$match": bson.M{
			"_id":       id,
			"deletedat": bson.M{"$exists": false},
		}},
		{
			"$lookup": bson.M{
				"from":         "units",
				"localField":   "unitcode",
				"foreignField": "unitcode",
				"as":           "units",
			},
		},
		{
			"$replaceRoot": bson.M{
				"newRoot": bson.M{
					"$mergeObjects": []interface{}{
						bson.M{"$arrayElemAt": []interface{}{"$units", 0}},
						"$$ROOT",
					},
				},
			},
		},
		{
			"$project": bson.M{
				"units": 0,
			},
		},
		{"$limit": 1},
	}

	err := repo.pst.Aggregate(models.InventoryDoc{}, pipeline, &findDocList)

	if err != nil {
		return models.InventoryDoc{}, err
	}

	if len(findDocList) < 1 {
		return models.InventoryDoc{}, errors.New("document not found")
	}

	return findDocList[0], nil
}

func (repo InventoryRepository) FindByGuid(shopID string, guid string) (models.InventoryDoc, error) {

	findDocList := []models.InventoryDoc{}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"shopid":    shopID,
				"guidfixed": guid,
				"deletedat": bson.M{"$exists": false},
			},
		},
		repo.unitLookupQuery(shopID),
		repo.unitUnwindQuery(),
		{"$limit": 1},
	}

	err := repo.pst.Aggregate(models.InventoryDoc{}, pipeline, &findDocList)

	if err != nil {
		return models.InventoryDoc{}, err
	}

	if len(findDocList) < 1 {
		return models.InventoryDoc{}, errors.New("document not found")
	}
	return findDocList[0], nil
}

func (repo InventoryRepository) FindByItemCode(shopID string, itemCode string) (models.InventoryDoc, error) {

	findDocList := []models.InventoryDoc{}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"shopid":    shopID,
				"itemcode":  itemCode,
				"deletedat": bson.M{"$exists": false},
			},
		},
		repo.unitLookupQuery(shopID),
		repo.unitUnwindQuery(),
		{"$limit": 1},
	}

	err := repo.pst.Aggregate(models.InventoryDoc{}, pipeline, &findDocList)

	if err != nil {
		return models.InventoryDoc{}, err
	}

	if len(findDocList) < 1 {
		return models.InventoryDoc{}, errors.New("document not found")
	}
	return findDocList[0], nil
}

func (repo InventoryRepository) FindPage(shopID string, filters map[string]interface{}, q string, page int, limit int) ([]models.InventoryInfo, mongopagination.PaginationData, error) {

	filterQuery := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"itemcode": q},
			bson.M{"name1": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}

	if len(filters) > 0 {
		for xkey, xval := range filters {
			filterQuery[xkey] = xval
		}

	}

	matchQuery := bson.M{"$match": filterQuery}

	aggData, err := repo.pst.AggregatePage(models.InventoryInfo{}, limit, page, matchQuery, repo.unitLookupQuery(shopID), repo.unitUnwindQuery())

	if err != nil {
		return []models.InventoryInfo{}, mongopagination.PaginationData{}, err
	}

	docList, err := mogoutil.AggregatePageDecode[models.InventoryInfo](aggData)

	if err != nil {
		return []models.InventoryInfo{}, mongopagination.PaginationData{}, err
	}

	return docList, aggData.Pagination, nil
}

func (repo InventoryRepository) FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.InventoryDeleteActivity, mongopagination.PaginationData, error) {

	docList := []models.InventoryDeleteActivity{}
	pagination, err := repo.pst.FindPage(&models.InventoryInfo{}, limit, page, bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$gte": lastUpdatedDate},
	}, &docList)

	if err != nil {
		return []models.InventoryDeleteActivity{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo InventoryRepository) FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.InventoryActivity, mongopagination.PaginationData, error) {

	matchQuery := bson.M{
		"$match": bson.M{
			"shopid":    shopID,
			"deletedat": bson.M{"$not": bson.M{"$gte": lastUpdatedDate}},
			"$or": []interface{}{
				bson.M{"createdat": bson.M{"$gte": lastUpdatedDate}},
				bson.M{"updatedat": bson.M{"$gte": lastUpdatedDate}},
			},
		},
	}

	aggData, err := repo.pst.AggregatePage(models.InventoryActivity{}, limit, page, matchQuery, repo.unitLookupQuery(shopID), repo.unitUnwindQuery())

	if err != nil {
		return []models.InventoryActivity{}, mongopagination.PaginationData{}, err
	}

	docList, err := mogoutil.AggregatePageDecode[models.InventoryActivity](aggData)

	if err != nil {
		return []models.InventoryActivity{}, mongopagination.PaginationData{}, err
	}

	return docList, aggData.Pagination, nil
}

func (repo InventoryRepository) FindByItemGuid(shopID string, itemguid string) (models.InventoryDoc, error) {

	findDocList := []models.InventoryDoc{}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"shopid":    shopID,
				"itemguid":  itemguid,
				"deletedat": bson.M{"$exists": false},
			},
		},
		repo.unitLookupQuery(shopID),
		repo.unitUnwindQuery(),
		{"$limit": 1},
	}

	err := repo.pst.Aggregate(models.InventoryDoc{}, pipeline, &findDocList)

	if err != nil {
		return models.InventoryDoc{}, err
	}

	if len(findDocList) > 0 {
		return models.InventoryDoc{}, nil
	}

	return findDocList[0], nil
}

func (repo InventoryRepository) FindByItemGuidList(shopID string, guidList []string) ([]models.InventoryDoc, error) {

	findDocList := []models.InventoryDoc{}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"shopid":    shopID,
				"itemguid":  bson.M{"$in": guidList},
				"deletedat": bson.M{"$exists": false},
			},
		},
		repo.unitLookupQuery(shopID),
		repo.unitUnwindQuery(),
	}

	err := repo.pst.Aggregate(models.InventoryDoc{}, pipeline, &findDocList)

	if err != nil {
		return []models.InventoryDoc{}, err
	}
	return findDocList, nil
}

func (repo InventoryRepository) FindByItemBarcode(shopID string, barcode string) (models.InventoryDoc, error) {

	findDocList := []models.InventoryDoc{}

	pipeline := []bson.M{
		{
			"$match": bson.M{
				"shopid":    shopID,
				"deletedat": bson.M{"$exists": false},
				"barcodes":  bson.M{"$exists": true},
				"$and": []interface{}{
					bson.M{"$or": []interface{}{bson.M{
						"barcodes.barcode": barcode,
						"barcode":          barcode,
					}}},
				},
			},
		},
		repo.unitLookupQuery(shopID),
		repo.unitUnwindQuery(),
		{"$limit": 1},
	}

	err := repo.pst.Aggregate(models.InventoryDoc{}, pipeline, &findDocList)

	if err != nil {
		return models.InventoryDoc{}, err
	}

	if len(findDocList) < 1 {
		return models.InventoryDoc{}, nil
	}

	return findDocList[0], nil
}

func (repo InventoryRepository) unitLookupQuery(shopID string) bson.M {
	return bson.M{
		"$lookup": bson.M{
			"from": "units",
			"let":  bson.M{"unitcode": "$unitcode"},
			"pipeline": []interface{}{
				bson.M{"$match": bson.M{
					"shopid": shopID,
					"$expr":  bson.M{"$eq": []string{"$$unitcode", "$unitcode"}},
				}},
			},
			"as": "unit",
		},
	}
}

func (repo InventoryRepository) unitUnwindQuery() bson.M {
	return bson.M{
		"$unwind": bson.M{
			"path":                       "$unit",
			"preserveNullAndEmptyArrays": true,
		},
	}
}

// func (repo InventoryRepository) unitReplaceQuery() bson.M {
// 	return bson.M{
// 		"$replaceRoot": bson.M{
// 			"newRoot": bson.M{
// 				"$mergeObjects": []interface{}{
// 					bson.M{"$arrayElemAt": []interface{}{"$units", 0}},
// 					"$$ROOT",
// 				},
// 			},
// 		},
// 	}
// }

// func (repo InventoryRepository) unitProjectQuery() bson.M {
// 	return bson.M{"$project": bson.M{
// 		"units": 0,
// 	}}
// }
