package purchase

import (
	"smlcloudplatform/internal/microservice"
	"smlcloudplatform/pkg/models"

	paginate "github.com/gobeam/mongo-go-pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IPurchaseRepository interface {
	Create(doc models.PurchaseDoc) (primitive.ObjectID, error)
	Update(shopID string, guid string, doc models.PurchaseDoc) error
	Delete(shopID string, guid string, username string) error
	FindByGuid(shopID string, guid string) (models.PurchaseDoc, error)
	FindPage(shopID string, q string, page int, limit int) ([]models.PurchaseInfo, paginate.PaginationData, error)
	FindItemsByGuidPage(guid string, shopID string, q string, page int, limit int) ([]models.PurchaseInfo, paginate.PaginationData, error)
}

type PurchaseRepository struct {
	pst microservice.IPersisterMongo
}

func NewPurchaseRepository(pst microservice.IPersisterMongo) PurchaseRepository {
	return PurchaseRepository{
		pst: pst,
	}
}

func (repo PurchaseRepository) Create(doc models.PurchaseDoc) (primitive.ObjectID, error) {
	idx, err := repo.pst.Create(&models.PurchaseDoc{}, doc)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return idx, nil
}

func (repo PurchaseRepository) Update(shopID string, guid string, doc models.PurchaseDoc) error {
	filterDoc := map[string]interface{}{
		"shopid":    shopID,
		"guidfixed": guid,
	}
	err := repo.pst.UpdateOne(&models.PurchaseDoc{}, filterDoc, doc)
	if err != nil {
		return err
	}
	return nil
}

func (repo PurchaseRepository) Delete(shopID string, guid string, username string) error {
	err := repo.pst.SoftDelete(&models.PurchaseDoc{}, username, bson.M{"guidfixed": guid, "shopid": shopID})
	if err != nil {
		return err
	}
	return nil
}

func (repo PurchaseRepository) FindByGuid(shopID string, guid string) (models.PurchaseDoc, error) {
	doc := &models.PurchaseDoc{}
	err := repo.pst.FindOne(&models.PurchaseDoc{}, bson.M{"shopid": shopID, "guidfixed": guid, "deletedat": bson.M{"$exists": false}}, doc)
	if err != nil {
		return *doc, err
	}
	return *doc, nil
}

func (repo PurchaseRepository) FindPage(shopID string, q string, page int, limit int) ([]models.PurchaseInfo, paginate.PaginationData, error) {

	docList := []models.PurchaseInfo{}
	pagination, err := repo.pst.FindPage(&models.PurchaseInfo{}, limit, page, bson.M{
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
		return []models.PurchaseInfo{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}

func (repo PurchaseRepository) FindItemsByGuidPage(guid string, shopID string, q string, page int, limit int) ([]models.PurchaseInfo, paginate.PaginationData, error) {

	docList := []models.PurchaseInfo{}
	pagination, err := repo.pst.FindPage(&models.PurchaseInfo{}, limit, page, bson.M{
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
		return []models.PurchaseInfo{}, paginate.PaginationData{}, err
	}

	return docList, pagination, nil
}
