package stockadjustment

import (
	"smlcloudplatform/internal/microservice"
	micromodels "smlcloudplatform/internal/microservice/models"
	"smlcloudplatform/pkg/transaction/stockadjustment/models"

	"github.com/userplant/mongopagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IStockAdjustmentRepository interface {
	Create(doc models.StockAdjustmentDoc) (primitive.ObjectID, error)
	Update(shopID string, guid string, doc models.StockAdjustmentDoc) error
	Delete(shopID string, guid string, username string) error
	FindByGuid(shopID string, guid string) (models.StockAdjustmentDoc, error)
	FindPage(shopID string, pageable micromodels.Pageable) ([]models.StockAdjustmentInfo, mongopagination.PaginationData, error)
	FindItemsByGuidPage(shopID string, guid string, pageable micromodels.Pageable) ([]models.StockAdjustmentInfo, mongopagination.PaginationData, error)
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

func (repo StockAdjustmentRepository) FindPage(shopID string, pageable micromodels.Pageable) ([]models.StockAdjustmentInfo, mongopagination.PaginationData, error) {
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
	docList := []models.StockAdjustmentInfo{}
	pagination, err := repo.pst.FindPage(&models.StockAdjustmentInfo{}, filterQueries, pageable, &docList)

	if err != nil {
		return []models.StockAdjustmentInfo{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo StockAdjustmentRepository) FindItemsByGuidPage(shopID string, guid string, pageable micromodels.Pageable) ([]models.StockAdjustmentInfo, mongopagination.PaginationData, error) {
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
	docList := []models.StockAdjustmentInfo{}
	pagination, err := repo.pst.FindPage(&models.StockAdjustment{}, filterQueries, pageable, &docList)

	if err != nil {
		return []models.StockAdjustmentInfo{}, mongopagination.PaginationData{}, err
	}

	return docList, pagination, nil
}
