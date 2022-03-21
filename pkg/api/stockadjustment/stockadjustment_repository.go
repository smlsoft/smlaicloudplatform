package stockadjustment

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IStockAdjustmentRepository interface {
	Create(doc models.StockAdjustment) (primitive.ObjectID, error)
	Update(guid string, doc models.StockAdjustment) error
	Delete(guid string, shopID string) error
	FindByGuid(guid string, shopID string) (models.StockAdjustment, error)
	FindPage(shopID string, q string, page int, limit int) ([]models.StockAdjustment, paginate.PaginationData, error)
	FindItemsByGuidPage(guid string, shopID string, q string, page int, limit int) ([]models.StockAdjustment, paginate.PaginationData, error)
}

type StockAdjustmentRepository struct {
	pst microservice.IPersisterMongo
}

func NewStockAdjustmentRepository(pst microservice.IPersisterMongo) StockAdjustmentRepository {
	return StockAdjustmentRepository{
		pst: pst,
	}
}

func (repo StockAdjustmentRepository) Create(doc models.StockAdjustment) (primitive.ObjectID, error) {
	idx, err := repo.pst.Create(&models.StockAdjustment{}, doc)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return idx, nil
}

func (repo StockAdjustmentRepository) Update(guid string, doc models.StockAdjustment) error {
	err := repo.pst.UpdateOne(&models.StockAdjustment{}, "guidFixed", guid, doc)
	if err != nil {
		return err
	}
	return nil
}

func (repo StockAdjustmentRepository) Delete(guid string, shopID string) error {
	err := repo.pst.SoftDelete(&models.StockAdjustment{}, bson.M{"guidFixed": guid, "shopID": shopID})
	if err != nil {
		return err
	}
	return nil
}

func (repo StockAdjustmentRepository) FindByGuid(guid string, shopID string) (models.StockAdjustment, error) {
	doc := &models.StockAdjustment{}
	err := repo.pst.FindOne(&models.StockAdjustment{}, bson.M{"shopID": shopID, "guidFixed": guid, "deleted": false}, doc)
	if err != nil {
		return *doc, err
	}
	return *doc, nil
}

func (repo StockAdjustmentRepository) FindPage(shopID string, q string, page int, limit int) ([]models.StockAdjustment, paginate.PaginationData, error) {

	docList := []models.StockAdjustment{}
	pagination, err := repo.pst.FindPage(&models.StockAdjustment{}, limit, page, bson.M{
		"shopID":  shopID,
		"deleted": false,
		"$or": []interface{}{
			bson.M{"guidFixed": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &docList)

	if err != nil {
		return []models.StockAdjustment{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo StockAdjustmentRepository) FindItemsByGuidPage(guid string, shopID string, q string, page int, limit int) ([]models.StockAdjustment, paginate.PaginationData, error) {

	docList := []models.StockAdjustment{}
	pagination, err := repo.pst.FindPage(&models.StockAdjustment{}, limit, page, bson.M{
		"shopID":    shopID,
		"guidFixed": guid,
		"deleted":   false,
		"$or": []interface{}{
			bson.M{"items.itemSku": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &docList)

	if err != nil {
		return []models.StockAdjustment{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}
