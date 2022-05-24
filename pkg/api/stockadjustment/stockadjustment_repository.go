package stockadjustment

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IStockAdjustmentRepository interface {
	Create(doc models.StockAdjustmentDoc) (primitive.ObjectID, error)
	Update(shopID string, guid string, doc models.StockAdjustmentDoc) error
	Delete(shopID string, guid string, username string) error
	FindByGuid(shopID string, guid string) (models.StockAdjustmentDoc, error)
	FindPage(shopID string, q string, page int, limit int) ([]models.StockAdjustmentInfo, paginate.PaginationData, error)
	FindItemsByGuidPage(shopID string, guid string, q string, page int, limit int) ([]models.StockAdjustmentInfo, paginate.PaginationData, error)
}

type StockAdjustmentRepository struct {
	pst microservice.IPersisterMongo
}

func NewStockAdjustmentRepository(pst microservice.IPersisterMongo) StockAdjustmentRepository {
	return StockAdjustmentRepository{
		pst: pst,
	}
}

func (repo StockAdjustmentRepository) Create(doc models.StockAdjustmentDoc) (primitive.ObjectID, error) {
	idx, err := repo.pst.Create(&models.StockAdjustmentDoc{}, doc)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return idx, nil
}

func (repo StockAdjustmentRepository) Update(shopID string, guid string, doc models.StockAdjustmentDoc) error {
	filterDoc := map[string]interface{}{
		"shopid":    shopID,
		"guidfixed": guid,
	}
	err := repo.pst.UpdateOne(&models.StockAdjustmentDoc{}, filterDoc, doc)
	if err != nil {
		return err
	}
	return nil
}

func (repo StockAdjustmentRepository) Delete(shopID string, guid string, username string) error {
	err := repo.pst.SoftDelete(&models.StockAdjustmentDoc{}, username, bson.M{"guidfixed": guid, "shopid": shopID})
	if err != nil {
		return err
	}
	return nil
}

func (repo StockAdjustmentRepository) FindByGuid(shopID string, guid string) (models.StockAdjustmentDoc, error) {
	doc := &models.StockAdjustmentDoc{}
	err := repo.pst.FindOne(&models.StockAdjustmentDoc{}, bson.M{"shopid": shopID, "guidfixed": guid, "deletedat": bson.M{"$exists": false}}, doc)
	if err != nil {
		return *doc, err
	}
	return *doc, nil
}

func (repo StockAdjustmentRepository) FindPage(shopID string, q string, page int, limit int) ([]models.StockAdjustmentInfo, paginate.PaginationData, error) {

	docList := []models.StockAdjustmentInfo{}
	pagination, err := repo.pst.FindPage(&models.StockAdjustmentInfo{}, limit, page, bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"guidfixed": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &docList)

	if err != nil {
		return []models.StockAdjustmentInfo{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo StockAdjustmentRepository) FindItemsByGuidPage(shopID string, guid string, q string, page int, limit int) ([]models.StockAdjustmentInfo, paginate.PaginationData, error) {

	docList := []models.StockAdjustmentInfo{}
	pagination, err := repo.pst.FindPage(&models.StockAdjustment{}, limit, page, bson.M{
		"shopid":    shopID,
		"guidfixed": guid,
		"deletedat": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"items.itemSku": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &docList)

	if err != nil {
		return []models.StockAdjustmentInfo{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}
