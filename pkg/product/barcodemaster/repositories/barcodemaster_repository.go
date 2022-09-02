package repositories

import (
	"errors"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/product/barcodemaster/models"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IBarcodeMasterRepository interface {
	CreateInBatch(inventories []models.BarcodeMasterDoc) error
	Create(barcodemaster models.BarcodeMasterDoc) (string, error)
	Update(shopID string, guid string, barcodemaster models.BarcodeMasterDoc) error
	Delete(shopID string, guid string, username string) error
	FindByItemCodeGuid(shopID string, itemCodeGuidList []string) ([]models.BarcodeMasterItemGuid, error)
	FindByID(id primitive.ObjectID) (models.BarcodeMasterDoc, error)
	FindByGuid(shopID string, guid string) (models.BarcodeMasterDoc, error)
	FindPage(shopID string, q string, page int, limit int) ([]models.BarcodeMasterInfo, paginate.PaginationData, error)
	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.BarcodeMasterDeleteActivity, paginate.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.BarcodeMasterActivity, paginate.PaginationData, error)
	FindByItemGuid(shopId string, itemguid string) (models.BarcodeMasterDoc, error)
	FindByItemGuidList(shopID string, guidList []string) ([]models.BarcodeMasterDoc, error)
	FindByItemBarcode(shopId string, barcode string) (models.BarcodeMasterDoc, error)
}

type BarcodeMasterRepository struct {
	pst microservice.IPersisterMongo
}

func NewBarcodeMasterRepository(pst microservice.IPersisterMongo) BarcodeMasterRepository {
	return BarcodeMasterRepository{
		pst: pst,
	}
}

func (repo BarcodeMasterRepository) CreateInBatch(inventories []models.BarcodeMasterDoc) error {
	var tempList []interface{}

	for _, inv := range inventories {
		tempList = append(tempList, inv)
	}

	err := repo.pst.CreateInBatch(&models.BarcodeMasterDoc{}, tempList)

	if err != nil {
		return err
	}
	return nil
}

func (repo BarcodeMasterRepository) Create(barcodemaster models.BarcodeMasterDoc) (string, error) {
	idx, err := repo.pst.Create(&models.BarcodeMasterDoc{}, barcodemaster)

	if err != nil {
		return "", err
	}
	return idx.Hex(), nil
}

func (repo BarcodeMasterRepository) Update(shopID string, guid string, barcodemaster models.BarcodeMasterDoc) error {

	filterDoc := map[string]interface{}{
		"shopid":    shopID,
		"guidfixed": guid,
	}

	err := repo.pst.UpdateOne(&models.BarcodeMasterDoc{}, filterDoc, barcodemaster)

	if err != nil {
		return err
	}

	return nil
}

func (repo BarcodeMasterRepository) Delete(shopID string, guid string, username string) error {

	err := repo.pst.SoftDeleteLastUpdate(&models.BarcodeMasterDoc{}, username, bson.M{"guidfixed": guid, "shopid": shopID})

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	return nil
}

func (repo BarcodeMasterRepository) FindByItemCodeGuid(shopID string, itemCodeGuidList []string) ([]models.BarcodeMasterItemGuid, error) {

	findDoc := []models.BarcodeMasterItemGuid{}
	err := repo.pst.Find(&models.BarcodeMasterItemGuid{}, bson.M{"shopid": shopID, "itemguid": bson.M{"$in": itemCodeGuidList}, "deletedat": bson.M{"$exists": false}}, &findDoc)

	if err != nil {
		return []models.BarcodeMasterItemGuid{}, err
	}
	return findDoc, nil
}

func (repo BarcodeMasterRepository) FindByID(id primitive.ObjectID) (models.BarcodeMasterDoc, error) {

	findDoc := &models.BarcodeMasterDoc{}
	err := repo.pst.FindOne(&models.BarcodeMasterDoc{}, bson.M{"_id": id}, findDoc)

	if err != nil {
		return models.BarcodeMasterDoc{}, err
	}

	if !findDoc.DeletedAt.IsZero() {
		return models.BarcodeMasterDoc{}, errors.New("document not found")
	}

	return *findDoc, nil
}

func (repo BarcodeMasterRepository) FindByGuid(shopID string, guid string) (models.BarcodeMasterDoc, error) {

	findDoc := &models.BarcodeMasterDoc{}
	err := repo.pst.FindOne(&models.BarcodeMasterDoc{}, bson.M{"shopid": shopID, "guidfixed": guid, "deletedat": bson.M{"$exists": false}}, findDoc)

	if err != nil {
		return models.BarcodeMasterDoc{}, err
	}
	return *findDoc, nil
}

func (repo BarcodeMasterRepository) FindPage(shopID string, q string, page int, limit int) ([]models.BarcodeMasterInfo, paginate.PaginationData, error) {

	docList := []models.BarcodeMasterInfo{}
	pagination, err := repo.pst.FindPage(&models.BarcodeMasterInfo{}, limit, page, bson.M{
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
		return []models.BarcodeMasterInfo{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo BarcodeMasterRepository) FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.BarcodeMasterDeleteActivity, paginate.PaginationData, error) {

	docList := []models.BarcodeMasterDeleteActivity{}
	pagination, err := repo.pst.FindPage(&models.BarcodeMasterInfo{}, limit, page, bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$gte": lastUpdatedDate},
	}, &docList)

	if err != nil {
		return []models.BarcodeMasterDeleteActivity{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo BarcodeMasterRepository) FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.BarcodeMasterActivity, paginate.PaginationData, error) {

	docList := []models.BarcodeMasterActivity{}
	pagination, err := repo.pst.FindPage(&models.BarcodeMasterInfo{}, limit, page, bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$not": bson.M{"$gte": lastUpdatedDate}},
		"$or": []interface{}{
			bson.M{"createdat": bson.M{"$gte": lastUpdatedDate}},
			bson.M{"updatedat": bson.M{"$gte": lastUpdatedDate}},
		},
	}, &docList)

	if err != nil {
		return []models.BarcodeMasterActivity{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo BarcodeMasterRepository) FindByItemGuid(shopID string, itemguid string) (models.BarcodeMasterDoc, error) {

	findDoc := models.BarcodeMasterDoc{}
	err := repo.pst.FindOne(&models.BarcodeMasterDoc{}, bson.M{"shopid": shopID, "itemguid": itemguid}, &findDoc)

	if err != nil {
		return models.BarcodeMasterDoc{}, err
	}
	return findDoc, nil
}

func (repo BarcodeMasterRepository) FindByItemGuidList(shopID string, guidList []string) ([]models.BarcodeMasterDoc, error) {

	findDoc := []models.BarcodeMasterDoc{}
	err := repo.pst.Find(&models.BarcodeMasterDoc{}, bson.M{"shopid": shopID, "itemguid": bson.M{"$in": guidList}}, &findDoc)

	if err != nil {
		return []models.BarcodeMasterDoc{}, err
	}
	return findDoc, nil
}

func (repo BarcodeMasterRepository) FindByItemBarcode(shopID string, barcode string) (models.BarcodeMasterDoc, error) {

	findDoc := &models.BarcodeMasterDoc{}
	err := repo.pst.FindOne(&models.BarcodeMasterDoc{}, bson.M{"shopid": shopID, "barcode": barcode, "deletedat": bson.M{"$exists": false}}, findDoc)

	if err != nil {
		return models.BarcodeMasterDoc{}, err
	}
	return *findDoc, nil
}
