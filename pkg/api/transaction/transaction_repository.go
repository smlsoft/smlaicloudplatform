package transaction

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ITransactionRepository interface {
	Create(trans models.TransactionDoc) (primitive.ObjectID, error)
	Update(guid string, trans models.TransactionDoc) error
	Delete(guid string, shopID string, username string) error
	FindByGuid(guid string, shopID string) (models.TransactionDoc, error)
	FindPage(shopID string, q string, page int, limit int) ([]models.TransactionInfo, paginate.PaginationData, error)
	FindItemsByGuidPage(guid string, shopID string, q string, page int, limit int) ([]models.TransactionInfo, paginate.PaginationData, error)
}

type TransactionRepository struct {
	pst microservice.IPersisterMongo
}

func NewTransactionRepository(pst microservice.IPersisterMongo) TransactionRepository {
	return TransactionRepository{
		pst: pst,
	}
}

func (repo TransactionRepository) Create(trans models.TransactionDoc) (primitive.ObjectID, error) {
	idx, err := repo.pst.Create(&models.TransactionDoc{}, trans)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return idx, nil
}

func (repo TransactionRepository) Update(guid string, trans models.TransactionDoc) error {
	err := repo.pst.UpdateOne(&models.TransactionDoc{}, "guidfixed", guid, trans)
	if err != nil {
		return err
	}
	return nil
}

func (repo TransactionRepository) Delete(guid string, shopID string, username string) error {
	err := repo.pst.SoftDelete(&models.TransactionDoc{}, username, bson.M{"guidfixed": guid, "shopid": shopID})
	if err != nil {
		return err
	}
	return nil
}

func (repo TransactionRepository) FindByGuid(guid string, shopID string) (models.TransactionDoc, error) {
	trans := &models.TransactionDoc{}
	err := repo.pst.FindOne(
		&models.TransactionDoc{},
		bson.M{"shopid": shopID, "guidfixed": guid, "deletedat": bson.M{"$exists": false}},
		trans,
	)
	if err != nil {
		return *trans, err
	}
	return *trans, nil
}

func (repo TransactionRepository) FindPage(shopID string, q string, page int, limit int) ([]models.TransactionInfo, paginate.PaginationData, error) {

	transList := []models.TransactionInfo{}
	pagination, err := repo.pst.FindPage(&models.TransactionInfo{}, limit, page, bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"guidfixed": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &transList)

	if err != nil {
		return []models.TransactionInfo{}, paginate.PaginationData{}, err
	}

	return transList, pagination, nil
}

func (repo TransactionRepository) FindItemsByGuidPage(guid string, shopID string, q string, page int, limit int) ([]models.TransactionInfo, paginate.PaginationData, error) {

	transList := []models.TransactionInfo{}
	pagination, err := repo.pst.FindPage(&models.Transaction{}, limit, page, bson.M{
		"shopid":    shopID,
		"guidfixed": guid,
		"deletedat": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"items.itemSku": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &transList)

	if err != nil {
		return []models.TransactionInfo{}, paginate.PaginationData{}, err
	}

	return transList, pagination, nil
}
