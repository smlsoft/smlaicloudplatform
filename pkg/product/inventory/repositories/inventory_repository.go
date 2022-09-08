package repositories

import (
	"errors"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/inventory/models"
	"time"

	"smlcloudplatform/pkg/utils/mogoutil"

	paginate "github.com/gobeam/mongo-go-pagination"
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
	FindPage(shopID string, q string, page int, limit int) ([]models.InventoryInfo, paginate.PaginationData, error)
	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.InventoryDeleteActivity, paginate.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.InventoryActivity, paginate.PaginationData, error)
	FindByItemGuid(shopId string, itemguid string) (models.InventoryDoc, error)
	FindByItemGuidList(shopID string, guidList []string) ([]models.InventoryDoc, error)
	FindByItemBarcode(shopId string, barcode string) (models.InventoryDoc, error)
	FindByBarcodes(shopID string, barcodes []string) ([]models.InventoryDoc, error)
}

type InventoryRepository struct {
	pst microservice.IPersisterMongo
}

func NewInventoryRepository(pst microservice.IPersisterMongo) InventoryRepository {
	return InventoryRepository{
		pst: pst,
	}
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
		"shopid":           shopID,
		"barcodes":         bson.M{"$exists": true},
		"barcodes.barcode": bson.M{"$in": barcodes},
		"deletedat":        bson.M{"$exists": false},
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
		repo.unitLookupQuery(),
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

func (repo InventoryRepository) FindPage(shopID string, q string, page int, limit int) ([]models.InventoryInfo, paginate.PaginationData, error) {

	matchQuery := bson.M{"$match": bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"guidfixed": q},
			bson.M{"name1": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}}

	aggData, err := repo.pst.AggregatePage(models.InventoryInfo{}, limit, page, matchQuery, repo.unitLookupQuery(), repo.unitUnwindQuery())

	if err != nil {
		return []models.InventoryInfo{}, paginate.PaginationData{}, err
	}

	docList, err := mogoutil.AggregatePageDecode[models.InventoryInfo](aggData)

	if err != nil {
		return []models.InventoryInfo{}, paginate.PaginationData{}, err
	}

	return docList, aggData.Pagination, nil
}

func (repo InventoryRepository) FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.InventoryDeleteActivity, paginate.PaginationData, error) {

	docList := []models.InventoryDeleteActivity{}
	pagination, err := repo.pst.FindPage(&models.InventoryInfo{}, limit, page, bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$gte": lastUpdatedDate},
	}, &docList)

	if err != nil {
		return []models.InventoryDeleteActivity{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo InventoryRepository) FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.InventoryActivity, paginate.PaginationData, error) {

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

	aggData, err := repo.pst.AggregatePage(models.InventoryActivity{}, limit, page, matchQuery, repo.unitLookupQuery(), repo.unitUnwindQuery())

	if err != nil {
		return []models.InventoryActivity{}, paginate.PaginationData{}, err
	}

	docList, err := mogoutil.AggregatePageDecode[models.InventoryActivity](aggData)

	if err != nil {
		return []models.InventoryActivity{}, paginate.PaginationData{}, err
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
		repo.unitLookupQuery(),
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
		repo.unitLookupQuery(),
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
				"barcode":   barcode,
				"deletedat": bson.M{"$exists": false},
			},
		},
		repo.unitLookupQuery(),
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

func (repo InventoryRepository) unitLookupQuery() bson.M {
	return bson.M{
		"$lookup": bson.M{
			"from":         "units",
			"localField":   "unitcode",
			"foreignField": "unitcode",
			"as":           "unit",
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
