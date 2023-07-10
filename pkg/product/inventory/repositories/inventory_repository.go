package repositories

import (
	"context"
	"errors"
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/product/inventory/models"
	"smlcloudplatform/pkg/repositories"
	"smlcloudplatform/pkg/utils/mogoutil"
	"time"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IInventoryRepository interface {
	CreateInBatch(ctx context.Context, docs []models.InventoryDoc) error
	Create(ctx context.Context, doc models.InventoryDoc) (string, error)
	Update(ctx context.Context, shopID string, guid string, inventory models.InventoryDoc) error
	Delete(ctx context.Context, shopID string, guid string, username string) error
	FindByItemCodeGuid(ctx context.Context, shopID string, itemCodeGuidList []string) ([]models.InventoryItemGuid, error)
	FindByID(ctx context.Context, id primitive.ObjectID) (models.InventoryDoc, error)
	FindByGuid(ctx context.Context, shopID string, guid string) (models.InventoryDoc, error)
	FindByItemCode(ctx context.Context, shopID string, itemCode string) (models.InventoryDoc, error)
	FindPage(ctx context.Context, shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.InventoryInfo, mongopagination.PaginationData, error)
	FindByItemGuid(ctx context.Context, shopId string, itemguid string) (models.InventoryDoc, error)
	FindByItemGuidList(ctx context.Context, shopID string, guidList []string) ([]models.InventoryDoc, error)
	FindByItemBarcode(ctx context.Context, shopId string, barcode string) (models.InventoryDoc, error)
	FindByBarcodes(ctx context.Context, shopID string, barcodes []string) ([]models.InventoryDoc, error)

	FindDeletedPage(ctx context.Context, shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.InventoryDeleteActivity, mongopagination.PaginationData, error)
	FindCreatedOrUpdatedPage(ctx context.Context, shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.InventoryActivity, mongopagination.PaginationData, error)
	FindDeletedStep(ctx context.Context, shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.InventoryDeleteActivity, error)
	FindCreatedOrUpdatedStep(ctx context.Context, shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageableStep micromodels.PageableStep) ([]models.InventoryActivity, error)
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

func (repo InventoryRepository) CreateInBatch(ctx context.Context, inventories []models.InventoryDoc) error {
	var tempList []interface{}

	for _, inv := range inventories {
		tempList = append(tempList, inv)
	}

	err := repo.pst.CreateInBatch(ctx, &models.InventoryDoc{}, tempList)

	if err != nil {
		return err
	}
	return nil
}

func (repo InventoryRepository) Create(ctx context.Context, inventory models.InventoryDoc) (string, error) {
	idx, err := repo.pst.Create(ctx, &models.InventoryDoc{}, inventory)

	if err != nil {
		return "", err
	}
	return idx.Hex(), nil
}

func (repo InventoryRepository) Update(ctx context.Context, shopID string, guid string, inventory models.InventoryDoc) error {

	filterDoc := map[string]interface{}{
		"shopid":    shopID,
		"guidfixed": guid,
		"deletedat": bson.M{"$exists": false},
	}

	err := repo.pst.UpdateOne(ctx, &models.InventoryDoc{}, filterDoc, inventory)

	if err != nil {
		return err
	}

	return nil
}

func (repo InventoryRepository) Delete(ctx context.Context, shopID string, guid string, username string) error {

	err := repo.pst.SoftDeleteLastUpdate(ctx, &models.InventoryDoc{}, username, bson.M{"guidfixed": guid, "shopid": shopID})

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	return nil
}

func (repo InventoryRepository) FindByBarcodes(ctx context.Context, shopID string, barcodes []string) ([]models.InventoryDoc, error) {

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

	err := repo.pst.Find(ctx, &models.InventoryDoc{}, filters, &findDoc)

	if err != nil {
		return []models.InventoryDoc{}, err
	}
	return findDoc, nil
}

func (repo InventoryRepository) FindByItemCodeGuid(ctx context.Context, shopID string, itemCodeGuidList []string) ([]models.InventoryItemGuid, error) {

	findDoc := []models.InventoryItemGuid{}
	err := repo.pst.Find(ctx,
		&models.InventoryItemGuid{},
		bson.M{"shopid": shopID, "itemguid": bson.M{"$in": itemCodeGuidList}, "deletedat": bson.M{"$exists": false}},
		&findDoc,
	)

	if err != nil {
		return []models.InventoryItemGuid{}, err
	}
	return findDoc, nil
}

func (repo InventoryRepository) FindByID(ctx context.Context, id primitive.ObjectID) (models.InventoryDoc, error) {

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

	err := repo.pst.Aggregate(ctx, models.InventoryDoc{}, pipeline, &findDocList)

	if err != nil {
		return models.InventoryDoc{}, err
	}

	if len(findDocList) < 1 {
		return models.InventoryDoc{}, errors.New("document not found")
	}

	return findDocList[0], nil
}

func (repo InventoryRepository) FindByGuid(ctx context.Context, shopID string, guid string) (models.InventoryDoc, error) {

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

	err := repo.pst.Aggregate(ctx, models.InventoryDoc{}, pipeline, &findDocList)

	if err != nil {
		return models.InventoryDoc{}, err
	}

	if len(findDocList) < 1 {
		return models.InventoryDoc{}, errors.New("document not found")
	}
	return findDocList[0], nil
}

func (repo InventoryRepository) FindByItemCode(ctx context.Context, shopID string, itemCode string) (models.InventoryDoc, error) {

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

	err := repo.pst.Aggregate(ctx, models.InventoryDoc{}, pipeline, &findDocList)

	if err != nil {
		return models.InventoryDoc{}, err
	}

	if len(findDocList) < 1 {
		return models.InventoryDoc{}, errors.New("document not found")
	}
	return findDocList[0], nil
}

func (repo InventoryRepository) FindPage(ctx context.Context, shopID string, filters map[string]interface{}, pageable micromodels.Pageable) ([]models.InventoryInfo, mongopagination.PaginationData, error) {

	filterQuery := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"itemcode": pageable.Query},
			bson.M{"name1": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + pageable.Query + ".*",
				Options: "i",
			}}},
		},
	}

	if len(filters) > 0 {
		for xkey, xval := range filters {
			filterQuery[xkey] = xval
		}

	}

	matchQuery := bson.M{"$match": filterQuery}

	aggData, err := repo.pst.AggregatePage(ctx, models.InventoryInfo{}, pageable, matchQuery, repo.unitLookupQuery(shopID), repo.unitUnwindQuery())

	if err != nil {
		return []models.InventoryInfo{}, mongopagination.PaginationData{}, err
	}

	docList, err := mogoutil.AggregatePageDecode[models.InventoryInfo](aggData)

	if err != nil {
		return []models.InventoryInfo{}, mongopagination.PaginationData{}, err
	}

	return docList, aggData.Pagination, nil
}

