package inventory

import (
	"errors"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"
	"time"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IInventoryRepository interface {
	CreateInBatch(inventories []models.InventoryDoc) error
	Create(inventory models.InventoryDoc) (string, error)
	Update(guid string, inventory models.InventoryDoc) error
	Delete(shopID string, guid string, username string) error
	FindByID(id primitive.ObjectID) (models.InventoryDoc, error)
	FindByGuid(shopID string, guid string) (models.InventoryDoc, error)
	FindPage(shopID string, q string, page int, limit int) ([]models.InventoryInfo, paginate.PaginationData, error)
	FindDeletedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.InventoryDeleteActivity, paginate.PaginationData, error)
	FindCreatedOrUpdatedPage(shopID string, lastUpdatedDate time.Time, page int, limit int) ([]models.InventoryActivity, paginate.PaginationData, error)
	FindByItemGuid(itemguid string, shopId string) (models.InventoryDoc, error)
	FindByItemBarcode(shopId string, barcode string) (models.InventoryDoc, error)
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

func (repo InventoryRepository) Update(guid string, inventory models.InventoryDoc) error {

	err := repo.pst.UpdateOne(&models.InventoryDoc{}, "guidfixed", guid, inventory)

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

func (repo InventoryRepository) FindByID(id primitive.ObjectID) (models.InventoryDoc, error) {

	findDoc := &models.InventoryDoc{}
	err := repo.pst.FindOne(&models.InventoryDoc{}, bson.M{"_id": id}, findDoc)

	if err != nil {
		return models.InventoryDoc{}, err
	}

	if !findDoc.DeletedAt.IsZero() {
		return models.InventoryDoc{}, errors.New("document not found")
	}

	return *findDoc, nil
}

func (repo InventoryRepository) FindByGuid(shopID string, guid string) (models.InventoryDoc, error) {

	findDoc := &models.InventoryDoc{}
	err := repo.pst.FindOne(&models.InventoryDoc{}, bson.M{"shopid": shopID, "guidfixed": guid, "deletedat": bson.M{"$exists": false}}, findDoc)

	if err != nil {
		return models.InventoryDoc{}, err
	}
	return *findDoc, nil
}

func (repo InventoryRepository) FindPage(shopID string, q string, page int, limit int) ([]models.InventoryInfo, paginate.PaginationData, error) {

	docList := []models.InventoryInfo{}
	pagination, err := repo.pst.FindPage(&models.InventoryInfo{}, limit, page, bson.M{
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
		return []models.InventoryInfo{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
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

	docList := []models.InventoryActivity{}
	pagination, err := repo.pst.FindPage(&models.InventoryInfo{}, limit, page, bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$not": bson.M{"$gte": lastUpdatedDate}},
		"$or": []interface{}{
			bson.M{"createdat": bson.M{"$gte": lastUpdatedDate}},
			bson.M{"updatedat": bson.M{"$gte": lastUpdatedDate}},
		},
	}, &docList)

	if err != nil {
		return []models.InventoryActivity{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo InventoryRepository) FindByItemGuid(itemguid string, shopID string) (models.InventoryDoc, error) {

	findDoc := &models.InventoryDoc{}
	err := repo.pst.FindOne(&models.InventoryDoc{}, bson.M{"shopid": shopID, "itemguid": itemguid, "deletedat": bson.M{"$exists": false}}, findDoc)

	if err != nil {
		return models.InventoryDoc{}, err
	}
	return *findDoc, nil
}

func (repo InventoryRepository) FindByItemBarcode(shopID string, barcode string) (models.InventoryDoc, error) {

	findDoc := &models.InventoryDoc{}
	err := repo.pst.FindOne(&models.InventoryDoc{}, bson.M{"shopid": shopID, "barcode": barcode, "deletedat": bson.M{"$exists": false}}, findDoc)

	if err != nil {
		return models.InventoryDoc{}, err
	}
	return *findDoc, nil
}
