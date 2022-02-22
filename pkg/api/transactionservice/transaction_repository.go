package transactionservice

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ITransactionRepository interface {
	Create(trans models.Transaction) (string, error)
	Update(transGuid string, trans models.Transaction) error
	Delete(transGuid string) error
	FindByGuid(transGuid string, merchantId string) (models.Transaction, error)
	FindPage(merchantId string, q string, page int, limit int) ([]models.Transaction, paginate.PaginationData, error)
	FindItemsByGuidPage(transGuid string, merchantId string, q string, page int, limit int) ([]models.Transaction, paginate.PaginationData, error)
}

type TransactionRepository struct {
	pst microservice.IPersisterMongo
}

func NewTransactionRepository(pst microservice.IPersisterMongo) ITransactionRepository {
	return &TransactionRepository{
		pst: pst,
	}
}

func (repo *TransactionRepository) Create(trans models.Transaction) (string, error) {
	idx, err := repo.pst.Create(&models.Transaction{}, trans)
	if err != nil {
		return "", err
	}

	return idx.Hex(), nil
}

func (repo *TransactionRepository) Update(transGuid string, trans models.Transaction) error {
	err := repo.pst.UpdateOne(&models.Transaction{}, "guidFixed", transGuid, trans)
	if err != nil {
		return err
	}
	return nil
}

func (repo *TransactionRepository) Delete(transGuid string) error {
	err := repo.pst.SoftDeleteByID(&models.Transaction{}, transGuid)
	if err != nil {
		return err
	}
	return nil
}

func (repo *TransactionRepository) FindByGuid(merchantId string, transGuid string) (models.Transaction, error) {
	trans := &models.Transaction{}
	err := repo.pst.FindOne(&models.Transaction{}, bson.M{"merchantId": merchantId, "guidFixed": transGuid, "deleted": false}, trans)
	if err != nil {
		return *trans, err
	}
	return *trans, nil
}

func (repo *TransactionRepository) FindPage(merchantId string, q string, page int, limit int) ([]models.Transaction, paginate.PaginationData, error) {

	transList := []models.Transaction{}
	pagination, err := repo.pst.FindPage(&models.Inventory{}, limit, page, bson.M{
		"merchantId": merchantId,
		"deleted":    false,
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

func (repo *TransactionRepository) FindItemsByGuidPage(transGuid string, merchantId string, q string, page int, limit int) ([]models.Transaction, paginate.PaginationData, error) {

	transList := []models.Transaction{}
	pagination, err := repo.pst.FindPage(&models.Inventory{}, limit, page, bson.M{
		"merchantId": merchantId,
		"guidFixed":  transGuid,
		"deleted":    false,
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
