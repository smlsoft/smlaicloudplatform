package transaction

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ITransactionRepository interface {
	Create(trans models.Transaction) (primitive.ObjectID, error)
	Update(guid string, trans models.Transaction) error
	Delete(guid string, shopID string) error
	FindByGuid(guid string, shopID string) (models.Transaction, error)
	FindPage(shopID string, q string, page int, limit int) ([]models.Transaction, paginate.PaginationData, error)
	FindItemsByGuidPage(guid string, shopID string, q string, page int, limit int) ([]models.Transaction, paginate.PaginationData, error)
}

type TransactionRepository struct {
	pst microservice.IPersisterMongo
}

func NewTransactionRepository(pst microservice.IPersisterMongo) TransactionRepository {
	return TransactionRepository{
		pst: pst,
	}
}

func (repo TransactionRepository) Create(trans models.Transaction) (primitive.ObjectID, error) {
	idx, err := repo.pst.Create(&models.Transaction{}, trans)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return idx, nil
}

func (repo TransactionRepository) Update(guid string, trans models.Transaction) error {
	err := repo.pst.UpdateOne(&models.Transaction{}, "guidFixed", guid, trans)
	if err != nil {
		return err
	}
	return nil
}

func (repo TransactionRepository) Delete(guid string, shopID string) error {
	err := repo.pst.SoftDelete(&models.Transaction{}, bson.M{"guidFixed": guid, "shopID": shopID})
	if err != nil {
		return err
	}
	return nil
}

func (repo TransactionRepository) FindByGuid(guid string, shopID string) (models.Transaction, error) {
	trans := &models.Transaction{}
	err := repo.pst.FindOne(&models.Transaction{}, bson.M{"shopID": shopID, "guidFixed": guid, "deleted": false}, trans)
	if err != nil {
		return *trans, err
	}
	return *trans, nil
}

func (repo TransactionRepository) FindPage(shopID string, q string, page int, limit int) ([]models.Transaction, paginate.PaginationData, error) {

	transList := []models.Transaction{}
	pagination, err := repo.pst.FindPage(&models.Transaction{}, limit, page, bson.M{
		"shopID":  shopID,
		"deleted": false,
		"$or": []interface{}{
			bson.M{"guidFixed": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &transList)

	if err != nil {
		return []models.Transaction{}, paginate.PaginationData{}, err
	}

	return transList, pagination, nil
}

func (repo TransactionRepository) FindItemsByGuidPage(guid string, shopID string, q string, page int, limit int) ([]models.Transaction, paginate.PaginationData, error) {

	transList := []models.Transaction{}
	pagination, err := repo.pst.FindPage(&models.Transaction{}, limit, page, bson.M{
		"shopID":    shopID,
		"guidFixed": guid,
		"deleted":   false,
		"$or": []interface{}{
			bson.M{"items.itemSku": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &transList)

	if err != nil {
		return []models.Transaction{}, paginate.PaginationData{}, err
	}

	return transList, pagination, nil
}
