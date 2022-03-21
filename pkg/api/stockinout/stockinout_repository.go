package stockinout

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IStockInOutRepository interface {
	Create(doc models.StockInOut) (primitive.ObjectID, error)
	Update(guid string, doc models.StockInOut) error
	Delete(guid string, shopID string) error
	FindByGuid(guid string, shopID string) (models.StockInOut, error)
	FindPage(shopID string, q string, page int, limit int) ([]models.StockInOut, paginate.PaginationData, error)
	FindItemsByGuidPage(guid string, shopID string, q string, page int, limit int) ([]models.StockInOut, paginate.PaginationData, error)
}

type StockInOutRepository struct {
	pst microservice.IPersisterMongo
}

func NewStockInOutRepository(pst microservice.IPersisterMongo) StockInOutRepository {
	return StockInOutRepository{
		pst: pst,
	}
}

func (repo StockInOutRepository) Create(doc models.StockInOut) (primitive.ObjectID, error) {
	idx, err := repo.pst.Create(&models.StockInOut{}, doc)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return idx, nil
}

func (repo StockInOutRepository) Update(guid string, doc models.StockInOut) error {
	err := repo.pst.UpdateOne(&models.StockInOut{}, "guidFixed", guid, doc)
	if err != nil {
		return err
	}
	return nil
}

func (repo StockInOutRepository) Delete(guid string, shopID string) error {
	err := repo.pst.SoftDelete(&models.StockInOut{}, bson.M{"guidFixed": guid, "shopID": shopID})
	if err != nil {
		return err
	}
	return nil
}

func (repo StockInOutRepository) FindByGuid(guid string, shopID string) (models.StockInOut, error) {
	doc := &models.StockInOut{}
	err := repo.pst.FindOne(&models.StockInOut{}, bson.M{"shopID": shopID, "guidFixed": guid, "deleted": false}, doc)
	if err != nil {
		return *doc, err
	}
	return *doc, nil
}

func (repo StockInOutRepository) FindPage(shopID string, q string, page int, limit int) ([]models.StockInOut, paginate.PaginationData, error) {

	docList := []models.StockInOut{}
	pagination, err := repo.pst.FindPage(&models.StockInOut{}, limit, page, bson.M{
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
		return []models.StockInOut{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo StockInOutRepository) FindItemsByGuidPage(guid string, shopID string, q string, page int, limit int) ([]models.StockInOut, paginate.PaginationData, error) {

	docList := []models.StockInOut{}
	pagination, err := repo.pst.FindPage(&models.StockInOut{}, limit, page, bson.M{
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
		return []models.StockInOut{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}
