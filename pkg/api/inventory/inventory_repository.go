package inventory

import (
	"errors"
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IInventoryRepository interface {
	Create(inventory models.InventoryDoc) (string, error)
	Update(guid string, inventory models.InventoryDoc) error
	Delete(shopID string, guid string, username string) error
	FindByID(id primitive.ObjectID) (models.InventoryDoc, error)
	FindByGuid(shopID string, guid string) (models.InventoryDoc, error)
	FindPage(shopID string, q string, page int, limit int) ([]models.InventoryInfo, paginate.PaginationData, error)
}

type InventoryRepository struct {
	pst microservice.IPersisterMongo
}

func NewInventoryRepository(pst microservice.IPersisterMongo) InventoryRepository {
	return InventoryRepository{
		pst: pst,
	}
}

func (repo InventoryRepository) Create(inventory models.InventoryDoc) (string, error) {
	idx, err := repo.pst.Create(&models.InventoryDoc{}, inventory)

	if err != nil {
		return "", err
	}
	return idx.Hex(), nil
}

func (repo InventoryRepository) Update(guid string, inventory models.InventoryDoc) error {

	err := repo.pst.UpdateOne(&models.InventoryDoc{}, "guidFixed", guid, inventory)

	if err != nil {
		return err
	}

	return nil
}

func (repo InventoryRepository) Delete(shopID string, guid string, username string) error {

	err := repo.pst.SoftDelete(&models.InventoryDoc{}, username, bson.M{"guidFixed": guid, "shopID": shopID})

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
	err := repo.pst.FindOne(&models.InventoryDoc{}, bson.M{"shopID": shopID, "guidFixed": guid, "deletedAt": bson.M{"$exists": false}}, findDoc)

	if err != nil {
		return models.InventoryDoc{}, err
	}
	return *findDoc, nil
}

func (repo InventoryRepository) FindPage(shopID string, q string, page int, limit int) ([]models.InventoryInfo, paginate.PaginationData, error) {

	docList := []models.InventoryInfo{}
	pagination, err := repo.pst.FindPage(&models.InventoryInfo{}, limit, page, bson.M{
		"shopID":    shopID,
		"deletedAt": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"guidFixed": q},
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
