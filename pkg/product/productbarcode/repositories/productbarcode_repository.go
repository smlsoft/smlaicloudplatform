package repositories

import (
	"errors"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/productbarcode/models"
	"smlcloudplatform/pkg/utils/mogoutil"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IProductBarcodeRepository interface {
	CreateInBatch(inventories []models.ProductBarcodeDoc) error
	Create(productbarcode models.ProductBarcodeDoc) (string, error)
	Update(shopID string, guid string, productbarcode models.ProductBarcodeDoc) error
	Delete(shopID string, guid string, username string) error
	FindByItemCodeGuid(shopID string, itemCodeGuidList []string) ([]models.ProductBarcodeItemGuid, error)
	FindByID(id primitive.ObjectID) (models.ProductBarcodeDoc, error)
	FindByGuid(shopID string, guid string) (models.ProductBarcodeDoc, error)
	FindPage(shopID string, q string, page int, limit int) ([]models.ProductBarcodeInfo, paginate.PaginationData, error)
	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.ProductBarcodeDeleteActivity, paginate.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.ProductBarcodeActivity, paginate.PaginationData, error)
	FindByItemGuid(shopId string, itemguid string) (models.ProductBarcodeDoc, error)
	FindByItemGuidList(shopID string, guidList []string) ([]models.ProductBarcodeDoc, error)
	FindByItemBarcode(shopId string, barcode string) (models.ProductBarcodeDoc, error)
	FindByBarcodes(shopID string, barcodes []string) ([]models.ProductBarcodeDoc, error)
}

type ProductBarcodeRepository struct {
	pst microservice.IPersisterMongo
}

func NewProductBarcodeRepository(pst microservice.IPersisterMongo) ProductBarcodeRepository {
	return ProductBarcodeRepository{
		pst: pst,
	}
}

func (repo ProductBarcodeRepository) CreateInBatch(inventories []models.ProductBarcodeDoc) error {
	var tempList []interface{}

	for _, inv := range inventories {
		tempList = append(tempList, inv)
	}

	err := repo.pst.CreateInBatch(&models.ProductBarcodeDoc{}, tempList)

	if err != nil {
		return err
	}
	return nil
}

func (repo ProductBarcodeRepository) Create(productbarcode models.ProductBarcodeDoc) (string, error) {
	idx, err := repo.pst.Create(&models.ProductBarcodeDoc{}, productbarcode)

	if err != nil {
		return "", err
	}
	return idx.Hex(), nil
}

func (repo ProductBarcodeRepository) Update(shopID string, guid string, productbarcode models.ProductBarcodeDoc) error {

	filterDoc := map[string]interface{}{
		"shopid":    shopID,
		"guidfixed": guid,
	}

	err := repo.pst.UpdateOne(&models.ProductBarcodeDoc{}, filterDoc, productbarcode)

	if err != nil {
		return err
	}

	return nil
}

func (repo ProductBarcodeRepository) Delete(shopID string, guid string, username string) error {

	err := repo.pst.SoftDeleteLastUpdate(&models.ProductBarcodeDoc{}, username, bson.M{"guidfixed": guid, "shopid": shopID})

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	return nil
}

func (repo ProductBarcodeRepository) FindByItemCodeGuid(shopID string, itemCodeGuidList []string) ([]models.ProductBarcodeItemGuid, error) {

	findDoc := []models.ProductBarcodeItemGuid{}
	err := repo.pst.Find(&models.ProductBarcodeItemGuid{}, bson.M{"shopid": shopID, "itemguid": bson.M{"$in": itemCodeGuidList}, "deletedat": bson.M{"$exists": false}}, &findDoc)

	if err != nil {
		return []models.ProductBarcodeItemGuid{}, err
	}
	return findDoc, nil
}

func (repo ProductBarcodeRepository) FindByID(id primitive.ObjectID) (models.ProductBarcodeDoc, error) {

	findDoc := &models.ProductBarcodeDoc{}
	err := repo.pst.FindOne(&models.ProductBarcodeDoc{}, bson.M{"_id": id}, findDoc)

	if err != nil {
		return models.ProductBarcodeDoc{}, err
	}

	if !findDoc.DeletedAt.IsZero() {
		return models.ProductBarcodeDoc{}, errors.New("document not found")
	}

	return *findDoc, nil
}

func (repo ProductBarcodeRepository) FindByGuid(shopID string, guid string) (models.ProductBarcodeDoc, error) {

	findDoc := &models.ProductBarcodeDoc{}
	err := repo.pst.FindOne(&models.ProductBarcodeDoc{}, bson.M{"shopid": shopID, "guidfixed": guid, "deletedat": bson.M{"$exists": false}}, findDoc)

	if err != nil {
		return models.ProductBarcodeDoc{}, err
	}
	return *findDoc, nil
}