func (repo InventoryRepository) FindDeletedPage(ctx context.Context, shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.InventoryDeleteActivity, mongopagination.PaginationData, error) {

	filterQueries := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$gte": lastUpdatedDate},
	}

	docList := []models.InventoryDeleteActivity{}
	pagination, err := repo.pst.FindPage(ctx, &models.InventoryInfo{}, filterQueries, pageable, &docList)

	if err != nil {
		return []models.InventoryDeleteActivity{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo InventoryRepository) FindCreatedOrUpdatedPage(ctx context.Context, shopID string, lastUpdatedDate time.Time, extraFilters map[string]interface{}, pageable micromodels.Pageable) ([]models.InventoryActivity, mongopagination.PaginationData, error) {

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

	aggData, err := repo.pst.AggregatePage(ctx, models.InventoryActivity{}, pageable, matchQuery, repo.unitLookupQuery(shopID), repo.unitUnwindQuery())

	if err != nil {
		return []models.InventoryActivity{}, mongopagination.PaginationData{}, err
	}

	docList, err := mogoutil.AggregatePageDecode[models.InventoryActivity](aggData)

	if err != nil {
		return []models.InventoryActivity{}, mongopagination.PaginationData{}, err
	}

	return docList, aggData.Pagination, nil
}

func (repo InventoryRepository) FindByItemGuid(ctx context.Context, shopID string, itemguid string) (models.InventoryDoc, error) {

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

	err := repo.pst.Aggregate(ctx, models.InventoryDoc{}, pipeline, &findDocList)

	if err != nil {
		return models.InventoryDoc{}, err
	}

	if len(findDocList) > 0 {
		return models.InventoryDoc{}, nil
	}

	return findDocList[0], nil
}

func (repo InventoryRepository) FindByItemGuidList(ctx context.Context, shopID string, guidList []string) ([]models.InventoryDoc, error) {

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

	err := repo.pst.Aggregate(ctx, models.InventoryDoc{}, pipeline, &findDocList)

	if err != nil {
		return []models.InventoryDoc{}, err
	}
	return findDocList, nil
}

func (repo InventoryRepository) FindByItemBarcode(ctx context.Context, shopID string, barcode string) (models.InventoryDoc, error) {

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

	err := repo.pst.Aggregate(ctx, models.InventoryDoc{}, pipeline, &findDocList)

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
