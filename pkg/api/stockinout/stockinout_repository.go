package stockinout

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IStockInOutRepository interface {
	Create(doc models.StockInOutDoc) (primitive.ObjectID, error)
	Update(guid string, doc models.StockInOutDoc) error
	Delete(guid string, shopID string, username string) error
	FindByGuid(guid string, shopID string) (models.StockInOutDoc, error)
	FindPage(shopID string, q string, page int, limit int) ([]models.StockInOutInfo, paginate.PaginationData, error)
	FindItemsByGuidPage(guid string, shopID string, q string, page int, limit int) ([]models.StockInOutInfo, paginate.PaginationData, error)
}

type StockInOutRepository struct {
	pst microservice.IPersisterMongo
}

func NewStockInOutRepository(pst microservice.IPersisterMongo) StockInOutRepository {
	return StockInOutRepository{
		pst: pst,
	}
}

func (repo StockInOutRepository) Create(doc models.StockInOutDoc) (primitive.ObjectID, error) {
	idx, err := repo.pst.Create(&models.StockInOutDoc{}, doc)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return idx, nil
}

func (repo StockInOutRepository) Update(guid string, doc models.StockInOutDoc) error {
	err := repo.pst.UpdateOne(&models.StockInOutDoc{}, "guidFixed", guid, doc)
	if err != nil {
		return err
	}
	return nil
}

func (repo StockInOutRepository) Delete(guid string, shopID string, username string) error {
	err := repo.pst.SoftDelete(&models.StockInOutDoc{}, username, bson.M{"guidFixed": guid, "shopID": shopID})
	if err != nil {
		return err
	}
	return nil
}

func (repo StockInOutRepository) FindByGuid(guid string, shopID string) (models.StockInOutDoc, error) {
	doc := &models.StockInOutDoc{}
	err := repo.pst.FindOne(&models.StockInOutDoc{}, bson.M{"shopID": shopID, "guidFixed": guid, "deletedAt": bson.M{"$exists": false}}, doc)
	if err != nil {
		return *doc, err
	}
	return *doc, nil
}

func (repo StockInOutRepository) FindPage(shopID string, q string, page int, limit int) ([]models.StockInOutInfo, paginate.PaginationData, error) {

	docList := []models.StockInOutInfo{}
	pagination, err := repo.pst.FindPage(&models.StockInOutInfo{}, limit, page, bson.M{
		"shopID":    shopID,
		"deletedAt": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"guidFixed": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &docList)

	if err != nil {
		return []models.StockInOutInfo{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo StockInOutRepository) FindItemsByGuidPage(guid string, shopID string, q string, page int, limit int) ([]models.StockInOutInfo, paginate.PaginationData, error) {

	docList := []models.StockInOutInfo{}
	pagination, err := repo.pst.FindPage(&models.StockInOutInfo{}, limit, page, bson.M{
		"shopID":    shopID,
		"guidFixed": guid,
		"deletedAt": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"items.itemSku": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + q + ".*",
				Options: "",
			}}},
		},
	}, &docList)

	if err != nil {
		return []models.StockInOutInfo{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}