func (repo ProductBarcodeRepository) FindPage(shopID string, q string, page int, limit int) ([]models.ProductBarcodeInfo, paginate.PaginationData, error) {

	docList := []models.ProductBarcodeInfo{}
	pagination, err := repo.pst.FindPage(&models.ProductBarcodeInfo{}, limit, page, bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"guidfixed": q},
			bson.M{"name1": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &docList)

	if err != nil {
		return []models.ProductBarcodeInfo{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo ProductBarcodeRepository) FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.ProductBarcodeDeleteActivity, paginate.PaginationData, error) {

	docList := []models.ProductBarcodeDeleteActivity{}
	repo.pst.AggregatePage(shopID, limit, page)
	pagination, err := repo.pst.FindPage(&models.ProductBarcodeInfo{}, limit, page, bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$gte": lastUpdatedDate},
	}, &docList)

	if err != nil {
		return []models.ProductBarcodeDeleteActivity{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo ProductBarcodeRepository) FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.ProductBarcodeActivity, paginate.PaginationData, error) {

	matchQuery1 := bson.M{
		"$match": bson.M{"shopid": shopID,
			"deletedat": bson.M{"$not": bson.M{"$gte": lastUpdatedDate}},
			"$or": []interface{}{
				bson.M{"createdat": bson.M{"$gte": lastUpdatedDate}},
				bson.M{"updatedat": bson.M{"$gte": lastUpdatedDate}},
			},
		},
	}

	unwindBarcodeQuery2 := bson.M{
		"$unwind": bson.M{"path": "$barcodes",
			"preserveNullAndEmptyArrays": true,
		},
	}

	projectDocQuery3 := bson.M{
		"$project": bson.M{
			"barcodedetail": "$barcodes",
			"doc":           "$$ROOT",
			"unitcode":      "$barcodes.unitcode",
		},
	}

	lookupUnitQuery4 := bson.M{
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

	unwindUnitQuery5 := bson.M{
		"$unwind": bson.M{
			"path":                       "$unit",
			"preserveNullAndEmptyArrays": true,
		},
	}

	replaceRootDocQuery6 := bson.M{
		"$replaceRoot": bson.M{
			"newRoot": bson.M{
				"$mergeObjects": []string{
					"$doc",
					"$$ROOT",
				},
			},
		},
	}

	projectUnuseQuery7 := bson.M{
		"$project": bson.M{
			"barcodes": 0,
			"doc":      0,
		},
	}

	aggData, err := repo.pst.AggregatePage(models.ProductBarcodeActivity{}, limit, page,
		matchQuery1,
		unwindBarcodeQuery2,
		projectDocQuery3,
		lookupUnitQuery4,
		unwindUnitQuery5,
		replaceRootDocQuery6,
		projectUnuseQuery7,
	)

	if err != nil {
		return []models.ProductBarcodeActivity{}, paginate.PaginationData{}, err
	}

	docList, err := mogoutil.AggregatePageDecode[models.ProductBarcodeActivity](aggData)

	if err != nil {
		return []models.ProductBarcodeActivity{}, paginate.PaginationData{}, err
	}

	return docList, aggData.Pagination, nil
}

func (repo ProductBarcodeRepository) FindByBarcodes(shopID string, barcodes []string) ([]models.ProductBarcodeDoc, error) {

	findDoc := []models.ProductBarcodeDoc{}

	filters := bson.M{
		"shopid":           shopID,
		"barcodes":         bson.M{"$exists": true},
		"barcodes.barcode": bson.M{"$in": barcodes},
	}

	err := repo.pst.Find(&models.ProductBarcodeDoc{}, filters, &findDoc)

	if err != nil {
		return []models.ProductBarcodeDoc{}, err
	}
	return findDoc, nil
}

func (repo ProductBarcodeRepository) FindByItemGuid(shopID string, itemguid string) (models.ProductBarcodeDoc, error) {

	findDoc := models.ProductBarcodeDoc{}
	err := repo.pst.FindOne(&models.ProductBarcodeDoc{}, bson.M{"shopid": shopID, "itemguid": itemguid}, &findDoc)

	if err != nil {
		return models.ProductBarcodeDoc{}, err
	}
	return findDoc, nil
}

func (repo ProductBarcodeRepository) FindByItemGuidList(shopID string, guidList []string) ([]models.ProductBarcodeDoc, error) {

	findDoc := []models.ProductBarcodeDoc{}
	err := repo.pst.Find(&models.ProductBarcodeDoc{}, bson.M{"shopid": shopID, "itemguid": bson.M{"$in": guidList}}, &findDoc)

	if err != nil {
		return []models.ProductBarcodeDoc{}, err
	}
	return findDoc, nil
}

func (repo ProductBarcodeRepository) FindByItemBarcode(shopID string, barcode string) (models.ProductBarcodeDoc, error) {

	findDoc := &models.ProductBarcodeDoc{}
	err := repo.pst.FindOne(&models.ProductBarcodeDoc{}, bson.M{"shopid": shopID, "barcode": barcode, "deletedat": bson.M{"$exists": false}}, findDoc)

	if err != nil {
		return models.ProductBarcodeDoc{}, err
	}
	return *findDoc, nil
}
