package stockinout

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/transaction/stockinout/models"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IStockInOutRepository interface {
	Create(doc models.StockInOutDoc) (primitive.ObjectID, error)
	Update(shopID string, guid string, doc models.StockInOutDoc) error
	Delete(shopID string, guid string, username string) error
	FindByGuid(shopID string, guid string) (models.StockInOutDoc, error)
	FindPage(shopID string, pageable micromodels.Pageable) ([]models.StockInOutInfo, mongopagination.PaginationData, error)
	FindItemsByGuidPage(shopID string, guid string, pageable micromodels.Pageable) ([]models.StockInOutInfo, mongopagination.PaginationData, error)
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

func (repo StockInOutRepository) Update(shopID string, guid string, doc models.StockInOutDoc) error {
	filterDoc := map[string]interface{}{
		"shopid":    shopID,
		"guidfixed": guid,
	}
	err := repo.pst.UpdateOne(&models.StockInOutDoc{}, filterDoc, doc)
	if err != nil {
		return err
	}
	return nil
}

func (repo StockInOutRepository) Delete(shopID string, guid string, username string) error {
	err := repo.pst.SoftDelete(&models.StockInOutDoc{}, username, bson.M{"guidfixed": guid, "shopid": shopID})
	if err != nil {
		return err
	}
	return nil
}

func (repo StockInOutRepository) FindByGuid(shopID string, guid string) (models.StockInOutDoc, error) {
	doc := &models.StockInOutDoc{}
	err := repo.pst.FindOne(&models.StockInOutDoc{}, bson.M{"shopid": shopID, "guidfixed": guid, "deletedat": bson.M{"$exists": false}}, doc)
	if err != nil {
		return *doc, err
	}
	return *doc, nil
}

func (repo StockInOutRepository) FindPage(shopID string, pageable micromodels.Pageable) ([]models.StockInOutInfo, mongopagination.PaginationData, error) {
	filterQueries := bson.M{
		"shopid":    shopID,
		"deletedat": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"guidfixed": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + pageable.Query + ".*",
				Options: "",
			}}},
		},
	}
	docList := []models.StockInOutInfo{}
	pagination, err := repo.pst.FindPage(&models.StockInOutInfo{}, filterQueries, pageable, &docList)

	if err != nil {
		return []models.StockInOutInfo{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo StockInOutRepository) FindItemsByGuidPage(shopID string, guid string, pageable micromodels.Pageable) ([]models.StockInOutInfo, mongopagination.PaginationData, error) {
	filterQueries := bson.M{
		"shopid":    shopID,
		"guidfixed": guid,
		"deletedat": bson.M{"$exists": false},
		"$or": []interface{}{
			bson.M{"items.itemSku": bson.M{"$regex": primitive.Regex{
				Pattern: ".*" + pageable.Query + ".*",
				Options: "",
			}}},
		},
	}
	docList := []models.StockInOutInfo{}
	pagination, err := repo.pst.FindPage(&models.StockInOutInfo{}, filterQueries, pageable, &docList)

	if err != nil {
		return []models.StockInOutInfo{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}
